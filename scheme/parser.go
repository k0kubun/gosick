// Parser is a type to analyse scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

import (
	"fmt"
)

type Parser struct {
	*Lexer
}

func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

func (p *Parser) Parse(parent Object) Object {
	p.ensureAvailability()
	return p.parseObject(parent)
}

func (p *Parser) parseObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseBlock(parent)
	case '\'':
		return p.parseQuotedObject(parent)
	case IntToken:
		return NewNumber(token, parent)
	case IdentifierToken:
		return NewVariable(token, parent)
	case BooleanToken:
		return NewBoolean(token, parent)
	case StringToken:
		return NewString(token[1:len(token)-1], parent)
	default:
		return nil
	}
	return nil
}

func (p *Parser) parseBlock(parent Object) Object {
	switch p.PeekToken() {
	case ")":
		p.NextToken()
		return Null
	case "define":
		p.NextToken()
		return p.parseDefinition(parent)
	case "quote":
		p.NextToken()
		object := p.parseQuotedList(parent)
		if !object.isList() || object.(*Pair).ListLength() != 1 {
			compileError("syntax-error: malformed quote")
		}
		quotedObject := object.(*Pair).Car
		quotedObject.setParent(parent)
		return quotedObject
	case "lambda":
		p.NextToken()
		return p.parseProcedure(parent)
	case "let", "let*", "letrec":
		p.NextToken()
		return p.parseLet(parent)
	case "set!":
		p.NextToken()
		set := NewSet(parent)
		object := p.parseList(set)
		if !object.isList() || object.(*Pair).ListLength() != 2 {
			compileError("syntax-error: malformed set!")
		}

		set.variable = object.(*Pair).ElementAt(0)
		set.value = object.(*Pair).ElementAt(1)
		return set
	case "if":
		p.NextToken()
		return p.parseIf(parent)
	case "cond":
		p.NextToken()
		return p.parseCond(parent)
	case "and":
		p.NextToken()
		return p.parseAnd(parent)
	case "or":
		p.NextToken()
		return p.parseOr(parent)
	case "begin":
		p.NextToken()
		return p.parseBegin(parent)
	case "do":
		p.NextToken()
		return p.parseDo(parent)
	}

	return p.parseApplication(parent)
}

// This function returns *Pair of first object and list from second.
// Scanner position ends with the next of close parentheses.
func (p *Parser) parseList(parent Object) Object {
	pair := NewPair(parent)
	pair.Car = p.parseObject(pair)
	if pair.Car == nil {
		return pair
	}
	pair.Cdr = p.parseList(pair).(*Pair)
	return pair
}

func (p *Parser) parseApplication(parent Object) Object {
	application := NewApplication(parent)
	application.procedure = p.parseObject(application)
	application.arguments = p.parseList(application)

	return application
}

func (p *Parser) parseProcedure(parent Object) Object {
	if p.TokenType() == '(' {
		p.NextToken()
	} else {
		compileError("syntax-error: malformed lambda")
	}

	procedure := new(Procedure)
	procedure.arguments = p.parseList(procedure)
	procedure.body = p.parseList(procedure)
	procedure.generateFunction(parent)
	return procedure
}

func (p *Parser) parseLet(parent Object) Object {
	if p.TokenType() == '(' {
		p.NextToken()
	} else {
		compileError("syntax-error: malformed let")
	}

	application := NewApplication(parent)
	procedure := new(Procedure)

	procedureArguments := NewPair(procedure)
	applicationArguments := NewPair(application)

	argumentSets := p.parseList(application)
	for _, set := range argumentSets.(*Pair).Elements() {
		if !set.isApplication() || set.(*Application).arguments.(*Pair).ListLength() != 1 {
			compileError("syntax-error: malformed let")
		}

		procedureArguments.Append(set.(*Application).procedure)
		applicationArguments.Append(set.(*Application).arguments.(*Pair).ElementAt(0))
	}

	procedure.arguments = procedureArguments
	procedure.body = p.parseList(application)
	procedure.generateFunction(parent)

	application.arguments = applicationArguments
	application.procedure = procedure
	return application
}

func (p *Parser) parseDefinition(parent Object) Object {
	definition := &Definition{ObjectBase: ObjectBase{parent: parent}}

	object := p.parseList(definition)
	if !object.isList() || object.(*Pair).ListLength() != 2 {
		runtimeError("Compile Error: syntax-error: (define)")
	}
	list := object.(*Pair)

	definition.variable = list.ElementAt(0).(*Variable)
	definition.variable.setParent(definition)

	definition.value = list.ElementAt(1)
	definition.value.setParent(definition)

	return definition
}

func (p *Parser) parseIf(parent Object) Object {
	ifStatement := NewIf(parent)
	ifStatement.condition = p.parseObject(ifStatement)
	ifStatement.trueBody = p.parseObject(ifStatement)
	if p.PeekToken() == ")" {
		p.NextToken()
		ifStatement.falseBody = undef
	} else {
		ifStatement.falseBody = p.parseObject(ifStatement)
		if p.NextToken() != ")" {
			compileError("syntax-error: malformed if")
		}
	}

	return ifStatement
}

func (p *Parser) parseCond(parent Object) Object {
	cond := NewCond(parent)

	caseExists := false
	elseExists := false
	for {
		// judge case continues or not
		firstToken := p.NextToken()
		if firstToken == ")" {
			break
		} else if firstToken != "(" {
			compileError("syntax-error: bad clause in cond")
		}

		// parse list body
		if elseExists {
			compileError("syntax-error: 'else' clause followed by more clauses")
		}
		switch p.PeekToken() {
		case "else":
			p.NextToken()
			cond.elseBody = p.parseList(cond)
			elseExists = true
		case ")":
			compileError("syntax-error: bad clause in cond")
		default:
			caseBody := p.parseList(cond)
			if !caseBody.isList() || caseBody.(*Pair).ListLength() == 0 {
				compileError("syntax-error: bad clause in cond3")
			}
			cond.cases = append(cond.cases, caseBody)
		}
		caseExists = true
	}

	if !caseExists {
		compileError("syntax-error: at least one clause is required for cond")
	}
	return cond
}

func (p *Parser) parseAnd(parent Object) Object {
	andStatement := NewAnd(parent)
	andStatement.body = p.parseList(andStatement)
	return andStatement
}

func (p *Parser) parseOr(parent Object) Object {
	orStatement := NewOr(parent)
	orStatement.body = p.parseList(orStatement)
	return orStatement
}

func (p *Parser) parseBegin(parent Object) Object {
	begin := NewBegin(parent)
	begin.body = p.parseList(begin)
	return begin
}

func (p *Parser) parseDo(parent Object) Object {
	do := NewDo(parent)

	// parse iterators
	if p.NextToken() != "(" {
		compileError("syntax-error: malformed do")
	}
	do.iterators = p.parseIterators(do)

	// parse test and a body for the case test is true
	if p.NextToken() != "(" {
		compileError("syntax-error: malformed do")
	}
	do.testBody = p.parseList(do)
	if do.testBody.(*Pair).ListLength() == 0 {
		compileError("syntax-error: malformed do")
	}

	// parse a body for the case test is false
	do.continueBody = p.parseList(do)
	return do
}

func (p *Parser) parseIterators(parent Object) []*Iterator {
	iterators := []*Iterator{}
	for {
		// check first is '('
		firstToken := p.NextToken()
		if firstToken == ")" {
			break
		} else if firstToken != "(" {
			compileError("syntax-error: malformed do")
		}

		// get element list and assert their number
		elementList := p.parseList(parent)
		if !elementList.isList() || elementList.(*Pair).ListLength() < 2 {
			compileError("syntax-error: malformed do")
		} else if elementList.(*Pair).ListLength() > 3 {
			compileError("syntax-error: bad update expr in do")
		}

		iterator := NewIterator(parent)

		// get variable
		iterator.variable = elementList.(*Pair).ElementAt(0)
		iterator.variable.setParent(iterator)

		// get value
		iterator.value = elementList.(*Pair).ElementAt(1)
		iterator.value.setParent(iterator)

		// get update
		if elementList.(*Pair).ListLength() == 3 {
			iterator.update = elementList.(*Pair).ElementAt(2)
			iterator.update.setParent(iterator)
		}

		iterators = append(iterators, iterator)
	}

	return iterators
}

func (p *Parser) parseQuotedObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(parent)
	case IntToken:
		return NewNumber(token, parent)
	case IdentifierToken:
		return NewSymbol(token)
	case BooleanToken:
		return NewBoolean(token, parent)
	case ')':
		return nil
	default:
		runtimeError("unterminated quote")
		return nil
	}
}

func (p *Parser) parseQuotedList(parent Object) Object {
	pair := NewPair(parent)
	pair.Car = p.parseQuotedObject(pair)
	if pair.Car == nil {
		return pair
	}
	pair.Cdr = p.parseQuotedList(pair).(*Pair)
	return pair
}

func (p *Parser) ensureAvailability() {
	// Error message will be printed by interpreter
	recover()
}

func compileError(format string, a ...interface{}) {
	runtimeError("Compile Error: "+format, a...)
}

func runtimeError(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func dumpObject(object interface{}) {
	fmt.Printf("%#v\n", object)
}
