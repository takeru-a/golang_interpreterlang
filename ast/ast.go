package ast

import "github.com/takeru-a/golang_interpreterlang/token"

type Node interface{
	TokenLiteral() string
}

type Statement interface{
	Node
	statementNode()
}

type Expression interface{
	Node
	expressionNode()
}

type Program struct{
	Statements []Statement
}

// 再帰する
func (p *Program) TokenLiteral() string{
	if len(p.Statements) > 0{
		return p.Statements[0].TokenLiteral()
	}else{
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name *Indetifier
	Value Expression
}

func (ls *LetStatement) statementNode(){}
func (ls *LetStatement) TokenLiteral() string {return ls.Token.Literal}

type Indetifier struct{
	Token token.Token
	Value string
}

func (i *Indetifier) expressionNode() {}
func (i *Indetifier) TokenLiteral() string {return i.Token.Literal}