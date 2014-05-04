// Environment has variable bindings.
// Interpreter has one Environment global variable for top-level environment.
// And each let block and procedure has Environment to hold its scope's variable binding.

package scheme

type Environment struct {
	ObjectBase
	parent  *Environment
	binding *Binding
}

type Binding map[string]*Procedure

var TopLevel = Environment{
	parent:  nil,
	binding: builtinProcedures,
}

var builtinProcedures = Binding{
	"": nil,
}

func NewEnvironment() *Environment {
	return &Environment{}
}
