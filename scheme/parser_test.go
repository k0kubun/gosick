package scheme

import (
	"reflect"
	"testing"
)

type parserTest struct {
	source string
	result Object
}

var parserTests = []parserTest{
	{"1", NewNumber(1)},
	{"-2", NewNumber(-2)},
	{"()", Null},
	{"'12", NewNumber(12)},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		i := NewInterpreter(test.source)
		i.Peek()
		object := i.Parse(nil)

		if !reflect.DeepEqual(object, test.result) {
			t.Errorf(
				"%s:\n  Expected:\n    %#v\n  Got:\n    %#v\n",
				test.source,
				test.result,
				object,
			)
			return
		}
	}
}
