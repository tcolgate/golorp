// Copyright 2016 Tristan Colgate-McFarlane
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package golorp

import (
	"fmt"

	"github.com/tcolgate/golorp/term"
)

// Cell implements an interface for items that can be stored on the heap
type Cell interface {
	IsCell()
	fmt.Stringer
}

// RefCell is a a heap cell containing a reference to another cell
type RefCell struct {
	Ptr int
}

// IsCell marks RefCell as a valid heap Cell
func (RefCell) IsCell() {
}

func (c RefCell) String() string {
	return fmt.Sprintf("REF %d", c.Ptr)
}

// StrCell is a structure header cell
type StrCell struct {
	Ptr int
}

// IsCell marks StrCell as a valid heap Cell
func (StrCell) IsCell() {
}

func (c StrCell) String() string {
	return fmt.Sprintf("STR %d", c.Ptr)
}

// FuncCell is not tagged in WAM-Book, but we need a type
type FuncCell struct {
	Atom term.Atom
	n    int
}

// IsCell marks FuncCell as a valid heap Cell
func (FuncCell) IsCell() {
}

func (c FuncCell) String() string {
	return fmt.Sprintf("%s/%d", c.Atom, c.n)
}

// HeapCells is a utility type to format a slice of
// cells as a heap
type HeapCells []Cell

func (cs HeapCells) String() string {
	str := ""
	for i, c := range cs {
		str += fmt.Sprintf("%d %s\n", i, c)
	}
	return str
}

// RegCells is a utility type to format a slice of
// cells as a X registers
type RegCells []Cell

func (cs RegCells) String() string {
	str := ""
	for i, c := range cs {
		str += fmt.Sprintf("X%d = %s\n", i, c)
	}
	return str
}

type PDL []int

func (p PDL) isEmpty() bool {
	if len(p) > 0 {
		return false
	}
	return true
}

func (p PDL) push(a int) {
	p = append(p, a)
}

func (p PDL) pop() int {
	a := p[len(p)-1]
	p = p[:len(p)-1]
	return a
}

// Machine hods the state of our WAM
type Machine struct {
	// We use to stop on failure
	Finished bool
	Failed   bool
	// M0
	Heap       []Cell
	XRegisters []Cell
	Mode       InstructionMode
	HReg       int
	SReg       int
	PDL        PDL

	// M1
	Code   []Cell
	Labels map[string]int

	PReg int

	// M2
	AndStack []Environment
	EReg     int
	CPReg    int

	// M3 - Prolog
	OrStack []ChoicePoint
	Trail   []Cell
	BReg    int
	TRReg   int
	GBReg   int

	// Optimisations
}

func NewMachine() *Machine {
	return &Machine{
		Heap:       make([]Cell, 20),
		XRegisters: make([]Cell, 10),
	}
}

func (m *Machine) String() string {
	str := fmt.Sprintf("Finished %v Failed: %v\n\n", m.Finished, m.Failed)
	str += fmt.Sprintf("H: %d S: %d\n", m.HReg, m.SReg)
	str += fmt.Sprintf("Mode %v\n", m.Mode)
	str += "X Registers:\n"
	str += fmt.Sprintf("%s\n", RegCells(m.XRegisters))
	str += "Heap:\n"
	str += fmt.Sprintf("%s\n", HeapCells(m.Heap))
	return str
}

type Environment Cell
type ChoicePoint Cell

type Instruction int
type InstructionMode int

const (
	InvalidMode InstructionMode = iota // zero value not sure
	Read
	Write
)

type machineFunc func(*Machine) (machineFunc, string)
type CodeCell struct {
	fn machineFunc
	string
}

type CodeCells []CodeCell

func (cs CodeCells) String() string {
	str := ""
	for _, i := range cs {
		str += fmt.Sprintf("%s\n", i.string)
	}
	return str
}

func (m *Machine) run(cs []CodeCell) {
	for _, c := range cs {
		fmt.Printf("%s\n", c.string)
		c.fn(m)
		if m.Finished {
			return
		}
	}
	m.Finished = true
}

// I0 - M0 insutrctions for L0

func PutStructure(fn term.Atom, n, xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		m.Heap[m.HReg] = StrCell{m.HReg + 1}
		m.Heap[m.HReg+1] = FuncCell{fn, n}
		m.XRegisters[xi] = m.Heap[m.HReg]
		m.HReg = m.HReg + 2
		return nil, ""
	}, fmt.Sprintf("put_structure %s/%d X%d", fn, n, xi)
}

func SetVariable(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		m.Heap[m.HReg] = RefCell{m.HReg}
		m.XRegisters[xi] = m.Heap[m.HReg]
		m.HReg = m.HReg + 1
		return nil, ""
	}, fmt.Sprintf("set_variable X%d", xi)
}

func SetValue(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		m.Heap[m.HReg] = m.XRegisters[xi]
		m.HReg = m.HReg + 1
		return nil, ""
	}, fmt.Sprintf("set_value X%d", xi)
}

func GetStructure(fn term.Atom, n, xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		addr := m.deref(m.XRegisters[xi].(StrCell).Ptr)

		cc := m.Heap[addr]
		switch c := cc.(type) {
		case RefCell:
			m.Heap[m.HReg] = StrCell{m.HReg + 1}
			m.Heap[m.HReg+1] = FuncCell{fn, n}
			m.XRegisters[xi] = m.Heap[m.HReg]
			m.bind(addr, m.HReg)
			m.HReg = m.HReg + 2
			m.Mode = Write
		case StrCell:
			if tc, ok := m.Heap[c.Ptr].(FuncCell); ok && tc.Atom == fn {
				fmt.Printf("IN HERE 1 %v %v %v\n", ok, tc, fn)
				m.SReg = c.Ptr + 1
				m.Mode = Read
			} else {
				fmt.Printf("IN HERE 2 %v %v %v\n", ok, tc, fn)
				m.Finished = true
				m.Failed = true
			}
		default:
			fmt.Printf("IN HERE 3 %#v\n", c)
			m.Finished = true
			m.Failed = true
		}
		return nil, ""
	}, fmt.Sprintf("get_structure %s/%d X%d", fn, n, xi)
}

func UnifyVariable(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		switch m.Mode {
		case Read:
			m.XRegisters[xi] = m.Heap[m.SReg]
		case Write:
			m.Heap[m.HReg] = RefCell{m.HReg}
			m.XRegisters[xi] = m.Heap[m.HReg]
			m.HReg = m.HReg + 1
		default:
			panic(fmt.Errorf("invalid read/write mode %v", m.Mode))
		}
		return nil, ""
	}, fmt.Sprintf("unify_variable X%d", xi)
}

func UnifyValue(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		switch m.Mode {
		case Read:
			m.unify(xi, m.SReg)
		case Write:
			m.Heap[m.HReg] = m.XRegisters[xi]
			m.HReg = m.HReg + 1
		default:
			panic(fmt.Errorf("invalid read/write mode %v", m.Mode))
		}
		return nil, ""
	}, fmt.Sprintf("unify_value X%d", xi)
}

func Unify(a1, a2 int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("unify A%d, A%d", a1, a2)
}

func Bind(a1, a2 int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("bind A%d A%d", a1, a2)
}

func (m *Machine) deref(a int) int {
	for {
		switch c := m.Heap[a].(type) {
		case RefCell:
			if c.Ptr == a { // Ref Cell points to itself, unbound
				return a
			}
			a = c.Ptr
		default:
			return a
		}
	}
}

func (m *Machine) bind(a, b int) {
}

func (m *Machine) unify(a1, a2 int) {
	m.PDL.push(a1)
	m.PDL.push(a2)
	fail := false
	for !(m.PDL.isEmpty() || fail) {
		d1 := m.deref(m.PDL.pop())
		d2 := m.deref(m.PDL.pop())
		if d1 != d2 {
			t1, t2 := m.Heap[d1], m.Heap[d2]
			_, ok1 := t1.(RefCell)
			_, ok2 := t2.(RefCell)
			if ok1 || ok2 {
				m.bind(d1, d2)
			} else {
				v1, ok1 := t1.(StrCell)
				v2, ok2 := t2.(StrCell)
				if !(ok1 && !ok2) {
					panic("Wrong cell type")
				}
				f1, ok1 := m.Heap[v1.Ptr].(FuncCell)
				f2, ok2 := m.Heap[v2.Ptr].(FuncCell)
				if !(ok1 && !ok2) {
					panic("Wrong cell type")
				}
				if f1.Atom == f2.Atom && f1.n == f2.n {
					for i := 0; i < f1.n; i++ {
						m.PDL.push(v1.Ptr + i)
						m.PDL.push(v2.Ptr + i)
					}
				} else {
					fail = true
				}
			}
		}
	}
}

func Call() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("call")
}

func Proceeed() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("proceed")
}

func PutVariable() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("put_variable")
}

func PutValue() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("put_value")
}

func GetVariable() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("get_variable")
}

func GetValue() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("get_value")
}

// L2
func Allocate() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("allocate")
}

func Deallocate() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("deallocate")
}

// L3 - Prolog
func TryMeElse() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("try_me_else")
}

func RetryMeElse() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("retry_me_else")
}

func TrustMe() (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("trust_me")
}

// Optimisations
