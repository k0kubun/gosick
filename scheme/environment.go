// Environment has variable bindings.
// Interpreter has one Environment global variable for top-level environment.
// And each let block and procedure has Environment to hold its scope's variable binding.

package scheme

import (
	"log"
)

type Environment struct {
	parent  *Environment
	binding Binding
}

type Binding map[string]Object

var TopLevel = Environment{
	parent:  nil,
	binding: builtinProcedures,
}

func NewEnvironment() *Environment {
	return &Environment{}
}

func (e *Environment) Bind(identifier string, value Object) {
	e.binding[identifier] = value
}

// Search procedure which is binded with given variable from environment,
// and invoke the procedure with given arguments.
func (e *Environment) invokeProcedure(variable, arguments Object) Object {
	if variable == nil {
		log.Fatal("Invoked procedure for <nil> variable.")
	}
	identifier := variable.(*Variable).identifier
	procedure := e.binding[identifier].(*Procedure)
	if procedure == nil {
		log.Printf("Unbound variable: %s\n", identifier)
		return nil
	}
	return procedure.invoke(arguments)
}

func (e *Environment) scopedBinding() Binding {
	return Binding{}
}
