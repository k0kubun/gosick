// Interpreter is a scheme source code interpreter.
// It owns a role of API for executing scheme program.
// Interpreter embeds Parser and delegates syntactic analysis to it.

package scheme

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"
)

type Interpreter struct {
	*Parser
	closure *Closure
}

func NewInterpreter(source string) *Interpreter {
	i := &Interpreter{
		Parser: NewParser(source),
		closure: &Closure{
			ObjectBase:   ObjectBase{parent: nil},
			localBinding: defaultBinding(),
		},
	}
	i.loadBuiltinLibrary("builtin")
	return i
}

// Load new source code with current environment
func (i *Interpreter) ReloadSourceCode(source string) {
	i.Parser = NewParser(source)
}

func (i *Interpreter) PrintResult(dumpAST bool) {
	results := i.EvalSource(dumpAST)
	if dumpAST {
		fmt.Printf("\n*** Result ***\n")
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func (i *Interpreter) EvalSource(dumpAST bool) (results []string) {
	defer func() {
		if err := recover(); err != nil {
			results = append(results, fmt.Sprintf("*** ERROR: %s", err))
		}
	}()

	for i.Peek() != scanner.EOF {
		expression := i.Parser.Parse(i.closure)
		if dumpAST {
			fmt.Printf("\n*** AST ***\n")
			i.DumpAST(expression, 0)
		}

		if expression == nil {
			return
		}
		results = append(results, expression.Eval().String())
	}
	return
}

func (i *Interpreter) DumpAST(object Object, indentLevel int) {
	if object == nil {
		return
	}
	switch object.(type) {
	case *Application:
		i.printWithIndent("Application", indentLevel)
		i.DumpAST(object.(*Application).procedure, indentLevel+1)
		i.DumpAST(object.(*Application).arguments, indentLevel+1)
	case *Pair:
		pair := object.(*Pair)
		if pair.Car == nil && pair.Cdr == nil {
			i.printWithIndent("()", indentLevel)
			return
		}
		i.printWithIndent("Pair", indentLevel)
		i.DumpAST(pair.Car, indentLevel+1)
		i.DumpAST(pair.Cdr, indentLevel+1)
	case *String:
		i.printWithIndent(fmt.Sprintf("String(%s)", object), indentLevel)
	case *Number:
		i.printWithIndent(fmt.Sprintf("Number(%s)", object), indentLevel)
	case *Boolean:
		i.printWithIndent(fmt.Sprintf("Boolean(%s)", object), indentLevel)
	case *Variable:
		i.printWithIndent(fmt.Sprintf("Variable(%s)", object.(*Variable).identifier), indentLevel)
	}
}

func (i *Interpreter) printWithIndent(text string, indentLevel int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indentLevel), text)
}

func (i *Interpreter) loadBuiltinLibrary(name string) {
	originalParser := i.Parser

	buffer, err := ioutil.ReadFile(i.libraryPath(name))
	if err != nil {
		log.Fatal(err)
	}
	i.Parser = NewParser(string(buffer))
	i.EvalSource(false)

	i.Parser = originalParser
}

func (i *Interpreter) libraryPath(name string) string {
	return filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"k0kubun",
		"gosick",
		"lib",
		name+".scm",
	)
}
