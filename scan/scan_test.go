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
