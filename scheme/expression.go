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
		fmt.Errorf(err.Error())
		return ""
	}
	return result.String()
}

func (e *Expression) Eval() (Type, error) {
	if e.isParenthesized() {
		expressions := strings.Split(e.parenthesesStrippedExpression(), " ")
		function := &Expression{expression: expressions[0]}

		childs := []*Expression{}
		for _, expression := range expressions[1:] {
			childs = append(childs, &Expression{expression: expression})
		}

		return function.evalFunction(childs)
	} else {
		if matched, _ := regexp.MatchString("[0-9]*", e.expression); matched {
			return NewConst(e.expression), nil
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid or unexpected token: %s\n", e.expression))
		}
	}
}

func (e *Expression) evalFunction(args []*Expression) (Type, error) {
	switch e.expression {
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
	reg, err := regexp.Compile("^\\(")
	if err != nil {
		log.Fatal(err)
	}
	e.expression = reg.ReplaceAllString(e.expression, "")

	reg, err = regexp.Compile("\\)$")
	if err != nil {
		log.Fatal(err)
	}
	e.expression = reg.ReplaceAllString(e.expression, "")

	return e.expression
}
