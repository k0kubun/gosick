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
	evalTest("12", "12"),
	evalTest("()", "()"),
	evalTest("#f #t", "#f", "#t"),
	evalTest("1234567890", "1234567890"),

	evalTest("'12", "12"),
	evalTest("'hello", "hello"),
	evalTest("'#f", "#f"),
	evalTest("'#t", "#t"),
	evalTest("'(1)", "(1)"),
	evalTest("'(  1   2   3  )", "(1 2 3)"),
	evalTest("'( 1 ( 2 3 ) )", "(1 (2 3))"),

	evalTest("(quote 12)", "12"),
	evalTest("(quote hello)", "hello"),
	evalTest("(quote #f)", "#f"),
	evalTest("(quote #t)", "#t"),
	evalTest("(quote  ( 1 (3) 4 ))", "(1 (3) 4)"),

	evalTest("\"\"", "\"\""),
	evalTest("\"hello\"", "\"hello\""),

	evalTest("(+)", "0"),
	evalTest("(- 1)", "1"),
	evalTest("(*)", "1"),
	evalTest("(/ 1)", "1"),

	evalTest("(+ 1 20 300 4000)", "4321"),
	evalTest(" ( + 1 2 3 ) ", "6"),
	evalTest("(+ 1 (+ 2 3) (+ 3 4))", "13"),
	evalTest("(- 3(- 2 3)(+ 3 0))", "1"),
	evalTest("(*(* 3 3)3)", "27"),
	evalTest("(/ 100(/ 4 2))", "50"),
	evalTest("(+ (* 100 3) (/(- 4 2) 2))", "301"),

	evalTest("(= 2 1)", "#f"),
	evalTest("(= (* 100 3) 300)", "#t"),

	evalTest("(not #f)", "#t"),
	evalTest("(not #t)", "#f"),
	evalTest("(not (number? ()))", "#t"),
	evalTest("(not 1)", "#f"),
	evalTest("(not ())", "#f"),

	evalTest("(car '(1))", "1"),
	evalTest("(cdr '(1))", "()"),
	evalTest("(car '(1 2))", "1"),
	evalTest("(cdr '(1 2))", "(2)"),

	evalTest("(list)", "()"),
	evalTest("(list 1 2 3)", "(1 2 3)"),
	evalTest("(cdr (list 1 2 3))", "(2 3)"),

	evalTest("(string-append)", "\"\""),
	evalTest("(string-append \"a\" \" \" \"b\")", "\"a b\""),

	evalTest("(string->symbol \"a\")", "a"),
	evalTest("(symbol->string 'a)", "\"a\""),

	evalTest("(number? 100", "#t"),
	evalTest("(number? (+ 3(* 2 8)))", "#t"),
	evalTest("(number? #t)", "#f"),
	evalTest("(number? ())", "#f"),

	evalTest("(procedure? 1)", "#f"),
	evalTest("(procedure? +)", "#t"),

	evalTest("(boolean? 1)", "#f"),
	evalTest("(boolean? ())", "#f"),
	evalTest("(boolean? #t)", "#t"),
	evalTest("(boolean? #f)", "#t"),
	evalTest("(boolean? (null? 1))", "#t"),

	evalTest("(pair? 1)", "#f"),
	evalTest("(pair? ())", "#f"),
	evalTest("(pair? '(1 2 3))", "#t"),

	evalTest("(list? 1)", "#f"),
	evalTest("(list? ())", "#t"),
	evalTest("(list? '(1 2 3))", "#t"),

	evalTest("(symbol? 1)", "#f"),
	evalTest("(symbol? 'hello)", "#t"),

	evalTest("(string? 1)", "#f"),
	evalTest("(string? \"\")", "#t"),
	evalTest("(string? \"hello\")", "#t"),

	evalTest("(null? 1)", "#f"),
	evalTest("(null? ())", "#t"),

	evalTest("(define x 1) x", "x", "1"),
	evalTest("(define x (+ 1 3)) x", "x", "4"),
	evalTest("(define x 1) (define y 2) (define z 3) (+ x (* y z))", "x", "y", "z", "7"),
	evalTest("(define x 1) (define x 2) x", "x", "x", "2"),

	evalTest("(lambda (x) x)", "#<closure #f>"),
}

var evalErrorTests = []evalErrorTest{
	{"(1)", "invalid application"},
	{"hello", "Unbound variable: hello"},
	{"(quote)", "Compile Error: syntax-error: malformed quote"},
	{"(define)", "Compile Error: syntax-error: (define)"},

	{"(-)", "Compile Error: procedure requires at least 1 argument"},
	{"(/)", "Compile Error: procedure requires at least 1 argument"},
	{"(number?)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},
	{"(null?)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},
	{"(null? 1 2)", "Compile Error: wrong number of arguments: number? requires 1, but got 2"},
	{"(not)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},

	{"(+ 1 #t)", "Compile Error: number required, but got #t"},
	{"(- #t)", "Compile Error: number required, but got #t"},
	{"(* ())", "Compile Error: number required, but got ()"},
	{"(/ '(1 2 3))", "Compile Error: number required, but got (1 2 3)"},

	{"(string-append #f)", "Compile Error: string required, but got #f"},
	{"(string-append 1)", "Compile Error: string required, but got 1"},

	{"(string->symbol)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},
	{"(string->symbol 'hello)", "Compile Error: string required, but got hello"},
	{"(symbol->string)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},
	{"(symbol->string \"\")", "Compile Error: symbol required, but got \"\""},

	{"(car ())", "Compile Error: pair required, but got ()"},
	{"(cdr ())", "Compile Error: pair required, but got ()"},
	{"(car)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},
	{"(cdr)", "Compile Error: wrong number of arguments: number? requires 1, but got 0"},
}

func evalTest(source string, results ...string) interpreterTest {
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
