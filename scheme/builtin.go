// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"fmt"
)

var builtinProcedures = Binding{
	"+":          NewProcedure(plus),
	"-":          NewProcedure(minus),
	"*":          NewProcedure(multiply),
	"/":          NewProcedure(divide),
	"=":          NewProcedure(equal),
	"number?":    NewProcedure(isNumber),
	"null?":      NewProcedure(isNull),
	"procedure?": NewProcedure(isProcedure),
}

func assertArgumentsMinimum(arguments Object, minimum int) {
	if !arguments.IsList() {
		panic("Compile Error: proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() < minimum {
		panic(fmt.Sprintf("Compile Error: procedure requires at least %d argument", minimum))
	}
}

func assertArgumentsEqual(arguments Object, length int) {
	if !arguments.IsList() {
		panic("Compile Error: proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() != length {
		panic(fmt.Sprintf("Compile Error: wrong number of arguments: number? requires %d, but got %d",
			length, arguments.(*Pair).ListLength()))
	}
}

func plus(arguments Object) Object {
	assertArgumentsMinimum(arguments, 0)

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
	assertArgumentsMinimum(arguments, 1)

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
	assertArgumentsMinimum(arguments, 0)

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
	assertArgumentsMinimum(arguments, 1)

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

func equal(arguments Object) Object {
	assertArgumentsMinimum(arguments, 2)

	list := arguments.(*Pair)
	length := list.ListLength()
	firstNumber := list.ElementAt(0).Eval().(*Number)
	for i := 1; i < length; i++ {
		if firstNumber.value != list.ElementAt(i).(*Number).value {
			return NewBoolean(false)
		}
	}
	return NewBoolean(true)
}

func isNumber(arguments Object) Object {
	assertArgumentsEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsNumber())
}

func isNull(arguments Object) Object {
	assertArgumentsEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	switch object.(type) {
	case *Pair:
		return NewBoolean(object.(*Pair).IsEmpty())
	default:
		return NewBoolean(false)
	}
}

func isProcedure(arguments Object) Object {
	assertArgumentsEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsProcedure())
}
