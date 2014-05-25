// A type for builtin functions

package scheme

import (
	"fmt"
)

type Subroutine struct {
	ObjectBase
	function func(Object) Object
}

func NewSubroutine(function func(Object) Object) *Subroutine {
	return &Subroutine{function: function}
}

func (s *Subroutine) String() string {
	return fmt.Sprintf("#<subr %s>", s.Bounder())
}

func (s *Subroutine) Eval() Object {
	return s
}

func (s *Subroutine) Invoke(argument Object) Object {
	return s.function(argument)
}

func (s *Subroutine) isProcedure() bool {
	return true
}
