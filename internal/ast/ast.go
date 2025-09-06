package ast

import (
	"bytes"
	"strings"

	"ixion/internal/token"
)

// Node is the base interface for all AST nodes.
type Node interface {
	// TokenLiteral returns the literal value of the token associated with the node.
	TokenLiteral() string
	String() string
}

// Statement represents a statement in the language.
// Statements do not produce values.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the language.
// Expressions produce values.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST.
// It is a collection of statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
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

// --- Statements ---

// VarStatement represents a variable declaration statement.
// e.g., var x = 5;
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Text }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ConstStatement represents a constant declaration statement.
// e.g., const x = 5;
type ConstStatement struct {
	Token token.Token // the token.CONST token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Text }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")
	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement.
// e.g., return 10;
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Text }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents a statement that consists of a single expression.
// e.g., x + 10;
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Text }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block of statements.
// e.g., { statement1; statement2; }
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Text }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// PrintStatement represents a print statement.
// e.g., print "hello world";
type PrintStatement struct {
	Token token.Token // the 'print' token
	Value Expression
}

func (ps *PrintStatement) statementNode()       {}
func (ps *PrintStatement) TokenLiteral() string { return ps.Token.Text }
func (ps *PrintStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + "(")
	if ps.Value != nil {
		out.WriteString(ps.Value.String())
	}
	out.WriteString(");")
	return out.String()
}

// ForStatement represents a for loop statement.
// e.g., for (x < y) { ... }
type ForStatement struct {
	Token     token.Token // the 'for' token
	Condition Expression
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Text }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for")
	out.WriteString(fs.Condition.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// --- Expressions ---

// Identifier represents an identifier (e.g., a variable name).
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Text }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Text }
func (il *IntegerLiteral) String() string       { return il.Token.Text }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Text }
func (sl *StringLiteral) String() string       { return sl.Token.Text }

// PrefixExpression represents a unary operation.
// e.g., -15
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Text }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents a binary operation.
// e.g., 5 + 5
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Text }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// FunctionLiteral represents a function definition.
// e.g., fn(x, y) { x + y; }
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Text }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents a function call.
// e.g., add(2, 3)
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Text }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
