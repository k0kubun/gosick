package scheme

import (
	"fmt"
)

type Actor struct {
	ObjectBase
	functions    map[string]func([]Object)
	receiver     chan []Object
	localBinding Binding
}

func NewActor(parent Object) *Actor {
	actor := &Actor{
		ObjectBase:   ObjectBase{parent: parent},
		functions:    make(map[string]func([]Object)),
		receiver:     make(chan []Object),
		localBinding: make(Binding),
	}
	actor.localBinding["self"] = actor
	return actor
}

func (a *Actor) Eval() Object {
	return a
}

func (a *Actor) Invoke(argument Object) Object {
	assertListMinimum(argument, 1)
	elements := argument.(*Pair).Elements()
	switch elements[0].(type) {
	case *Variable:
		switch elements[0].(*Variable).identifier {
		case "start":
			go a.Start()
		case "!":
			a.receiver <- elements[1:]
		default:
			runtimeError("unexpected method for actor: %s", elements[0].(*Variable).identifier)
		}
	}
	return undef
}

func (a *Actor) Start() {
	for {
		select {
		case received := <-a.receiver:
			if len(received) == 0 {
				continue
			}
			assertObjectType(received[0], "string")

			a.functions[received[0].(*String).text](received[1:])
		}
	}
}

func (a *Actor) String() string {
	if a.Bounder() == nil {
		return "#<actor #f>"
	}
	return fmt.Sprintf("#<actor %s>", a.Bounder())
}

func (a *Actor) tryDefine(variable Object, object Object) {
	if variable.isVariable() {
		a.localBinding[variable.(*Variable).identifier] = object
	}
}

func (a *Actor) define(identifier string, object Object) {
	a.localBinding[identifier] = object
}

func (a *Actor) set(identifier string, object Object) {
	if a.localBinding[identifier] == nil {
		if a.parent == nil {
			runtimeError("symbol not defined")
		} else {
			a.parent.set(identifier, object)
		}
	} else {
		a.localBinding[identifier] = object
	}
}

func (a *Actor) binding() Binding {
	return a.localBinding
}
