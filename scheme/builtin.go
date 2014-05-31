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

func DefaultBinding() Binding {
	binding := make(Binding)
	for key, value := range builtinProcedures {
		binding[key] = value
	}
	for key, value := range builtinSyntaxes {
		binding[key] = value
	}
	return binding
}

func plusProc(arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	sum := 0
	for _, number := range numbers {
		sum += number.(*Number).value
	}
	return NewNumber(sum)
}

func minusProc(arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	difference := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		difference -= number.(*Number).value
	}
	return NewNumber(difference)
}

func multiplyProc(arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	product := 1
	for _, number := range numbers {
		product *= number.(*Number).value
	}
	return NewNumber(product)
}

func divideProc(arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	quotient := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		quotient /= number.(*Number).value
	}
	return NewNumber(quotient)
}

func equalProc(arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a == b })
}

func lessThanProc(arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a < b })
}

func lessEqualProc(arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a <= b })
}

func greaterThanProc(arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a > b })
}

func greaterEqualProc(arguments Object) Object {
	return compareNumbers(arguments, func(a, b int) bool { return a >= b })
}

func isNumberProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isNumber() })
}

func isProcedureProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isProcedure() })
}

func isBooleanProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isBoolean() })
}

func isPairProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isPair() })
}

func isListProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isList() })
}

func isSymbolProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isSymbol() })
}

func isStringProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isString() })
}

func consProc(arguments Object) Object {
	assertListEqual(arguments, 2)
	objects := evaledObjects(arguments.(*Pair).Elements())

	return &Pair{
		ObjectBase: ObjectBase{parent: arguments.Parent()},
		Car:        objects[0],
		Cdr:        objects[1],
	}
}

func carProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Car
}

func cdrProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Cdr
}

func listProc(arguments Object) Object {
	return arguments
}

func setCarProc(arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Car = object
	return undef
}

func setCdrProc(arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Cdr = object
	return undef
}

func lengthProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	list := arguments.(*Pair).ElementAt(0).Eval()
	assertListMinimum(list, 0)

	return NewNumber(list.(*Pair).ListLength())
}

func memqProc(arguments Object) Object {
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

func lastProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	list := arguments.(*Pair).ElementAt(0).Eval()
	if !list.isPair() {
		runtimeError("pair required: %s", list)
	}
	assertListMinimum(list, 1)

	elements := list.(*Pair).Elements()
	return elements[len(elements)-1].Eval()
}

func appendProc(arguments Object) Object {
	assertListMinimum(arguments, 0)
	elements := evaledObjects(arguments.(*Pair).Elements())

	appendedList := NewPair(arguments)
	for _, element := range elements {
		appendedList = appendedList.AppendList(element)
	}

	return appendedList
}

func stringAppendProc(arguments Object) Object {
	assertListMinimum(arguments, 0)

	stringObjects := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(stringObjects, "string")

	texts := []string{}
	for _, stringObject := range stringObjects {
		texts = append(texts, stringObject.(*String).text)
	}
	return NewString(strings.Join(texts, ""))
}

func symbolToStringProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "symbol")
	return NewString(object.(*Symbol).identifier)
}

func stringToSymbolProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewSymbol(object.(*String).text)
}

func stringToNumberProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewNumber(object.(*String).text)
}

func numberToStringProc(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "number")
	return NewString(object.(*Number).value)
}

func areIdentical(a Object, b Object) bool {
	if typeName(a) != typeName(b) {
		return false
	}

	switch a.(type) {
	case *Number:
		return a.(*Number).value == b.(*Number).value
	case *Boolean:
		return a.(*Boolean).value == b.(*Boolean).value
	default:
		return a == b
	}
}

func areEqual(a Object, b Object) bool {
	if a == nil {
		return true
	}
	if typeName(a) != typeName(b) {
		return false
	} else if areIdentical(a, b) {
		return true
	}

	switch a.(type) {
	case *Pair:
		return areEqual(a.(*Pair).Car, b.(*Pair).Car) && areEqual(a.(*Pair).Cdr, b.(*Pair).Cdr)
	default:
		return false
	}
}

func areSameList(a Object, b Object) bool {
	if typeName(a) != typeName(b) {
		return false
	}

	switch a.(type) {
	case *Pair:
		return areSameList(a.(*Pair).Car, b.(*Pair).Car) && areSameList(a.(*Pair).Cdr, b.(*Pair).Cdr)
	default:
		return areIdentical(a, b)
	}
}

func isEqProc(arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areIdentical(objects[0], objects[1]))
}

func isNeqProc(arguments Object) Object {
	return NewBoolean(!isEqProc(arguments).(*Boolean).value)
}

func isEqualProc(arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areEqual(objects[0], objects[1]))
}

func loadProc(arguments Object) Object {
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

func writeProc(arguments Object) Object {
	assertListEqual(arguments, 1) // TODO: accept output port

	object := arguments.(*Pair).ElementAt(0).Eval()
	fmt.Printf("%s", object)
	return undef
}

func printProc(arguments Object) Object {
	assertListEqual(arguments, 1) // TODO: accept output port

	object := arguments.(*Pair).ElementAt(0).Eval()
	fmt.Printf("%s\n", object)
	return undef
}
