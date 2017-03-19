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
	"bytes"
	"fmt"
	"testing"

	"github.com/tcolgate/golorp/context"
	"github.com/tcolgate/golorp/scan"
)

type test struct {
	name string
	src  string
	exp  string
}

var tests = []test{
	{"clause0", `likes.`, `("likes"/0 [])`},
	{"clause1", `1 + 2.`, `("+"/2 [(number 1) (number 2)])`},
	{"clause1", `2 / 3.`, `("/"/2 [(number 2) (number 3)])`},
	{"clause2", `print(1 + 2 + 3 + 4 + 5).`, `("print"/1 [("+"/2 [("+"/2 [("+"/2 [("+"/2 [(number 1) (number 2)]) (number 3)]) (number 4)]) (number 5)])])`},
	{"clause2", `1 + (2 * 3).`, `("+"/2 [(number 1) ("*"/2 [(number 2) (number 3)])])`},
	{"clause3", `-2.`, `("-"/1 [(number 2)])`},
	{"clause4", `likes(sam).`, `("likes"/1 [("sam"/0 [])])`},
	{"clause5", `likes(sam,Food).`, `("likes"/2 [("sam"/0 []) (var Food)])`},
	{"clause6", `likes(sam,orange).`, `("likes"/2 [("sam"/0 []) ("orange"/0 [])])`},
	{"clause7", `likes(sam,_).`, `("likes"/2 [("sam"/0 []) (var _)])`},
	{"clause8", `likes/2(sam,__thing).`, `("likes/2"/2 [("sam"/0 []) (var __thing)])`},
	{"clause9", `likes/2(sam,Thing) :- yummy(Thing).`, `(":-"/2 [("likes/2"/2 [("sam"/0 []) (var Thing)]) ("yummy"/1 [(var Thing)])])`},
	{"clause10", `eatenChocs(tristan,1000000).`, `("eatenChocs"/2 [("tristan"/0 []) (number 1e+06)])`},
	{"clause11", `eatenChocs(tristan + 4,1000000).`, `("eatenChocs"/2 [("+"/2 [("tristan"/0 []) (number 4)]) (number 1e+06)])`},
}

func TestNew(t *testing.T) {
	var ctx context.Context
	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			s := scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.src)))
			p := New("file.pl", s)

			t0, _ := p.NextTerm()
			str := fmt.Sprintf("%v", t0)
			if str != st.exp {
				t.Fatalf("\nexpected: %#v\ngot: %#v", st.exp, str)
			}
		})
	}
}
