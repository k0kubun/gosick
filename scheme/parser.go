package scheme

type Parser struct {
	Lexer
}

func (p Parser) Parse() *Object {
	token := p.ReadToken()
	return nil
}
