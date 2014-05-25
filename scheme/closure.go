// Closure is an object returned by lambda.
// It has references for closures scoped when it was generated.

package scheme

type Closure struct {
	ObjectBase
	localBinding Binding
}

func NewClosure(parent Object) *Closure {
	return &Closure{ObjectBase: ObjectBase{parent: parent}, localBinding: make(Binding)}
}

// This method is for define syntax form.
// Define a local variable in the most inner closure.
func (c *Closure) define(identifier string, object Object) {
	c.localBinding[identifier] = object
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

func (c *Closure) boundedObject(identifier string) Object {
	boundedObject := c.localBinding[identifier]
	if boundedObject != nil {
		return boundedObject
	}

	if c.parent == nil {
		runtimeError("unbound variable: %s", identifier)
	}
	return c.parent.boundedObject(identifier)
}

func (c *Closure) binding() Binding {
	return c.localBinding
}
func (c *Closure) scopedBinding() Binding {
	return c.localBinding
}
