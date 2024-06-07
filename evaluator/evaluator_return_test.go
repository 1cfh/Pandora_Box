package evaluator

import "testing"

func TestReturnStatements(t *testing.T) {
	testsRet := []struct {
		input    string
		expected int64
	}{
		//{"return 10;", 10},
		//{"return 10; 9;", 10},
		//{"return 2*5; 9", 10},
		//{"9;return 2 *5;9", 10},
		{
			`
	if(10>11){
		if(6>3){
			return 10;
		}
		return 1;
	}else{
		return 2;
	}
`,
			2,
		},
	}

	for _, tt := range testsRet {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
