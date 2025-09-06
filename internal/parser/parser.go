package parser

import (
	"fmt"
	"strconv"

	"ixion/internal/ast"
	"ixion/internal/token"
)

// Operator precedence
const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:    ASSIGN,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.DIV:       PRODUCT,
	token.MUL:       PRODUCT,
	token.LPAREN:    CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	tokens *token.Tokens

	currentPos int
	curToken   token.Token
	peekToken  token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(tokens *token.Tokens) *Parser {
	p := &Parser{
		tokens: tokens,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER_LITERAL, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING_LITERAL, p.parseStringLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.FN, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t.String(), p.peekToken.Type.String())
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	if p.currentPos < len(*p.tokens) {
		p.peekToken = (*p.tokens)[p.currentPos]
	} else {
		p.peekToken = token.New(token.EOF, "")
	}
	p.currentPos++
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.CONST:
		// TODO: Implement parseConstStatement
		return nil
	case token.RETURN:
		return p.parseReturnStatement()
	case token.PRINT:
		return p.parsePrintStatement()
	case token.FN:
		// This could be a function declaration or a function literal assigned to a variable.
		// For now, assume it's a function declaration if followed by an identifier.
		if p.peekTokenIs(token.IDENT) {
			return p.parseFunctionDeclaration()
		}
		// If not a declaration, it will be handled as an expression statement
		fallthrough
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		// Error already added by expectPeek
		return stmt // Return partially constructed statement to avoid nil panic
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Text}

	// Optional type annotation
	if p.peekToken.IsType() {
		p.nextToken() // Advance to the type token
		stmt.Type = &ast.TypeLiteral{Token: p.curToken, Value: p.curToken.Text}
	}

	if !p.expectPeek(token.ASSIGN) {
		// Error already added by expectPeek
		return stmt // Return partially constructed statement to avoid nil panic
	}

	p.nextToken() // Advance past ASSIGN

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // Advance past RETURN

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return stmt
	}

	p.nextToken() // Advance past LPAREN

	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return stmt
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // Advance past LBRACE

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if !p.curTokenIs(token.RBRACE) {
		p.errors = append(p.errors, "expected '}' at end of block statement")
		return nil
	}

	return block
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	fnDecl := &ast.FunctionDeclaration{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	fnDecl.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Text}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fnDecl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Optional return type
	if p.peekToken.IsType() {
		p.nextToken() // Advance to the type token
		fnDecl.ReturnType = &ast.TypeLiteral{Token: p.curToken, Value: p.curToken.Text}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fnDecl.Body = p.parseBlockStatement()

	return fnDecl
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
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
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Text}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Text, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Text)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Text}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Text,
	}

	p.nextToken() // Advance past the operator

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Text,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken() // Advance past the operator

	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // Advance past LPAREN

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseFunctionParameters() []*ast.FunctionParameter {
	parameters := []*ast.FunctionParameter{}

	if p.peekTokenIs(token.RPAREN) {
		return parameters
	}

	p.nextToken() // Advance past LPAREN or COMMA

	param := &ast.FunctionParameter{Token: p.curToken}
	if !p.curTokenIs(token.IDENT) {
		p.errors = append(p.errors, "expected identifier for function parameter")
		return nil
	}
	param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Text}

	// Optional type annotation
	if p.peekToken.IsType() {
		p.nextToken() // Advance to the type token
		param.Type = &ast.TypeLiteral{Token: p.curToken, Value: p.curToken.Text}
	}
	parameters = append(parameters, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Advance past COMMA
		p.nextToken() // Advance to the next parameter's identifier

		param := &ast.FunctionParameter{Token: p.curToken}
		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, "expected identifier for function parameter")
			return nil
		}
		param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Text}

		// Optional type annotation
		if p.peekToken.IsType() {
			p.nextToken() // Advance to the type token
			param.Type = &ast.TypeLiteral{Token: p.curToken, Value: p.curToken.Text}
		}
		parameters = append(parameters, param)
	}

	return parameters
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Optional return type
	if p.peekToken.IsType() {
		p.nextToken() // Advance to the type token
		lit.ReturnType = &ast.TypeLiteral{Token: p.curToken, Value: p.curToken.Text}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() // Advance past RPAREN
		return args
	}

	p.nextToken() // Advance past LPAREN or COMMA

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Advance past COMMA
		p.nextToken() // Advance to the next argument
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	exp := &ast.AssignmentExpression{Token: p.curToken, Left: left}

	precedence := p.curPrecedence()
	p.nextToken() // Advance past ASSIGN

	exp.Value = p.parseExpression(precedence)

	return exp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t.String())
	p.errors = append(p.errors, msg)
}
