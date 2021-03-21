package ast_test

import (
	"github.com/YReshetko/rash-lang/ast"
	"github.com/YReshetko/rash-lang/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: tokens.Token{Type: tokens.LET, Literal: "let"},
				Name:  &ast.Identifier{
					Token: tokens.Token{Type: tokens.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: tokens.Token{Type: tokens.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	assert.Equal(t, "let myVar = anotherVar;", program.String())
}
