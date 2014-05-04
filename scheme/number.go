// Number is a scheme number object, which is expressed by number literal.

package scheme

import (
	"fmt"
	"log"
	"strconv"
)

type Number struct {
	ObjectBase
	value int
}

func NewNumber(argument interface{}) *Number {
	var value int
	var err error

	switch argument.(type) {
	case int:
		value = argument.(int)
	case string:
		value, err = strconv.Atoi(argument.(string))
		if err != nil {
			log.Fatal(fmt.Sprintf("String conversion %s to integer failed.", argument.(string)))
		}
	default:
		log.Fatal("Unexpected argument type for NewNumber()")
	}

	return &Number{value: value}
}

func (n *Number) String() string {
	return strconv.Itoa(n.value)
}
