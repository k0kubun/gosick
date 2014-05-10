// This file is for statements by syntax form, such as set!

package scheme

type Set struct {
	ObjectBase
	variable Object
	value    Object
}

func NewSet(parent Object) *Set {
	return &Set{ObjectBase: ObjectBase{parent: parent}}
}

func (s *Set) Eval() Object {
	variable := s.variable.Eval()
	if !variable.isVariable() {
		variable = variable.Bounder()
		if variable == nil {
			compileError("syntax-error: malformed set!")
		}
	}

	value := s.value.Eval()
	s.ObjectBase.updateBinding(variable.(*Variable).identifier, value)
	return NewSymbol("#<undef>")
}
