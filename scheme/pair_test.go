package scheme

import (
	"testing"
)

type lengthTest struct {
	source string
	length int
}

var lengthTests = []lengthTest{
	{"()", 0},
}

func TestListLength(t *testing.T) {
	for _, test := range lengthTests {
		p := NewParser(test.source)

		p.Peek()
		parsedObject := p.Parse()
		list := parsedObject.(*Pair)

		actual := list.ListLength()
		if actual != test.length {
			t.Errorf("%s => %d; want %d", test.source, actual, test.length)
		}
	}
}
