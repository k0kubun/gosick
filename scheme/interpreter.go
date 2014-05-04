package scheme

import (
	"fmt"
	"text/scanner"
)

type Interpreter struct {
	*Parser
}

func NewInterpreter(source string) *Interpreter {
	interpreter := &Interpreter{NewParser(source)}
	return interpreter
}

func (i *Interpreter) IndentLevel() int {
	return 0
}

func (i *Interpreter) Eval() {
	for i.Peek() != scanner.EOF {
		expression := i.Parser.Parse()

		if expression == nil {
			return
		}
		fmt.Println(expression.String())
	}
}
