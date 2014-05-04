// Number is a scheme number object, which is expressed by number literal.

package scheme

import (
	"fmt"
	"log"
	"strconv"
)

type Number struct {
	ObjectBase
	number int
}

func NewNumber(numberText string) *Number {
	number, err := strconv.Atoi(numberText)
	if err != nil {
		log.Fatal(fmt.Sprintf("String conversion %s to integer failed.", numberText))
	}
	return &Number{number: number}
}

func (n *Number) String() string {
	return strconv.Itoa(n.number)
}
