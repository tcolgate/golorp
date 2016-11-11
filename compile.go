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
	code := []CodeCell{}

	seen := map[int]bool{}
	ts := make(chan streamToken)
	go flattenq(ts, q)

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

	for _, i := range code {
		fmt.Printf("Q CODE %#v\n", i.string)
	}

	code = []CodeCell{}
	seen = map[int]bool{}
	ts = make(chan streamToken)
	go flattenp(ts, p)

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

	for _, i := range code {
		fmt.Printf("P CODE %#v\n", i.string)
	}
	return code
}

func flattenq(ts chan<- streamToken, t term.Term) {
	defer close(ts)

	regs := map[int]term.Term{}
	invregs := map[term.Term]int{}
	vars := map[term.Variable]int{}

	var assign func(p term.Term)
	assign = func(p term.Term) {
		switch t := p.(type) {
		case *term.Callable:
			for _, at := range t.Args() {
				switch t := at.(type) {
				case term.Variable:
					if xi, ok := vars[t]; !ok {
						regs[len(regs)] = t
						vars[t] = len(regs) - 1
						invregs[at] = len(regs) - 1
						continue
					} else {
						invregs[at] = xi
					}
				case *term.Callable:
					regs[len(regs)] = t
					invregs[at] = len(regs) - 1
				default:
				}
			}

		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	var output func(p term.Term)
	output = func(p term.Term) {
		xi, ok := invregs[p]
		if !ok {
			panic("unknown term")
		}
		switch t := p.(type) {
		case *term.Callable:
			fn, argc := t.Functor()
			ts <- streamToken{
				fn: fmt.Sprintf("%s/%d ", fn, argc),
				xi: xi,
			}
			for _, at := range t.Args() {
				xi, ok := invregs[at]
				if !ok {
					panic("unknown term")
				}
				ts <- streamToken{
					xi: xi,
				}
			}
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	regs[len(regs)] = t
	invregs[t] = len(regs) - 1
	term.WalkDepthFirst(assign, output, t)
}

func flattenp(ts chan<- streamToken, t term.Term) {
	defer close(ts)

	regs := map[int]term.Term{}
	invregs := map[term.Term]int{}
	vars := map[term.Variable]int{}

	var assign func(p term.Term)
	assign = func(p term.Term) {
		switch t := p.(type) {
		case *term.Callable:
			for _, at := range t.Args() {
				switch t := at.(type) {
				case term.Variable:
					if xi, ok := vars[t]; !ok {
						regs[len(regs)] = t
						vars[t] = len(regs) - 1
						invregs[at] = len(regs) - 1
						continue
					} else {
						invregs[at] = xi
					}
				case *term.Callable:
					regs[len(regs)] = t
					invregs[at] = len(regs) - 1
				default:
				}
			}
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	var output func(p term.Term)
	output = func(p term.Term) {
		xi, ok := invregs[p]
		if !ok {
			panic("unknown term")
		}
		switch t := p.(type) {
		case *term.Callable:
			fn, argc := t.Functor()
			ts <- streamToken{
				fn: fmt.Sprintf("%s/%d ", fn, argc),
				xi: xi,
			}
			for _, at := range t.Args() {
				xi, ok := invregs[at]
				if !ok {
					panic("unknown term")
				}
				ts <- streamToken{
					xi: xi,
				}
			}
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	regs[len(regs)] = t
	invregs[t] = len(regs) - 1
	term.WalkDepthFirst(assign, nil, t)

	term.WalkDepthFirst(output, nil, t)
}
