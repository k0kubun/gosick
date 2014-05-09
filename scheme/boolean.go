// Boolean is a type for scheme bool objects, such as #f, #t.

package scheme

type Boolean struct {
	ObjectBase
	value bool
}

func NewBoolean(value interface{}, options ...Object) (boolean *Boolean) {
	switch value.(type) {
	case bool:
		boolean = &Boolean{value: value.(bool)}
	case string:
		if value.(string) == "#t" {
			boolean = &Boolean{value: true}
		} else if value.(string) == "#f" {
			boolean = &Boolean{value: false}
		} else {
			compileError("Unexpected value for NewBoolean")
		}
	default:
		return nil
	}

	if len(options) > 0 {
		boolean.parent = options[0]
	}
	return
}

func (b *Boolean) Eval() Object {
	return b
}

func (b *Boolean) String() string {
	if b.value == true {
		return "#t"
	} else {
		return "#f"
	}
}

func (b *Boolean) isBoolean() bool {
	return true
}
