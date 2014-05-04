package scheme

import (
	"testing"
)

type parserTest struct {
	source string
	result string
}

var parserTests = []parserTest{
	{"()", "()"},
	{"10", "10"},
	{"123456789", "123456789"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		parsedObject := NewParser(test.source).Parse()
		if parsedObject == nil {
			t.Errorf("%s => <nil>; want %s", test.source, test.result)
			return
		}
		actual := parsedObject.String()
		if actual != test.result {
			t.Errorf("%s => %s; want %s", test.source, actual, test.result)
		}
	}
}