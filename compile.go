package golorp

import (
	"fmt"

	"github.com/tcolgate/golorp/term"
)

// Compile a single query, and a program
func compileL0(q, p term.Term) []CodeCell {
	code := []CodeCell{}
	qregs := assignReg(q)

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
					if n, ok := seen[t]; ok {
						inst, str := SetValue(n)
						code = append(code, CodeCell{inst, str})
						continue
					}
					seen[t] = i
					inst, str := SetVariable(i)
					code = append(code, CodeCell{inst, str})
				case *term.Callable:
					inst, str := PutStructure(t.Functor())
					code = append(code, CodeCell{inst, str})
				}
			}
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	return code
}

func assignReg(t term.Term) map[int]term.Term {
	regs := map[int]term.Term{}
	vars := map[term.Variable]int{}

	var assign func(p term.Term)
	assign = func(p term.Term) {
		switch t := p.(type) {
		case *term.Callable:
			for _, at := range t.Args() {
				switch t := at.(type) {
				case term.Variable:
					if _, ok := vars[t]; !ok {
						regs[len(regs)] = t
						vars[t] = len(regs) - 1
					}
				case *term.Callable:
					regs[len(regs)] = t
				default:
				}
			}

			for _, at := range t.Args() {
				switch t := at.(type) {
				case *term.Callable:
					term.WalkDepthFirst(assign, t)
				default:
				}
			}

		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	regs[len(regs)] = t
	term.WalkDepthFirst(assign, t)

	return regs
}
