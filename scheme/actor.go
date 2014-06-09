package scheme

import (
	"fmt"
)

type Actor struct {
	ObjectBase
	receiver map[string]func(Object)
}

func NewActor() *Actor {
	return &Actor{}
}

func (a *Actor) Eval() Object {
	return a
}

func (a *Actor) Invoke(argument Object) Object {
	return undef
}

func (a *Actor) String() string {
	if a.Bounder() == nil {
		return "#<actor #f>"
	}
	return fmt.Sprintf("#<actor %s>", a.Bounder())
}
