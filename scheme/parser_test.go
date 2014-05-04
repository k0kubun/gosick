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
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		actual := NewParser(test.source).Parse().String()
		if actual != test.result {
			t.Errorf("%s => %s; want %s", test.source, actual, test.result)
		}
	}
}
