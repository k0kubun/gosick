package scheme

type Object interface {
	String() string
}

type ObjectBase struct {
}

func (o *ObjectBase) String() string {
	return "This type's Eval() is not implemented yet."
}
