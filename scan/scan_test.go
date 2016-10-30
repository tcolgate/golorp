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

package scan

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tcolgate/golorp/context"
)

type test struct {
	name string
	src  string
	exp  []Token
}

var tests = []test{
	{"single", "/* This is a test */", []Token{Token{Type: 3, Line: 1, Text: "/* This is a test */"}}},
	{"line", "% This is a test", []Token{Token{Type: 3, Line: 1, Text: "% This is a test"}}},
	{"mixed", `% This is a test
% This is a test
% This is a test
/* This is also %
  a test */`, []Token{Token{Type: 3, Line: 1, Text: "% This is a test\n"}, Token{Type: 3, Line: 1, Text: "% This is a test\n"}, Token{Type: 3, Line: 1, Text: "% This is a test\n"}, Token{Type: 3, Line: 1, Text: "/* This is also %\n  a test */"}}},
	{"atom0", `cheese`, []Token{Token{Type: 5, Line: 1, Text: "cheese"}}},
	{"atom0", `cheese123`, []Token{Token{Type: 5, Line: 1, Text: "cheese123"}}},
	{"atom0", `cheeseAndSalami`, []Token{Token{Type: 5, Line: 1, Text: "cheeseAndSalami"}}},
	{"atom0", `cheese_a_thing`, []Token{Token{Type: 5, Line: 1, Text: "cheese_a_thing"}}},
	{"atom1", `'this atom'`, []Token{Token{Type: 5, Line: 1, Text: "'this atom'"}}},
	{"atom2", `'this \' atom'`, []Token{Token{Type: 5, Line: 1, Text: "'this \\' atom'"}}},
	{"variable0", `X`, []Token{Token{Type: 8, Line: 1, Text: "X"}}},
	{"variable1", `Food`, []Token{Token{Type: 8, Line: 1, Text: "Food"}}},
	{"cluase0", `likes(sam,Food).`, []Token{Token{Type: 6, Line: 1, Text: "likes"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 5, Line: 1, Text: "sam"}, Token{Type: 16, Line: 1, Text: ","}, Token{Type: 8, Line: 1, Text: "Food"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 15, Line: 1, Text: "."}}},
	{"cluase1", `likes(sam,orange).`, []Token{Token{Type: 6, Line: 1, Text: "likes"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 5, Line: 1, Text: "sam"}, Token{Type: 16, Line: 1, Text: ","}, Token{Type: 5, Line: 1, Text: "orange"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 15, Line: 1, Text: "."}}},
	{"cluase2", `likes(sam,_).`, []Token{Token{Type: 6, Line: 1, Text: "likes"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 5, Line: 1, Text: "sam"}, Token{Type: 16, Line: 1, Text: ","}, Token{Type: 9, Line: 1, Text: "_"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 15, Line: 1, Text: "."}}},
	{"cluase2", `likes/2(sam,__thing).`, []Token{Token{Type: 6, Line: 1, Text: "likes/2"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 5, Line: 1, Text: "sam"}, Token{Type: 16, Line: 1, Text: ","}, Token{Type: 8, Line: 1, Text: "__thing"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 15, Line: 1, Text: "."}}},
	{"cluase4", `likes/2(sam,Thing) :- yummy(Thing).`, []Token{Token{Type: 6, Line: 1, Text: "likes/2"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 5, Line: 1, Text: "sam"}, Token{Type: 16, Line: 1, Text: ","}, Token{Type: 8, Line: 1, Text: "Thing"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 7, Line: 1, Text: ":-"}, Token{Type: 6, Line: 1, Text: "yummy"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 8, Line: 1, Text: "Thing"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 15, Line: 1, Text: "."}}},
	{"cluase5", `eatenChocs(tristan,1000000).`, []Token{Token{Type: 6, Line: 1, Text: "eatenChocs"}, Token{Type: 13, Line: 1, Text: "("}, Token{Type: 5, Line: 1, Text: "tristan"}, Token{Type: 16, Line: 1, Text: ","}, Token{Type: 4, Line: 1, Text: "1000000"}, Token{Type: 14, Line: 1, Text: ")"}, Token{Type: 15, Line: 1, Text: "."}}},
}

func TestNew(t *testing.T) {
	for _, st := range tests {
		func(st test) {
			t.Run(st.name, func(t *testing.T) {
				var ctx context.Context
				s := New(ctx, "test.pl", bytes.NewBuffer([]byte(st.src)))

				ts := []Token{}
				for {
					l := s.Next()
					if l.Type == EOF {
						break
					}
					ts = append(ts, l)
				}
				if !reflect.DeepEqual(st.exp, ts) {
					t.Fatalf("expected: %#v, got: %#v", st.exp, ts)
				}
			})
		}(st)
	}
}
