// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	builtinProcedures = Binding{
		"+":              NewSubroutine(plusSubr),
		"-":              NewSubroutine(minusSubr),
		"*":              NewSubroutine(multiplySubr),
		"/":              NewSubroutine(divideSubr),
		"=":              NewSubroutine(equalSubr),
		"<":              NewSubroutine(lessThanSubr),
		"<=":             NewSubroutine(lessEqualSubr),
		">":              NewSubroutine(greaterThanSubr),
		">=":             NewSubroutine(greaterEqualSubr),
		"append":         NewSubroutine(appendSubr),
		"boolean?":       NewSubroutine(isBooleanSubr),
		"car":            NewSubroutine(carSubr),
		"cdr":            NewSubroutine(cdrSubr),
		"cons":           NewSubroutine(consSubr),
		"eq?":            NewSubroutine(isEqSubr),
		"equal?":         NewSubroutine(isEqualSubr),
		"last":           NewSubroutine(lastSubr),
		"length":         NewSubroutine(lengthSubr),
		"list":           NewSubroutine(listSubr),
		"list?":          NewSubroutine(isListSubr),
		"load":           NewSubroutine(loadSubr),
		"memq":           NewSubroutine(memqSubr),
		"neq?":           NewSubroutine(isNeqSubr),
		"number?":        NewSubroutine(isNumberSubr),
		"number->string": NewSubroutine(numberToStringSubr),
		"pair?":          NewSubroutine(isPairSubr),
		"print":          NewSubroutine(printSubr),
		"procedure?":     NewSubroutine(isProcedureSubr),
		"set-car!":       NewSubroutine(setCarSubr),
		"set-cdr!":       NewSubroutine(setCdrSubr),
		"string?":        NewSubroutine(isStringSubr),
		"string-append":  NewSubroutine(stringAppendSubr),
		"string->number": NewSubroutine(stringToNumberSubr),
		"string->symbol": NewSubroutine(stringToSymbolSubr),
		"symbol?":        NewSubroutine(isSymbolSubr),
		"symbol->string": NewSubroutine(symbolToStringSubr),
		"write":          NewSubroutine(writeSubr),
	}
)

func carSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Car
}

func cdrSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Cdr
}

func consSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)
	objects := evaledObjects(arguments.(*Pair).Elements())

	return &Pair{
		ObjectBase: ObjectBase{parent: arguments.Parent()},
		Car:        objects[0],
		Cdr:        objects[1],
	}
}

func divideSubr(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	quotient := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		quotient /= number.(*Number).value
	}
	return NewNumber(quotient)
}

func equalSubr(s *Subroutine, arguments Object) Object {
	return s.compareNumbers(arguments, func(a, b int) bool { return a == b })
}

func greaterThanSubr(s *Subroutine, arguments Object) Object {
	return s.compareNumbers(arguments, func(a, b int) bool { return a > b })
}

func greaterEqualSubr(s *Subroutine, arguments Object) Object {
	return s.compareNumbers(arguments, func(a, b int) bool { return a >= b })
}

func lengthSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	list := arguments.(*Pair).ElementAt(0).Eval()
	assertListMinimum(list, 0)

	return NewNumber(list.(*Pair).ListLength())
}

func lessEqualSubr(s *Subroutine, arguments Object) Object {
	return s.compareNumbers(arguments, func(a, b int) bool { return a <= b })
}

func lessThanSubr(s *Subroutine, arguments Object) Object {
	return s.compareNumbers(arguments, func(a, b int) bool { return a < b })
}

func listSubr(s *Subroutine, arguments Object) Object {
	return arguments
}

func memqSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	searchObject := arguments.(*Pair).ElementAt(0).Eval()
	list := arguments.(*Pair).ElementAt(1).Eval()

	for {
		switch list.(type) {
		case *Pair:
			if areIdentical(list.(*Pair).Car, searchObject) {
				return list
			}
		default:
			break
		}

		if list = list.(*Pair).Cdr; list == nil {
			break
		}
	}
	return NewBoolean(false)
}

func minusSubr(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	difference := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		difference -= number.(*Number).value
	}
	return NewNumber(difference)
}

func multiplySubr(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	product := 1
	for _, number := range numbers {
		product *= number.(*Number).value
	}
	return NewNumber(product)
}

func lastSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	list := arguments.(*Pair).ElementAt(0).Eval()
	if !list.isPair() {
		runtimeError("pair required: %s", list)
	}
	assertListMinimum(list, 1)

	elements := list.(*Pair).Elements()
	return elements[len(elements)-1].Eval()
}

func appendSubr(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)
	elements := evaledObjects(arguments.(*Pair).Elements())

	appendedList := NewPair(arguments)
	for _, element := range elements {
		appendedList = appendedList.AppendList(element)
	}

	return appendedList
}

func numberToStringSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "number")
	return NewString(object.(*Number).value)
}

func isBooleanSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isBoolean() })
}

func isEqSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areIdentical(objects[0], objects[1]))
}

func isEqualSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areEqual(objects[0], objects[1]))
}

func isListSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isList() })
}

func isNeqSubr(s *Subroutine, arguments Object) Object {
	return NewBoolean(!isEqSubr(s, arguments).(*Boolean).value)
}

func isNumberSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isNumber() })
}

func isPairSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isPair() })
}

func isProcedureSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isProcedure() })
}

func isSymbolSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isSymbol() })
}

func isStringSubr(s *Subroutine, arguments Object) Object {
	return s.booleanByFunc(arguments, func(object Object) bool { return object.isString() })
}

func loadSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")

	buffer, err := ioutil.ReadFile(object.(*String).text)
	if err != nil {
		runtimeError("cannot find \"%s\"", object.(*String).text)
		return nil
	}

	parser := NewParser(string(buffer))
	for parser.Peek() != EOF {
		expression := parser.Parse(arguments.Parent())
		if expression != nil {
			expression.Eval()
		}
	}

	return NewBoolean(true)
}

func plusSubr(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	sum := 0
	for _, number := range numbers {
		sum += number.(*Number).value
	}
	return NewNumber(sum)
}

func printSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1) // TODO: accept output port

	object := arguments.(*Pair).ElementAt(0).Eval()
	fmt.Printf("%s\n", object)
	return undef
}

func setCarSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Car = object
	return undef
}

func setCdrSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Cdr = object
	return undef
}

func stringAppendSubr(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)

	stringObjects := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(stringObjects, "string")

	texts := []string{}
	for _, stringObject := range stringObjects {
		texts = append(texts, stringObject.(*String).text)
	}
	return NewString(strings.Join(texts, ""))
}

func stringToNumberSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewNumber(object.(*String).text)
}

func symbolToStringSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "symbol")
	return NewString(object.(*Symbol).identifier)
}

func stringToSymbolSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewSymbol(object.(*String).text)
}

func writeSubr(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1) // TODO: accept output port

	object := arguments.(*Pair).ElementAt(0).Eval()
	fmt.Printf("%s", object)
	return undef
}
