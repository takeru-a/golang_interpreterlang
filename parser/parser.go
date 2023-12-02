package parser

import (
	"fmt"

	"github.com/takeru-a/golang_interpreterlang/ast"
	"github.com/takeru-a/golang_interpreterlang/lexer"
	"github.com/takeru-a/golang_interpreterlang/token"
)

// 優先度 数字の大きさ　優先度高
const (
	_ int = iota
	LOWEST
	EQALS           // ==
	LESSGREATER     // >  <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -x !x
	CALL            // myFunction(x)
)

type Parser struct {
	l         *lexer.Lexer //字句解析器
	errors    []string     // errormessage
	curToken  token.Token  //現在のトーク
	peekToken token.Token  //次のトークン

	prefixParseFns map[token.TokenType]prefixParseFn // 前置構文解析関数
	infixParseFns  map[token.TokenType]infixParseFn  // 中置構文解析関数
}

type (
	prefixParseFn func() ast.Expression               // 前置構文解析関数 -1, !x
	infixParseFn  func(ast.Expression) ast.Expression // 中置構文解析関数  1 + 1, 1 * 1
)

// 前置構文解析関数の登録
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// 中置構文解析関数の登録
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// curToken, peekTokenを初期化
	p.nextToken()
	p.nextToken()

	// 前置構文解析関数の設定
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.INDENT, p.parseIdentifier)

	return p
}

// 識別子
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Indetifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

// errormessageを追加
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// tokenを更新する
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// tokenに応じて、適切なステートメントの構文解析関数を呼び出す
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// 終端に届くまでforを回す
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// 文の構文解析
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// letの構文解析
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	// もしletの次のtokenのTypeがINDENTでなければ
	if !p.expectPeek(token.INDENT) {
		return nil
	}
	stmt.Name = &ast.Indetifier{Token: p.curToken, Value: p.curToken.Literal}

	// letの次の次のTokenが = でなければ　(前のexpectPeekでnextTokenされている)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// Semicolonまで読み飛ばす
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

//現在のTokenと指定したTokenとのType比較
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

//次のTokenと指定したTokenとのType比較
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

//次のTokenが望むものか
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// returnの構文解析
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// セミコロンまで飛ばす
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 識別子式
func (p *Parser) parseExpression(precedence int) ast.Expression {

	// tokenに応じた関数
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	// !x, -xなど識別子の前置されているものを返却
	leftExp := prefix()

	return leftExp
}

// 識別子式の構文解析
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return stmt;
}
