// Interpreter is a scheme source code interpreter.
// It owns a role of API for executing scheme program.
// Interpreter embeds Parser and delegates syntactic analysis to it.

package scheme

import (
	"fmt"
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

		if expression == nil {
			return
		}
		fmt.Println(expression.String())
	}
}
