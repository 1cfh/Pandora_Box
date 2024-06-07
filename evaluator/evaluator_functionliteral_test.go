package evaluator

import (
	"Pandora_Box/object"
	"testing"
)

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2;};"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	// 参数列表的长度
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	// 参数的名字
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameters is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x+2)"

	// 代码块
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}

}

func TestFunctionApplication(t *testing.T) {
	funcTests := []struct {
		input    string
		expected int64
	}{
		{
			"let identity = fn(x){x;};identity(5);",
			5,
		},
		{
			"let identity = fn(x){ return x;};identity(5);",
			5,
		},
		{
			"let double = fn(x){ x*2;}; double(5);",
			10,
		},
		{
			"let add = fn(x,y){ x + y;}; add(5, 5);",
			10,
		},
		{
			"let add = fn(x,y){ x + y;}; add(5+5, add(5, 5));",
			20,
		},
		{
			"fn(x){ x; }(5)",
			5,
		},
	}

	for _, tt := range funcTests {
		val := testEval(tt.input)
		// fmt.Println(val)
		testIntegerObject(t, val, tt.expected)
	}

}
