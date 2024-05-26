package ast

import (
	"Pandora_Box/token"
	"bytes"
)

// Notes: golang中实现接口是隐式的

// Node 必须实现TokenLiteral方法
type Node interface {
	TokenLiteral() string // 用于调试和测试
	String() string
}

/*
	Statement和Expression继承了Node接口
	Statement
*/
// Statement 语句接口
type Statement interface {
	Node
	statementNode()
}

// Expression 表达式接口
type Expression interface {
	Node
	expressionNode()
}

/*
Program 中实现了Statement接口,
作用是作为AST的根节点, 其中的Statement数组代表了语句序列
*/
type Program struct {
	Statements []Statement
}

// TokenLiteral Program需要实现Node中的TokenLiteral方法
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// String 可以理解为这是一种魔术方法
func (p *Program) String() string {
	var out bytes.Buffer // 创建缓冲区

	for _, s := range p.Statements {
		out.WriteString(s.String()) // 将每个statement写入到缓冲区
	}
	return out.String()
}

/*
LetStatement 中必须是实现Expression

	let x = 5;
	let x = 5 * 5;
*/
type LetStatement struct {
	Token token.Token // token.LET 词法单元
	Name  *Identifier // 变量的标识符
	Value Expression  // 表达式
}

// TokenLiteral 返回LetStatement对象中的Token Literal
// 实现Node接口
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// statementNode
// 实现Statement接口
func (ls *LetStatement) statementNode() {

}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	/* let identifier = xxx; */

	// write let token
	out.WriteString(ls.TokenLiteral() + " ")

	// write the identifier
	out.WriteString(ls.Name.String())

	out.WriteString(" = ")

	// write expression
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	// ;
	out.WriteString(";")

	return out.String()
}

/*
ReturnStatement return语句的statement节点
*/
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

/*
ExpressionStatement 针对于表达式声明的statement节点
*/
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {

}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/*
Identifier 变量标识符结构体
*/
type Identifier struct {
	Token token.Token // token.IDENT 词法单元
	Value string
}

// TokenLiteral 实现了TokenLiteral
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {

}

// String
func (i *Identifier) String() string {
	return i.Value
}

/*
IntegerLiteral
*/

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {

}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

/*
PrefixExpression 前缀表达式
*/
type PrefixExpression struct {
	Token    token.Token // 前缀词法单元
	Operator string      // 前缀表达式中的运算符
	Right    Expression  // 表达式
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String 作为打印对象的魔术方法
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
InfixExpression 中缀表达式
*/
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
