package evaluator

import (
	"Pandora_Box/object"
	"testing"
)

func TestEvalBooleanExpression(t *testing.T) {
	testsBoolean := []struct {
		input    string
		expected bool
	}{
		//{"true", true},
		//{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1==2", false},
		{"1==1", false},
		{"1!=1", false},
		{"1!=2", true},
		{"1>1", false},
		{"1<1", false},
		{"true==true", true},
		{"false==false", true},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range testsBoolean {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}

}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%T, want=%t", result.Value, expected)
		return false
	}

	return true
}
