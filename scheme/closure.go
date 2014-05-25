// Closure is an object returned by lambda.
// It has references for closures scoped when it was generated.

package scheme

type Closure struct {
	ObjectBase
	localBinding Binding
	function     func(Object) Object
}

func NewClosure(parent Object) *Closure {
	return &Closure{ObjectBase: ObjectBase{parent: parent}, localBinding: make(Binding)}
}

func (c *Closure) String() string {
	return "#<closure #f>"
}

func (c *Closure) Invoke(argument Object) Object {
	return c.function(argument)
}

func (c *Closure) isClosure() bool {
	return true
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

func (c *Closure) binding() Binding {
	return c.localBinding
}

func (c *Closure) scopedBinding() Binding {
	return c.localBinding
}
