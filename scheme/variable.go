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
	object := v.content()
	if object == nil {
		runtimeError("unbound variable: %s", v.identifier)
	}
	object.setBounder(v)
	return object
}

func (v *Variable) String() string {
	return v.identifier
}

func (v *Variable) content() Object {
	return v.boundedObject(v.identifier)
}

func (v *Variable) isVariable() bool {
	return true
}
