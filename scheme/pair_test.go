package scheme

import (
	"testing"
)

type lengthTest struct {
	source string
	length int
}

type listTest struct {
	source string
	result bool
}

var lengthTests = []lengthTest{
	{"()", 0},
	{"'((1 2) 3)", 2},
	{"'(1 2 3)", 3},
}

var listTests = []listTest{
	{"'()", true},
	{"'(1 2 3)", true},
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

func TestIsList(t *testing.T) {
	for _, test := range listTests {
		p := NewParser(test.source)

		p.Peek()
		parsedObject := p.Parse()
		list := parsedObject.(*Pair)

		actual := list.IsList()
		if actual != test.result {
			t.Errorf("%s => %s; want %s", test.source, actual, test.result)
		}
	}
}
