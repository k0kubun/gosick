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
		}

		return p.parseApplication(parent)
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
	application := &Application{
		ObjectBase: ObjectBase{parent: parent},
	}
	application.procedure = p.parseObject(application)
	application.arguments = p.parseList(application)

	return application
}

func (p *Parser) parseProcedure(parent Object) Object {
	if p.TokenType() == '(' {
		p.NextToken()
		procedure := new(Procedure)
		procedure.arguments = p.parseList(procedure)
		procedure.body = p.parseList(procedure)
		procedure.generateFunction(parent)
		return procedure
	} else {
		runtimeError("Not implemented yet.")
		return nil
	}
}

func (p *Parser) parseDefinition(parent Object) Object {
	definition := &Definition{ObjectBase: ObjectBase{parent: parent}}

	object := p.parseList(definition)
	if !object.isList() || object.(*Pair).ListLength() != 2 {
		runtimeError("Compile Error: syntax-error: (define)")
	}
	list := object.(*Pair)

	definition.variable = list.ElementAt(0).(*Variable)
	definition.variable.setParent(definition) // before this, variable's parent is list.

	definition.value = list.ElementAt(1)
	definition.value.setParent(definition) // before this, value's parent is list.

	return definition
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
