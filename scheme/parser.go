package scheme

import (
	"fmt"
	"log"
	"text/scanner"
)

type Parser struct {
	Lexer
}

func (p Parser) Parse() *Object {
	token := p.Peek()
	fmt.Printf("%c\n", token)
	switch token {
	case ')':
		log.Fatal("Unexpected token: ')'")
	case '\'':
	case '(':
	case scanner.EOF:
		return nil
	}
	return nil
}

func (p Parser) parseListBody() (*Cons, error) {
	if p.Peek() == ')' {
	}
	return nil, nil
}
