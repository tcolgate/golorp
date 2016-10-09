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
)

type test struct {
	name string
	src  string
	exp  []Term
}

var tests = []test{
	{"clause0", `likes.`, []Term{}},
	{"clause1", `likes(sam).`, []Term{}},
	{"clause2", `likes(sam,Food).`, []Term{}},
	{"clause3", `likes(sam,orange).`, []Term{}},
	{"clause4", `likes(sam,_).`, []Term{}},
	{"clause5", `likes/2(sam,__thing).`, []Term{}},
	//{"clause6", `likes/2(sam,Thing) :- yummy(Thing).`, []Term{}},
	{"clause7", `eatenChocs(tristan,1000000).`, []Term{}},
}

func TestNew(t *testing.T) {
	var ctx context.Context
	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			log.Println(st.name)

			s := scan.New(ctx, "file.pl", bytes.NewBuffer([]byte(st.src)))
			p := New("file.pl", s)

			for {
				t, err := p.NextTerm()
				log.Println(t)
				log.Println(err)
				if err != nil {
					break
				}
			}
		})
	}
}
