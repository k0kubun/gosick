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

type If struct {
	ObjectBase
	condition Object
	trueBody  Object
	falseBody Object
}

func NewIf(parent Object) *If {
	return &If{ObjectBase: ObjectBase{parent: parent}}
}

func (i *If) Eval() Object {
	result := i.condition.Eval()
	if result.isBoolean() && result.(*Boolean).value {
		return i.trueBody.Eval()
	} else {
		return i.falseBody.Eval()
	}
}

type Cond struct {
	ObjectBase
	cases    []Object
	elseBody Object
}

func NewCond(parent Object) *Cond {
	return &Cond{ObjectBase: ObjectBase{parent: parent}}
}

func (c *Cond) Eval() Object {
	for _, caseBody := range c.cases {
		elements := caseBody.(*Pair).Elements()
		lastResult := elements[0].Eval()

		if lastResult.isBoolean() && lastResult.(*Boolean).value == false {
			continue
		}

		for _, element := range elements {
			lastResult = element.Eval()
		}
		return lastResult
	}

	if c.elseBody == nil {
		return NewSymbol("#<undef>")
	}

	elements := c.elseBody.(*Pair).Elements()
	lastResult := Object(NewSymbol("#<undef>"))
	for _, element := range elements {
		lastResult = element.Eval()
	}
	return lastResult
}

type And struct {
	ObjectBase
	body Object
}

func NewAnd(parent Object) *And {
	return &And{ObjectBase: ObjectBase{parent: parent}}
}

func (a *And) Eval() Object {
	lastResult := Object(NewBoolean(true))
	for _, object := range a.body.(*Pair).Elements() {
		lastResult = object.Eval()
		if lastResult.isBoolean() && lastResult.(*Boolean).value == false {
			return NewBoolean(false)
		}
	}
	return lastResult
}

type Or struct {
	ObjectBase
	body Object
}

func NewOr(parent Object) *Or {
	return &Or{ObjectBase: ObjectBase{parent: parent}}
}

func (o *Or) Eval() Object {
	lastResult := Object(NewBoolean(false))
	for _, object := range o.body.(*Pair).Elements() {
		lastResult = object.Eval()
		if !lastResult.isBoolean() || lastResult.(*Boolean).value != false {
			return lastResult
		}
	}
	return lastResult
}

type Begin struct {
	ObjectBase
	body Object
}

func NewBegin(parent Object) *Begin {
	return &Begin{ObjectBase: ObjectBase{parent: parent}}
}

func (b *Begin) Eval() Object {
	lastResult := Object(NewSymbol("#<undef>"))
	for _, object := range b.body.(*Pair).Elements() {
		lastResult = object.Eval()
	}
	return lastResult
}
