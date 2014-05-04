// Environment has variable bindings.
// Interpreter has one Environment global variable for top-level environment.
// And each let block and procedure has Environment to hold its scope's variable binding.

package scheme

import (
	"log"
)

type Environment struct {
	ObjectBase
	parent  *Environment
	binding Binding
}

type Binding map[string]*Procedure

var TopLevel = Environment{
	parent:  nil,
	binding: builtinProcedures,
}

func NewEnvironment() *Environment {
	return &Environment{}
}

// Search procedure which is binded with given variable from environment,
// and invoke the procedure with given arguments.
func (e *Environment) invokeProcedure(variable, arguments Object) Object {
	if variable == nil {
		log.Fatal("Invoked procedure for <nil> variable.")
	}
	procedure := TopLevel.binding[variable.(*Variable).identifier]
	return procedure.invoke(arguments)
}
