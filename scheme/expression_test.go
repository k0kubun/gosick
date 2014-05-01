package scheme

import (
	"testing"
)

func TestExpression(t *testing.T) {
	expectedResultByExpression := map[string]string{
		"1":        "1",
		"20":       "20",
		"10000000": "10000000",
	}

	for expression, expectedResult := range expectedResultByExpression {
		actual := NewExpression(expression).String()

		if expectedResult != actual {
			t.Errorf("Expected: %v, Got: %v", expectedResult, actual)
			return
		}
	}
}
