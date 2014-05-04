package scheme

import (
	"log"
	"strings"
	"text/scanner"
)

type Parser struct {
	*Lexer
}

func NewParser(source string) *Parser {
	parser := &Parser{&Lexer{}}
	parser.Init(strings.NewReader(source))
	return parser
}

func (p *Parser) Parse() Object {
	token := p.Next()
	switch token {
	case ')':
		log.Print("Unexpected token: ')'")
		return nil
	case '\'':
	case '(':
		return p.parseListBody()
	case scanner.EOF:
		return nil
	}
	return nil
}

// This function returns only *Cons.
// But returns Object because if this method returns nil which is not interface type,
// Parse()'s result cannot be judged as nil.
// To avoid such a situation, this method's return value is Object.
func (p *Parser) parseListBody() Object {
	if p.Peek() == ')' {
		p.Next()
		return new(Cons)
	}

	car := p.Parse()
	if car == nil {
		log.Print("Unsupported flow (maybe incomplete source or unexpected expression)")
		return nil
	}
	cdr := p.parseListBody().(*Cons)
	return &Cons{Car: car, Cdr: cdr}
}
