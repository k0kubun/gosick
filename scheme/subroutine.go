// A type for builtin functions

package scheme

import (
	"fmt"
)

type Subroutine struct {
	ObjectBase
	function func(*Subroutine, Object) Object
}

func NewSubroutine(function func(*Subroutine, Object) Object) *Subroutine {
	return &Subroutine{function: function}
}

func (s *Subroutine) String() string {
	return fmt.Sprintf("#<subr %s>", s.Bounder())
}

func (s *Subroutine) Eval() Object {
	return s
}

func (s *Subroutine) Invoke(argument Object) Object {
	return s.function(s, argument)
}

func (s *Subroutine) isProcedure() bool {
	return true
}

func booleanByFunc(arguments Object, typeCheckFunc func(Object) bool) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(typeCheckFunc(object))
}

func compareNumbers(arguments Object, compareFunc func(int, int) bool) Object {
	assertListMinimum(arguments, 2)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	oldValue := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		if !compareFunc(oldValue, number.(*Number).value) {
			return NewBoolean(false)
		}
		oldValue = number.(*Number).value
	}
	return NewBoolean(true)
}
