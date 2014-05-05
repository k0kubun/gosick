// Parser is a type to analyse scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

import (
	"log"
)

type Parser struct {
	*Lexer
}

func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

func (p *Parser) Parse() Object {
	return p.parseObject(&TopLevel)
}

func (p *Parser) parseObject(environment *Environment) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		peekToken := p.PeekToken()
		if p.TokenType() == ')' {
			p.NextToken()
			return new(Pair)
		} else if peekToken == "define" {
			return p.parseDefinition(environment)
		} else if peekToken == "quote" {
			p.NextToken()
			object := p.parseQuotedList(environment)
			if !object.IsList() || object.(*Pair).ListLength() != 1 {
				log.Print("Error: syntax-error: malformed quote")
				return nil
			}
			return object.(*Pair).Car
		}

		return p.parseApplication(environment)
	case '\'':
		return p.parseQuotedObject(environment)
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewVariable(token, environment)
	case BooleanToken:
		return NewBoolean(token)
	default:
		return nil
	}
	return nil
}

// This function returns *Pair of first object and list from second.
// Returns value is Object because if a method returns nil which is not
// interface type, the method's result cannot be judged as nil.
func (p *Parser) parseList(environment *Environment) Object {
	car := p.parseObject(environment)
	if car == nil {
		return new(Pair)
	}
	cdr := p.parseList(environment).(*Pair)
	return &Pair{Car: car, Cdr: cdr}
}

func (p *Parser) parseApplication(environment *Environment) Object {
	firstObject := p.parseObject(environment)
	if firstObject == nil {
		log.Print("Unexpected flow: procedure application car is nil")
		return nil
	}

	list := p.parseList(environment)
	if list == nil {
		log.Print("Unexpected flow: procedure application cdr is nil")
		return nil
	}
	return &Application{
		procedureVariable: firstObject,
		arguments:         list,
		environment:       environment,
	}
}

func (p *Parser) parseDefinition(environment *Environment) Object {
	p.NextToken() // skip "define"

	object := p.parseList(environment)
	if !object.IsList() || object.(*Pair).ListLength() != 2 {
		log.Print("Compile Error: syntax-error: (define)")
		return nil
	}

	list := object.(*Pair)
	variable := list.ElementAt(0).(*Variable)
	value := list.ElementAt(1)

	return &Definition{
		environment: environment,
		variable:    variable,
		value:       value,
	}
}

func (p *Parser) parseQuotedObject(environment *Environment) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(environment)
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewSymbol(token)
	case BooleanToken:
		return NewBoolean(token)
	default:
		return nil
	}
}

func (p *Parser) parseQuotedList(environment *Environment) Object {
	car := p.parseQuotedObject(environment)
	if car == nil {
		return new(Pair)
	}
	cdr := p.parseQuotedList(environment).(*Pair)
	return &Pair{Car: car, Cdr: cdr}
}
