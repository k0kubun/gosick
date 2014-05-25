// This file is for statements by syntax form, such as set!

package scheme

import (
	"fmt"
)

var (
	builtinSyntaxes = Binding{
		"and":    NewSyntax(andSyntax),
		"begin":  NewSyntax(beginSyntax),
		"cond":   NewSyntax(condSyntax),
		"define": NewSyntax(defineSyntax),
		"do":     NewSyntax(doSyntax),
		"if":     NewSyntax(ifSyntax),
		"lambda": NewSyntax(lambdaSyntax),
		"or":     NewSyntax(orSyntax),
		"quote":  NewSyntax(quoteSyntax),
		"set!":   NewSyntax(setSyntax),
	}
)

type Syntax struct {
	ObjectBase
	function func(*Syntax, Object) Object
}

func NewSyntax(function func(*Syntax, Object) Object) *Syntax {
	return &Syntax{ObjectBase: ObjectBase{parent: nil}, function: function}
}

func (s *Syntax) Invoke(arguments Object) Object {
	return s.function(s, arguments)
}

func (s *Syntax) String() string {
	return fmt.Sprintf("#<syntax %s>", s.Bounder())
}

func (s *Syntax) isSyntax() bool {
	return true
}

func (s *Syntax) malformedError() {
	syntaxError("malformed %s: %s", s.Bounder(), s.Bounder().Parent())
}

func (s *Syntax) assertListEqual(arguments Object, length int) {
	if !arguments.isList() || arguments.(*Pair).ListLength() != length {
		s.malformedError()
	}
}

func (s *Syntax) assertListMinimum(arguments Object, minimum int) {
	if !arguments.isList() || arguments.(*Pair).ListLength() < minimum {
		s.malformedError()
	}
}

func (s *Syntax) assertListRange(arguments Object, lengthRange []int) {
	if !arguments.isList() {
		s.malformedError()
	}

	for _, length := range lengthRange {
		if length == arguments.(*Pair).ListLength() {
			return
		}
	}
	s.malformedError()
}

func andSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 0)

	lastResult := Object(NewBoolean(true))
	for _, object := range arguments.(*Pair).Elements() {
		lastResult = object.Eval()
		if lastResult.isBoolean() && lastResult.(*Boolean).value == false {
			return NewBoolean(false)
		}
	}
	return lastResult
}

func beginSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 0)

	lastResult := undef
	for _, object := range arguments.(*Pair).Elements() {
		lastResult = object.Eval()
	}
	return lastResult
}

func condSyntax(s *Syntax, arguments Object) Object {
	if arguments.isApplication() {
		arguments = NewList(arguments.Parent(), arguments)
	}
	s.assertListMinimum(arguments, 0)
	if arguments.(*Pair).ListLength() == 0 {
		syntaxError("at least one clause is required for cond")
	}
	elements := arguments.(*Pair).Elements()

	// First: syntax check
	elseExists := false
	for _, element := range elements {
		if elseExists {
			syntaxError("'else' clause followed by more clauses")
		} else if element.isApplication() && element.(*Application).procedure.isVariable() &&
			element.(*Application).procedure.(*Variable).identifier == "else" {
			elseExists = true
		}

		if element.isNull() || !element.isApplication() {
			syntaxError("bad clause in cond")
		}
	}

	// Second: eval cases
	for _, element := range elements {
		lastResult := undef
		application := element.(*Application)

		isElse := application.procedure.isVariable() && application.procedure.(*Variable).identifier == "else"
		if !isElse {
			lastResult = application.procedure.Eval()
		}

		// first element is 'else' or not '#f'
		if isElse || !lastResult.isBoolean() || lastResult.(*Boolean).value == true {
			for _, object := range application.arguments.(*Pair).Elements() {
				lastResult = object.Eval()
			}
			return lastResult
		}
	}
	return undef
}

func defineSyntax(s *Syntax, arguments Object) Object {
	s.assertListEqual(arguments, 2)
	elements := arguments.(*Pair).Elements()

	if !elements[0].isVariable() {
		s.malformedError()
	}
	variable := elements[0].(*Variable)
	s.Bounder().define(variable.identifier, elements[1].Eval())

	return NewSymbol(variable.identifier)
}

// FIXME: This implementation is extremely dirty...
func doSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 2)
	elements := arguments.(*Pair).Elements()

	// Insert closure betweetn application and its parent
	application := arguments.Parent()
	closure := NewClosure(application.Parent())
	application.setParent(closure)

	// Parse iterator list and define first variable
	iteratorList := elements[0]
	if iteratorList.isApplication() {
		iteratorList = iteratorList.(*Application).toList()
	}
	s.assertListMinimum(iteratorList, 0)
	iteratorBodies := iteratorList.(*Pair).Elements()
	for _, iteratorBody := range iteratorBodies {
		if iteratorBody.isApplication() {
			iteratorBody = iteratorBody.(*Application).toList()
		}
		s.assertListMinimum(iteratorBody, 2)
		if iteratorBody.(*Pair).ListLength() > 3 {
			compileError("bad update expr in %s: %s", s.Bounder(), s.Bounder().Parent())
		}

		variable := iteratorBody.(*Pair).ElementAt(0)
		value := iteratorBody.(*Pair).ElementAt(1)
		if variable.isVariable() {
			closure.localBinding[variable.(*Variable).identifier] = value.Eval()
		}
	}

	// eval test ->
	//   true: eval testBody and returns its result
	//  false: eval continueBody, eval iterator's update
	testBody := elements[1]
	if testBody.isApplication() {
		testBody = testBody.(*Application).toList()
	}
	s.assertListMinimum(testBody, 1)
	testElements := testBody.(*Pair).Elements()
	continueElements := elements[2:]
	for {
		testResult := testElements[0].Eval()
		if !testResult.isBoolean() || testResult.(*Boolean).value == true {
			for _, element := range testElements[1:] {
				testResult = element.Eval()
			}
			return testResult
		} else {
			// eval continueBody
			for _, element := range continueElements {
				element.Eval()
			}

			// update iterators
			for _, iteratorBody := range iteratorBodies {
				if iteratorBody.isApplication() {
					iteratorBody = iteratorBody.(*Application).toList()
				}
				iteratorElements := iteratorBody.(*Pair).Elements()
				if len(iteratorElements) == 3 {
					variable := iteratorElements[0]
					if variable.isVariable() {
						closure.localBinding[variable.(*Variable).identifier] = iteratorElements[2].Eval()
					}
				}
			}
		}
	}
	return undef
}

func ifSyntax(s *Syntax, arguments Object) Object {
	s.assertListRange(arguments, []int{2, 3})
	elements := arguments.(*Pair).Elements()

	result := elements[0].Eval()
	if result.isBoolean() && !result.(*Boolean).value {
		if len(elements) == 3 {
			return elements[2].Eval()
		} else {
			return undef
		}
	} else {
		return elements[1].Eval()
	}
}

func lambdaSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 1)
	elements := arguments.(*Pair).Elements()

	// Insert closure between application and its parent
	application := arguments.Parent()
	closure := NewClosure(application.Parent())
	application.setParent(closure)

	// Parse argument list
	argumentList := elements[0]
	if argumentList.isApplication() {
		argumentList = argumentList.(*Application).toList()
	}
	s.assertListMinimum(argumentList, 0)
	variables := argumentList.(*Pair).Elements()

	// generate function
	closure.function = func(givenArguments Object) Object {
		// assert given arguments
		s.assertListMinimum(givenArguments, 0)
		givenElements := givenArguments.(*Pair).Elements()
		if len(variables) != len(givenElements) {
			compileError("wrong number of arguments: requires %d, but got %d", len(variables), len(givenElements))
		}

		// define arguments to local scope
		for index, variable := range variables {
			object := givenElements[index].Eval()
			if variable.isVariable() {
				closure.localBinding[variable.(*Variable).identifier] = object
			}
		}

		// returns last eval result
		lastResult := undef
		for _, element := range elements[1:] {
			lastResult = element.Eval()
		}
		return lastResult
	}
	return closure
}

func orSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 0)

	lastResult := Object(NewBoolean(false))
	for _, object := range arguments.(*Pair).Elements() {
		lastResult = object.Eval()
		if !lastResult.isBoolean() || lastResult.(*Boolean).value != false {
			return lastResult
		}
	}
	return lastResult
}

func quoteSyntax(s *Syntax, arguments Object) Object {
	s.assertListEqual(arguments, 1)
	object := arguments.(*Pair).ElementAt(0)

	p := NewParser(object.String())
	p.Peek()
	return p.parseQuotedObject(s.Bounder())
}

func setSyntax(s *Syntax, arguments Object) Object {
	s.assertListEqual(arguments, 2)
	elements := arguments.(*Pair).Elements()

	variable := elements[0]
	if !variable.isVariable() {
		s.malformedError()
	}
	value := elements[1].Eval()
	s.Bounder().set(variable.(*Variable).identifier, value)
	return value
}
