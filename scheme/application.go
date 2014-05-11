// Application is a type to express scheme procedure application.
// Application has a procedure and its argument as list which consists
// of Pair.

package scheme

type Application struct {
	ObjectBase
	procedure Object
	arguments Object
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
	result := a.Eval()
	if result == nil {
		compileError("Procedure appication returns nil")
	}
	return result.String()
}

func (a *Application) applyProcedure() Object {
	evaledObject := a.procedure.Eval()
	if !evaledObject.isProcedure() {
		runtimeError("invalid application")
	}
	procedure := evaledObject.(*Procedure)
	return procedure.Invoke(a.arguments)
}

func (a *Application) isApplication() bool {
	return true
}
