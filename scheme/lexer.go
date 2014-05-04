// Lexer is the abbreviation for Lexical Analyzer.
// Lexer consists of Scanner and Tokenizer, and it owns
// a role of analyzing tokens.
// And this is used by Parser for syntactic analysis.
//
// Package text/scanner has both of them, so Lexer uses
// their customized version.

package scheme

import (
	"fmt"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
}

func (l *Lexer) NextToken() Object {
	text := ""
	switch l.Scan() {
	case '(', ')', '\'', scanner.EOF:
		fmt.Println("Unexpected flow")
	case scanner.Int:
		return NewNumber(l.TokenText())
	case '-':
		fmt.Println("Unexpected")
	case scanner.String:
		text = l.TokenText()
	default:
		text = l.TokenText()
	}
	fmt.Println(text)
	return nil
}
