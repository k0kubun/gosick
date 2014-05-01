package scheme

import (
	"errors"
	"fmt"
	"regexp"
)

type Const SchemeType

func NewConst(expression string) *Const {
	return &Const{
		expression: expression,
	}
}

func (c *Const) String() string {
	return c.expression
}

func (c *Const) Eval() (Type, error) {
	if matched, _ := regexp.MatchString("[0-9]*", c.expression); matched {
		return NewNumber(c.expression), nil
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid or unexpected token: %s\n", c.expression))
	}
}
