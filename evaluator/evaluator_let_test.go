package evaluator

import "testing"

func TestLetStatements(t *testing.T) {
	letTests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5*5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range letTests {
		letVal := testEval(tt.input)
		testIntegerObject(t, letVal, tt.expected)
	}

}
