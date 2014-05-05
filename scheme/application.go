// Application is a type to express scheme procedure application.
// Application has a procedure and its argument as list which consists
// of Pair.

package scheme

import (
	"log"
)

type Application struct {
	ObjectBase
	procedureVariable Object
	arguments         Object // expect *Pair
	environment       *Environment
}

func (a *Application) String() string {
	result := a.applyProcedure()
	if result == nil {
		log.Fatal("Procedure appication returns nil")
	}
	return result.String()
}

func (a *Application) applyProcedure() Object {
	if a.environment == nil {
		log.Fatal("Procedure does not have environment")
	}
	return a.environment.invokeProcedure(a.procedureVariable, a.arguments)
}

func (a *Application) IsApplication() bool {
	return true
}
