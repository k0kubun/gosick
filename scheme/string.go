// String is a type for scheme string object, which is
// expressed like "string".

package scheme

type String struct {
	ObjectBase
	text string
}

func NewString(text string) *String {
	return &String{text: text}
}
