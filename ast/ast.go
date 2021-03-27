package ast

import (
	"bytes"
	"fmt"
	"github.com/YReshetko/rash-lang/tokens"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
	StackLine() string
}

type Statement interface {
	Node
	statementNode() // We need the function to define type unambiguously, there is no needed any real implementations
}

type Expression interface {
	Node
	expressionNode() // We need the function to define type unambiguously, there is no needed any real implementations
}

type Declaration interface {
	Node
	declarationNode() //  We need the function to define type unambiguously, there is no needed any real implementations
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}
func (p *Program) String() string {
	out := bytes.Buffer{}
	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}
func (p *Program) StackLine() string {
	return ""
}

type LetStatement struct {
	Token tokens.Token // LET token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) statementNode()       {}
func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatement) String() string {
	out := bytes.Buffer{}
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (l *LetStatement) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", l.Token.FileName, l.Token.LineNumber)
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode() ()   {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	if i == nil {
		return ""
	}
	return i.Value
}
func (i *Identifier) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", i.Token.FileName, i.Token.LineNumber)
}

type ReturnStatement struct {
	Token tokens.Token // RETURN token
	Value Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(r.TokenLiteral() + " ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (r *ReturnStatement) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", r.Token.FileName, r.Token.LineNumber)
}

type ExpressionStatement struct {
	Token      tokens.Token // The first token in the expression
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}
func (e *ExpressionStatement) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", e.Token.FileName, e.Token.LineNumber)
}

type DeclarationStatement struct {
	Token       tokens.Token // The first token in the expression
	Declaration Declaration
}

func (d *DeclarationStatement) statementNode()       {}
func (d *DeclarationStatement) TokenLiteral() string { return d.Token.Literal }
func (d *DeclarationStatement) String() string {
	if d.Declaration != nil {
		return d.Declaration.String()
	}
	return ""
}
func (d *DeclarationStatement) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", d.Token.FileName, d.Token.LineNumber)
}

type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }
func (i *IntegerLiteral) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", i.Token.FileName, i.Token.LineNumber)
}

type StringLiteral struct {
	Token tokens.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }
func (s *StringLiteral) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", s.Token.FileName, s.Token.LineNumber)
}

type BooleanLiteral struct {
	Token tokens.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }
func (b *BooleanLiteral) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", b.Token.FileName, b.Token.LineNumber)
}

type ArrayLiteral struct {
	Token    tokens.Token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	out := bytes.Buffer{}
	elems := make([]string, len(a.Elements))
	for i, elem := range a.Elements {
		elems[i] = elem.String()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}
func (a *ArrayLiteral) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", a.Token.FileName, a.Token.LineNumber)
}

type HashLiteral struct {
	Token    tokens.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	out := bytes.Buffer{}
	pairs := make([]string, len(h.Pairs))
	i := 0
	for k, v := range h.Pairs {
		pairs[i] = k.String() + ":"+v.String()
		i++
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
func (h *HashLiteral) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", h.Token.FileName, h.Token.LineNumber)
}

type PrefixExpression struct {
	Token    tokens.Token // Prefix operator, eg. !
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	out := bytes.Buffer{}

	out.WriteString("(")
	out.WriteString(p.Operator)
	if p.Right != nil {
		out.WriteString(p.Right.String())
	}
	out.WriteString(")")

	return out.String()
}
func (p *PrefixExpression) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", p.Token.FileName, p.Token.LineNumber)
}

type InfixExpression struct {
	Token    tokens.Token // Operator, eg. +
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	out := bytes.Buffer{}

	out.WriteString("(")
	if i.Left != nil {
		out.WriteString(i.Left.String())
	}
	if i.Operator == "." {
		out.WriteString(i.Operator)
	} else {
		out.WriteString(" " + i.Operator + " ")
	}

	if i.Right != nil {
		out.WriteString(i.Right.String())
	}
	out.WriteString(")")

	return out.String()
}
func (i *InfixExpression) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", i.Token.FileName, i.Token.LineNumber)
}

type IfExpression struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	out := bytes.Buffer{}

	out.WriteString("if ")
	out.WriteString(i.Condition.String() + " ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}
func (i *IfExpression) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", i.Token.FileName, i.Token.LineNumber)
}

type BlockStatement struct {
	Token      tokens.Token // start block literal -> {
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	out := bytes.Buffer{}
	for _, v := range b.Statements {
		out.WriteString(v.String())
	}
	return out.String()
}
func (b *BlockStatement) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", b.Token.FileName, b.Token.LineNumber)
}

type FunctionLiteral struct {
	Token      tokens.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionLiteral) String() string {
	out := bytes.Buffer{}

	out.WriteString(f.Token.Literal + "(")

	params := make([]string, len(f.Parameters))
	for i, parameter := range f.Parameters {
		params[i] = parameter.Value
	}

	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}
func (f *FunctionLiteral) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", f.Token.FileName, f.Token.LineNumber)
}

type CallExpression struct {
	Token     tokens.Token // Token for (
	Function  Expression   // Identifier or Function literal
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	out := bytes.Buffer{}

	out.WriteString(c.Function.String())
	out.WriteString("(")

	args := make([]string, len(c.Arguments))
	for i, argument := range c.Arguments {
		args[i] = argument.String()
	}

	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
func (c *CallExpression) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", c.Token.FileName, c.Token.LineNumber)
}

type IncludeDeclaration struct {
	Token   tokens.Token
	Alias   *Identifier
	Include *StringLiteral
}

func (c *IncludeDeclaration) declarationNode()     {}
func (c *IncludeDeclaration) TokenLiteral() string { return c.Token.Literal }
func (c *IncludeDeclaration) String() string {
	out := bytes.Buffer{}

	out.WriteString("# ")
	if c.Alias != nil {
		out.WriteString(c.Alias.Value + " ")
	}
	out.WriteString("\"")
	out.WriteString(c.Include.Value)
	out.WriteString("\";")

	return out.String()
}
func (c *IncludeDeclaration) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", c.Token.FileName, c.Token.LineNumber)
}

type IndexExpression struct {
	Token tokens.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
	out := bytes.Buffer{}

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("]")
	out.WriteString(")")

	return out.String()
}
func (i *IndexExpression) StackLine() string {
	return fmt.Sprintf("file: %s; line: %d", i.Token.FileName, i.Token.LineNumber)
}
