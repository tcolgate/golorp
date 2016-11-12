package golorp

import (
	"fmt"

	"github.com/tcolgate/golorp/term"
)

type streamToken struct {
	fn string // functor name
	vn string // variable name, for named vars
	xi int    // The register to use
}

// Compile a single query, and a program
func compileL0(q, p term.Term) []CodeCell {
	qcode := compileL0Query(q)
	for _, i := range qcode {
		fmt.Printf("Q CODE %#v\n", i.string)
	}

	pcode := compileL0Program(p)
	for _, i := range pcode {
		fmt.Printf("P CODE %#v\n", i.string)
	}

	return append(qcode, pcode...)
}

// Compile a single query, and a program
func compileL0Query(t term.Term) []CodeCell {
	code := []CodeCell{}

	seen := map[int]bool{}
	ts := make(chan streamToken)
	go flattenq(ts, t)

	for ft := range ts {
		switch {
		case ft.fn != "":
			inst, str := PutStructure(term.Atom(ft.fn), ft.xi)
			seen[ft.xi] = true
			code = append(code, CodeCell{inst, str})
		case ft.fn == "":
			if _, ok := seen[ft.xi]; ok {
				inst, str := SetValue(ft.xi)
				code = append(code, CodeCell{inst, str})
				continue
			}
			seen[ft.xi] = true
			inst, str := SetVariable(ft.xi)
			code = append(code, CodeCell{inst, str})
		default:
			panic("unknown term m0 type")
		}
	}

	return code
}

// Compile a single l0 program term
func compileL0Program(t term.Term) []CodeCell {
	code := []CodeCell{}

	seen := map[int]bool{}
	ts := make(chan streamToken)
	go flattenp(ts, t)

	for ft := range ts {
		switch {
		case ft.fn != "":
			inst, str := GetStructure(term.Atom(ft.fn), ft.xi)
			seen[ft.xi] = true
			code = append(code, CodeCell{inst, str})
		case ft.fn == "":
			if _, ok := seen[ft.xi]; ok {
				inst, str := UnifyValue(ft.xi)
				code = append(code, CodeCell{inst, str})
				continue
			}
			seen[ft.xi] = true
			inst, str := UnifyVariable(ft.xi)
			code = append(code, CodeCell{inst, str})
		default:
			panic("unknown term m0 type")
		}
	}

	return code
}

type tokeningCtx struct {
	regs    map[int]term.Term
	invregs map[term.Term]int
	vars    map[term.Variable]int
	ts      chan<- streamToken
}

func newTokeningCtx(ts chan<- streamToken) *tokeningCtx {
	return &tokeningCtx{
		regs:    map[int]term.Term{},
		invregs: map[term.Term]int{},
		vars:    map[term.Variable]int{},
		ts:      ts,
	}
}

func flattenq(ts chan<- streamToken, t term.Term) {
	defer close(ts)

	ctx := newTokeningCtx(ts)

	ctx.regs[len(ctx.regs)] = t
	ctx.invregs[t] = len(ctx.regs) - 1

	term.WalkDepthFirst(ctx.assign, ctx.tokenize, t)
}

func flattenp(ts chan<- streamToken, t term.Term) {
	defer close(ts)

	ctx := newTokeningCtx(ts)

	ctx.regs[len(ctx.regs)] = t
	ctx.invregs[t] = len(ctx.regs) - 1
	term.WalkDepthFirst(ctx.assign, nil, t)

	term.WalkDepthFirst(ctx.tokenize, nil, t)
}

func (ctx *tokeningCtx) assign(p term.Term) {
	switch t := p.(type) {
	case *term.Callable:
		for _, at := range t.Args() {
			switch t := at.(type) {
			case term.Variable:
				if xi, ok := ctx.vars[t]; !ok {
					ctx.regs[len(ctx.regs)] = t
					ctx.vars[t] = len(ctx.regs) - 1
					ctx.invregs[at] = len(ctx.regs) - 1
					continue
				} else {
					ctx.invregs[at] = xi
				}
			case *term.Callable:
				ctx.regs[len(ctx.regs)] = t
				ctx.invregs[at] = len(ctx.regs) - 1
			default:
			}
		}
	case term.Variable:
	default:
		panic("unknown term m0 type")
	}
}

func (ctx *tokeningCtx) tokenize(p term.Term) {
	xi, ok := ctx.invregs[p]
	if !ok {
		panic("unknown term")
	}
	switch t := p.(type) {
	case *term.Callable:
		fn, argc := t.Functor()
		ctx.ts <- streamToken{
			fn: fmt.Sprintf("%s/%d ", fn, argc),
			xi: xi,
		}
		for _, at := range t.Args() {
			xi, ok := ctx.invregs[at]
			if !ok {
				panic("unknown term")
			}
			ctx.ts <- streamToken{
				xi: xi,
			}
		}
	case term.Variable:
	default:
		panic("unknown term m0 type")
	}
}
