package parser

import (
	"errors"
	"log"

	ls "github.com/Piyush01Bhatt/interpreter_go/internal/scanner"
)

// Grammar to parse
// program        → declaration* EOF
// declaration    → classDecl | funDecl | varDecl | statement
// classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}"
// funDecl        → "fun" function
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";"
// statement      → exprStmt | ifStmt | printStmt | returnStmt | whileStmt | block
// exprStmt       → expression ";"
// expression     → equality
// equality       → comparison ( ( "!=" | "==" ) comparison )*
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
// term           → factor ( ( "-" | "+" ) factor )*
// factor         → unary ( ( "/" | "*" ) unary )*
// unary          → ( "!" | "-" ) unary
//                | primary
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")"
//                | IDENTIFIER

type Parser struct {
	tokens  []ls.Token
	current int
}

func NewParser(tokens []ls.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	var stmts []Stmt
	for !p.isAtEnd() {
		stmt := p.declaration()
		stmts = append(stmts, stmt)
	}
	return stmts
}

func (p *Parser) declaration() Stmt {
	if p.match(ls.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	if p.match(ls.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) varDeclaration() Stmt {
	name, err := p.consume(ls.IDENTIFIER, "expect variable name")
	if err != nil {
		log.Fatal(err)
	}

	var initializer Expr

	if p.match(ls.EQUAL) {
		initializer = p.ParseExpression()
	}

	_, err = p.consume(ls.SEMICOLON, "expect ';' after expression")
	if err != nil {
		log.Fatal(err)
	}

	return &VarStmt{
		Name: &name,
		Expr: initializer,
	}
}

func (p *Parser) printStatement() Stmt {
	expr := p.ParseExpression()
	_, err := p.consume(ls.SEMICOLON, "expect ';' after expression")
	if err != nil {
		log.Fatal(err)
	}
	return &PrintStmt{
		Expr: expr,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.ParseExpression()
	_, err := p.consume(ls.SEMICOLON, "expect ';' after expression")
	if err != nil {
		log.Fatal(err)
	}
	return &ExpressionStmt{
		Expr: expr,
	}
}

func (p *Parser) ParseExpression() Expr {
	return p.expression()
}

// expression -> equality

func (p *Parser) expression() Expr {
	return p.equality()
}

// equality  → comparison ( ( "!=" | "==" ) comparison )*

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(ls.BANG_EQUAL, ls.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{
			Left:     expr,
			Operator: &operator,
			Right:    right,
		}
	}
	return expr
}

// comparison  → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(ls.GREATER, ls.GREATER_EQUAL, ls.LESS, ls.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{
			Left:     expr,
			Operator: &operator,
			Right:    right,
		}
	}
	return expr
}

// term  → factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(ls.MINUS, ls.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{
			Left:     expr,
			Operator: &operator,
			Right:    right,
		}
	}
	return expr
}

// factor → unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(ls.SLASH, ls.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{
			Left:     expr,
			Operator: &operator,
			Right:    right,
		}
	}
	return expr
}

// unary          → ( "!" | "-" ) unary
//
//	| primary
func (p *Parser) unary() Expr {
	if p.match(ls.BANG, ls.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &Unary{
			Operator: &operator,
			Right:    right,
		}
	}
	return p.primary()
}

// primary  → NUMBER | STRING | "true" | "false" | "nil"
//
//	| "(" expression ")"
func (p *Parser) primary() Expr {
	if p.match(ls.NUMBER, ls.STRING, ls.TRUE, ls.FALSE, ls.NIL) {
		token := p.previous()
		literal := token.Literal
		var value *Value

		switch token.Type {
		case ls.NUMBER:
			value = NewFloatValue(literal.(float64))
		case ls.STRING:
			value = NewStringValue(literal.(string))
		case ls.TRUE, ls.FALSE:
			value = NewBoolValue(literal.(bool))
		default:
			value = nil
		}
		return &Literal{
			Value: value,
		}
	}

	if p.match(ls.IDENTIFIER) {
		return &Variable{
			Name: p.previous().Lexeme,
		}
	}

	var expr Expr

	if p.match(ls.LEFT_PAREN) {
		expr = p.expression()
		_, err := p.consume(ls.RIGHT_PAREN, "expect ')' after expression")
		if err != nil {
			log.Fatal(err)
		}
	}

	return expr
}

// utilities
// match for tokens
func (p *Parser) match(tokens ...ls.TokenType) bool {
	for _, tokenType := range tokens {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() ls.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(tokenType ls.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ls.EOF
}

func (p *Parser) peek() ls.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() ls.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType ls.TokenType, message string) (ls.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return ls.Token{}, errors.New(message)
}
