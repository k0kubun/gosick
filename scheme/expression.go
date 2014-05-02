package scheme

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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
		log.Fatal(err.Error())
		return ""
	}
	return result.String()
}

func (e *Expression) Eval() (Type, error) {
	if e.isParenthesized() {
		expressions := strings.Split(e.parenthesesStrippedExpression(), " ")
		function := &Expression{expression: expressions[0]}

		args := []*Expression{}
		for _, expression := range expressions[1:] {
			args = append(args, &Expression{expression: expression})
		}

		return function.evalFunction(args)
	} else {
		if matched, _ := regexp.MatchString("[0-9]*", e.expression); matched {
			value, err := NewConst(e.expression).Eval()
			if err != nil {
				return nil, err
			}
			return value, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid or unexpected token: %s\n", e.expression))
		}
	}
}

func (e *Expression) evalFunction(args []*Expression) (Type, error) {
	switch e.expression {
	case "number?":
		if len(args) == 1 {
			arg, err := args[0].Eval()
			if err != nil {
				return nil, err
			}
			return NewBoolean(arg.IsNumber()), nil
		} else {
			return nil, errors.New(
				fmt.Sprintf("Wrong number of arguments: %s requires 1, but got %d", e.expression, len(args)),
			)
		}
	case "+":
		sum := 0
		for _, arg := range args {
			value, err := arg.Eval()
			if err != nil {
				return nil, err
			}
			number, _ := strconv.Atoi(value.String())
			sum += number
		}

		return NewNumber(sum), nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid or unexpected function: %s\n", e.expression))
	}
}

func (e *Expression) isParenthesized() bool {
	matched, _ := regexp.MatchString("^\\(.*\\)$", e.expression)
	return matched
}

func (e *Expression) parenthesesStrippedExpression() string {
	expression := strings.TrimPrefix(e.expression, "(")
	expression = strings.TrimSuffix(expression, ")")
	return expression
}
