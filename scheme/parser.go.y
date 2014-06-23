%{
// Parser is a type to analyse scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme
%}

%union{
	object Object
}

%type<object> program
%type<object> expr
%type<object> const

%token<object> NUMBER

%%

program:
	expr
		{
			$$ = $1
			if l, ok := yylex.(*Lexer); ok {
				l.result = NewNumber(1)
			}
		}

expr:
	const
		{
			$$ = $1
		}

const:
	NUMBER
		{
			$$ = NewNumber(1)
		}

%%

type Parser struct {
	*Lexer
}

func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

func (p *Parser) Parse(parent Object) Object {
	p.ensureAvailability()
	if yyParse(p.Lexer) != 0 {
		panic("parse error")
	}
	p.result.setParent(parent)
	return p.result
}

func (p *Parser) parseObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseApplication(parent)
	case '\'':
		return p.parseSingleQuote(parent)
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

// This is for parsing syntax sugar '*** => (quote ***)
func (p *Parser) parseSingleQuote(parent Object) Object {
	if len(p.PeekToken()) == 0 {
		runtimeError("unterminated quote")
	}
	application := NewApplication(parent)
	application.procedure = NewVariable("quote", application)
	application.arguments = NewList(application, p.parseObject(application))
	return application
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
	if p.PeekToken() == ")" {
		p.NextToken()
		return Null
	}
	application := NewApplication(parent)
	application.procedure = p.parseObject(application)
	application.arguments = p.parseList(application)

	return application
}

func (p *Parser) parseQuotedObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(parent)
	case '\'':
		return p.parseSingleQuote(parent)
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
