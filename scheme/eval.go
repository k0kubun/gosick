package scheme

import (
	"strings"
)

// This function takes string and interprets it as scheme program in the top-level environment.
func Eval(program string) (Type, error) {
	normalizedProgram := normalizeExpression(program)
	return NewExpression(normalizedProgram).Eval()
}

func normalizeExpression(expression string) string {
	strippedExpression := strings.Trim(expression, " 　	")
	return strippedExpression
}
