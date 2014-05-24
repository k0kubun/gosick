package scheme

import (
	"testing"
)

type parserTest struct {
	source  string
	results []string
}

var parserTests = []parserTest{
	parseTest(" x ", "x"),
	parseTest("'( 1 )", "'(1)"),
	parseTest("''''hello", "''''hello"),
	parseTest("( x 1 )", "(x 1)"),
	parseTest("( x ( 1 ) )", "(x (1))"),
}

func parseTest(source string, results ...string) parserTest {
	return parserTest{source: source, results: results}
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		i := NewInterpreter(test.source)
		parseResults := []string{}
		for i.Peek() != EOF {
			object := i.Parse(i)
			if object != nil {
				parseResults = append(parseResults, object.String())
			}
		}

		for i := 0; i < len(test.results); i++ {
			expect := test.results[i]
			actual := parseResults[i]
			if actual != expect {
				t.Errorf("%s =>\n got: %s\nwant: %s", test.source, actual, expect)
			}
		}
	}
}
