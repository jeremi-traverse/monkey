package parser

import (
	"testing"

	"github.com/jeremi-traverse/monkey/ast"
	"github.com/jeremi-traverse/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("progam.Statements does not contain 3 statments. got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, expectedName string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("Statement isn't a LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != expectedName {
		t.Errorf("letStmt.Name not %q. got=%q", expectedName, letStmt.Name)
		return false
	}

	if letStmt.Name.TokenLiteral() != expectedName {
		t.Errorf("letStmt.Name.TokenLiteral() not %q. got=%q",
			expectedName,
			letStmt.Name.TokenLiteral())
	}

	return true
}
