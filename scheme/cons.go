package scheme

type Cons struct {
	ObjectBase
	Car Object
	Cdr *Cons
}

func (c *Cons) String() string {
	if c.Car == nil && c.Cdr == nil {
		return "()"
	} else {
		return "Not implemented."
	}
}
