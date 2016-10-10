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
	"log"

	"github.com/tcolgate/golorp/scan"
)

// This code owes a lot to golog

type Term struct {
}

func (p *Parser) NextTerm() (*Term, error) {
	t, err := p.readTerm(1200)
	if err != nil {
		return nil, err
	}

	nl := p.next()
	if nl.Type != scan.Stop {
		return nil, fmt.Errorf("unterminated term")
	}
	return t, nil
}

// readTerm reads the next available term, assuming the
// priority of the current term is prio
// T: 'var' |
//    'fun' |
//    'fun' '(' B ')' |
//    'fun' T |
//    T 'fun' |
//    T 'fun' T |
//    '(' T ')'
// B: T | T ',' B
//
// T: 'var' |
//    'atom' |
//    'fun' '(' B ')' |
//    'prefixop' T |
//    T 'postfixop' |
//    T 'infixop' T |
//    '(' T ')'
//    '[]'
//    '[' L ']'
// B: T | T ',' B
func (p *Parser) readTerm(pri int) (*Term, error) {
	log.Println("reading term")
Loop:
	for {
		l := p.next()
		switch l.Type {
		case scan.EOF:
			break Loop
		case scan.Comment:
			continue

		case scan.Variable:
			return p.readRest(0, pri, &Term{})

		case scan.Number:
			return p.readRest(0, pri, &Term{})

		case scan.Atom, scan.SpecialAtom, scan.Comma:
			opp, argp, ok := p.operators.Prefix(l.Text)
			if ok && opp <= pri {
				t0, err := p.readTerm(argp)
				if err == nil {
					return p.readRest(opp, pri, t0)
				}
			}

			return p.readRest(0, pri, &Term{})

		case scan.FunctorAtom:
			tb := p.next()
			if tb.Type != scan.LeftParen {
				panic(fmt.Errorf("functor atom without leftParen should be impossible"))
			}

			fargs, err := p.readFunctorArgs()
			if err != nil {
				return nil, err
			}

			tb = p.peek()
			if tb.Type != scan.RightParen {
				return nil, fmt.Errorf("Unterminated functor arguments")
			}
			p.next() // discard ')'
			log.Printf("args: %v", fargs)
			return p.readRest(0, pri, &Term{})

		case scan.LeftParen:
			t0, err := p.readTerm(1200)
			if err != nil {
				return nil, err
			}

			t1 := p.peek()
			if t1.Type != scan.RightParen {
				return nil, fmt.Errorf("Unterminated right parenthesis")
			}
			p.next() // discard ')'
			return p.readRest(0, pri, t0)
		case scan.EmptyList:
		case scan.LeftBrack:
		case scan.Unbound:
			return p.readRest(0, pri, &Term{})
		default:
			return nil, fmt.Errorf("syntax error")
		}
	}

	return nil, fmt.Errorf("premaature end of stream")
}

// readRest reads the remaining terms.
// restTerm lt ->
//   postfixTerm restTerm
//   infixTerm term restTerm
func (p *Parser) readRest(lpri, pri int, lt *Term) (*Term, error) {
	l := p.peek()
	switch l.Type {
	case scan.Atom, scan.Comma, scan.SpecialAtom:
		loppri, oppri, roppri, ok := p.operators.Infix(l.Text)
		if ok && pri >= oppri && lpri <= loppri {
			log.Printf("INFIX %#v %#v %#v %#v %#v", l.Text, lpri, pri, oppri, loppri)
			p.next() // consume the token
			t0, err := p.readTerm(roppri)
			if err != nil {
				return nil, err
			}
			return p.readRest(oppri, pri, t0)
		}
		opppri, argpri, ok := p.operators.Postfix(l.Text)
		if ok && opppri <= pri && lpri <= argpri {
			log.Printf("POSTFIX %#v %#v %#v %#v %#v", l.Text, lpri, pri, oppri, loppri)
			p.next() // consume the token
			t0 := &Term{}
			return p.readRest(opppri, pri, t0)
		}
	}
	return lt, nil
}

func (p *Parser) readListItems() (*Term, error) {
	return &Term{}, nil
}

func (p *Parser) readFunctorArgs() ([]*Term, error) {
	log.Println("Reading args")
	fargs := []*Term{}

	for {
		lt := p.peek()
		if lt.Type == scan.RightParen {
			break
		}
		log.Println("Reading arg")

		t, err := p.readTerm(999)
		if err != nil {
			log.Println("Got error from readTerm")
			return nil, err
		}
		fargs = append(fargs, t)
		lt = p.peek()
		log.Printf("peeked %#v\n", lt)
		if lt.Type == scan.RightParen {
			break
		}
		if lt.Type != scan.Comma {
			return nil, fmt.Errorf("invalid functor arguments, tok %#v", lt)
		}
		// discard comma
		p.next()
	}

	return fargs, nil
}
