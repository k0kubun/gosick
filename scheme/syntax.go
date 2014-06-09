// This file is for statements by syntax form, such as set!

package scheme

import (
	"fmt"
)

var (
	builtinSyntaxes = Binding{
		"actor":        NewSyntax(actorSyntax),
		"and":          NewSyntax(andSyntax),
		"begin":        NewSyntax(beginSyntax),
		"cond":         NewSyntax(condSyntax),
		"define":       NewSyntax(defineSyntax),
		"define-macro": NewSyntax(defineMacroSyntax),
		"do":           NewSyntax(doSyntax),
		"if":           NewSyntax(ifSyntax),
		"lambda":       NewSyntax(lambdaSyntax),
		"let":          NewSyntax(letSyntax),
		"let*":         NewSyntax(letStarSyntax),
		"letrec":       NewSyntax(letrecSyntax),
		"or":           NewSyntax(orSyntax),
		"quote":        NewSyntax(quoteSyntax),
		"set!":         NewSyntax(setSyntax),
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

// Returns elements in list object with type assertion (syntax form specific error message)
// Assertion is minimum
func (s *Syntax) elementsMinimum(list Object, minimum int) []Object {
	if list.isApplication() {
		list = list.(*Application).toList()
	}
	s.assertListMinimum(list, minimum)
	return list.(*Pair).Elements()
}

// Returns elements in list object with type assertion (syntax form specific error message)
// Assertion is equal
func (s *Syntax) elementsExact(list Object, value int) []Object {
	if list.isApplication() {
		list = list.(*Application).toList()
	}
	s.assertListEqual(list, value)
	return list.(*Pair).Elements()
}

// Eval all given objects and returns last object's eval result.
// When 'objects' is empty, returns #<undef>.
func evalAll(objects []Object) Object {
	lastResult := undef
	for _, object := range objects {
		lastResult = object.Eval()
	}
	return lastResult
}

func actorSyntax(s *Syntax, arguments Object) Object {
	elements := s.elementsMinimum(arguments, 0)

	// Insert over the application to override scope
	application := arguments.Parent()
	actor := NewActor(application.Parent())
	application.setParent(actor)

	for _, element := range elements {
		caseElements := s.elementsMinimum(element, 1)
		caseArguments := s.elementsMinimum(caseElements[0], 1)
		assertObjectType(caseArguments[0], "string")

		actor.functions[caseArguments[0].(*String).text] = func(objects []Object) {
			if len(caseArguments[1:]) != len(objects) {
				runtimeError("invalid message argument length: requires %d, but got %d", len(caseArguments[1:]), len(objects))
			}

			for index, variable := range caseArguments[1:] {
				actor.tryDefine(variable, objects[index].Eval())
			}
			evalAll(caseElements[1:])
		}
	}

	return actor
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
	elements := s.elementsMinimum(arguments, 0)
	return evalAll(elements)
}

func condSyntax(s *Syntax, arguments Object) Object {
	elements := s.elementsMinimum(arguments, 0)
	if len(elements) == 0 {
		syntaxError("at least one clause is required for cond")
	}

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
	elements := s.elementsExact(arguments, 2)

	if elements[0].isVariable() {
		variable := elements[0].(*Variable)
		s.Bounder().define(variable.identifier, elements[1].Eval())

		return NewSymbol(variable.identifier)
	} else if elements[0].isApplication() {
		closure := WrapClosure(arguments)

		defineElements := s.elementsMinimum(elements[0], 1)
		funcName := defineElements[0]
		closure.DefineFunction(s, defineElements[1:], elements[1:])

		if funcName.isVariable() {
			s.Bounder().define(funcName.(*Variable).identifier, closure)
			return funcName
		}
	}
	return syntaxError("%s", s.Bounder().Parent())
}

func defineMacroSyntax(s *Syntax, arguments Object) Object {
	elements := s.elementsMinimum(arguments, 2)

	macroElements := s.elementsMinimum(elements[0], 1)
	assertObjectType(macroElements[0], "variable")

	macro := NewMacro()
	s.Bounder().define(macroElements[0].(*Variable).identifier, macro)
	return undef
}

func doSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())

	// Parse iterator list and define first variable
	elements := s.elementsMinimum(arguments, 2)
	iteratorBodies := s.elementsMinimum(elements[0], 0)
	for _, iteratorBody := range iteratorBodies {
		iteratorElements := s.elementsMinimum(iteratorBody, 2)
		if len(iteratorElements) > 3 {
			compileError("bad update expr in %s: %s", s.Bounder(), s.Bounder().Parent())
		}
		closure.tryDefine(iteratorElements[0], iteratorElements[1].Eval())
	}

	// eval test ->
	//   true: eval testBody and returns its result
	//  false: eval continueBody, eval iterator's update
	testElements := s.elementsMinimum(elements[1], 1)
	for {
		testResult := testElements[0].Eval()
		if !testResult.isBoolean() || testResult.(*Boolean).value == true {
			for _, element := range testElements[1:] {
				testResult = element.Eval()
			}
			return testResult
		} else {
			// eval continueBody
			evalAll(elements[2:])

			// update iterators
			for _, iteratorBody := range iteratorBodies {
				iteratorElements := s.elementsMinimum(iteratorBody, 2)
				if len(iteratorElements) == 3 {
					closure.tryDefine(iteratorElements[0], iteratorElements[2].Eval())
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
	closure := WrapClosure(arguments.Parent())

	elements := s.elementsMinimum(arguments, 1)
	variables := s.elementsMinimum(elements[0], 0)
	closure.DefineFunction(s, variables, elements[1:])
	return closure
}

func letSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())
	elements := s.elementsMinimum(arguments, 1)

	// define arguments to local scope
	variables := []Object{}
	results := []Object{}
	for _, argumentElement := range s.elementsMinimum(elements[0], 0) {
		variableElements := s.elementsExact(argumentElement, 2)

		variableElements[1].setParent(closure.Parent())
		variables = append(variables, variableElements[0])
		results = append(results, variableElements[1].Eval())
	}

	for index, variable := range variables {
		closure.tryDefine(variable, results[index])
	}

	// eval body
	return evalAll(elements[1:])
}

func letStarSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())
	elements := s.elementsMinimum(arguments, 1)

	// define arguments to local scope
	for _, argumentElement := range s.elementsMinimum(elements[0], 0) {
		variableElements := s.elementsExact(argumentElement, 2)
		variable, result := variableElements[0], variableElements[1].Eval()

		closure.tryDefine(variable, result)
		result.setParent(closure.Parent())
	}

	// eval body
	return evalAll(elements[1:])
}

func letrecSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())
	elements := s.elementsMinimum(arguments, 1)

	// define arguments to local scope
	variables := []Object{}
	results := []Object{}
	for _, argumentElement := range s.elementsMinimum(elements[0], 0) {
		variableElements := s.elementsExact(argumentElement, 2)
		variables = append(variables, variableElements[0])

		results = append(results, variableElements[1].Eval())
	}

	for index, variable := range variables {
		closure.tryDefine(variable, results[index])
	}

	// eval body
	return evalAll(elements[1:])
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
	elements := s.elementsExact(arguments, 2)

	variable := elements[0]
	if !variable.isVariable() {
		s.malformedError()
	}
	value := elements[1].Eval()
	s.Bounder().set(variable.(*Variable).identifier, value)
	return value
}
