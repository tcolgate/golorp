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
func compileL0(q, p term.Term) CodeCells {
	qcode := compileL0Query(q)
	pcode := compileL0Program(p)

	return CodeCells(append(qcode, pcode...))
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

func (ctx *tokeningCtx) assignReg(t term.Term) int {
	next := len(ctx.regs)
	ctx.regs[next] = t
	ctx.invregs[t] = next

	return next
}

func flattenq(ts chan<- streamToken, t term.Term) {
	defer close(ts)

	ctx := newTokeningCtx(ts)

	ctx.assignReg(t)
	term.WalkDepthFirst(ctx.assign, ctx.tokenize, t)
}

func flattenp(ts chan<- streamToken, t term.Term) {
	defer close(ts)

	ctx := newTokeningCtx(ts)

	ctx.assignReg(t)
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
					xi := ctx.assignReg(at)
					ctx.vars[t] = xi
					continue
				} else {
					ctx.invregs[at] = xi
				}
			case *term.Callable:
				ctx.assignReg(at)
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
