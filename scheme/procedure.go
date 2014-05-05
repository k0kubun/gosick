// Procedure is a type for scheme procedure, which is expressed
// by lambda syntax form, like (lambda (x) x)
// When procedure has free variable, free variable must be binded when
// procedure is generated.
// So all Procedures have variable binding by Environment type (when there is
// no free variable, Procedure has Environment which is empty).

package scheme

import (
	"log"
)

type Procedure struct {
	ObjectBase
	environment *Environment
	function    func(Object) Object
}

var builtinProcedures = Binding{
	"+": NewProcedure(plus),
	"-": NewProcedure(minus),
	"*": NewProcedure(multiply),
	"/": NewProcedure(divide),
}

func NewProcedure(function func(Object) Object) *Procedure {
	return &Procedure{
		environment: nil,
		function:    function,
	}
}

func (p *Procedure) invoke(argument Object) Object {
	return p.function(argument)
}

func plus(arguments Object) Object {
	sum := 0
	for arguments != nil {
		pair := arguments.(*Pair)
		if pair == nil {
			break
		}
		if car := pair.EvaledCar(); car != nil {
			number := car.(*Number)
			sum += number.value
		}
		arguments = pair.Cdr
	}
	return NewNumber(sum)
}

func minus(arguments Object) Object {
	switch arguments.(type) {
	case *Pair:
		pair := arguments.(*Pair)
		if pair.IsEmpty() {
			log.Print("procedure requires at least one argument: (-)")
			return nil
		}

		difference := pair.EvaledCar().(*Number).value
		list := pair.Cdr
		for {
			if list == nil {
				break
			}
			if car := list.EvaledCar(); car != nil {
				number := car.(*Number)
				difference -= number.value
			}
			list = list.Cdr
		}
		return NewNumber(difference)
	default:
		log.Print("procedure requires at least one argument: (-)")
		return nil
	}
}

func multiply(arguments Object) Object {
	product := 1
	for arguments != nil {
		pair := arguments.(*Pair)
		if pair == nil {
			break
		}
		if car := pair.EvaledCar(); car != nil {
			number := car.(*Number)
			product *= number.value
		}
		arguments = pair.Cdr
	}
	return NewNumber(product)
}

func divide(arguments Object) Object {
	switch arguments.(type) {
	case *Pair:
		pair := arguments.(*Pair)
		if pair.IsEmpty() {
			log.Print("procedure requires at least one argument: (/)")
			return nil
		}

		quotient := pair.EvaledCar().(*Number).value
		list := pair.Cdr
		for {
			if list == nil {
				break
			}
			if car := list.EvaledCar(); car != nil {
				number := car.(*Number)
				quotient /= number.value
			}
			list = list.Cdr
		}
		return NewNumber(quotient)
	default:
		log.Print("procedure requires at least one argument: (/)")
		return nil
	}
}
