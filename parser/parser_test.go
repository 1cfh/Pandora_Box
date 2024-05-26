package parser

import (
	"Pandora_Box/ast"
	"Pandora_Box/lexer"
	"fmt"
	"testing"
)

type testStatementStruct struct {
	expectedIdentifier string
}

// TestLetStatements 测试函数以Test开头
func TestLetStatements(t *testing.T) {
	input := `
	let x = ;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)
	// golang中使用计数引用进行回收, 所以将地址返回给program时, 函数中声明的ast.program仍然存在
	// 获得AST
	program := p.ParseProgram()

	// check errors
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []testStatementStruct{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
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

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}

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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}

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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string // 完整的表达式
		leftValue  int64  // 左表达式
		operator   string // 运算符
		rightValue int64  // 右表达式
	}{
		// 测试用例1
		{"(1+3) > 2", (1 + 3), ">", 2},
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
	opTests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a)*b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a+b+c",
			"((a+b)+c)",
		},
	}

	for _, tt := range opTests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}

	}

}
