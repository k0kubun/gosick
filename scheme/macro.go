package scheme

import (
	"fmt"
)

type Macro struct {
	ObjectBase
}

func NewMacro() *Macro {
	return &Macro{}
}

func (m *Macro) Eval() Object {
	return m
}

func (m *Macro) String() string {
	if m.Bounder() == nil {
		return "#<macro #f>"
	}
	return fmt.Sprintf("#<macro %s>", m.Bounder())
}
