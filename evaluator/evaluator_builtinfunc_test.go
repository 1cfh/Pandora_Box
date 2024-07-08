package evaluator

import (
	"Pandora_Box/object"
	"testing"
)

func TestBuiltinFunctions(t *testing.T) {
	// len内建函数的测试用例
	tests := []struct {
		input    string
		expected interface{}
	}{
		//{`len("")`, 0},
		//{`len("123456")`, 6},
		//{`len("Hello World")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		//{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errobj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if errobj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errobj.Message)
			}
		}
	}

}
