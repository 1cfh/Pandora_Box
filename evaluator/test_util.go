package evaluator

import (
	"Pandora_Box/lexer"
	"Pandora_Box/object"
	"Pandora_Box/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)       // input为输入的源码字符串, 生成Token序列
	p := parser.New(l)          // 根据生成的Token序列 , 返回Parser
	program := p.ParseProgram() // 调用Parser的ParseProgram的方法生成抽象语法树

	return Eval(program) // 解析抽象语法树
}
