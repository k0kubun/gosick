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
		log.Fatal("Unexpected token: ')'")
	case '\'':
	case '(':
		return p.parseListBody()
	case scanner.EOF:
		return nil
	}
	return nil
}

func (p *Parser) parseListBody() Object {
	if p.Peek() == ')' {
		p.Next()
		return new(Cons)
	}
	return nil
}
