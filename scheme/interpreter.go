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
	ObjectBase
	*Parser
	topLevel Binding
}

func NewInterpreter(source string) *Interpreter {
	return &Interpreter{Parser: NewParser(source), topLevel: BuiltinProcedures()}
}

// Load new source code with current environment
func (i *Interpreter) ReloadSourceCode(source string) {
	i.Parser = NewParser(source)
}

func (i *Interpreter) PrintResult(dumpAST bool) {
	for _, result := range i.EvalSource(dumpAST) {
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
		expression := i.Parser.Parse(i)
		if dumpAST {
			fmt.Printf("\n*** AST ***\n")
			i.DumpAST(expression, 0)
			fmt.Printf("\n*** Result ***\n")
		}

		if expression == nil {
			return
		}
		results = append(results, expression.Eval().String())
	}
	return
}

func (i *Interpreter) bind(identifier string, object Object) {
	i.topLevel[identifier] = object
}

func (i *Interpreter) binding() Binding {
	return i.topLevel
}

func (i *Interpreter) scopedBinding() Binding {
	return i.topLevel
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
	case *Definition:
		i.printWithIndent("Definition", indentLevel)
		i.DumpAST(object.(*Definition).variable, indentLevel+1)
		i.DumpAST(object.(*Definition).value, indentLevel+1)
	case *Procedure:
		i.printWithIndent("Procedure", indentLevel)
		i.DumpAST(object.(*Procedure).arguments, indentLevel+1)
		i.DumpAST(object.(*Procedure).body, indentLevel+1)
	case *Set:
		i.printWithIndent("Set", indentLevel)
		i.DumpAST(object.(*Set).variable, indentLevel+1)
		i.DumpAST(object.(*Set).value, indentLevel+1)
	}
}

func (i *Interpreter) printWithIndent(text string, indentLevel int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indentLevel), text)
}
