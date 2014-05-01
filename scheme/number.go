package scheme

type Number SchemeType

func NewNumber(expression string) *Number {
	return &Number{
		expression: expression,
	}
}

func (n *Number) String() string {
	return n.expression
}
