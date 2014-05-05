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
		if p.TokenType() == ')' {
			p.NextToken()
			return new(Pair)
		}

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
	case ')':
		return nil
	case EOF:
		return nil
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewVariable(token)
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
