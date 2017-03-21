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
	"testing"

	"github.com/tcolgate/golorp/context"
	"github.com/tcolgate/golorp/parse"
	"github.com/tcolgate/golorp/scan"
	"github.com/tcolgate/golorp/term"
)

type testL0 struct {
	name string
	q    string
	p    string
	res  string
}

var testsL0 = []testL0{
	{"query0",
		`p(Z,h(Z,W),f(W)).`,
		`p(f(X),h(Y,f(a)),Y).`,
		`put_structure (atom h)/2 X2
set_variable X1
set_variable X4
put_structure (atom f)/1 X3
set_value X4
put_structure (atom p)/3 X0
set_value X1
set_value X2
set_value X3
get_structure (atom p)/3 X0
unify_variable X1
unify_variable X2
unify_variable X3
get_structure (atom f)/1 X1
unify_variable X4
get_structure (atom h)/2 X2
unify_value X3
unify_variable X5
get_structure (atom f)/1 X5
unify_variable X6
get_structure (atom a)/0 X6
`,
	},
}

func TestCompileL0(t *testing.T) {
	var ctx context.Context
	for _, st := range testsL0 {
		t.Run(st.name, func(t *testing.T) {
			s := scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.q)))
			qt := parse.New("file.pl", s)

			s = scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.p)))
			pt := parse.New("file.pl", s)

			q, _ := qt.NextTerm()
			p, _ := pt.NextTerm()

			cs := compileL1(q, []term.Term{p})
			if cs.String() != st.res {
				t.Fatalf("expected: %s, got: %s", st.res, cs)
			}
		})
	}
}
