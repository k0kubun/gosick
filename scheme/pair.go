// Pair is a type which is generated by cons procedure.
// Pair has two pointers, which are named car and cdr.
//
// List is expressed by linked list of Pair.
// And procedure application has list which consists of Pair
// as its arguments.

package scheme

import (
	"fmt"
	"strings"
)

type Pair struct {
	ObjectBase
	Car Object
	Cdr Object
}

func NewNull(parent Object) *Pair {
	return &Pair{ObjectBase: ObjectBase{parent: parent}, Car: nil, Cdr: nil}
}

func (p *Pair) Eval() Object {
	return p
}

func (p *Pair) String() string {
	if p.isNull() {
		return "()"
	} else if p.isList() {
		length := p.ListLength()
		tokens := []string{}
		for i := 0; i < length; i++ {
			tokens = append(tokens, p.ElementAt(i).Eval().String())
		}
		return fmt.Sprintf("(%s)", strings.Join(tokens, " "))
	} else {
		return fmt.Sprintf("(%s . %s)", p.Car, p.Cdr)
	}
}

func (p *Pair) isNull() bool {
	return p.Car == nil && p.Cdr == nil
}

func (p *Pair) isPair() bool {
	return !p.isNull()
}

func (p *Pair) isList() bool {
	pair := p

	for {
		if pair.isNull() {
			return true
		}
		switch pair.Cdr.(type) {
		case *Pair:
			pair = pair.Cdr.(*Pair)
		default:
			return false
		}
	}
	return false
}

func (p *Pair) Elements() []Object {
	elements := []Object{}

	pair := p
	for {
		if pair.Car == nil {
			break
		} else {
			elements = append(elements, pair.Car)
		}
		pair = pair.Cdr.(*Pair)
	}
	return elements
}

func (p *Pair) ElementAt(index int) Object {
	return p.Elements()[index]
}

func (p *Pair) ListLength() int {
	if p.isNull() {
		return 0
	} else {
		return p.Cdr.(*Pair).ListLength() + 1
	}
}

func (p *Pair) Append(object Object) *Pair {
	assertListMinimum(p, 0)
	assertListMinimum(object, 0)

	listTail := p
	for {
		if listTail.isNull() {
			break
		} else {
			listTail = listTail.Cdr.(*Pair)
		}
	}

	listTail.Car = object.(*Pair).Car
	listTail.Cdr = object.(*Pair).Cdr
	return p
}
