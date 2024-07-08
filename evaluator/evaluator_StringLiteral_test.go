package evaluator

import (
	"Pandora_Box/object"
	"testing"
)

func TestStringLiteral(t *testing.T) {
	input := `"Hello world!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello world!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}

}

func TestStringInfixExpression(t *testing.T) {
	//input := `"Hello" + " " + "World!"`
	//evaluated := testEval(input)
	//
	//str, ok := evaluated.(*object.String)
	//if !ok {
	//	t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	//}
	//
	//if str.Value != "Hello World!" {
	//	t.Errorf("String has wrong value. got=%q", str.Value)
	//}

	input := `"Hello" == "Hello"`
	evaluated := testEval(input)

	boolean, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+v)", evaluated, evaluated)
	}

	if !boolean.Value {
		t.Errorf("String has wrong value. got=%t", boolean.Value)
	}

	input2 := `"Hello" != "World"`
	evaluated2 := testEval(input2)

	boolean2, ok2 := evaluated2.(*object.Boolean)
	if !ok2 {
		t.Fatalf("object is not Boolean. got=%T (%+v)", evaluated, evaluated)
	}
	if !boolean2.Value {
		t.Errorf("String has wrong value. got=%t", boolean.Value)
	}

}
