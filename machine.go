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

type Cell interface {
	IsCell()
}

type RefCell interface {
	Cell
	Ptr() Cell
	Name() string
}

type StructurCell interface {
	Cell
	Ptr() Cell
}

// FunctorCell is not tagged in WAM-Book, but we need a type
type FunctorCell interface {
	Cell
	Functor() (term.Atom, int)
}

type Machine struct {
	// M0
	Heap       []Cell
	XRegisters []int
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

// I0 - M0 insutrctions for L0

func PutStructure(fn term.Atom, xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("put_structure %s X%d", fn, xi)
}

func SetVariable(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
		return nil, ""
	}, fmt.Sprintf("set_variable X%d", xi)
}

func SetValue(xi int) (machineFunc, string) {
	return func(m *Machine) (machineFunc, string) {
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

func deref(m *Machine, xi int) Cell {
	return nil
}

// L1
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
