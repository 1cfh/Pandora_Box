package evaluator

import (
	"fmt"
	"testing"
)

func TestClosures(t *testing.T) {
	//	input := `
	//let newAdder = fn(x){
	//	fn(y) { x };
	//};
	//let addTwo = newAdder(2);
	//addTwo(3);
	//`

	input2 := `
let applyFunc = fn(a, b, func) { func(a, b) };
applyFunc(2, 2, add);

`

	val := testEval(input2)
	fmt.Println(val)
	// testIntegerObject(t, val, 4)
}
