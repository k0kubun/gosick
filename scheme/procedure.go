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
	// Create new local environment which has the same binding with procedure generated place
	localEnvironment := &Environment{parent: nil, binding: environment.ScopedBinding()}

	function := func(givenArguments Object) Object {
		if !arguments.IsList() || !givenArguments.IsList() {
			runtimeError("Given non-list arguments")
		}

		// assert arguments count
		expectedLength := arguments.(*Pair).ListLength()
		actualLength := givenArguments.(*Pair).ListLength()
		if expectedLength != actualLength {
			compileError("wrong number of arguments: #f requires %d, but got %d", expectedLength, actualLength)
		}

		// bind arguments to local scope
		parameters := arguments.(*Pair).Elements()
		objects := evaledObjects(givenArguments.(*Pair).Elements())
		for i, parameter := range parameters {
			if parameter.IsVariable() {
				localEnvironment.Bind(parameter.(*Variable).identifier, objects[i])
			}
		}

		// returns last eval result
		var returnValue Object
		elements := body.(*Pair).Elements()
		for _, element := range elements {
			returnValue = element.Eval()
		}
		return returnValue
	}

	return &Procedure{
		environment: localEnvironment,
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
