package golorp

import (
	"fmt"
	"strings"

	"github.com/tcolgate/golorp/term"
)

// Compile a single query, and a program
func compileL0(q, p term.Term) []machineFunc {
	qregs := assignReg(q)

	fmt.Printf("REGS: %s\n", qregs)

	qfl := flattenBottomUp(qregs)

	fmt.Printf("FLAT: %s\n", qfl)

	insts := []string{}
	funs := map[*term.Callable]int{}
	vars := map[term.Variable]int{}

	for _, qt := range qfl {
		switch t := qt.(type) {
		case *term.Callable:
			insts = append(insts, fmt.Sprintf("STR %d", len(insts)+1))
			fn, arity := t.Functor()
			insts = append(insts, fmt.Sprintf("%s/%d", fn, arity))
			funs[t] = len(insts) - 1
			for _, at := range t.Args() {
				switch t := at.(type) {
				case term.Variable:
					if i, ok := vars[t]; ok {
						insts = append(insts, fmt.Sprintf("REF %v", i))
						continue
					}
					vars[t] = len(insts)
					insts = append(insts, fmt.Sprintf("REF %v", vars[t]))
				case *term.Callable:
					if _, ok := funs[t]; !ok {
						panic("Function not serialized yet")
					}
					insts = append(insts, fmt.Sprintf("STR %v", funs[t]))
				}
			}
		default:
			panic("unknown term m0 type")
		}
	}

	fmt.Println(strings.Join(insts, "\n"))

	return nil
}

func assignReg(t term.Term) []term.Term {
	regs := []term.Term{}
	vars := map[term.Variable]int{}

	assign := func(p term.Term) {
		switch t := p.(type) {
		case *term.Callable:
			for _, at := range t.Args() {
				switch t := at.(type) {
				case term.Variable:
					if _, ok := vars[t]; !ok {
						regs = append(regs, t)
						vars[t] = len(regs) - 1
					}
				case *term.Callable:
					regs = append(regs, t)
				default:
				}
			}
		case term.Variable:
		default:
			panic("unknown term m0 type")
		}
	}

	regs = append(regs, t)
	term.WalkDepthFirst(assign, t)

	return regs
}

func flattenBottomUp(ts []term.Term) []term.Term {
	out := []term.Term{}
	for i := len(ts) - 1; i >= 0; i-- {
		switch ts[i].(type) {
		case term.Variable:
		case *term.Callable:
			out = append(out, ts[i])
		default:
			panic("unknown term m0 type")
		}
	}
	return out
}
