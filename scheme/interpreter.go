// Interpreter is a scheme source code interpreter.
// It owns a role of API for executing scheme program.
// Interpreter embeds Parser and delegates syntactic analysis to it.

package scheme

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Interpreter struct {
	*Parser
	topLevel *Environment
}

func NewInterpreter(source string) *Interpreter {
	interpreter := &Interpreter{Parser: NewParser(source)}
	interpreter.topLevel = &Environment{
		parent:  nil,
		binding: BuiltinProcedures(),
	}
	return interpreter
}

func (i *Interpreter) PrintResult(dumpAST bool) {
	for _, result := range i.Eval(dumpAST) {
		fmt.Println(result)
	}
}

func (i *Interpreter) Eval(dumpAST bool) (results []string) {
	defer func() {
		if err := recover(); err != nil {
			results = append(results, fmt.Sprintf("*** ERROR: %s", err))
		}
	}()

	for i.Peek() != scanner.EOF {
		expression := i.Parser.Parse(i.topLevel)
		if dumpAST {
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
		i.DumpAST(object.(*Application).procedureVariable, indentLevel+1)
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
	case *Number:
		i.printWithIndent(fmt.Sprintf("Number(%s)", object), indentLevel)
	case *Boolean:
		i.printWithIndent(fmt.Sprintf("Boolean(%s)", object), indentLevel)
	case *Variable:
		i.printWithIndent(fmt.Sprintf("Variable(%s)", object.(*Variable).identifier), indentLevel)
	case *Definition:
		i.printWithIndent("Definition", indentLevel)
		i.DumpAST(object.(*Definition).variable, indentLevel+1)
		i.DumpAST(object.(*Definition).value, indentLevel+1)
	case *Procedure:
		i.printWithIndent("Procedure", indentLevel)
		i.DumpAST(object.(*Procedure).arguments, indentLevel+1)
		i.DumpAST(object.(*Procedure).body, indentLevel+1)
	}
}

func (i *Interpreter) printWithIndent(text string, indentLevel int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indentLevel), text)
}
