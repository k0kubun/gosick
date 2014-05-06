// Object and ObjectBase is an abstract class for all scheme expressions.
// A return value of a method which returns scheme object is Object.
// And ObjectBase has Object's implementation of String().

package scheme

type Object interface {
	Eval() Object
	String() string
	IsNumber() bool
	IsBoolean() bool
	IsProcedure() bool
	IsNull() bool
	IsPair() bool
	IsList() bool
	IsSymbol() bool
	IsApplication() bool
}

type ObjectBase struct {
}

func (o *ObjectBase) Eval() Object {
	panic("This object's Eval() is not implemented yet.")
}

func (o *ObjectBase) String() string {
	panic("This object's String() is not implemented yet.")
}

func (o *ObjectBase) IsNumber() bool {
	return false
}

func (o *ObjectBase) IsBoolean() bool {
	return false
}

func (o *ObjectBase) IsProcedure() bool {
	return false
}

func (o *ObjectBase) IsNull() bool {
	return false
}

func (o *ObjectBase) IsPair() bool {
	return false
}

func (o *ObjectBase) IsList() bool {
	return false
}

func (o *ObjectBase) IsSymbol() bool {
	return false
}

func (o *ObjectBase) IsApplication() bool {
	return false
}
