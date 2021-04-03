package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/tokens"
)

func TestNextToken_Simple(t *testing.T) {
	input := "+=;)(}{,."
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.PLUS, "+"},
		{tokens.ASSIGN, "="},
		{tokens.SEMICOLON, ";"},
		{tokens.RPAREN, ")"},
		{tokens.LPAREN, "("},
		{tokens.RBRACE, "}"},
		{tokens.LBRACE, "{"},
		{tokens.COMMA, ","},
		{tokens.DOT, "."},
		{tokens.EOF, ""},
	}

	l := lexer.New(input, "non-file")

	for _, v := range tests {
		next := l.NextToken()
		assert.Equal(t, v.expectedLiteral, next.Literal)
		assert.Equal(t, v.expectedType, next.Type)
	}
}

func TestNextToken_Program(t *testing.T) {
	input := `
	# sys "sys";

	let five = 5;
	let ten = 10;

	let add = fn(a, b) {
		a + b;
	};

	let result = add(five, ten);
	"foobar";
	"foo bar";
	"hello \"world\"";
	"hello \n\t \"world\"";
	[1, 2];
	{"foo":"bar", true: 195};
	for (a < b) {};

	let tenFifteen = 10.15;
`
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.HASH, "#"},
		{tokens.IDENT, "sys"},
		{tokens.STRING, "sys"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "five"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "5"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "ten"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "add"},
		{tokens.ASSIGN, "="},
		{tokens.FUNCTION, "fn"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "a"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "b"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.IDENT, "a"},
		{tokens.PLUS, "+"},
		{tokens.IDENT, "b"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "result"},
		{tokens.ASSIGN, "="},
		{tokens.IDENT, "add"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "five"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "ten"},
		{tokens.RPAREN, ")"},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, "foobar"},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, "foo bar"},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, `hello "world"`},
		{tokens.SEMICOLON, ";"},
		{tokens.STRING, `hello 
	 "world"`},
		{tokens.SEMICOLON, ";"},
		{tokens.LBRACKET, "["},
		{tokens.INT, "1"},
		{tokens.COMMA, ","},
		{tokens.INT, "2"},
		{tokens.RBRACKET, "]"},
		{tokens.SEMICOLON, ";"},
		{tokens.LBRACE, "{"},
		{tokens.STRING, "foo"},
		{tokens.COLON, ":"},
		{tokens.STRING, "bar"},
		{tokens.COMMA, ","},
		{tokens.TRUE, "true"},
		{tokens.COLON, ":"},
		{tokens.INT, "195"},
		{tokens.RBRACE, "}"},
		{tokens.SEMICOLON, ";"},
		//for (a < b) {};
		{tokens.FOR, "for"},
		{tokens.LPAREN, "("},
		{tokens.IDENT, "a"},
		{tokens.LT, "<"},
		{tokens.IDENT, "b"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RBRACE, "}"},
		{tokens.SEMICOLON, ";"},
		//let tenFifteen = 10.15;
		{tokens.LET, "let"},
		{tokens.IDENT, "tenFifteen"},
		{tokens.ASSIGN, "="},
		{tokens.DOUBLE, "10.15"},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, ""},
	}

	l := lexer.New(input, "non-file")

	for _, v := range tests {
		next := l.NextToken()
		assert.Equal(t, v.expectedLiteral, next.Literal)
		require.Equal(t, v.expectedType, next.Type)
	}
}

func TestNextToken_LogicOperators(t *testing.T) {
	input := `
	!*/-5
	10 < 20 > 5
`
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.BANG, "!"},
		{tokens.ASTERISK, "*"},
		{tokens.SLASH, "/"},
		{tokens.MINUS, "-"},
		{tokens.INT, "5"},
		{tokens.INT, "10"},
		{tokens.LT, "<"},
		{tokens.INT, "20"},
		{tokens.GT, ">"},
		{tokens.INT, "5"},
		{tokens.EOF, ""},
	}

	l := lexer.New(input, "non-file")

	for _, v := range tests {
		next := l.NextToken()
		assert.Equal(t, v.expectedLiteral, next.Literal)
		assert.Equal(t, v.expectedType, next.Type)
	}
}

func TestNextToken_Bool(t *testing.T) {
	input := `
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
`
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.IF, "if"},
		{tokens.LPAREN, "("},
		{tokens.INT, "5"},
		{tokens.LT, "<"},
		{tokens.INT, "10"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.TRUE, "true"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.ELSE, "else"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.FALSE, "false"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.EOF, ""},
	}

	l := lexer.New(input, "non-file")

	for _, v := range tests {
		next := l.NextToken()
		assert.Equal(t, v.expectedLiteral, next.Literal)
		assert.Equal(t, v.expectedType, next.Type)
	}
}

func TestNextToken_Equal_NotEqual(t *testing.T) {
	input := `
	10 == 10;
	9 != 10;
`
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.INT, "10"},
		{tokens.EQ, "=="},
		{tokens.INT, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.INT, "9"},
		{tokens.NOT_EQ, "!="},
		{tokens.INT, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, ""},
	}

	l := lexer.New(input, "non-file")

	for _, v := range tests {
		next := l.NextToken()
		assert.Equal(t, v.expectedLiteral, next.Literal)
		assert.Equal(t, v.expectedType, next.Type)
	}
}

func TestDottedIdentifiers(t *testing.T) {
	input := `
			let a = sys.PI;
			sys.time();
`
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.LET, "let"},
		{tokens.IDENT, "a"},
		{tokens.ASSIGN, "="},
		{tokens.IDENT, "sys"},
		{tokens.DOT, "."},
		{tokens.IDENT, "PI"},
		{tokens.SEMICOLON, ";"},
		{tokens.IDENT, "sys"},
		{tokens.DOT, "."},
		{tokens.IDENT, "time"},
		{tokens.LPAREN, "("},
		{tokens.RPAREN, ")"},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, ""},
	}

	l := lexer.New(input, "non-file")

	for _, v := range tests {
		next := l.NextToken()
		assert.Equal(t, v.expectedLiteral, next.Literal)
		assert.Equal(t, v.expectedType, next.Type)
	}

}
