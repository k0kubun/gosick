// Symbol is a type to express scheme symbol object, which
// is expressed like 'symbol or (quote symbol).

package scheme

var (
	symbols = make(map[string]*Symbol)
	undef   = Object(&Symbol{identifier: "#<undef>"}) // FIXME: this should not be symbol
)

type Symbol struct {
	ObjectBase
	identifier string
}

func NewSymbol(identifier string) *Symbol {
	if symbols[identifier] == nil {
		symbols[identifier] = &Symbol{ObjectBase: ObjectBase{parent: nil}, identifier: identifier}
	}
	return symbols[identifier]
}

func (s *Symbol) Eval() Object {
	return s
}

func (s *Symbol) String() string {
	return s.identifier
}

func (s *Symbol) isSymbol() bool {
	return true
}
