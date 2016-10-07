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
	"log"
	"testing"

	"github.com/tcolgate/golorp/context"
)

type test struct {
	name string
	src  string
	exp  []Token
}

var tests = []test{
	{"single", "/* This is a test */", []Token{}},
	{"line", "% This is a test", []Token{}},
	{"mixed", `% This is a test
% This is a test
% This is a test
/* This is also %
  a test */`, []Token{}},
	{"atom0", `cheese`, []Token{}},
	{"atom0", `cheese123`, []Token{}},
	{"atom0", `cheeseAndSalami`, []Token{}},
	{"atom0", `cheese_a_thing`, []Token{}},
	{"atom1", `'this atom'`, []Token{}},
	{"atom2", `'this \' atom'`, []Token{}},
	{"variable0", `X`, []Token{}},
	{"variable1", `Food`, []Token{}},
	{"cluase0", `likes(sam,Food).`, []Token{}},
	{"cluase1", `likes(sam,orange).`, []Token{}},
	{"cluase2", `likes(sam,_).`, []Token{}},
	{"cluase2", `likes/2(sam,__thing).`, []Token{}},
	{"cluase4", `likes/2(sam,Thing) :- yummy(Thing).`, []Token{}},
	{"cluase5", `eatenChocs(tristan,1000000).`, []Token{}},
}

func TestNew(t *testing.T) {
	var ctx context.Context
	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			log.Println(st.name)

			s := New(ctx, "file.pl", bytes.NewBuffer([]byte(st.src)))

			for {
				l := s.Next()
				log.Println(l)
				if l.Type == EOF {
					return
				}
			}
		})
	}
}
