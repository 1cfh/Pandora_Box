package parser

import (
	"Pandora_Box/ast"
	"Pandora_Box/lexer"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testStatementStruct struct {
	expectedIdentifier string
}

// TestLetStatements 测试函数以Test开头
func TestLetStatements(t *testing.T) {

	// original version
	//input := `
	//let x = ;
	//let y = 10;
	//let foobar = 838383;
	//`
	//l := lexer.New(input)
	//p := New(l)
	//// golang中使用计数引用进行回收, 所以将地址返回给program时, 函数中声明的ast.program仍然存在
	//// 获得AST
	//program := p.ParseProgram()
	//
	//// check errors
	//checkParserErrors(t, p)
	//
	//if program == nil {
	//	t.Fatalf("ParseProgram() returned nil")
	//}
	//if len(program.Statements) != 3 {
	//	t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	//}
	//
	//tests := []testStatementStruct{
	//	{"x"},
	//	{"y"},
	//	{"foobar"},
	//}
	//
	//for i, tt := range tests {
	//	stmt := program.Statements[i]
	//	if !testLetStatement(t, stmt, tt.expectedIdentifier) {
	//		return
	//	}
	//}

	// new version
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{
			input:              "let x = 5;",
			expectedIdentifier: "x",
			expectedValue:      5,
		},
		{
			input:              "let y = true;",
			expectedIdentifier: "y",
			expectedValue:      true,
		},
		{
			input:              "let foobar = y;",
			expectedIdentifier: "foobar",
			expectedValue:      "y",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}

	}

}

// testLetStatement 测试函数
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	// 检查let语句
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	/*
		此处为Type Assertion类型断言的用法
		因为此处为接口实现
		需要判断传入的值为哪个接口
	*/
	letStmt, ok := s.(*ast.LetStatement)

	// 检查是否为LetStatement
	if !ok {
		t.Errorf("s not *ast.Letstatement. got=%T", s)
		return false
	}

	//
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%T", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993 332;
`
	// 生成词法分析器
	l := lexer.New(input)
	// 生成语法分析器
	p := New(l)

	// 语法分析
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement) // assert
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// 使用testIdentifier辅助函数进行替代如下代码
	ok = testIdentifier(t, stmt.Expression, input)

	if !ok {
		t.Errorf("Test Error")
		os.Exit(-1)
	}

	//ident, ok := stmt.Expression.(*ast.Identifier)
	//if !ok {
	//	t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	//}
	//
	//if ident.Value != "foobar" {
	//	t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	//}
	//
	//if ident.TokenLiteral() != "foobar" {
	//	t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	//}

}

func TestIntegerLiteralExpression(t *testing.T) {

	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral.got=%T", stmt.Expression)
	}

	// 使用辅助函数进行替换
	val, err := strconv.ParseInt(strings.ReplaceAll(input, ";", ""), 10, 64)
	if err != nil {
		t.Errorf("ParseInt Error")
		os.Exit(-1)
	}
	ok = testIntegerLiteral(t, stmt.Expression, val)
	if !ok {
		t.Errorf("Test Error")
		os.Exit(-1)
	}

	//literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	//
	//if literal.Value != 5 {
	//	t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	//}
	//
	//if literal.TokenLiteral() != "5" {
	//	t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	//}

}

func TestParsingPrefixExpressions(t *testing.T) {
	// 前缀表达式测试用例(一元表达式)
	prefixTests := []struct {
		input        string // 完整输入
		operator     string // 前缀运算符
		integerValue int64  // 表达式
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	// 遍历所有的测试用例
	for _, tt := range prefixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		// 测试statement
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		// 测试表达式(包括运算符)
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}

	}

}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string // 完整的表达式
		leftValue  int64  // 左表达式
		operator   string // 运算符
		rightValue int64  // 右表达式
	}{
		// 测试用例1
		//{"(1+3) > 2", (1 + 3), ">", 2},
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		//ok = testInfixExpression(t, stmt.Expression, stmt.Expression.(*ast.InfixExpression).Left, stmt.Expression.(*ast.InfixExpression).Operator, stmt.Expression.(*ast.InfixExpression).Operator)
		//if !ok {
		//	t.Errorf("Test Error")
		//	os.Exit(-1)
		//}

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			t.Fatal("error: exp.Left")
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	//opTests := []struct {
	//	input    string
	//	expected string
	//}{
	//	{
	//		"1+2+3",
	//		"((1+2)+3)",
	//	},
	//	{
	//		"-a * b",
	//		"((-a)*b)",
	//	},
	//	{
	//		"!-a",
	//		"(!(-a))",
	//	},
	//	{
	//		"a+b+c",
	//		"((a+b)+c)",
	//	},
	//}

	opTests := []struct {
		input    string
		expected string
	}{
		{
			input:    "1+(2+3)+4",
			expected: "((1+(2+3))+4)",
		},
		{
			input:    "a + add(b*c)+c",
			expected: "((a+add((b*c)))+c)",
		},
		{
			input:    "add(1, 1+2, add(1,2))",
			expected: "add(1, (1+2), add(1, 2))",
		},
	}

	for _, tt := range opTests {

		l := lexer.New(tt.input) // 构建词法解析器
		p := New(l)              // 执行词法解析器 利用其中的nextToken逐个解析成词法单元

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}

	}

}

func TestBooleanExpression(t *testing.T) {
	booleanTest := []struct {
		strvalue string
		boolVal  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range booleanTest {
		l := lexer.New(tt.strvalue)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		// ok = testBooleanLiteral(t, stmt.Expression, booleanTest)
		ok = testLiteralExpression(t, stmt.Expression, tt.boolVal)
		if !ok {
			t.Fatal("test Boolean Literal Error")
		}
	}

}

func TestIfExpression(t *testing.T) {
	input := `if (x<y) {x}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil.got=%+v", exp.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if(x<y){x}else{y}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// 检测else部分 -- 是否为一个block
	//if exp.Alternative == nil {
	//	t.Errorf("exp.Alternative.Statements was not nil.got=%+v", exp.Alternative)
	//}

	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil.got=%+v", exp.Alternative)
	}

}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x,y){x+y;}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	paramsTests := []struct {
		input          string
		expectedParams []string
	}{
		{
			input:          "fn() {};",
			expectedParams: []string{},
		},
		{
			input:          "fn(x) {}",
			expectedParams: []string{"x"},
		},
		{
			input:          "fn(x,y,z) {}",
			expectedParams: []string{"x", "y", "z"},
		},
	}

	for _, tt := range paramsTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		// 检测形参长度
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, but got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		//
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}

	}

}

func TestCallExpressionParsing(t *testing.T) {
	callTests := "add(1, 2 * 3 ,4+5);"

	l := lexer.New(callTests)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of argument. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)

}

// testInfixExpression 测试 InfixExpression
func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%S)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

// testIntegerLiteral 测试IntegerLiteral的辅助函数
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	// 断言判断是否为IntegerLiteral
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

/*
testLiteralExpression 测试Literal汇总函数
*/
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// testIdentifier 测试标识符
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	// boolVal, err := strconv.ParseBool(value)

	//if err != nil {
	//	t.Fatal("ParseBool Error")
	//}

	if boolean.Value != value {
		t.Errorf("boolean.Value %t. got=%t", value, boolean.Value)
		return false
	}

	boolTokenLiteral := strconv.FormatBool(value)

	if boolean.TokenLiteral() != boolTokenLiteral {
		t.Errorf("boolean.TokenLiteral not %s. got=%s", boolTokenLiteral, boolean.TokenLiteral())
		return false
	}

	return true
}
