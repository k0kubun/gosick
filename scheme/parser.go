// Parser is a type to analyse scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

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

func (p *Parser) parseBlock(parent Object) Object {
	switch p.PeekToken() {
	case ")":
		p.NextToken()
		return Null
	case "lambda":
		p.NextToken()
		return p.parseProcedure(parent)
	case "let", "let*", "letrec":
		p.NextToken()
		return p.parseLet(parent)
	case "do":
		p.NextToken()
		return p.parseDo(parent)
	}

	return p.parseApplication(parent)
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
	application := NewApplication(parent)
	application.procedure = p.parseObject(application)
	application.arguments = p.parseList(application)

	return application
}

func (p *Parser) parseProcedure(parent Object) Object {
	if p.TokenType() == '(' {
		p.NextToken()
	} else {
		syntaxError("malformed lambda")
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
		syntaxError("malformed let")
	}

	application := NewApplication(parent)
	procedure := new(Procedure)

	procedureArguments := NewPair(procedure)
	applicationArguments := NewPair(application)

	argumentSets := p.parseList(application)
	for _, set := range argumentSets.(*Pair).Elements() {
		if !set.isApplication() || set.(*Application).arguments.(*Pair).ListLength() != 1 {
			syntaxError("malformed let")
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

func (p *Parser) parseDo(parent Object) Object {
	do := NewDo(parent)

	// parse iterators
	if p.NextToken() != "(" {
		syntaxError("malformed do")
	}
	do.iterators = p.parseIterators(do)

	// parse test and a body for the case test is true
	if p.NextToken() != "(" {
		syntaxError("malformed do")
	}
	do.testBody = p.parseList(do)
	if do.testBody.(*Pair).ListLength() == 0 {
		syntaxError("malformed do")
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
			syntaxError("malformed do")
		}

		// get element list and assert their number
		elementList := p.parseList(parent)
		if !elementList.isList() || elementList.(*Pair).ListLength() < 2 {
			syntaxError("malformed do")
		} else if elementList.(*Pair).ListLength() > 3 {
			syntaxError("bad update expr in do")
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
