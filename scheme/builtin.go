// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	builtinProcedures = Binding{
		"+":              NewSubroutine(plusProc),
		"-":              NewSubroutine(minusProc),
		"*":              NewSubroutine(multiplyProc),
		"/":              NewSubroutine(divideProc),
		"=":              NewSubroutine(equalProc),
		"<":              NewSubroutine(lessThanProc),
		"<=":             NewSubroutine(lessEqualProc),
		">":              NewSubroutine(greaterThanProc),
		">=":             NewSubroutine(greaterEqualProc),
		"append":         NewSubroutine(appendProc),
		"boolean?":       NewSubroutine(isBooleanProc),
		"car":            NewSubroutine(carProc),
		"cdr":            NewSubroutine(cdrProc),
		"cons":           NewSubroutine(consProc),
		"eq?":            NewSubroutine(isEqProc),
		"equal?":         NewSubroutine(isEqualProc),
		"last":           NewSubroutine(lastProc),
		"length":         NewSubroutine(lengthProc),
		"list":           NewSubroutine(listProc),
		"list?":          NewSubroutine(isListProc),
		"load":           NewSubroutine(loadProc),
		"memq":           NewSubroutine(memqProc),
		"neq?":           NewSubroutine(isNeqProc),
		"number?":        NewSubroutine(isNumberProc),
		"number->string": NewSubroutine(numberToStringProc),
		"pair?":          NewSubroutine(isPairProc),
		"print":          NewSubroutine(printProc),
		"procedure?":     NewSubroutine(isProcedureProc),
		"set-car!":       NewSubroutine(setCarProc),
		"set-cdr!":       NewSubroutine(setCdrProc),
		"string?":        NewSubroutine(isStringProc),
		"string-append":  NewSubroutine(stringAppendProc),
		"string->number": NewSubroutine(stringToNumberProc),
		"string->symbol": NewSubroutine(stringToSymbolProc),
		"symbol?":        NewSubroutine(isSymbolProc),
		"symbol->string": NewSubroutine(symbolToStringProc),
		"write":          NewSubroutine(writeProc),
	}
)

func plusProc(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	sum := 0
	for _, number := range numbers {
		sum += number.(*Number).value
	}
	return NewNumber(sum)
}

func minusProc(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	difference := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		difference -= number.(*Number).value
	}
	return NewNumber(difference)
}

func multiplyProc(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	product := 1
	for _, number := range numbers {
		product *= number.(*Number).value
	}
	return NewNumber(product)
}

func divideProc(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	quotient := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		quotient /= number.(*Number).value
	}
	return NewNumber(quotient)
}

func equalProc(s *Subroutine, arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a == b })
}

func lessThanProc(s *Subroutine, arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a < b })
}

func lessEqualProc(s *Subroutine, arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a <= b })
}

func greaterThanProc(s *Subroutine, arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a > b })
}

func greaterEqualProc(s *Subroutine, arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a >= b })
}

func isNumberProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isNumber() })
}

func isProcedureProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isProcedure() })
}

func isBooleanProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isBoolean() })
}

func isPairProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isPair() })
}

func isListProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isList() })
}

func isSymbolProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isSymbol() })
}

func isStringProc(s *Subroutine, arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isString() })
}

func consProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)
	objects := evaledObjects(arguments.(*Pair).Elements())

	return &Pair{
		ObjectBase: ObjectBase{parent: arguments.Parent()},
		Car:        objects[0],
		Cdr:        objects[1],
	}
}

func carProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Car
}

func cdrProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Cdr
}

func listProc(s *Subroutine, arguments Object) Object {
	return arguments
}

func setCarProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Car = object
	return undef
}

func setCdrProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Cdr = object
	return undef
}

func lengthProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	list := arguments.(*Pair).ElementAt(0).Eval()
	assertListMinimum(list, 0)

	return NewNumber(list.(*Pair).ListLength())
}

func memqProc(s *Subroutine, arguments Object) Object {
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

func lastProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	list := arguments.(*Pair).ElementAt(0).Eval()
	if !list.isPair() {
		runtimeError("pair required: %s", list)
	}
	assertListMinimum(list, 1)

	elements := list.(*Pair).Elements()
	return elements[len(elements)-1].Eval()
}

func appendProc(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)
	elements := evaledObjects(arguments.(*Pair).Elements())

	appendedList := NewPair(arguments)
	for _, element := range elements {
		appendedList = appendedList.AppendList(element)
	}

	return appendedList
}

func stringAppendProc(s *Subroutine, arguments Object) Object {
	assertListMinimum(arguments, 0)

	stringObjects := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(stringObjects, "string")

	texts := []string{}
	for _, stringObject := range stringObjects {
		texts = append(texts, stringObject.(*String).text)
	}
	return NewString(strings.Join(texts, ""))
}

func symbolToStringProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "symbol")
	return NewString(object.(*Symbol).identifier)
}

func stringToSymbolProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewSymbol(object.(*String).text)
}

func stringToNumberProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewNumber(object.(*String).text)
}

func numberToStringProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "number")
	return NewString(object.(*Number).value)
}

func isEqProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areIdentical(objects[0], objects[1]))
}

func isNeqProc(s *Subroutine, arguments Object) Object {
	return NewBoolean(!isEqProc(s, arguments).(*Boolean).value)
}

func isEqualProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areEqual(objects[0], objects[1]))
}

func loadProc(s *Subroutine, arguments Object) Object {
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

func writeProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1) // TODO: accept output port

	object := arguments.(*Pair).ElementAt(0).Eval()
	fmt.Printf("%s", object)
	return undef
}

func printProc(s *Subroutine, arguments Object) Object {
	assertListEqual(arguments, 1) // TODO: accept output port

	object := arguments.(*Pair).ElementAt(0).Eval()
	fmt.Printf("%s\n", object)
	return undef
}
