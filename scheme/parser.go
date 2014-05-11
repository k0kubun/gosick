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
	peekToken := p.PeekToken()
	if p.TokenType() == ')' {
		p.NextToken()
		return NewNull(parent)
	} else if peekToken == "define" {
		p.NextToken()
		return p.parseDefinition(parent)
	} else if peekToken == "quote" {
		p.NextToken()
		object := p.parseQuotedList(parent)
		if !object.isList() || object.(*Pair).ListLength() != 1 {
			compileError("syntax-error: malformed quote")
		}
		quotedObject := object.(*Pair).Car
		quotedObject.setParent(parent)
		return quotedObject
	} else if peekToken == "lambda" {
		p.NextToken()
		return p.parseProcedure(parent)
	} else if peekToken == "let" || peekToken == "let*" || peekToken == "letrec" {
		p.NextToken()
		return p.parseLet(parent)
	} else if peekToken == "set!" {
		p.NextToken()
		set := NewSet(parent)
		object := p.parseList(set)
		if !object.isList() || object.(*Pair).ListLength() != 2 {
			compileError("syntax-error: malformed set!")
		}

		set.variable = object.(*Pair).ElementAt(0)
		set.value = object.(*Pair).ElementAt(1)
		return set
	} else if peekToken == "if" {
		p.NextToken()
		return p.parseIf(parent)
	}

	return p.parseApplication(parent)
}

// This function returns *Pair of first object and list from second.
// Returns value is Object because if a method returns nil which is not
// interface type, the method's result cannot be judged as nil.
func (p *Parser) parseList(parent Object) Object {
	pair := NewNull(parent)
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

	procedureArguments := NewNull(procedure)
	applicationArguments := NewNull(application)

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
		ifStatement.falseBody = NewSymbol("#<undef>")
	} else {
		ifStatement.falseBody = p.parseObject(ifStatement)
		if p.NextToken() != ")" {
			compileError("syntax-error: malformed if")
		}
	}

	return ifStatement
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
		return NewSymbol(token, parent)
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
	pair := NewNull(parent)
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
