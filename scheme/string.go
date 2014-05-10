// String is a type for scheme string object, which is
// expressed like "string".

package scheme

import (
	"fmt"
)

type String struct {
	ObjectBase
	text string
}

func NewString(object interface{}, options ...Object) *String {
	text := ""
	switch object.(type) {
	case string:
		text = object.(string)
	case int:
		text = fmt.Sprintf("%d", object.(int))
	default:
		runtimeError("Unexpected conversion")
	}
	if len(options) > 0 {
		return &String{ObjectBase: ObjectBase{parent: options[0]}, text: text}
	} else {
		return &String{text: text}
	}
}

func (s *String) Eval() Object {
	return s
}

func (s *String) String() string {
	return fmt.Sprintf("\"%s\"", s.text)
}

func (s *String) isString() bool {
	return true
}
