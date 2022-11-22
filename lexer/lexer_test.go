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

func (l *Lexer) readIdentifier ()string{
	position := l.position
	for isLetter(l.ch){
		l.readChar()
	}
	return l.input[position: l.position]
}

// 空白を読み飛ばす
func (l *Lexer) skipWhitespace(){
	for l.ch == ' ' || l.ch =='\t' || l.ch == '\n' || l.ch == '\r'{
		l.readChar()
	}
}

func (l *Lexer) readNumber() string{
	position:= l.position
	for isDigit(l.ch){
		l.readChar()
	}
	return l.input[position: l.position]
}

func newToken(tokenType token.TokenType, ch byte)token.Token{
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool{
	// ASCII Code 文字の定義
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool{
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) NextToken() token.Token{
	var tok token.Token
	l.skipWhitespace()
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
	default:
		if isLetter(l.ch){
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}else if isDigit(l.ch){
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}else{
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func TestNextToken(t *testing.T){
	input := `let a = 5;
			  let b = 10;
			  let add = fn(x, y){
				x + y;
			  };
			  let result = add(a, b);
			  `

	tests := []struct{
		//期待する型・文字
		expectedType token.TokenType
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
		{token.ASSIGN,"="},
		{token.INDENT, "add"},
		{token.LPAREN, "("},
		{token.INDENT, "a"},
		{token.COMMA, ","},
		{token.INDENT, "b"},
		{token.RPAREN, ")"},
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