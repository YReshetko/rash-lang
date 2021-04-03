package parser

import (
	"fmt"
	"github.com/YReshetko/rash-lang/ast"
	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/tokens"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currToken tokens.Token
	peekToken tokens.Token

	errors []string

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: map[tokens.TokenType]prefixParseFn{},
		infixParseFns:  map[tokens.TokenType]infixParseFn{},
	}

	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
	p.registerPrefix(tokens.INT, p.parseIntegerLiteral)
	p.registerPrefix(tokens.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.BANG, p.parsePrefixExpression)
	p.registerPrefix(tokens.MINUS, p.parsePrefixExpression)
	p.registerPrefix(tokens.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(tokens.IF, p.parseIfExpression)
	p.registerPrefix(tokens.FOR, p.parseForExpression)
	p.registerPrefix(tokens.LET, p.parseLetExpression)
	p.registerPrefix(tokens.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(tokens.STRING, p.parseStringLiteral)
	p.registerPrefix(tokens.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(tokens.LBRACE, p.parseHashLiteral)

	p.registerInfix(tokens.EQ, p.parseInfixExpression)
	p.registerInfix(tokens.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(tokens.LT, p.parseInfixExpression)
	p.registerInfix(tokens.GT, p.parseInfixExpression)
	p.registerInfix(tokens.PLUS, p.parseInfixExpression)
	p.registerInfix(tokens.MINUS, p.parseInfixExpression)
	p.registerInfix(tokens.SLASH, p.parseInfixExpression)
	p.registerInfix(tokens.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tokens.LPAREN, p.parseCallExpression)
	p.registerInfix(tokens.DOT, p.parseInfixExpression)
	p.registerInfix(tokens.LBRACKET, p.parseInfixIndexExpression)
	p.registerInfix(tokens.ASSIGN, p.parseInfixExpression)

	// Call twice to set current and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(tokens.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.HASH:
		return p.parseIncludeDeclarationStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() ast.Statement {
	defer untrace(trace("parseLetStatement"))
	statement := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeekToken(tokens.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeekToken(tokens.ASSIGN) {
		return nil
	}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseLetExpression() ast.Expression {
	return p.parseLetStatement().(*ast.LetStatement)
}

func (p *Parser) parseReturnStatement() ast.Statement {
	defer untrace(trace("parseReturnStatement"))
	statement := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()

	if p.currTokenIs(tokens.SEMICOLON) {
		return statement
	}

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseIncludeDeclarationStatement() ast.Statement {
	defer untrace(trace("parseReturnStatement"))
	statement := &ast.DeclarationStatement{
		Token: p.currToken,
	}
	if !p.expectPeekToken(tokens.IDENT) {
		return nil
	}

	alias := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeekToken(tokens.STRING) {
		return nil
	}

	include := &ast.StringLiteral{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	statement.Declaration = &ast.IncludeDeclaration{
		Token:   statement.Token,
		Alias:   alias,
		Include: include,
	}
	return statement
}

const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	EQUAL       // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	DOT         // .
	CALL        // myFunc(x)
	INDEX       //array[index]
)

var precedences = map[tokens.TokenType]int{
	tokens.ASSIGN:   ASSIGN,
	tokens.EQ:       EQUAL,
	tokens.NOT_EQ:   EQUAL,
	tokens.LT:       LESSGREATER,
	tokens.GT:       LESSGREATER,
	tokens.PLUS:     SUM,
	tokens.MINUS:    SUM,
	tokens.SLASH:    PRODUCT,
	tokens.ASTERISK: PRODUCT,
	tokens.LPAREN:   CALL,
	tokens.DOT:      DOT,
	tokens.LBRACKET: INDEX,
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	defer untrace(trace("parseExpressionStatement"))
	statement := &ast.ExpressionStatement{Token: p.currToken}
	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(tokens.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)

	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(lit.TokenLiteral(), 0, 64)
	if err != nil {
		message := fmt.Sprintf("expected integer literal on line %d instead of %s", p.currToken.LineNumber, p.currToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	defer untrace(trace("parseStringLiteral"))
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	defer untrace(trace("parseBooleanLiteral"))
	return &ast.BooleanLiteral{Token: p.currToken, Value: p.currTokenIs(tokens.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	exp := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	exp := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	prec := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(prec)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	defer untrace(trace("parseGroupedExpression"))
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeekToken(tokens.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	defer untrace(trace("parseIfExpression"))
	exp := &ast.IfExpression{
		Token: p.currToken,
	}

	if !p.expectPeekToken(tokens.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeekToken(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeekToken(tokens.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if !p.peekTokenIs(tokens.ELSE) {
		return exp
	}
	p.nextToken()
	if !p.expectPeekToken(tokens.LBRACE) {
		return nil
	}

	exp.Alternative = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseForExpression() ast.Expression {
	defer untrace(trace("parseForExpression"))
	exp := &ast.ForExpression{
		Token: p.currToken,
	}

	if !p.expectPeekToken(tokens.LPAREN) {
		return nil
	}
	// Parse for params expressions
	args := p.parseForArguments()

	if !p.expectPeekToken(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeekToken(tokens.LBRACE) {
		return nil
	}

	switch len(args) {
	case 1:
		exp.Condition = args[0]
	case 2:
		exp.Condition = args[0]
		exp.Complete = args[1]
	case 3:
		exp.Initial = args[0]
		exp.Condition = args[1]
		exp.Complete = args[2]
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseForArguments() []ast.Expression {
	var args []ast.Expression
	if p.peekTokenIs(tokens.RPAREN) {
		return args
	}
	p.nextToken()

	args = append(args, p.parseForArgument())
	if p.peekTokenIs(tokens.RPAREN) {
		return args
	}
	p.nextToken()

	args = append(args, p.parseForArgument())
	if p.peekTokenIs(tokens.RPAREN) {
		return args
	}
	p.nextToken()

	args = append(args, p.parseForArgument())
	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return args
}

func (p *Parser) parseForArgument() ast.Expression {
	arg := p.parseExpression(LOWEST)
	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}
	return arg
}


func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	defer untrace(trace("parseBlockStatement"))
	block := &ast.BlockStatement{
		Token:      p.currToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for !p.currTokenIs(tokens.RBRACE) && !p.currTokenIs(tokens.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	defer untrace(trace("parseFunctionLiteral"))
	fnLit := &ast.FunctionLiteral{
		Token: p.currToken,
	}

	if !p.expectPeekToken(tokens.LPAREN) {
		return nil
	}

	fnLit.Parameters = p.parseFunctionParameters()

	if !p.expectPeekToken(tokens.LBRACE) {
		return nil
	}

	fnLit.Body = p.parseBlockStatement()

	return fnLit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	defer untrace(trace("parseFunctionParameters"))
	idents := []*ast.Identifier{}
	p.nextToken()

	if p.currTokenIs(tokens.RPAREN) {
		return idents
	}

	ident := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
	idents = append(idents, ident)

	for p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
		p.nextToken()
		ident = &ast.Identifier{
			Token: p.currToken,
			Value: p.currToken.Literal,
		}
		idents = append(idents, ident)
	}

	if !p.expectPeekToken(tokens.RPAREN) {
		return nil
	}

	return idents
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	defer untrace(trace("parseCallExpression"))
	call := &ast.CallExpression{
		Token:    p.currToken,
		Function: function,
	}

	call.Arguments = p.parseExpressionList(tokens.RPAREN)
	return call
}

func (p *Parser) parseExpressionList(end tokens.TokenType) []ast.Expression {
	defer untrace(trace("parseExpressionList"))
	expressions := []ast.Expression{}
	p.nextToken()
	if p.currTokenIs(end) {
		return expressions
	}

	expressions = append(expressions, p.parseExpression(LOWEST))

	for p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
		p.nextToken()
		expressions = append(expressions, p.parseExpression(LOWEST))
	}

	if !p.expectPeekToken(end) {
		return nil
	}
	return expressions
}

func (p *Parser) parseInfixIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.currToken,
		Left:  left,
	}
	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeekToken(tokens.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.currToken,
		Elements: p.parseExpressionList(tokens.RBRACKET),
	}
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token: p.currToken,
		Pairs: map[ast.Expression]ast.Expression{},
	}
	for !p.peekTokenIs(tokens.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeekToken(tokens.COLON) {
			return nil
		}
		p.nextToken()
		hash.Pairs[key] = p.parseExpression(LOWEST)
		if !p.peekTokenIs(tokens.RBRACE) && !p.expectPeekToken(tokens.COMMA) {
			return nil
		}
	}
	if !p.expectPeekToken(tokens.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) currTokenIs(t tokens.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeekToken(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	prec, ok := precedences[p.peekToken.Type]
	if ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	prec, ok := precedences[p.currToken.Type]
	if ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected token %s on line %d; instead got %s", t, p.peekToken.LineNumber, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t tokens.TokenType) {
	msg := fmt.Sprintf("no prefix parse functions found for %s on line %d", t, p.currToken.LineNumber)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(t tokens.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t tokens.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}
