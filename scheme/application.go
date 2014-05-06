// Application is a type to express scheme procedure application.
// Application has a procedure and its argument as list which consists
// of Pair.

package scheme

type Application struct {
	ObjectBase
	procedureVariable Object
	arguments         Object // expect *Pair
	environment       *Environment
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
	if a.environment == nil {
		compileError("Procedure does not have environment")
	}
	return a.environment.invokeProcedure(a.procedureVariable, a.arguments)
}

func (a *Application) IsApplication() bool {
	return true
}
