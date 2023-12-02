package ast

import (
	"testing"
	"github.com/takeru-a/golang_interpreterlang/token"
)


func TestString(t *testing.T){
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Indetifier{
					Token: token.Token{Type: token.INDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Indetifier{
					Token: token.Token{Type: token.INDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;"{
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}