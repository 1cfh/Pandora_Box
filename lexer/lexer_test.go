package lexer

import (
	"Pandora_Box/token"
	"testing"
)

type LexicalUnit struct {
	expectedType    token.TokenType
	expectedLiteral string
}

// 测试用例
const (
	input1 = `=+(){},;`

	input2 = `let five = 5;
	let ten = 10;
	let add = fn(x,y){
		x + y;
	};
	let result = add(five, ten);
`
	input3 = `!-/*5;5 < 10 > 5;`

	input4 = `
	if ( 5 < 10 ){
		return true;
	}else{
		return false;
	}
`

	input5 = `
	10 == 10;
	10 != 9;
`
)

// 测试用例对应的Token列表
var (
	test1 = []LexicalUnit{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	test2 = []LexicalUnit{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}
	test3 = []LexicalUnit{
		// !-/*5;5 <10> 5;
		{token.EXCLAMATION, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	test4 = []LexicalUnit{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	test5 = []LexicalUnit{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
)

func TestNextToken(t *testing.T) {
	input := input5

	// 测试样例中首先对如下关键字及运算符进行测试
	//tests := []struct {
	//	expectedType    token.TokenType
	//	expectedLiteral string
	//}{
	//	{token.ASSIGN, "="},
	//	{token.PLUS, "+"},
	//	{token.LPAREN, "("},
	//	{token.RPAREN, ")"},
	//	{token.LBRACE, "{"},
	//	{token.RBRACE, "}"},
	//	{token.COMMA, ","},
	//	{token.SEMICOLON, ";"},
	//}

	// test the input string
	list := New(input)
	// i => index ; tt => TokenType
	for i, tt := range test5 {
		tok := list.NextToken()

		// 同时检测type及其符号形式
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong . expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong . expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}

		// fmt.Println(tok)

	}

}
