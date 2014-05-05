// Object and ObjectBase is an abstract class for all scheme expressions.
// A return value of a method which returns scheme object is Object.
// And ObjectBase has Object's implementation of String().

package scheme

type Object interface {
	String() string
	IsNumber() bool
	IsList() bool
	IsApplication() bool
}

type ObjectBase struct {
}

func (o *ObjectBase) String() string {
	return "This type's String() is not implemented yet."
}

func (o *ObjectBase) IsNumber() bool {
	return false
}

func (o *ObjectBase) IsList() bool {
	return false
}

func (o *ObjectBase) IsApplication() bool {
	return false
}
