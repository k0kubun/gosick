package scheme

import (
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
}

func (l Lexer) ReadToken() *Token {
	return nil
}
