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
}

// IsCell marks FuncCell as a valid heap Cell
func (FuncCell) IsCell() {
}

func (c FuncCell) String() string {
	return fmt.Sprintf("%s", c.Atom)
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
		str += fmt.Sprintf("X%d  = %s\n", i, c)
	}
	return str
}

// Machine hods the state of our WAM
type Machine struct {
	// M0
	Heap       []Cell
	XRegisters []Cell
	Mode       InstructionMode
	HReg       int
	SReg       int

	// M1
	Code   []Cell
	Labels map[string]int

	PDL  []*Cell
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
	str := "HEAP:\n"
	str += fmt.Sprintf("%s\n", HeapCells(m.Heap))
	str += "X Registers:\n"
	str += fmt.Sprintf("%s\n", RegCells(m.XRegisters))
	str += fmt.Sprintf("H: %d S: %d\n", m.HReg, m.SReg)

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
	}
}

// I0 - M0 insutrctions for L0

func PutStructure(fn term.Atom, xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		m.Heap[m.HReg] = StrCell{m.HReg + 1}
		m.Heap[m.HReg+1] = FuncCell{fn}
		m.XRegisters[xi] = m.Heap[m.HReg]
		m.HReg = m.HReg + 2
		return nil, ""
	}, fmt.Sprintf("put_structure %s X%d", fn, xi)
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

func GetStructure(fn term.Atom, xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("get_structure %s X%d", fn, xi)
}

func UnifyVariable(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("unify_variable X%d", xi)
}

func UnifyValue(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
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

func deref(m *Machine, xi int) int {
	for {
		switch c := m.Heap[xi].(type) {
		case RefCell:
			if c.Ptr == xi { // Ref Cell points to itself, unbound
				return xi
			}
			xi = c.Ptr
		default:
			return xi
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
