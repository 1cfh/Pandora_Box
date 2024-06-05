package evaluator

import (
	"Pandora_Box/object"
	"testing"
)

func TestIfElseExpression(t *testing.T) {
	testsIfElse := []struct {
		input    string
		expected interface{}
	}{
		//{"if (true) { 10 }", 10},
		//{"if (false) { 10 }", nil},
		//{"if (1) { 10 }", 10},
		{"if (1<2) { 10 }", 10},
		//{"if (1>2) { 10 }", nil},
		//{"if (1>2) { 10 } else { 20 }", 20},
		//{"if (1<2) { 10 } else { 20 }", 10},
	}

	for _, tt := range testsIfElse {
		evaluated := testEval(tt.input)
		intExpected, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(intExpected))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got =%T (%+v)", obj, obj)
		return false
	}
	return true
}
