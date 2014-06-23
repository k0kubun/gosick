package scheme

import (
	"testing"
)

type parserTest struct {
	source string
	result string
}

var parserTests = []parserTest{
	{"1", "1"},
	{"-2", "-2"},
	{"'12", "12"},
	{"()", "()"},
	{"'()", "()"},
	{"#f", "#f"},
	{"#t", "#t"},
	{"'#f", "#f"},
	{"'#t", "#t"},
	{"hello", "hello"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		i := NewInterpreter(test.source)
		i.Peek()
		object := i.Parse(nil)

		if object.String() != test.result {
			t.Errorf(
				"%s:\n  Expected:\n    %s\n  Got:\n    %s\n",
				test.source,
				test.result,
				object,
			)
			return
		}
	}
}
