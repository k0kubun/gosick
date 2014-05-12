// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func BuiltinProcedures() Binding {
	return Binding{
		"+":              builtinProcedure(plusProc),
		"-":              builtinProcedure(minusProc),
		"*":              builtinProcedure(multiplyProc),
		"/":              builtinProcedure(divideProc),
		"=":              builtinProcedure(equalProc),
		"<":              builtinProcedure(lessThanProc),
		"<=":             builtinProcedure(lessEqualProc),
		">":              builtinProcedure(greaterThanProc),
		">=":             builtinProcedure(greaterEqualProc),
		"number?":        builtinProcedure(isNumberProc),
		"null?":          builtinProcedure(isNullProc),
		"procedure?":     builtinProcedure(isProcedureProc),
		"boolean?":       builtinProcedure(isBooleanProc),
		"pair?":          builtinProcedure(isPairProc),
		"list?":          builtinProcedure(isListProc),
		"symbol?":        builtinProcedure(isSymbolProc),
		"string?":        builtinProcedure(isStringProc),
		"not":            builtinProcedure(notProc),
		"cons":           builtinProcedure(consProc),
		"car":            builtinProcedure(carProc),
		"cdr":            builtinProcedure(cdrProc),
		"set-car!":       builtinProcedure(setCarProc),
		"set-cdr!":       builtinProcedure(setCdrProc),
		"list":           builtinProcedure(listProc),
		"length":         builtinProcedure(lengthProc),
		"memq":           builtinProcedure(memqProc),
		"last":           builtinProcedure(lastProc),
		"append":         builtinProcedure(appendProc),
		"string-append":  builtinProcedure(stringAppendProc),
		"symbol->string": builtinProcedure(symbolToStringProc),
		"string->symbol": builtinProcedure(stringToSymbolProc),
		"number->string": builtinProcedure(numberToStringProc),
		"string->number": builtinProcedure(stringToNumberProc),
		"eq?":            builtinProcedure(isEqProc),
		"neq?":           builtinProcedure(isNeqProc),
		"equal?":         builtinProcedure(isEqualProc),
		"load":           builtinProcedure(loadProc),
	}
}

func builtinProcedure(function func(Object) Object) *Procedure {
	return &Procedure{function: function}
}

func assertListMinimum(arguments Object, minimum int) {
	if !arguments.isList() {
		compileError("proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() < minimum {
		compileError("procedure requires at least %d argument", minimum)
	}
}

func assertListEqual(arguments Object, length int) {
	if !arguments.isList() {
		compileError("proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() != length {
		compileError("wrong number of arguments: number? requires %d, but got %d",
			length, arguments.(*Pair).ListLength())
	}
}

func assertObjectsType(objects []Object, typeName string) {
	for _, object := range objects {
		assertObjectType(object, typeName)
	}
}

func typeName(object Object) string {
	switch object.(type) {
	case *Pair:
		if object.isNull() {
			return "null"
		} else {
			return "pair"
		}
	default:
		rawTypeName := fmt.Sprintf("%T", object)
		typeName := strings.Replace(rawTypeName, "*scheme.", "", 1)
		return strings.ToLower(typeName)
	}
}

func assertObjectType(object Object, assertType string) {
	if assertType != typeName(object) {
		compileError("%s required, but got %s", assertType, object)
	}
}

func evaledObjects(objects []Object) []Object {
	evaledObjects := []Object{}

	for _, object := range objects {
		evaledObjects = append(evaledObjects, object.Eval())
	}
	return evaledObjects
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

func isNullProc(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool { return object.isNull() })
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

func notProc(arguments Object) Object {
	return booleanByFunc(arguments,
		func(object Object) bool { return object.isBoolean() && !object.(*Boolean).value })
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
	return NewSymbol("#<undef>")
}

func setCdrProc(arguments Object) Object {
	assertListEqual(arguments, 2)

	object := arguments.(*Pair).ElementAt(1).Eval()
	pair := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(pair, "pair")

	pair.(*Pair).Cdr = object
	return NewSymbol("#<undef>")
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

	appendedList := NewNull(arguments)
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
	case *Symbol:
		return a.(*Symbol).identifier == b.(*Symbol).identifier
	case *Boolean:
		return a.(*Boolean).value == b.(*Boolean).value
	default:
		return false
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
