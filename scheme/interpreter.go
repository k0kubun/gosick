package scheme

import (
	"strings"
)

type Interpreter struct {
	Parser
}

func NewInterpreter(source string) *Interpreter {
	interpreter := new(Interpreter)
	interpreter.Init(strings.NewReader(source))
	return interpreter
}

func (i *Interpreter) IndentLevel() int {
	return 0
}

func (i *Interpreter) Eval() {
	i.Parse()
}
