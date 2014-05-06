// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"fmt"
	"strings"
)

var builtinProcedures = Binding{
	"+":             NewProcedure(plus),
	"-":             NewProcedure(minus),
	"*":             NewProcedure(multiply),
	"/":             NewProcedure(divide),
	"=":             NewProcedure(equal),
	"number?":       NewProcedure(isNumber),
	"null?":         NewProcedure(isNull),
	"procedure?":    NewProcedure(isProcedure),
	"boolean?":      NewProcedure(isBoolean),
	"pair?":         NewProcedure(isPair),
	"list?":         NewProcedure(isList),
	"symbol?":       NewProcedure(isSymbol),
	"string?":       NewProcedure(isString),
	"not":           NewProcedure(not),
	"car":           NewProcedure(car),
	"cdr":           NewProcedure(cdr),
	"list":          NewProcedure(list),
	"string-append": NewProcedure(stringAppend),
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
	for _, object := range objects {
		if typeName == "Number" && !object.IsNumber() {
			panic("Compile Error: number required")
		} else if typeName == "String" && !object.IsString() {
			panic("Compile Error: string required")
		}
	}
}

func assertPair(object Object) {
	switch object.(type) {
	case *Pair:
		if object.IsPair() {
			return
		}
	}
	panic("Compile Error: pair required")
}

func evaledObjects(objects []Object) []Object {
	evaledObjects := []Object{}

	for _, object := range objects {
		evaledObjects = append(evaledObjects, object.Eval())
	}
	return evaledObjects
}

func typeCheckResult(arguments Object, typeCheckFunc func(Object) bool) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(typeCheckFunc(object))
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
	return typeCheckResult(arguments, func(object Object) bool { return object.IsNumber() })
}

func isNull(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsNull() })
}

func isProcedure(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsProcedure() })
}

func isBoolean(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsBoolean() })
}

func isPair(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsPair() })
}

func isList(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsList() })
}

func isSymbol(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsSymbol() })
}

func isString(arguments Object) Object {
	return typeCheckResult(arguments, func(object Object) bool { return object.IsString() })
}

func not(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(object.IsBoolean() && !object.(*Boolean).value)
}

func car(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertPair(object)
	return object.(*Pair).Car
}

func cdr(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertPair(object)
	return object.(*Pair).Cdr
}

func list(arguments Object) Object {
	return arguments
}

func stringAppend(arguments Object) Object {
	assertListMinimum(arguments, 0)

	stringObjects := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(stringObjects, "String")

	texts := []string{}
	for _, stringObject := range stringObjects {
		texts = append(texts, stringObject.(*String).text)
	}
	return NewString(strings.Join(texts, ""))
}
