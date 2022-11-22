package lexer

import(
	"testing"
	"github.com/takeru-a/golang_interpreterlang/token"
)

type Lexer struct{
	input string
	position int //現在の位置
	readPosition int //次の文字
	ch byte // 検査中の文字
}

func New(input string) *Lexer{
	l := &Lexer{input: input}
	// 初期化
	l.readChar()
	return l
}

func (l *Lexer) readChar(){
	if l.readPosition >= len(l.input){
		// 終端 ASCII
		l.ch = 0
	}else{
		//検査する文字の指定
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte)token.Token{
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token{
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal =""
		tok.Type  = token.EOF
	}
	l.readChar()
	return tok
}

func TestNextToken(t *testing.T){
	input := `=+(){},;`

	tests := []struct{
		//期待する型・文字
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests{
		// tokenを取得
		tok := l.NextToken()
		// typeが合っているか
		if tok.Type != tt.expectedType{
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
			 i, tt.expectedType, tok.Type)
		}
		// Literal(文字) が合ってるか
		if tok.Literal != tt.expectedLiteral{
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
			i, tt.expectedLiteral, tok.Literal)
		}
	}
}