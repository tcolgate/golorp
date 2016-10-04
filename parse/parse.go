package parse

import (
	"fmt"

	"github.com/tcolgate/golorp/scan"
)

// Parser stores the state for the ivy parser.
type Parser struct {
	scanner  *scan.Scanner
	fileName string

	lineNum    int
	errorCount int // Number of errors.

	peekTok scan.Token
	curTok  scan.Token // most recent token from scanner

	operators map[string]Op
}

// NewParser returns a new parser that will read from the scanner.
// The context must have have been created by this package's NewContext function.
func NewParser(fileName string, scanner *scan.Scanner) *Parser {
	return &Parser{
		scanner:  scanner,
		fileName: fileName,

		operators: defaultOps, // TODO: Need a deep copy here
	}
}

func (p *Parser) next() scan.Token {
	return p.nextErrorOut(true)
}

// nextErrorOut accepts a flag whether to trigger a panic on error.
// The flag is set to false when we are draining input tokens in FlushToNewline.
func (p *Parser) nextErrorOut(errorOut bool) scan.Token {
	tok := p.peekTok
	if tok.Type != scan.EOF {
		p.peekTok = scan.Token{Type: scan.EOF}
	} else {
		tok = p.scanner.Next()
	}
	if tok.Type == scan.Error && errorOut {
		fmt.Errorf("%q", tok) // Need a local output writer
	}
	p.curTok = tok
	if tok.Type != scan.Newline {
		// Show the line number before we hit the newline.
		p.lineNum = tok.Line
	}
	return tok
}

func (p *Parser) peek() scan.Token {
	tok := p.peekTok
	if tok.Type != scan.EOF {
		return tok
	}
	p.peekTok = p.scanner.Next()
	return p.peekTok
}
