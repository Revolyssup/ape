package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Revolyssup/ape/token"
)

//Everything is a node in AST and has to implement a TokenLiteral method
type Node interface {
	TokenLiteral() string
	String() string // Return the exact string of code. Useful for debugging
}

//There are two types of node. Expression and Statement.

type Expression interface {
	Node
	expNode()
}

type Statement interface {
	Node
	stateNode()
}

//Our program is essentially a slice of statements.

//Root node
type Program struct {
	Statements []Statement
}

//Like other nodes, root node also implements a token literal method
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral() //Every further Node(statement/exp) will implement its tokenliteral
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

/**Expressions will have- Identifiers/ Integer Literals / Booleans **/
//Identifiers are token which hold some string like x,y,z...
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

/***Integer Literal*/
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

//String
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expNode() {}
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

//Object- key-value pairs
type ObjectLiteral struct {
	Token token.Token
	Value map[Expression]Expression
}

func (obj *ObjectLiteral) expNode() {}
func (obj *ObjectLiteral) TokenLiteral() string {
	return obj.Token.Literal
}
func (obj *ObjectLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for key, value := range obj.Value {
		out.WriteString(key.String() + ":" + value.String() + ",\n")
	}
	out.WriteString("}")
	return out.String()
}

//Array
type ArrayLiteral struct {
	Token token.Token
	Value []Expression
}

func (arr *ArrayLiteral) expNode() {}
func (arr *ArrayLiteral) TokenLiteral() string {
	return arr.Token.Literal
}

func (arr *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for _, ele := range arr.Value {
		out.WriteString(ele.String() + ",")
	}
	temp := strings.TrimSuffix(out.String(), ",") + "]"
	out.Reset()
	out.WriteString(temp)
	return out.String()
}

//Booleans
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expNode() {}
func (i *Boolean) TokenLiteral() string {
	return i.Token.Literal
}
func (b *Boolean) String() string {
	var out bytes.Buffer
	out.WriteString(b.TokenLiteral())
	return out.String()
}

/***LET STATEMENT****/
type LetStatement struct {
	Token token.Token //LET token
	Name  *Identifier
	Value Expression
}

//every statement has a method stateNode.
func (ls *LetStatement) stateNode() {}

//every statement is also a node and hence implements token literal method.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String() + " = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String() + ";")
	}

	return out.String()
}

/*****RETURN STATEMENT*******/

type ReturnStatement struct {
	Token       token.Token //RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) stateNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String() + ";")
	}

	return out.String()
}

/*************Expression Statement*******/

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) stateNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

//PREFIX
type PrefixExpression struct {
	Token           token.Token
	RightExpression Expression
	Operator        string
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) expNode() {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.RightExpression.String())
	out.WriteString(")")
	return out.String()
}

//INFIX
type InfixExpression struct {
	Token           token.Token
	LeftExpression  Expression
	RightExpression Expression
	Operator        string
}

func (ie *InfixExpression) expNode() {}
func (pe *InfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.LeftExpression.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.RightExpression.String())
	out.WriteString(")")
	return out.String()
}

//Block expressions are slice of statements ,nested under { []statements }
type BlockStatement struct {
	Token token.Token
	Stmts []Statement
}

func (bs *BlockStatement) stateNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range bs.Stmts {
		out.WriteString(stmt.String())
	}
	return out.String()
}

//If/Else expression
type IfExpression struct {
	Token     token.Token
	Condition Expression
	MainStmt  *BlockStatement
	AltStmt   *BlockStatement
}

func (ife *IfExpression) expNode() {}
func (ife *IfExpression) TokenLiteral() string {
	return ife.Token.Literal
}

func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	out.WriteString(" ")
	out.WriteString(ife.MainStmt.String())

	if ife.AltStmt != nil {
		out.WriteString(" else ")
		out.WriteString(ife.AltStmt.String())
	}
	return out.String()
}

//For expression
type ForExpression struct {
	Token     token.Token
	Condition Expression
	Stmt      *BlockStatement
}

func (fe *ForExpression) expNode() {}
func (fe *ForExpression) TokenLiteral() string {
	return fe.Token.Literal
}
func (fe *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	out.WriteString(fe.Condition.String())
	out.WriteString(" ")
	out.WriteString(fe.Stmt.String())
	return out.String()
}

//Function Literalss fn(params){body}
type FunctionLiteral struct {
	Token  token.Token //fn
	Params []*Identifier
	Body   *BlockStatement
}

func (fl *FunctionLiteral) expNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("fn")
	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")

	out.WriteString(fl.Body.String())
	return out.String()
}

//Function calls- <expression>(args). expression can be either an identifier pointing to a function literal or a function literal itself.And args is also expression.
//We can have nested funciton literals inside of out function call as arguments.

type FunctionCall struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (fc *FunctionCall) expNode() {}
func (fc *FunctionCall) TokenLiteral() string {
	return fc.Token.Literal
}

func (fc *FunctionCall) String() string {
	var out bytes.Buffer
	out.WriteString(fc.Function.String())
	out.WriteString("(")
	args := []string{}
	for _, arg := range fc.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

//Array Element
type ArrObjElement struct {
	Token token.Token //IDENT
	Name  Expression
	Index Expression
}

func (ae *ArrObjElement) expNode() {}
func (ae *ArrObjElement) TokenLiteral() string {
	return ae.Name.TokenLiteral()
}
func (ae *ArrObjElement) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.TokenLiteral())
	out.WriteString("[")
	out.WriteString(fmt.Sprint(ae.Index))
	out.WriteString("]")
	return out.String()
}
