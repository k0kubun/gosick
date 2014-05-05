package scheme

import (
	"fmt"
	"testing"
)

type lexerTest struct {
	source string
	result rune
}

var lexerTests = []lexerTest{
	{"100", IntToken},
	{"+", IdentifierToken},
	{"-", IdentifierToken},
	{"f2000", IdentifierToken},
	{"a0?!*/<=>:$%^&_~", IdentifierToken},
}

func TestLexer(t *testing.T) {
	for _, test := range lexerTests {
		l := NewLexer(test.source)

		actual := l.TokenType()
		if actual != test.result {
			t.Errorf("%s => %s; want %s", test.source, tokenTypeString(actual), tokenTypeString(test.result))
		}
	}
}

func tokenTypeString(tokenType rune) string {
	switch tokenType {
	case EOF:
		return "EOF"
	case IdentifierToken:
		return "IdentifierToken"
	case IntToken:
		return "IntToken"
	default:
		return fmt.Sprintf("%c", tokenType)
	}
}
