// Lexer is the abbreviation for Lexical Analyzer.
// Lexer consists of Scanner and Tokenizer, and it owns
// a role of analyzing tokens.
// And this is used by Parser for syntactic analysis.
//
// Package text/scanner has both of them, so Lexer uses
// their customized version.

package scheme

import (
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
}
