// Procedure is a type for scheme procedure, which is expressed
// by lambda syntax form, like (lambda (x) x)
// When procedure has free variable, free variable must be binded when
// procedure is generated.
// So all Procedures have variable binding by Environment type (when there is
// no free variable, Procedure has Environment which is empty).

package scheme

import (
	"fmt"
)

type Procedure struct {
	ObjectBase
	environment *Environment
	function    func(Object) Object
}

var builtinProcedures = Binding{
	"+":       NewProcedure(plus),
	"-":       NewProcedure(minus),
	"*":       NewProcedure(multiply),
	"/":       NewProcedure(divide),
	"number?": NewProcedure(isNumber),
}

func NewProcedure(function func(Object) Object) *Procedure {
	return &Procedure{
		environment: nil,
		function:    function,
	}
}

func (p *Procedure) Eval() Object {
	return p
}

func (p *Procedure) invoke(argument Object) Object {
	return p.function(argument)
}

//
// *** Builtin Procedures ***
//

func assertArgumentsMinimum(arguments Object, minimum int) bool {
	if !arguments.IsList() {
		panic("Compile Error: proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() < minimum {
		panic(fmt.Sprintf("Compile Error: procedure requires at least %d argument\n", minimum))
	}
	return true
}

func assertArgumentsEqual(arguments Object, length int) bool {
	if !arguments.IsList() {
		panic("Compile Error: proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() != length {
		panic(fmt.Sprintf("Compile Error: wrong number of arguments: number? requires %d, but got %d",
			length, arguments.(*Pair).ListLength()))
	}
	return true
}

func plus(arguments Object) Object {
	if !assertArgumentsMinimum(arguments, 0) {
		return nil
	}

	sum := 0
	for arguments != nil {
		pair := arguments.(*Pair)
		if pair == nil || pair.Car == nil {
			break
		}
		if car := pair.Car.Eval(); car != nil {
			number := car.(*Number)
			sum += number.value
		}
		arguments = pair.Cdr
	}
	return NewNumber(sum)
}

func minus(arguments Object) Object {
	if !assertArgumentsMinimum(arguments, 1) {
		return nil
	}

	pair := arguments.(*Pair)
	difference := pair.Car.Eval().(*Number).value
	list := pair.Cdr
	for {
		if list == nil || list.Car == nil {
			break
		}
		if car := list.Car.Eval(); car != nil {
			number := car.(*Number)
			difference -= number.value
		}
		list = list.Cdr
	}
	return NewNumber(difference)
}

func multiply(arguments Object) Object {
	if !assertArgumentsMinimum(arguments, 0) {
		return nil
	}

	product := 1
	for arguments != nil {
		pair := arguments.(*Pair)
		if pair == nil || pair.Car == nil {
			break
		}
		if car := pair.Car.Eval(); car != nil {
			number := car.(*Number)
			product *= number.value
		}
		arguments = pair.Cdr
	}
	return NewNumber(product)
}

func divide(arguments Object) Object {
	if !assertArgumentsMinimum(arguments, 1) {
		return nil
	}

	pair := arguments.(*Pair)
	quotient := pair.Car.Eval().(*Number).value
	list := pair.Cdr
	for {
		if list == nil || list.Car == nil {
			break
		}
		if car := list.Car.Eval(); car != nil {
			number := car.(*Number)
			quotient /= number.value
		}
		list = list.Cdr
	}
	return NewNumber(quotient)
}

func isNumber(arguments Object) Object {
	if !assertArgumentsEqual(arguments, 1) {
		return nil
	}
	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsNumber())
}
