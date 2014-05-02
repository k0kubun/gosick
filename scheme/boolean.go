package scheme

import (
	"log"
)

type Boolean struct {
	SchemeType
}

func NewBoolean(expression interface{}) *Boolean {
	switch expression.(type) {
	case string:
		return &Boolean{SchemeType{expression: expression.(string)}}
	case bool:
		if expression.(bool) {
			return NewBoolean("#t")
		} else {
			return NewBoolean("#f")
		}
	default:
		log.Fatal("Unexpected expression for NewBoolean()")
		return nil
	}
}

func (b *Boolean) String() string {
	return b.expression
}

// func (b *Boolean) IsNumber() bool {
// 	return false
// }
