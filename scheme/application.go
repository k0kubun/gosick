// Application is a type to express scheme procedure application.
// Application has a procedure and its argument as list which consists
// of Pair.

package scheme

type Application struct {
	ObjectBase
	procedure Object
	arguments Object
}

type Invoker interface {
	Invoke(Object) Object
}

func NewApplication(parent Object) *Application {
	return &Application{
		ObjectBase: ObjectBase{parent: parent},
	}
}

func (a *Application) Eval() Object {
	return a.applyProcedure()
}

func (a *Application) String() string {
	// Exceptional handling for special form: quote
	list := a.toList()
	firstObject := list.ElementAt(0)
	if firstObject.isVariable() && firstObject.(*Variable).content() == builtinSyntaxes["quote"] {
		if a.arguments.isNull() {
			return "(quote)"
		} else {
			return "'" + a.arguments.(*Pair).ElementAt(0).String()
		}
	}

	return list.String()
}

func (a *Application) toList() *Pair {
	list := NewPair(a.Parent())
	list.Car = a.procedure
	list.Car.setParent(list)
	list.Cdr = a.arguments
	list.Cdr.setParent(list)
	return list
}

func (a *Application) applyProcedure() Object {
	evaledObject := a.procedure.Eval()

	switch evaledObject.(type) {
	case Invoker:
		return evaledObject.(Invoker).Invoke(a.arguments)
	default:
		runtimeError("invalid application")
		return nil
	}
}

func (a *Application) isApplication() bool {
	return true
}
