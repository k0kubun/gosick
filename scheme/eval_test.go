package scheme

import (
	"testing"
)

func TestEval(t *testing.T) {
	expectedResultByExpression := map[string]string{
		" 1 ":      "1",
		" 　 20　 　": "20",
		" 　	300 　	": "300",
	}

	for expression, expectedResult := range expectedResultByExpression {
		result, err := Eval(expression)
		if err != nil {
			t.Errorf(err.Error())
		}
		actual := result.String()

		if expectedResult != actual {
			t.Errorf("Expected: %v, Got: %v", expectedResult, actual)
			return
		}
	}
}
