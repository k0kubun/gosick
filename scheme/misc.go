package scheme

import (
	"fmt"
	"strings"
)

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

func evaledObjects(objects []Object) []Object {
	evaledObjects := []Object{}

	for _, object := range objects {
		evaledObjects = append(evaledObjects, object.Eval())
	}
	return evaledObjects
}
