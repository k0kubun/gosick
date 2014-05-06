package scheme

import (
	"testing"
)

type interpreterTest struct {
	source  string
	results []string
}

type evalErrorTest struct {
	source  string
	message string
}

var interpreterTests = []interpreterTest{
	makeIT("12", "12"),
	makeIT("()", "()"),
	makeIT("#f #t", "#f", "#t"),
	makeIT("1234567890", "1234567890"),

	makeIT("(+)", "0"),
	makeIT("(- 1)", "1"),
	makeIT("(*)", "1"),
	makeIT("(/ 1)", "1"),

	makeIT("(+ 1 20 300 4000)", "4321"),
	makeIT(" ( + 1 2 3 ) ", "6"),
	makeIT("(+ 1 (+ 2 3) (+ 3 4))", "13"),
	makeIT("(- 3(- 2 3)(+ 3 0))", "1"),
	makeIT("(*(* 3 3)3)", "27"),
	makeIT("(/ 100(/ 4 2))", "50"),
	makeIT("(+ (* 100 3) (/(- 4 2) 2))", "301"),

	makeIT("(number? 100", "#t"),
	makeIT("(number? (+ 3(* 2 8)))", "#t"),
	makeIT("(number? #t)", "#f"),
	makeIT("(number? ())", "#f"),

	makeIT("(null? 1)", "#f"),
	makeIT("(null? ())", "#t"),

	makeIT("(define x 1) x", "x", "1"),
	makeIT("(define x (+ 1 3)) x", "x", "4"),
	makeIT("(define x 1) (define y 2) (define z 3) (+ x (* y z))", "x", "y", "z", "7"),
	makeIT("(define x 1) (define x 2) x", "x", "x", "2"),

	makeIT("'12", "12"),
	makeIT("'hello", "hello"),
	makeIT("'#f", "#f"),
	makeIT("'#t", "#t"),
	makeIT("'(1)", "(1)"),
	makeIT("'(  1   2   3  )", "(1 2 3)"),
	makeIT("'( 1 ( 2 3 ) )", "(1 (2 3))"),

	makeIT("(quote 12)", "12"),
	makeIT("(quote hello)", "hello"),
	makeIT("(quote #f)", "#f"),
	makeIT("(quote #t)", "#t"),
	makeIT("(quote  ( 1 (3) 4 ))", "(1 (3) 4)"),
}

var evalErrorTests = []evalErrorTest{
	{"(1)", "invalid application"},
	{"hello", "Unbound variable: hello"},
	{"(quote)", "Compile Error: syntax-error: malformed quote"},
	{"(define)", "Compile Error: syntax-error: (define)"},
	{"(null? 1 2)", "Compile Error: wrong number of arguments: number? requires 1, but got 2"},
}

func makeIT(source string, results ...string) interpreterTest {
	return interpreterTest{source: source, results: results}
}

func TestInterpreter(t *testing.T) {
	for _, test := range interpreterTests {
		p := NewParser(test.source)

		// I don't know what p.Peek() affects to p, but this test fails without p.Peek().
		p.Peek()

		for i := 0; i < len(test.results); i++ {
			result := test.results[i]

			parsedObject := p.Parse()
			if parsedObject == nil {
				t.Errorf("%s => <nil>; want %s", test.source, result)
				return
			}
			actual := parsedObject.String()
			if actual != result {
				t.Errorf("%s => %s; want %s", test.source, actual, result)
			}
		}
	}
}

func TestEvalError(t *testing.T) {
	for _, test := range evalErrorTests {
		assertError(t, test.source, test.message)
	}
}

func assertError(t *testing.T, source string, message string) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("\"%s\" did not panic\n want: %s\n", source, message)
		} else if err != message {
			t.Errorf("\"%s\" paniced\nwith: %s\nwant: %s\n", source, err, message)
		}
	}()

	p := NewParser(source)
	p.Peek()
	p.Parse().Eval()
}
