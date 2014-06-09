package scheme

import (
	"fmt"
)

type Actor struct {
	ObjectBase
	functions map[string]func([]Object)
	receiver  chan []Object
}

func NewActor() *Actor {
	return &Actor{
		functions: make(map[string]func([]Object)),
		receiver:  make(chan []Object),
	}
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
			if len(received) > 0 {
				assertObjectType(received[0], "string")
				a.functions[received[0].(*String).text](received[1:])
			}
		}
	}
}

func (a *Actor) String() string {
	if a.Bounder() == nil {
		return "#<actor #f>"
	}
	return fmt.Sprintf("#<actor %s>", a.Bounder())
}
