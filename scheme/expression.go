package scheme

import (
	"errors"
	"fmt"
	"regexp"
)

type Expression SchemeType

func NewExpression(expression string) *Expression {
	return &Expression{
		expression: expression,
	}
}

func (e *Expression) String() string {
	result, err := e.Eval()
	if err != nil {
		fmt.Errorf(err.Error())
		return ""
	}
	return result.String()
}

func (e *Expression) Eval() (Type, error) {
	if matched, _ := regexp.MatchString("[0-9]*", e.expression); matched {
		return NewConst(e.expression), nil
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid or unexpected token: %s\n", e.expression))
	}
}
