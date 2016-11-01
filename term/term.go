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

type Atom string

func (a Atom) String() string {
	return fmt.Sprintf("(atom %s)", string(a))
}

func (Atom) isTerm() {}

func NewAtom(s string) Term {
	return Atom(s)
}

type Number struct {
	n *big.Float
}

func (n *Number) String() string {
	return fmt.Sprintf("(number %v)", n.n)
}

func (*Number) isTerm() {}

func NewNumber(n *big.Float) Term {
	return &Number{n}
}

type TermList []Term

func (ts TermList) String() string {
	ss := []string{}
	for _, t := range ts {
		ss = append(ss, fmt.Sprintf("%v", t))
	}
	return strings.Join(ss, ",")
}

type Callable struct {
	fn   string
	args []Term
}

func (c *Callable) String() string {
	return fmt.Sprintf("(%q/%d %s)", c.fn, len(c.args), fmt.Sprintf("%s", c.args))
}

func (*Callable) isTerm() {}

func (c *Callable) Functor() (string, int) {
	return c.fn, len(c.args)
}

func (c *Callable) Args() []Term {
	return c.args
}

func NewCallable(fn string, args []Term) Term {
	return &Callable{fn, args}
}

type Variable string

func (v Variable) String() string {
	return fmt.Sprintf("(var %s)", string(v))
}

func (Variable) isTerm() {}

func NewVariable(vn string) Term {
	return Variable(vn)
}
