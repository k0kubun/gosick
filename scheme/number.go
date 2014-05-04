// Number is a scheme number object, which is expressed by number literal.

package scheme

type Number struct {
	number int
}

func NewNumber(number int) *Number {
	return &Number{number: number}
}
