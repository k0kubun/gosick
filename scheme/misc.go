package scheme

import (
	"fmt"
	"strings"
)

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
		compileError("wrong number of arguments: requires %d, but got %d",
			length, arguments.(*Pair).ListLength())
	}
}

func assertObjectsType(objects []Object, typeName string) {
	for _, object := range objects {
		assertObjectType(object, typeName)
	}
}

func assertObjectType(object Object, assertType string) {
	if assertType != typeName(object) {
		compileError("%s required, but got %s", assertType, object)
	}
}

func compileError(format string, a ...interface{}) Object {
	return runtimeError("Compile Error: "+format, a...)
}

func defaultBinding() Binding {
	binding := make(Binding)
	for key, value := range builtinProcedures {
		binding[key] = value
	}
	for key, value := range builtinSyntaxes {
		binding[key] = value
	}
	return binding
}

func evaledObjects(objects []Object) []Object {
	evaledObjects := []Object{}

	for _, object := range objects {
		evaledObjects = append(evaledObjects, object.Eval())
	}
	return evaledObjects
}

func runtimeError(format string, a ...interface{}) Object {
	panic(fmt.Sprintf(format, a...))
	return undef
}

func syntaxError(format string, a ...interface{}) Object {
	return compileError("syntax-error: "+format, a...)
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
