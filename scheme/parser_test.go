package scheme

import (
	"testing"
)

func testClosedScope(t *testing.T, closedObject Object, source string) {
	if closedObject.Parent() != nil {
		t.Errorf("Test Case: %s\n%T.Parent():\nwant: nil\n got:%#v", source, closedObject, closedObject.Parent())
	}
	testTreeStructure(t, closedObject, source)
}

func testParentRelationship(t *testing.T, parent Object, child Object, source string) {
	if parent != child.Parent() {
		t.Errorf("Test Case: %s\nParent: %T, Child: %T, Child.Parent()\n want: %#v\n got: %#v\n", source, parent, child, parent, child.Parent())
	}
	testTreeStructure(t, child, source)
}

func testTreeStructure(t *testing.T, object Object, source string) {
	if object == nil {
		return
	}
	switch object.(type) {
	case *Application:
		procedure := object.(*Application).procedure
		switch procedure.(type) {
		case *Variable:
			testParentRelationship(t, object, procedure, source)
		case *Procedure:
			testClosedScope(t, procedure, source)
		}
		testParentRelationship(t, object, object.(*Application).arguments, source)
	case *Pair:
		car := object.(*Pair).Car
		cdr := object.(*Pair).Cdr
		if car != nil {
			testParentRelationship(t, object, object.(*Pair).Car, source)
		}
		if cdr != nil {
			testParentRelationship(t, object, object.(*Pair).Cdr, source)
		}
	case *Definition:
		testParentRelationship(t, object, object.(*Definition).variable, source)
		testParentRelationship(t, object, object.(*Definition).value, source)
	case *Procedure:
		arguments := object.(*Procedure).arguments
		body := object.(*Procedure).body
		if arguments != nil {
			testParentRelationship(t, object, object.(*Procedure).arguments, source)
		}
		if body != nil {
			testParentRelationship(t, object, object.(*Procedure).body, source)
		}
	}
}

func TestParser(t *testing.T) {
	for _, test := range interpreterTests {
		i := NewInterpreter(test.source)
		parser := i.Parser
		parser.Peek()
		tree := parser.Parse(i)
		testTreeStructure(t, tree, test.source)
	}
}
