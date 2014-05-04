// Parser is a type to analyse scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

import (
	"log"
	"strings"
	"text/scanner"
)

type Parser struct {
	*Lexer
}

func NewParser(source string) *Parser {
	parser := &Parser{&Lexer{}}
	parser.Init(strings.NewReader(source))
	return parser
}

func (p *Parser) Parse() Object {
	switch p.Peek() {
	case ')':
		log.Print("Unexpected token: ')'")
		return nil
	case '(':
		p.Next()
		object := p.parseListBody()
		if object == nil {
			log.Fatal("Unexpected flow(parseListBody returns nil)")
			return nil
		}
		list := object.(*Pair)
		if list.Car == nil && list.Cdr == nil {
			return list
		}
		return &Application{
			procedureVariable: list.Car,
			arguments:         list.Cdr,
		}
	case scanner.EOF:
		return nil
	}

	object := p.NextToken()
	return object
}

// This function returns only *Pair.
// But returns Object because if this method returns nil which is not interface type,
// Parse()'s result cannot be judged as nil.
// To avoid such a situation, this method's return value is Object.
func (p *Parser) parseListBody() Object {
	if p.Peek() == ')' {
		p.Next()
		return new(Pair)
	}

	car := p.Parse()
	if car == nil {
		log.Print("Unsupported flow (maybe incomplete source or unexpected expression)")
		return nil
	}
	cdr := p.parseListBody().(*Pair)
	return &Pair{Car: car, Cdr: cdr}
}
