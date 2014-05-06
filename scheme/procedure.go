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
	arguments   Object
	body        Object
}

func NewProcedure(environment *Environment, arguments Object, body Object) *Procedure {
	function := func(Object) Object {
		var returnValue Object

		elements := body.(*Pair).Elements()
		for _, element := range elements {
			returnValue = element.Eval()
		}
		return returnValue
	}

	return &Procedure{
		environment: environment,
		function:    function,
		arguments:   arguments,
		body:        body,
	}
}

func (p *Procedure) String() string {
	return "#<closure #f>"
}

func (p *Procedure) Eval() Object {
	return p
}

func (p *Procedure) Invoke(argument Object) Object {
	return p.function(argument)
}

func (p *Procedure) IsProcedure() bool {
	return true
}
