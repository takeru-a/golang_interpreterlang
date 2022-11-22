package token

type TokenType string

type Token struct{
	Type TokenType
	Literal string
}
// 予約語
var keyword = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
}

// 予約語判定
func LookupIdent(ident string) TokenType{
	if tok, ok := keyword[ident]; ok{
		return tok
	}
	return INDENT
}

const (
	ILLEGAL = "ILLEGAL" //　未知なもの
	EOF = "EOF" //ファイル終端

	INDENT = "INDENT"
	INT = "INT"

	// 演算子
	ASSIGN = "="
	PLUS = "+"

	COMMA = ","
	SEMICOLON = ";"

	LPAREN ="("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET = "LET"
)