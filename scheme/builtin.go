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
	"boolean?":   NewProcedure(isBoolean),
	"not":        NewProcedure(not),
}

func assertListMinimum(arguments Object, minimum int) {
	if !arguments.IsList() {
		panic("Compile Error: proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() < minimum {
		panic(fmt.Sprintf("Compile Error: procedure requires at least %d argument", minimum))
	}
}

func assertListEqual(arguments Object, length int) {
	if !arguments.IsList() {
		panic("Compile Error: proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() != length {
		panic(fmt.Sprintf("Compile Error: wrong number of arguments: number? requires %d, but got %d",
			length, arguments.(*Pair).ListLength()))
	}
}

func assertObjectsType(objects []Object, typeName string) {
	if typeName == "Number" {
		for _, object := range objects {
			if !object.IsNumber() {
				panic("Compile Error: procedure expects arguments to be Number")
			}
		}
	}
}

func evaledObjects(objects []Object) []Object {
	evaledObjects := []Object{}

	for _, object := range objects {
		evaledObjects = append(evaledObjects, object.Eval())
	}
	return evaledObjects
}

func plus(arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "Number")

	sum := 0
	for _, number := range numbers {
		sum += number.(*Number).value
	}
	return NewNumber(sum)
}

func minus(arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "Number")

	difference := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		difference -= number.(*Number).value
	}
	return NewNumber(difference)
}

func multiply(arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "Number")

	product := 1
	for _, number := range numbers {
		product *= number.(*Number).value
	}
	return NewNumber(product)
}

func divide(arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "Number")

	quotient := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		quotient /= number.(*Number).value
	}
	return NewNumber(quotient)
}

func equal(arguments Object) Object {
	assertListMinimum(arguments, 2)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "Number")

	firstValue := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		if firstValue != number.(*Number).value {
			return NewBoolean(false)
		}
	}
	return NewBoolean(true)
}

func isNumber(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsNumber())
}

func isNull(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsPair() && object.(*Pair).IsEmpty())
}

func isProcedure(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsProcedure())
}

func isBoolean(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsBoolean())
}

func not(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsBoolean() && !object.(*Boolean).value)
}
