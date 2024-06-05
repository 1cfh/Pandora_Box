package evaluator

import (
	"Pandora_Box/object"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	testsInteger := []struct {
		input    string
		expected int64
	}{
		{"+5", 5},
		//{"5", 5},
		//{"10", 10},
		//{"5", 5},
		//{"50", 50},
		//{"500", 500},
		//{"-5", -5},
		//{"-50", -50},
		//{"-500", -500},
		//{"5+5+5-4", 11},
		//{"2*3*4/8", 3},
		//{"2*(2+3)", 10},
		//{"5*3-4", 11},
		//{"(5+6*2-3)/2", 7},
	}

	for _, tt := range testsInteger {
		evaluated := testEval(tt.input)              // 解析AST
		testIntegerObject(t, evaluated, tt.expected) // 测试结果
	}

}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}
