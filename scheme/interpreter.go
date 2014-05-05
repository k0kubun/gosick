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
}

func NewInterpreter(source string) *Interpreter {
	return &Interpreter{NewParser(source)}
}

func (i *Interpreter) IndentLevel() int {
	return 0
}

func (i *Interpreter) Eval() {
	for i.Peek() != scanner.EOF {
		expression := i.Parser.Parse()
		i.DumpAST(expression)

		if expression == nil {
			return
		}
		fmt.Println(expression.String())
	}
}

func (i *Interpreter) DumpAST(object Object) {
	i.dumpASTWithIndent(object, 0)
}

func (i *Interpreter) dumpASTWithIndent(object Object, indentLevel int) {
	if object == nil {
		return
	}
	switch object.(type) {
	case *Application:
		i.printWithIndent("Application", indentLevel)
		i.dumpASTWithIndent(object.(*Application).procedureVariable, indentLevel+1)
		i.dumpASTWithIndent(object.(*Application).arguments, indentLevel+1)
	case *Pair:
		pair := object.(*Pair)
		if pair.Car == nil && pair.Cdr == nil {
			i.printWithIndent("()", indentLevel)
			return
		}
		i.printWithIndent("Pair", indentLevel)
		i.dumpASTWithIndent(pair.Car, indentLevel+1)
		i.dumpASTWithIndent(pair.Cdr, indentLevel+1)
	case *Number:
		i.printWithIndent(fmt.Sprintf("Number(%d)", object.(*Number).value), indentLevel)
	case *Boolean:
		i.printWithIndent(fmt.Sprintf("Boolean(%s)", object.(*Boolean).String()), indentLevel)
	case *Variable:
		i.printWithIndent(fmt.Sprintf("Variable(%s)", object.(*Variable).identifier), indentLevel)
	}
}

func (i *Interpreter) printWithIndent(text string, indentLevel int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indentLevel), text)
}
