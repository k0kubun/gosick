// Number is a scheme number object, which is expressed by number literal.

package scheme

import (
	"strconv"
)

type Number struct {
	ObjectBase
	value int
}

func NewNumber(argument interface{}, options ...Object) *Number {
	var value int
	var err error

	switch argument.(type) {
	case int:
		value = argument.(int)
	case string:
		value, err = strconv.Atoi(argument.(string))
		if err != nil {
			runtimeError("String conversion %s to integer failed.", argument.(string))
		}
	default:
		runtimeError("Unexpected argument type for NewNumber()")
	}

	if len(options) > 0 {
		return &Number{ObjectBase: ObjectBase{parent: options[0]}, value: value}
	} else {
		return &Number{value: value}
	}
}

func (n *Number) Eval() Object {
	return n
}

func (n *Number) String() string {
	return strconv.Itoa(n.value)
}

func (n *Number) isNumber() bool {
	return true
}
