package scheme

type Definition struct {
	ObjectBase
	environment *Environment
	variable    *Variable
	value       Object
}

func (d *Definition) String() string {
	TopLevel.Bind(d.variable.identifier, d.value)
	return d.variable.identifier
}
