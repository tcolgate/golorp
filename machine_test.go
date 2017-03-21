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

package golorp

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/tcolgate/golorp/context"
	"github.com/tcolgate/golorp/parse"
	"github.com/tcolgate/golorp/scan"
	"github.com/tcolgate/golorp/term"
)

type mtestL0 struct {
	q    string
	p    string
	fail bool
}

var mtestsL0 = []mtestL0{
	/*
		{
			`p(Z,h(Z,W),f(W)).`,
			`p(f(X),h(Y,f(a)),Y).`,
			false,
		},
		{
			`p(Z,h(Z,W)).`,
			`p(f(X),h(Y,f(a)),Y).`,
			true,
		},
		{
			`p(Z,h(Z,a())).`,
			`p(Z,h(Z,Z)).`,
			false,
		},
		{
			`p(A,h(A,a(),D)).`,
			`p(a(),h(Z,B,B)).`,
			false,
		},
		{
			`a().`,
			`a().`,
			false,
		},
		{
			`a().`,
			`b().`,
			true,
		},
		{
			`a().`,
			`X.`,
			false,
		},
		{
			`a(X,b()).`,
			`a(Y,b()).`,
			false,
		},
		{
			`a(X,b()).`,
			`a(c(),Y).`,
			false,
		},
		{
			`a(X,b(),X).`,
			`a(b(),Y,Y).`,
			false,
		},
		{
			`X.`,
			`Y.`,
			false,
		},
		{
			`a(X,f(Y,Y),Y).`,
			`a(X,X,f(Y,Y)).`,
			false,
		},
		{
			`a(f(Z,Z),f(Z,Z,Z)).`,
			`a(Y,Y).`,
			true,
		},
	*/
	{
		`p(Z,h(Z,W),f(W)).`,
		`p(f(X),h(Y,f(a)),Y).`,
		true,
	},
}

func TestMachine0(t *testing.T) {
	var ctx context.Context
	for i, st := range mtestsL0 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s := scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.q)))
			qt := parse.New("file.pl", s)

			s = scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.p)))
			pt := parse.New("file.pl", s)

			q, _ := qt.NextTerm()
			p, _ := pt.NextTerm()

			cs := compileL1(q, []term.Term{p})
			m := NewMachine()

			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					fmt.Println(m.String())
					return
				}
				if !m.Finished {
					t.Fatalf("test failed, did not finish")
				}
			}()

			m.run(cs)

			if st.fail != m.Failed {
				fmt.Println(m.String())
				t.Fatalf("%s failed, expected %v, got %v", t.Name(), st.fail, m.Failed)
			}

			fmt.Printf("%s", m)
		})
	}
}
