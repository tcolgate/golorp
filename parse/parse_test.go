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
	"log"
	"testing"

	"github.com/tcolgate/golorp/context"
	"github.com/tcolgate/golorp/scan"
	"github.com/tcolgate/golorp/term"
)

type test struct {
	name string
	src  string
	exp  []term.Term
}

var tests = []test{
	{"clause0", `likes.`, []term.Term{}},
	{"clause1", `1 + 2.`, []term.Term{}},
	{"clause1", `2 / 3.`, []term.Term{}},
	{"clause2", `print(1 + 2 + 3 + 4 + 5).`, []term.Term{}},
	{"clause2", `1 + (2 * 3).`, []term.Term{}},
	{"clause3", `-2.`, []term.Term{}},
	{"clause4", `likes(sam).`, []term.Term{}},
	{"clause5", `likes(sam,Food).`, []term.Term{}},
	{"clause6", `likes(sam,orange).`, []term.Term{}},
	{"clause7", `likes(sam,_).`, []term.Term{}},
	{"clause8", `likes/2(sam,__thing).`, []term.Term{}},
	{"clause9", `likes/2(sam,Thing) :- yummy(Thing).`, []term.Term{}},
	{"clause10", `eatenChocs(tristan,1000000).`, []term.Term{}},
	{"clause11", `eatenChocs(tristan + 4,1000000).`, []term.Term{}},
}

func TestNew(t *testing.T) {
	var ctx context.Context
	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			s := scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.src)))
			p := New("file.pl", s)

			t0, err := p.NextTerm()
			log.Println(t0, err)
		})
	}
}
