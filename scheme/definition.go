package scheme

type Definition struct {
	ObjectBase
	environment *Environment
	variable    *Variable
	value       Object
}

func (d *Definition) Eval() Object {
	TopLevel.Bind(d.variable.identifier, d.value.Eval())
	return NewSymbol(d.variable.identifier)
}

func (d *Definition) String() string {
	return d.Eval().String()
}
