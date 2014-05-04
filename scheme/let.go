// Let is a type to express code block with local variable binding,
// which is written like (let ((x 1)) x).
// Each Let has Environment to hold local variable binding.

package scheme

type Let struct {
	ObjectBase
}

func NewLet() *Let {
	return &Let{}
}
