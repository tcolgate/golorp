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
	"log"
	"testing"

	"github.com/tcolgate/golorp/context"
	"github.com/tcolgate/golorp/parse"
	"github.com/tcolgate/golorp/scan"
)

type testL0 struct {
	name string
	q    string
	p    string
}

var testsL0 = []testL0{
	{"query0", `p(Z,h(Z,W),f(W)).`, `p(Z,h(Z,W),f(W)).`},
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

			log.Println(q)
			log.Println(p)

			compileL0(p, q)
		})
	}
}
