package scheme

import (
	"strconv"
)

type Number SchemeType

func NewNumber(expression interface{}) *Number {
	switch expression.(type) {
	case string:
		return &Number{expression: expression.(string)}
	case int:
		value := strconv.Itoa(expression.(int))
		return &Number{expression: value}
	default:
		panic("Caught unexpected flow")
	}
}

func (n *Number) String() string {
	return n.expression
}
