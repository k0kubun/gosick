package scheme

import (
	"io"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
}

func NewLexer(source io.Reader) *Lexer {
	lexer := Lexer{}
	lexer.Init(source)

	return &lexer
}
