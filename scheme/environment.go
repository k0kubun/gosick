// Environment has variable bindings.
// Interpreter has one Environment global variable for top-level environment.
// And each let block and procedure has Environment to hold its scope's variable binding.

package scheme

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

// Returns ultimate-ancestral environment.
// This returns virtual top level environment in closure,
// which is separated from TopLevel.
func (e *Environment) topLevel() *Environment {
	environment := e
	for environment.parent != nil {
		environment = environment.parent
	}
	return environment
}

// Search procedure which is binded with given variable from environment,
// and invoke the procedure with given arguments.
func (e *Environment) invokeProcedure(object, arguments Object) Object {
	if object == nil {
		runtimeError("Invoked procedure for <nil> variable.")
	}

	evaledObject := object.Eval()
	if !evaledObject.IsProcedure() {
		runtimeError("invalid application")
	}
	procedure := evaledObject.(*Procedure)
	return procedure.Invoke(arguments)
}

func (e *Environment) boundedObject(identifier string) Object {
	object := e.scopedBinding()[identifier]
	if object == nil {
		runtimeError("Unbound variable: %s", identifier)
	}
	return e.scopedBinding()[identifier]
}

func (e *Environment) scopedBinding() Binding {
	scopedBinding := make(map[string]Object)
	environment := e

	for environment != nil {
		for identifier, object := range environment.binding {
			if scopedBinding[identifier] == nil {
				scopedBinding[identifier] = object
			}
		}
		environment = e.parent
	}
	return scopedBinding
}
