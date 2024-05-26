package ast

import (
	"Pandora_Box/token"
	"testing"
)

// TestString test the let statement: let x = 1;
func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"}, // 词法单元
					Value: "x",                                          // Identifier的Value字段
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "1"},
					Value: "1",
				},
			},
		},
	}

	if program.String() != "let x = 1;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}

}
