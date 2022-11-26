package lexer

import (
	"testing"
	"github.com/takeru-a/golang_interpreterlang/token"
)

func TestNextToken(t *testing.T) {
	input := `let a = 5;
			  let b = 10;
			  let add = fn(x, y){
				x + y;
			  };
			  let result = add(a, b);
			  !-/*5;
			  5 < 10 > 5;
			  if (5 < 10){
				return true;
			  }else{
				return false;
			  }
			  10 == 10;
			  10 != 9;
			  `

	tests := []struct {
		//期待する型・文字
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.INDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.INDENT, "b"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.INDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.INDENT, "x"},
		{token.COMMA, ","},
		{token.INDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.INDENT, "x"},
		{token.PLUS, "+"},
		{token.INDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.INDENT, "result"},
		{token.ASSIGN, "="},
		{token.INDENT, "add"},
		{token.LPAREN, "("},
		{token.INDENT, "a"},
		{token.COMMA, ","},
		{token.INDENT, "b"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
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
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN,"return"},
		{token.TRUE, "true"},
		{token.SEMICOLON,";"},
		{token.RBRACE, "}"},
		{token.ELSE,"else"},
		{token.LBRACE, "{"},
		{token.RETURN,"return"},
		{token.FALSE, "false"},
		{token.SEMICOLON,";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ,"!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		// tokenを取得
		tok := l.NextToken()
		// typeが合っているか
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		// Literal(文字) が合ってるか
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
