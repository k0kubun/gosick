package scheme

type Interpreter struct {
	Parser
}

func NewInterpreter(expression string) *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) IndentLevel() int {
	return 0
}

func (i *Interpreter) Eval() {
}
