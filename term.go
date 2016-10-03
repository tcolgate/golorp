package golorp

import "math/big"

type Atom string

type Cell interface{}

type Ref interface {
	Cell
	Ptr() Cell
}

type Structure interface {
	Cell
	Functor() Atom
	Arity() int
	Subterms() []Cell
}

type Number interface {
	Cell
	Value() big.Rat
}
