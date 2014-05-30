// Closure is an object returned by lambda.
// It has references for closures scoped when it was generated.

package scheme

import (
	"fmt"
)

type Closure struct {
	ObjectBase
	localBinding Binding
	function     func(Object) Object
}

func NewClosure(parent Object) *Closure {
	return &Closure{ObjectBase: ObjectBase{parent: parent}, localBinding: make(Binding)}
}

// Cover the given object with a new closure.
// Insert this into tree structure between given object and its parent.
func WrapClosure(wrappedObject Object) *Closure {
	closure := NewClosure(wrappedObject.Parent())
	wrappedObject.setParent(closure)
	return closure
}

func (c *Closure) String() string {
	if c.Bounder() == nil {
		return "#<closure #f>"
	}
	return fmt.Sprintf("#<closure %s>", c.Bounder())
}

func (c *Closure) Invoke(argument Object) Object {
	return c.function(argument)
}

func (c *Closure) isClosure() bool {
	return true
}

func (c *Closure) isProcedure() bool {
	return true
}

// This method is for define syntax form.
// Define a local variable in the most inner closure.
func (c *Closure) define(identifier string, object Object) {
	c.localBinding[identifier] = object
}

// If variable is *Variable, define value.
func (c *Closure) tryDefine(variable Object, object Object) {
	if variable.isVariable() {
		c.localBinding[variable.(*Variable).identifier] = object
	}
}

// This method is for set! syntax form.
// Update most inner scoped closure's binding, otherwise raise error.
func (c *Closure) set(identifier string, object Object) {
	if c.localBinding[identifier] == nil {
		if c.parent == nil {
			runtimeError("symbol not defined")
		} else {
			c.parent.set(identifier, object)
		}
	} else {
		c.localBinding[identifier] = object
	}
}

func (c *Closure) binding() Binding {
	return c.localBinding
}
