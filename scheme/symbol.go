// Symbol is a type to express scheme symbol object, which
// is expressed like 'symbol or (quote symbol).

package scheme

type Symbol struct {
	ObjectBase
	identifier string
}

func NewSymbol(identifier string) *Symbol {
	return &Symbol{identifier: identifier}
}

func (s *Symbol) Eval() Object {
	return s
}

func (s *Symbol) String() string {
	return s.identifier
}
