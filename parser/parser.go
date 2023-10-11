package parser

import (
	"fmt"

	"github.com/jeremi-traverse/monkey/ast"
	"github.com/jeremi-traverse/monkey/lexer"
	"github.com/jeremi-traverse/monkey/token"
)

type Parser struct {
	l            *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	// Declares a new Program
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatment()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// Get the next statement
		p.nextToken()
	}
	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// let y = 5;
// -> LetStatement
//
//	Identifier(y)
//	Expression(5)
func (p *Parser) parseStatment() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {

		return nil
	}

	// Skipping until we get a semicolon token
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	fmt.Printf("last type %s\n", p.currentToken.Type)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	/*
		if !p.expectPeek(token.IDENT) {
			return nil
		}
	*/
	// Skipping until we get a semicolon token
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}

	p.addPeekError(t)
	return false
}

func (p *Parser) addPeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected token type %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
