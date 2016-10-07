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

package parse

import (
	"fmt"
	"io"

	"github.com/tcolgate/golorp/scan"
)

// This code owes a lot to golog

type Term struct {
}

func (p *Parser) NextTerm() (*Term, error) {
	return nil, io.EOF
}

// readTerm reads the next available term, assuming the
// priority of the current term is prio
func (p *Parser) readTerm(pri int) (*Term, error) {
	for {
		l := p.next()
		switch l.Type {
		case scan.Comment:
			continue
		case scan.Atom:
			opp, argp, ok := p.operators.Prefix(l.Text)
			if ok && opp <= pri {
				t0, err := p.readTerm(argp)
				if err == nil {
					return p.readRest(opp, pri, t0)
				}
			}
			fallthrough
		case scan.Variable:
			return &Term{}, nil
		case scan.FunctorAtom:
			continue
		case scan.Stop:
			return &Term{}, nil
		default:
			return nil, fmt.Errorf("unknown token")
		}
	}
}

// readRest reads the remaining terms.
func (p *Parser) readRest(lpri, pri int, lt *Term) (*Term, error) {
Loop:
	for {
		l := p.next()
		switch l.Type {
		case scan.Comment:
			continue
		case scan.Stop:
			break Loop
		default:
			return nil, fmt.Errorf("unknown token")
		}
	}

	return nil, io.EOF
}
