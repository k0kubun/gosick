package scheme

import (
	"fmt"
	"strings"
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
	{"(", '('},
	{")", ')'},
	{"'", '\''},

	{"100", IntToken},

	{"#f", BooleanToken},
	{"#t", BooleanToken},

	{"+", IdentifierToken},
	{"-", IdentifierToken},
	{"f2000", IdentifierToken},
	{"a0?!*/<=>:$%^&_~", IdentifierToken},

	{"\"a b\"", StringToken},
}

var tokenizeTests = []tokenizeTest{
	{"1", makeTokens("1")},
	{"#f", makeTokens("#f")},
	{"#t", makeTokens("#t")},

	{"(+ 1)", makeTokens("(,+,1,)")},
	{"(+ 1 (+ 1))", makeTokens("(,+,1,(,+,1,),)")},
	{"(+ (- 1)2)", makeTokens("(,+,(,-,1,),2,)")},
	{"(* (/ 1)2)", makeTokens("(,*,(,/,1,),2,)")},
	{"(number? 1)", makeTokens("(,number?,1,)")},

	{"((()", makeTokens("(,(,(,)")},
	{"()))", makeTokens("(,),),)")},
	{")))", makeTokens("),),)")},
	{"((()))(()", makeTokens("(,(,(,),),),(,(,)")},

	{"'2", makeTokens("',2")},
	{"'#f", makeTokens("',#f")},
	{"'#t", makeTokens("',#t")},
	{"'hello", makeTokens("',hello")},
	{"'(1 2 3)", makeTokens("',(,1,2,3,)")},
	{"'(1(2 3))", makeTokens("',(,1,(,2,3,),)")},

	{"\"a b\"", makeTokens("\"a b\"")},
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

func TestAllTokens(t *testing.T) {
	for _, test := range tokenizeTests {
		l := NewLexer(test.source)
		actual := l.AllTokens()

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
	case StringToken:
		return "StringToken"
	default:
		return fmt.Sprintf("%c", tokenType)
	}
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

func makeTokens(text string) []string {
	return strings.Split(text, ",")
}
