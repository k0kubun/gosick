package scheme

type Cons struct {
	ObjectBase
	Car *Object
	Cdr *Cons
}
