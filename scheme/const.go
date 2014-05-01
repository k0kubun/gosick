package scheme

type Const SchemeType

func NewConst(expression string) *Const {
	return &Const{
		expression: expression,
	}
}

func (c *Const) String() string {
	return c.expression
}
