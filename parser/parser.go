package parser

import (
	"fmt"
	"strconv"

	"github.com/takeru-a/golang_interpreterlang/ast"
	"github.com/takeru-a/golang_interpreterlang/lexer"
	"github.com/takeru-a/golang_interpreterlang/token"
)

// 優先度 数字の大きさ　優先度高
const (
	_ int = iota
	LOWEST
	EQALS       // ==
	LESSGREATER // >  <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x !x
	CALL        // myFunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQALS,
	token.NOT_EQ:   EQALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	//　中置構文解析関数の設定
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

// 次のトークンの優先度を返却する
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 現在のトークンの優先度を返却する
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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

// 識別子
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Indetifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 前置構文がないときのエラーメッセージを追加
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 識別子式
func (p *Parser) parseExpression(precedence int) ast.Expression {

	// tokenに応じた関数
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// !x, -x, 5, numなど識別子を返却
	leftExp := prefix()

	// 次のtokenが優先度が高い場合、中置構文解析関数を呼び出す
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		// 中置構文解析関数がない場合
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// 識別子式の構文解析
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 整数リテラルの構文解析
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// 前置表現の構文解析
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 中置演算子の構文解析
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedences := p.curPrecedence()
	p.nextToken()
	// 中置演算子の右側の式を構文解析
	expression.Right = p.parseExpression(precedences)

	return expression
}

// 真偽値リテラルの構文解析
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// グループ化された式
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// ifの構文
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// {}の中
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}