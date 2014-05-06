// Procedure is a type for scheme procedure, which is expressed
// by lambda syntax form, like (lambda (x) x)
// When procedure has free variable, free variable must be binded when
// procedure is generated.
// So all Procedures have variable binding by Environment type (when there is
// no free variable, Procedure has Environment which is empty).

package scheme

type Procedure struct {
	ObjectBase
	environment *Environment
	function    func(Object) Object
}

func NewProcedure(function func(Object) Object) *Procedure {
	return &Procedure{
		environment: nil,
		function:    function,
	}
}

func (p *Procedure) Eval() Object {
	return p
}

func (p *Procedure) invoke(argument Object) Object {
	return p.function(argument)
}
