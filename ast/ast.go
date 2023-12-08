package ast

import (
	"bytes"
	"github.com/takeru-a/golang_interpreterlang/token"
)

// 抽象構文木(abstract syntax tree)のノードのインターフェース
type Node interface {
	// トークンリテラルを返す
	// トークンの具体的な文字列表現
	TokenLiteral() string

	// nodeの内容を文字列で返す
	String() string
}

// 文のインターフェース
type Statement interface {
	Node
	statementNode()
}

// 式のインターフェース
type Expression interface {
	Node
	expressionNode()
}

// プログラム全体を表す
type Program struct {
	Statements []Statement
}

// 文の構文解析
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression  // 式
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// 式の文字列表現を返す
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// 再帰する
// トークンの種類によって、構文解析関数を呼び出す
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// プログラムの文字列表現を返す
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// let文の構文解析
type LetStatement struct {
	Token token.Token
	Name  *Indetifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// let x = 5;
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// return文の構文解析
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// return 5;
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// 式の構文解析
type Indetifier struct {
	Token token.Token
	Value string
}

func (i *Indetifier) expressionNode() {}

// プログラムが読み取る用
func (i *Indetifier) TokenLiteral() string { return i.Token.Literal }

// 人間が読む用
func (i *Indetifier) String() string { return i.Value }

// 整数リテラルの構文解析
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// 前置構文解析
type PrefixExpression struct {
	Token    token.Token // 前置トークン !, -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// 中置演算子の構文解析
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// 真偽値リテラルの構文解析
type Boolean struct {
	Token token.Token
	Value bool
}

func(b *Boolean) expressionNode() {}
func(b *Boolean) TokenLiteral() string { return b.Token.Literal }
func(b *Boolean) String() string { return b.Token.Literal }

// if文の構文解析
type IfExpression struct {
	Token token.Token   // 'if'トークン
	Condition Expression // 条件
	Consequence *BlockStatement  // ifの処理
	Alternative *BlockStatement  // elseの処理
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}