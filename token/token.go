package token

type TokenType string

// トークンの構造体
type Token struct {
	Type    TokenType   // トークンの種類
	Literal string      // 文字
}

// 予約語
var keyword = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// 予約語判定
func LookupIdent(ident string) TokenType {
	// 予約語ならば、その予約語のトークンを返す
	if tok, ok := keyword[ident]; ok {
		return tok
	}
	// 予約語でなければ、変数名として扱う
	return INDENT
}

const (
	ILLEGAL = "ILLEGAL" //　未知なもの
	EOF     = "EOF"     //ファイル終端

	INDENT = "INDENT" // 識別子 (add, x, yなど　変数、定数、関数の名前)
	INT    = "INT"

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="

	LT        = "<"
	GT        = ">"
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)
