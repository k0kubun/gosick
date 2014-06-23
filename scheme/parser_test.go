package scheme

import (
	"reflect"
	"testing"
)

type easyParserTest struct {
	source string
	result string
}

type deepParserTest struct {
	source string
	result Object
}

var easyParserTests = []easyParserTest{
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
	{"(+)", "(+)"},
	{"(- 1)", "(- 1)"},
}

var deepParserTests = []deepParserTest{
	{
		"(+)",
		func() Object {
			app := &Application{
				procedure: &Variable{identifier: "+"},
				arguments: Null,
			}
			app.procedure.setParent(app)
			return app
		}(),
	},
	{
		"(- 1)",
		func() Object {
			pair := &Pair{
				Cdr: Null,
			}
			pair.Car = NewNumber(1, pair)
			app := &Application{
				procedure: &Variable{identifier: "-"},
				arguments: pair,
			}
			app.procedure.setParent(app)
			pair.setParent(app)
			return app
		}(),
	},
}

func TestParser(t *testing.T) {
	for _, test := range easyParserTests {
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

	for _, test := range deepParserTests {
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
