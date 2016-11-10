package golorp

import (
	"fmt"

	"github.com/tcolgate/golorp/term"
)

// Compile a single query, and a program
func compileL0(q, p term.Term) []CodeCell {
	code := []CodeCell{}
	qregs, invqregs := assignReg(q)

	for i := 0; i < len(qregs); i++ {
		fmt.Printf("REGS: X%d = %s\n", i, qregs[i])
	}

	seen := map[term.Variable]int{}
	for i := len(qregs) - 1; i >= 0; i-- {
		qt := qregs[i]

		switch t := qt.(type) {
		case *term.Callable:
			fn, _ := t.Functor()
			inst, str := PutStructure(term.Atom(fn), i)
			code = append(code, CodeCell{inst, str})
			for _, at := range t.Args() {
				switch t := at.(type) {
				case term.Variable:
					var xi int
					var ok bool
					if xi, ok = invqregs[at]; !ok {
						panic("variable without assigned register")
					}
					if _, ok := seen[t]; ok {
						inst, str := SetValue(xi)
						code = append(code, CodeCell{inst, str})
						continue
					}
					seen[t] = i
					inst, str := SetVariable(xi)
					code = append(code, CodeCell{inst, str})
				case *term.Callable:
					fn, argc := t.Functor()
					inst, str := PutStructure(term.Atom(fn), argc)
					code = append(code, CodeCell{inst, str})
				}
			}
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	for _, i := range code {
		fmt.Printf("CODE %#v\n", i.string)
	}
	return code
}

func assignReg(t term.Term) (map[int]term.Term, map[term.Term]int) {
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
			fmt.Printf("X%d = ", xi)

			fn, argc := t.Functor()
			fmt.Printf("%s/%d ", fn, argc)
			for _, at := range t.Args() {
				xi, ok := invregs[at]
				if !ok {
					panic("unknown term")
				}
				fmt.Printf("X%d ", xi)
			}

			fmt.Printf("\n")
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}

	}

	regs[len(regs)] = t
	invregs[t] = len(regs) - 1
	term.WalkDepthFirst(assign, output, t)

	return regs, invregs
}
