package ast

import "github.com/jeremi-traverse/monkey/token"

type Node interface {
	TokenLiteral() string
}

// A Statement doesn't produce a value
type Statement interface {
	Node
	statementNode()
}

// An expression produce a value
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// Implements Statement interface
type Identifier struct {
	Token token.Token
	Value string
}

// Implements Statement interface
type LetStatement struct {
	Token token.Token
	Name  *Identifier // Identifier of the binding (left part)
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (i *Identifier) statementNode()        {}
func (ls *Identifier) TokenLiteral() string { return ls.Token.Literal }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
