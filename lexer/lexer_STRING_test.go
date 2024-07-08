package lexer

import (
	"Pandora_Box/token"
	"testing"
)

func TestNextToken_STRING(t *testing.T) {
	input := `
	"foobar"
	"foo bar"
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{
			token.STRING,
			"foobar",
		},
		{
			token.STRING,
			"foo bar",
		},
		{
			token.EOF,
			"",
		},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		// 同时检测type及其符号形式
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong . expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong . expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
