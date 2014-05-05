// Boolean is a type for scheme bool objects, such as #f, #t.

package scheme

import (
	"log"
)

type Boolean struct {
	ObjectBase
	value bool
}

func NewBoolean(value interface{}) *Boolean {
	switch value.(type) {
	case bool:
		return &Boolean{value: value.(bool)}
	case string:
		if value.(string) == "#t" {
			return &Boolean{value: true}
		} else if value.(string) == "#f" {
			return &Boolean{value: false}
		} else {
			log.Fatal("Unexpected value for NewBoolean")
		}
	}
	return nil
}

func (b *Boolean) String() string {
	if b.value == true {
		return "#t"
	} else {
		return "#f"
	}
}
