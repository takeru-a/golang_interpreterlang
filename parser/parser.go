package parser

import (
	"fmt"

	"github.com/takeru-a/golang_interpreterlang/ast"
	"github.com/takeru-a/golang_interpreterlang/lexer"
	"github.com/takeru-a/golang_interpreterlang/token"
)

type Parser struct{
	l *lexer.Lexer //字句解析器
	errors [] string
	curToken token.Token //現在のトークン
	peekToken token.Token //次のトークン
}

func New(l *lexer.Lexer) *Parser{
	p := &Parser{l: l, errors: []string{}}
	// curToken, peekTokenを初期化
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string{
	return p.errors
}

func (p *Parser) peekError(t token.TokenType){
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program{
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	// 終端に届くまでforを回す
	for p.curToken.Type != token.EOF{
		stmt := p.parseStatement()
		if stmt != nil{
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement{
	switch p.curToken.Type{
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// letの構文解析
func (p *Parser) parseLetStatement() *ast.LetStatement{
	stmt := &ast.LetStatement{Token: p.curToken}
	// もしletの次のtokenのTypeがINDENTでなければ
	if !p.expectPeek(token.INDENT){
		return nil
	}
	stmt.Name = &ast.Indetifier{Token: p.curToken, Value: p.curToken.Literal}

	// letの次の次のTokenが = でなければ　(前のexpectPeekでnextTokenされている)
	if !p.expectPeek(token.ASSIGN){
		return nil
	}
	// Semicolonまで読み飛ばす
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

//現在のTokenと指定したTokenとのType比較
func (p *Parser) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t
}
//次のTokenと指定したTokenとのType比較
func (p *Parser) peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}

//次のTokenが望むものか
func (p *Parser)expectPeek(t token.TokenType) bool{
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	}else{
		p.peekError(t)
		return false
	}
}