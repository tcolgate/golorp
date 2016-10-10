package term

import (
	"fmt"
	"math/big"
	"strings"
)

type Term interface {
	String() string
	isTerm()
}

type number struct {
	n *big.Float
}

func (n *number) String() string {
	return fmt.Sprintf("(number %v)", n.n)
}

func (*number) isTerm() {}

func NewNumber(n *big.Float) Term {
	return &number{n}
}

type TermList []Term

func (ts TermList) String() string {
	ss := []string{}
	for _, t := range ts {
		ss = append(ss, fmt.Sprintf("%v", t))
	}
	return strings.Join(ss, ",")
}

type callable struct {
	fn   string
	args []Term
}

func (c *callable) String() string {
	return fmt.Sprintf("(%q/%d %s)", c.fn, len(c.args), fmt.Sprintf("%s", c.args))
}

func (*callable) isTerm() {}

func (c *callable) Functor() string {
	return c.fn
}

func (c *callable) Arity() int {
	return len(c.args)
}

func (c *callable) Args() []Term {
	return c.args
}

func NewCallable(fn string, args []Term) Term {
	return &callable{fn, args}
}

type variable struct {
	vn string
}

func (v *variable) String() string {
	return fmt.Sprintf("(var %v)", v.vn)
}

func (*variable) isTerm() {}

func NewVariable(vn string) Term {
	return &variable{vn}
}
