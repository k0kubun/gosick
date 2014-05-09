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
			return object.(*Pair).Car
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
		return NewBoolean(token)
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
	pair.Cdr = p.parseList(parent).(*Pair)
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
		procedure.generateFunction(parent, p.parseList(procedure), p.parseList(procedure))
		return procedure
	} else {
		runtimeError("Not implemented yet.")
		return nil
	}
}

func (p *Parser) parseDefinition(parent Object) Object {
	object := p.parseList(parent)
	if !object.isList() || object.(*Pair).ListLength() != 2 {
		runtimeError("Compile Error: syntax-error: (define)")
	}

	list := object.(*Pair)
	variable := list.ElementAt(0).(*Variable)
	value := list.ElementAt(1)

	return &Definition{
		ObjectBase: ObjectBase{parent: parent},
		variable:   variable,
		value:      value,
	}
}

func (p *Parser) parseQuotedObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(parent)
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewSymbol(token, parent)
	case BooleanToken:
		return NewBoolean(token)
	default:
		return nil
	}
}

func (p *Parser) parseQuotedList(parent Object) Object {
	car := p.parseQuotedObject(parent)
	if car == nil {
		return NewNull(parent)
	}
	cdr := p.parseQuotedList(parent).(*Pair)
	return &Pair{ObjectBase: ObjectBase{parent: parent}, Car: car, Cdr: cdr}
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
