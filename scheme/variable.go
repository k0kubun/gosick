// Scheme's identifier is classified to a symbol or a variable.
// And this type owns a role to express a variable.
// Variable itself does not have a value for identifier,
// interpreter searches it from its code block scope by Variable's identifier.

package scheme

type Variable struct {
	ObjectBase
	identifier  string
	environment *Environment
}

func NewVariable(identifier string, environment *Environment) *Variable {
	return &Variable{
		identifier:  identifier,
		environment: environment,
	}
}

func (v *Variable) Eval() Object {
	return v.environment.boundedObject(v.identifier)
}

func (v *Variable) String() string {
	return v.Eval().String()
}

func (v *Variable) IsVariable() bool {
	return true
}
