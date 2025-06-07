package parse

import "github.com/codecrafters-io/interpreter-starter-go/app/pkg/scan"

type Parser struct {
	tokens  []scan.Token
	current int
}

func NewParser(tokens []scan.Token) Parser {
	return Parser{tokens: tokens}
}

func (v *Parser) Parse() Expression {
	return v.expression()
}

func (v *Parser) peek() scan.Token {
	return v.tokens[v.current]
}

func (v *Parser) previous() scan.Token {
	if v.current == 0 {
		panic("can't fetch previous")
	}
	return v.tokens[v.current-1]
}

func (v *Parser) isAtEnd() bool {
	return v.peek().TokenType == scan.EOF
}

func (v *Parser) advance() scan.Token {
	if !v.isAtEnd() {
		v.current += 1
	}
	return v.previous()
}

func (v *Parser) check(t scan.TokenType) bool {
	if v.isAtEnd() {
		return false
	}
	return v.peek().TokenType == t
}

func (v *Parser) match(types ...scan.TokenType) bool {
	for _, t := range types {
		if v.check(t) {
			v.advance()
			return true
		}
	}
	return false
}

func (v *Parser) expression() Expression {
	return v.equality()
}

func (v *Parser) equality() Expression {
	expr := v.comparison()

	for v.match(scan.BANG_EQUAL, scan.EQUAL_EQUAL) {
		operator := v.previous()
		right := v.comparison()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) comparison() Expression {
	expr := v.term()

	for v.match(scan.GREATER, scan.GREATER_EQUAL, scan.LESS, scan.LESS_EQUAL) {
		operator := v.previous()
		right := v.term()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) term() Expression {
	expr := v.factor()

	for v.match(scan.MINUS, scan.PLUS) {
		operator := v.previous()
		right := v.factor()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) factor() Expression {
	expr := v.unary()

	for v.match(scan.SLASH, scan.STAR) {
		operator := v.previous()
		right := v.unary()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) unary() Expression {
	if v.match(scan.BANG, scan.MINUS) {
		operator := v.previous()
		right := v.unary()
		return UnaryExpr{operator, right}
	}

	return v.primary()
}

func (v *Parser) primary() Expression {
	if v.match(scan.FALSE) {
		return LiteralExpr{scan.BooleanLiteral{false}}
	}
	if v.match(scan.TRUE) {
		return LiteralExpr{scan.BooleanLiteral{true}}
	}
	if v.match(scan.NUMBER, scan.STRING) {
		return LiteralExpr{v.previous().Literal}
	}
	if v.match(scan.NIL) {
		return LiteralExpr{scan.NullLiteral{}}
	}
	// if v.match(scan.LEFT_PAREN) { // todo: left to error part
	// 	expr := v.expression()

	// }
	return BinaryExpr{} // todo: what to return by default?
}

// func (v *Parser) consume(t scan.TokenType, message string) {

// }
