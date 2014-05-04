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
	token := l.Peek()
	fmt.Printf("Detected token: %c\n", token)
	text := ""
	switch token {
	case '(', ')', '\'', scanner.EOF:
		fmt.Println("Unexpected")
	case scanner.Int, scanner.Float:
		fmt.Println("Unexpected")
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
