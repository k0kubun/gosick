// Procedure is a type for scheme procedure, which is expressed
// by lambda syntax form, like (lambda (x) x)
// When procedure has free variable, free variable must be binded when
// procedure is generated.
// So all Procedures have variable binding by Environment type (when there is
// no free variable, Procedure has Environment which is empty).

package scheme

type Procedure struct {
	ObjectBase
	function     func(Object) Object
	arguments    Object
	body         Object
	localBinding Binding
}

func NewProcedure(function func(Object) Object) *Procedure {
	return &Procedure{function: function}
}

func (p *Procedure) generateFunction(parent Object) {
	// Create local binding for procedure
	localBinding := parent.scopedBinding()
	p.localBinding = localBinding

	p.function = func(givenArguments Object) Object {
		if !p.arguments.isList() || !givenArguments.isList() {
			runtimeError("Given non-list arguments")
		}

		// assert arguments count
		expectedLength := p.arguments.(*Pair).ListLength()
		actualLength := givenArguments.(*Pair).ListLength()
		if expectedLength != actualLength {
			compileError("wrong number of arguments: requires %d, but got %d", expectedLength, actualLength)
		}

		// bind arguments to local scope
		parameters := p.arguments.(*Pair).Elements()
		objects := evaledObjects(givenArguments.(*Pair).Elements())
		for i, parameter := range parameters {
			if parameter.isVariable() {
				localBinding[parameter.(*Variable).identifier] = objects[i]
			}
		}

		// returns last eval result
		var returnValue Object
		elements := p.body.(*Pair).Elements()
		for _, element := range elements {
			returnValue = element.Eval()
		}
		return returnValue
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

func (p *Procedure) isProcedure() bool {
	return true
}

func (p *Procedure) binding() Binding {
	return p.localBinding
}

func (p *Procedure) define(identifier string, object Object) {
	p.localBinding[identifier] = object
}
