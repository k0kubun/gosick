// Scheme's identifier is classified to a symbol or a variable.
// And this type owns a role to express a variable.
// Variable itself does not have a value for identifier,
// interpreter searches it from its code block scope by Variable's identifier.

package scheme

type Variable struct {
	ObjectBase
	identifier string
}

func NewVariable(identifier string, parent Object) *Variable {
	return &Variable{
		ObjectBase: ObjectBase{parent: parent},
		identifier: identifier,
	}
}

func (v *Variable) Eval() Object {
	object := v.boundedObject(v.identifier)
	if object == nil {
		runtimeError("Unbound variable: %s", v.identifier)
	}
	return object
}

func (v *Variable) String() string {
	return v.Eval().String()
}

func (v *Variable) isVariable() bool {
	return true
}
