package scheme

import (
	"fmt"
	"testing"
)

type tokenTypeTest struct {
	source string
	result rune
}

type tokenizeTest struct {
	source string
	result []string
}

var tokenTypeTests = []tokenTypeTest{
	{"100", IntToken},
	{"+", IdentifierToken},
	{"-", IdentifierToken},
	{"f2000", IdentifierToken},
	{"a0?!*/<=>:$%^&_~", IdentifierToken},
}

var tokenizeTests = []tokenizeTest{
	{"1", []string{"1"}},
	{"(+ 1)", []string{"(", "+", "1", ")"}},
	{"(+ 1 (+ 1))", []string{"(", "+", "1", "(", "+", "1", ")", ")"}},
	{"(+ (- 1)2)", []string{"(", "+", "(", "-", "1", ")", "2", ")"}},
}

func TestTokenType(t *testing.T) {
	for _, test := range tokenTypeTests {
		l := NewLexer(test.source)

		actual := l.TokenType()
		if actual != test.result {
			t.Errorf("%s => %s; want %s", test.source, tokenTypeString(actual), tokenTypeString(test.result))
		}
	}
}

func TestNextToken(t *testing.T) {
	for _, test := range tokenizeTests {
		actual := tokenizedStrings(test.source)
		if !areTheSameStrings(actual, test.result) {
			t.Errorf("%s => %s; want %s", test.source, actual, test.result)
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

func tokenizedStrings(source string) []string {
	l := NewLexer(source)
	tokens := []string{}

	for {
		token := l.NextToken()
		if token == "" {
			break
		}
		tokens = append(tokens, token)
	}

	return tokens
}

func areTheSameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
