package scheme

import (
	"fmt"
	"io/ioutil"
	"os"
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

type loadTest struct {
	loadSource string
	source     string
	message    string
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

	evalTest("(< 2 1)", "#f"),
	evalTest("(< 1 2)", "#t"),
	evalTest("(< 1 2 1)", "#f"),
	evalTest("(< 1 2 3)", "#t"),

	evalTest("(<= 2 1)", "#f"),
	evalTest("(<= 1 2)", "#t"),
	evalTest("(<= 2 2)", "#t"),
	evalTest("(<= 1 2 1)", "#f"),
	evalTest("(<= 1 2 3)", "#t"),
	evalTest("(<= 1 1 1)", "#t"),

	evalTest("(> 1 2)", "#f"),
	evalTest("(> 2 1)", "#t"),
	evalTest("(> 1 2 1)", "#f"),
	evalTest("(> 3 2 1)", "#t"),

	evalTest("(>= 1 2)", "#f"),
	evalTest("(>= 2 1)", "#t"),
	evalTest("(>= 2 2)", "#t"),
	evalTest("(>= 1 2 1)", "#f"),
	evalTest("(>= 3 2 1)", "#t"),
	evalTest("(>= 1 1 1)", "#t"),

	evalTest("(not #f)", "#t"),
	evalTest("(not #t)", "#f"),
	evalTest("(not (number? ()))", "#t"),
	evalTest("(not 1)", "#f"),
	evalTest("(not ())", "#f"),

	evalTest("(cons 1 2)", "(1 . 2)"),
	evalTest("(car (cons 1 2))", "1"),
	evalTest("(cons '(1 2) (list 1 2 3))", "((1 2) 1 2 3)"),
	evalTest("(cons (cons 1 2) 3)", "((1 . 2) . 3)"),
	evalTest("(car '(1))", "1"),
	evalTest("(cdr '(1))", "()"),
	evalTest("(car '(1 2))", "1"),
	evalTest("(cdr '(1 2))", "(2)"),

	evalTest("(list)", "()"),
	evalTest("(list 1 2 3)", "(1 2 3)"),
	evalTest("(cdr (list 1 2 3))", "(2 3)"),

	evalTest("(length ())", "0"),
	evalTest("(length '(1 2))", "2"),
	evalTest("(length (list 1 '(2 3) 4))", "3"),

	evalTest("(memq (car (cons 'b 'c)) '(a b c))", "(b c)"),
	evalTest("(memq 'd '(a b c))", "#f"),
	evalTest("(memq 'a (cons 'a 'b))", "(a . b)"),

	evalTest("(last '(1 2 3))", "3"),
	evalTest("(last (list 1 (+ 2 3)))", "5"),

	evalTest("(append)", "()"),
	evalTest("(append '(1))", "(1)"),
	evalTest("(append '(1 2) '(3 4))", "(1 2 3 4)"),

	evalTest("(string-append)", "\"\""),
	evalTest("(string-append \"a\" \" \" \"b\")", "\"a b\""),

	evalTest("(string->symbol \"a\")", "a"),
	evalTest("(symbol->string 'a)", "\"a\""),

	evalTest("(string->number \"1\")", "1"),
	evalTest("(number->string 1)", "\"1\""),

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

	evalTest("(eq? 1 1)", "#t"),
	evalTest("(eq? 1 2)", "#f"),
	evalTest("(eq? 1 #f)", "#f"),
	evalTest("(eq? #f #f)", "#t"),
	evalTest("(eq? 'foo 'foo)", "#t"),
	evalTest("(eq? \"foo\" \"foo\")", "#f"),
	evalTest("(eq? '(1 2) '(1 2))", "#f"),

	evalTest("(neq? 1 1)", "#f"),
	evalTest("(neq? 1 2)", "#t"),
	evalTest("(neq? 1 #f)", "#t"),
	evalTest("(neq? #f #f)", "#f"),
	evalTest("(neq? 'foo 'foo)", "#f"),
	evalTest("(neq? \"foo\" \"foo\")", "#t"),
	evalTest("(neq? '(1 2) '(1 2))", "#t"),

	evalTest("(equal? 1 1)", "#t"),
	evalTest("(equal? 1 2)", "#f"),
	evalTest("(equal? 1 #f)", "#f"),
	evalTest("(equal? #f #f)", "#t"),
	evalTest("(equal? 'foo 'foo)", "#t"),
	evalTest("(equal? \"foo\" \"foo\")", "#f"),
	evalTest("(equal? '(1 1) '(1 2))", "#f"),
	evalTest("(equal? '(1 2) '(1 2))", "#t"),

	evalTest("(define x 1) x", "x", "1"),
	evalTest("(define x (+ 1 3)) x", "x", "4"),
	evalTest("(define x 1) (define y 2) (define z 3) (+ x (* y z))", "x", "y", "z", "7"),
	evalTest("(define x 1) (define x 2) x", "x", "x", "2"),

	evalTest("(lambda (x) x)", "#<closure #f>"),
	evalTest("((lambda (x) 1) 2)", "1"),
	evalTest("((lambda (x y z) (+ 3 4) (- 4 1) ) 2 3 3)", "3"),
	evalTest("((lambda (x) (+ x x)) 1)", "2"),
	evalTest("(define x 1) ((lambda (y) (+ x y)) 2)", "x", "3"),
	evalTest("(define x 1) ((lambda (x) (+ x x)) 2)", "x", "4"),
	evalTest("((lambda (x) (define x 3) x) 2)", "3"),
	evalTest("((lambda (x y z) (* (+ x y) z)) 1 2 3)", "9"),
	evalTest("(define x (lambda (a) (* 2 a))) (define y (lambda (a) (* 3 a))) (define z (lambda (a b) (x a) (y b))) (* (x 3) (y 2) (z 4 5))", "x", "y", "z", "540"),

	evalTest("(define x 2) (set! x 3) x", "x", "#<undef>", "3"),
	evalTest("(define x 4) ((lambda (x) (set! x 3) x) 2) x", "x", "3", "4"),

	evalTest("(define x (cons 3 2)) (set-car! x 1) x", "x", "#<undef>", "(1 . 2)"),
	evalTest("(define x (cons 3 2)) (set-cdr! x 1) x", "x", "#<undef>", "(3 . 1)"),
	evalTest("((lambda (x) (set-car! x 1) x) (cons 2 3))", "(1 . 3)"),
	evalTest("((lambda (x) (set-cdr! x 1) x) (cons 2 3))", "(2 . 1)"),

	evalTest("(if #t 1 2)", "1"),
	evalTest("(if #f 1 2)", "2"),
	evalTest("(if (null? ()) 1 2)", "1"),
	evalTest("(if (null? 3) 1)", "#<undef>"),
	evalTest("(if (number? 3) 'num)", "num"),

	evalTest("(cond (#t))", "#t"),
	evalTest("(cond (()))", "()"),
	evalTest("(cond (else))", "#<undef>"),
	evalTest("(cond (#f 1) (#t 2) (else 3))", "2"),
	evalTest("(cond ((number? 3) 'hello) (else 'no))", "hello"),

	evalTest("(and)", "#t"),
	evalTest("(and #t 3)", "3"),
	evalTest("(and (number? 3) (boolean? #f))", "#t"),
	evalTest("(and (number? 3) (boolean? 3))", "#f"),

	evalTest("(or)", "#f"),
	evalTest("(or #f 3)", "3"),
	evalTest("(or 3)", "3"),
	evalTest("(or #f 3 #f)", "3"),
	evalTest("(or (number? 3) (boolean? 3))", "#t"),
	evalTest("(or (number? #f) (boolean? 3))", "#f"),

	evalTest("(begin)", "#<undef>"),
	evalTest("(begin 1 2 3)", "3"),
	evalTest("(begin (define x 2) (set! x 3) x) x", "3", "3"),

	evalTest("(do () (#t)))", "#t"),
	evalTest("(do ((i 1) (j 1)) (#t)))", "#t"),
	evalTest("(define x \"\") (do ((i 1 (+ i 1)) (j 1 (* j 2))) ((> i 3) x) (begin (set! x (string-append x (number->string i))) (set! x (string-append x (number->string j)))))", "x", "\"112234\""),
}

// let parsing break tree structure, so not to apply parser test
var letTests = []interpreterTest{
	evalTest("(let ((x 1)) x)", "1"),
	evalTest("(let ((x 1) (y 2)) (+ x y))", "3"),
	evalTest("(let* ((x 1)) x)", "1"),
	evalTest("(let* ((x 1) (y 2)) (+ x y))", "3"),
	evalTest("(letrec ((x 1)) x)", "1"),
	evalTest("(letrec ((x 1) (y 2)) (+ x y))", "3"),
}

var runtimeErrorTests = []interpreterTest{
	evalTest("(1)", "*** ERROR: invalid application"),
	evalTest("hello", "*** ERROR: Unbound variable: hello"),
	evalTest("((lambda (x) (define y 1) 1) 1) y", "1", "*** ERROR: Unbound variable: y"),
	evalTest("'1'", "1", "*** ERROR: unterminated quote"),
	evalTest("(last ())", "*** ERROR: pair required: ()"),
	evalTest("((lambda (x) (set! x 3) x) 2) x", "3", "*** ERROR: Unbound variable: x"),
}

var compileErrorTests = []interpreterTest{
	evalTest("(quote)", "*** ERROR: Compile Error: syntax-error: malformed quote"),
	evalTest("(define)", "*** ERROR: Compile Error: syntax-error: (define)"),

	evalTest("(-)", "*** ERROR: Compile Error: procedure requires at least 1 argument"),
	evalTest("(/)", "*** ERROR: Compile Error: procedure requires at least 1 argument"),
	evalTest("(number?)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),
	evalTest("(null?)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),
	evalTest("(null? 1 2)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 2"),
	evalTest("(not)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),

	evalTest("(+ 1 #t)", "*** ERROR: Compile Error: number required, but got #t"),
	evalTest("(- #t)", "*** ERROR: Compile Error: number required, but got #t"),
	evalTest("(* ())", "*** ERROR: Compile Error: number required, but got ()"),
	evalTest("(/ '(1 2 3))", "*** ERROR: Compile Error: number required, but got (1 2 3)"),

	evalTest("(string-append #f)", "*** ERROR: Compile Error: string required, but got #f"),
	evalTest("(string-append 1)", "*** ERROR: Compile Error: string required, but got 1"),

	evalTest("(string->symbol)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),
	evalTest("(string->symbol 'hello)", "*** ERROR: Compile Error: string required, but got hello"),
	evalTest("(symbol->string)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),
	evalTest("(symbol->string \"\")", "*** ERROR: Compile Error: symbol required, but got \"\""),
	evalTest("(string->number 1)", "*** ERROR: Compile Error: string required, but got 1"),
	evalTest("(number->string \"1\")", "*** ERROR: Compile Error: number required, but got \"1\""),

	evalTest("(car ())", "*** ERROR: Compile Error: pair required, but got ()"),
	evalTest("(cdr ())", "*** ERROR: Compile Error: pair required, but got ()"),
	evalTest("(car)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),
	evalTest("(cdr)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 1, but got 0"),

	evalTest("(length (cons 1 2))", "*** ERROR: Compile Error: proper list required for function application or macro use"),
	evalTest("(memq 'a '(a b c) 1)", "*** ERROR: Compile Error: wrong number of arguments: number? requires 2, but got 3"),
	evalTest("(append () 1 ())", "*** ERROR: Compile Error: proper list required for function application or macro use"),
	evalTest("(set! x 1 1)", "*** ERROR: Compile Error: syntax-error: malformed set!"),

	evalTest("(cond)", "*** ERROR: Compile Error: syntax-error: at least one clause is required for cond"),
	evalTest("(cond ())", "*** ERROR: Compile Error: syntax-error: bad clause in cond"),
	evalTest("(cond (#t) (else) ())", "*** ERROR: Compile Error: syntax-error: 'else' clause followed by more clauses"),

	evalTest("(do () ()))", "*** ERROR: Compile Error: syntax-error: malformed do"),
}

func evalTest(source string, results ...string) interpreterTest {
	return interpreterTest{source: source, results: results}
}

func runTests(t *testing.T, tests []interpreterTest) {
	for _, test := range tests {
		i := NewInterpreter(test.source)
		evalResults := i.EvalSource(false)

		for i := 0; i < len(test.results); i++ {
			expect := test.results[i]
			actual := evalResults[i]
			if actual != expect {
				t.Errorf("%s => %s; want %s", test.source, actual, expect)
			}
		}
	}
}

func TestInterpreter(t *testing.T) {
	runTests(t, interpreterTests)
	runTests(t, letTests)
	runTests(t, runtimeErrorTests)
	runTests(t, compileErrorTests)
}

func TestLoad(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "load_test")
	if err != nil {
		panic(err)
	}
	fileSource := "(define x 3)"
	ioutil.WriteFile(file.Name(), []byte(fileSource), os.ModeAppend)
	defer os.Remove(file.Name())

	source := fmt.Sprintf("(load \"%s\") x (load invalid)", file.Name())
	interpreter := NewInterpreter(source)
	expects := []string{"#t", "3", "*** ERROR: Unbound variable: invalid"}
	actuals := interpreter.EvalSource(false)

	for i := 0; i < len(actuals); i++ {
		expect := expects[i]
		actual := actuals[i]
		if actual != expect {
			t.Errorf("%s => %s; want %s", source, actual, expect)
		}
	}
}
