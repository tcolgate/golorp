package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tcolgate/golorp/context"
	"github.com/tcolgate/golorp/parse"
	"github.com/tcolgate/golorp/scan"
)

const qprompt = "?- "

func main() {
	flag.Parse()

	// Load databse file s from the command line

	for _, fn := range flag.Args() {
		f, err := os.Open(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: opening %s failed, %+v", fn, err)
			os.Exit(1)
		}

		var ctx context.Context
		s := scan.New(ctx, fn, bufio.NewReader(f))
		p := parse.New(fn, s)

		for {
			t, err := p.NextTerm()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: reading %s failed, %+v", fn, err)
				os.Exit(1)
			}

			fmt.Println("Got %v\n", t)
		}
	}

	// Process queries
	var ctx context.Context
	s := scan.New(ctx, "stdin", bufio.NewReader(os.Stdin))
	p := parse.New("stdin", s)
	for {
		fmt.Printf(qprompt)
		t0, err := p.NextTerm()
		if err == io.EOF {
			fmt.Println(os.Stderr, "got EOF\n", err)
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			break
		}
		fmt.Println("Got %v\n", t0)
	}
}
