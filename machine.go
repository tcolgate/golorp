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
	Functor() (Atom, int)
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

type machineFunc func(*Machine) machineFunc

// I0 - M0 insutrctions for L0

func PutStructure(fn Atom, xi int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func SetVariable(xi int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func SetValue(xi int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func GetStructure(fn Atom, xi int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func UnifyVariable(xi int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func UnifyValue(xi int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func Unify(a1, a2 int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func Bind(a1, a2 int) machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func deref(m *Machine, xi int) Cell {
	return nil
}

// L1
func Call() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func Proceeed() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func PutVariable() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func PutValue() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func GetVariable() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func GetValue() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

// L2
func Allocate() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func Deallocate() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

// L3 - Prolog
func TryMeElse() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func RetryMeElse() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

func TrustMe() machineFunc {
	return func(m *Machine) machineFunc {
		return nil
	}
}

// Optimisations
