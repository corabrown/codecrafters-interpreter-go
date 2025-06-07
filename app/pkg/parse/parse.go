package parse

import "github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"

type Parser struct {
	tokens  []data.Token
	current int
}

func NewParser(tokens []data.Token) Parser {
	return Parser{tokens: tokens}
}

func (v *Parser) Parse() data.Expression {
	return v.expression()
}

func (v *Parser) peek() data.Token {
	return v.tokens[v.current]
}

func (v *Parser) previous() data.Token {
	if v.current == 0 {
		panic("can't fetch previous")
	}
	return v.tokens[v.current-1]
}

func (v *Parser) isAtEnd() bool {
	return v.peek().TokenType == data.EOF
}

func (v *Parser) advance() data.Token {
	if !v.isAtEnd() {
		v.current += 1
	}
	return v.previous()
}

func (v *Parser) check(t data.TokenType) bool {
	if v.isAtEnd() {
		return false
	}
	return v.peek().TokenType == t
}

func (v *Parser) match(types ...data.TokenType) bool {
	for _, t := range types {
		if v.check(t) {
			v.advance()
			return true
		}
	}
	return false
}

func (v *Parser) consume(t data.TokenType, message string) data.Token {
	if v.check(t) {
		return v.advance()
	}

	return data.Token{}
}

func (v *Parser) expression() data.Expression {
	return v.equality()
}

func (v *Parser) equality() data.Expression {
	expr := v.comparison()

	for v.match(data.BANG_EQUAL, data.EQUAL_EQUAL) {
		operator := v.previous()
		right := v.comparison()
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) comparison() data.Expression {
	expr := v.term()

	for v.match(data.GREATER, data.GREATER_EQUAL, data.LESS, data.LESS_EQUAL) {
		operator := v.previous()
		right := v.term()
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) term() data.Expression {
	expr := v.factor()

	for v.match(data.MINUS, data.PLUS) {
		operator := v.previous()
		right := v.factor()
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) factor() data.Expression {
	expr := v.unary()

	for v.match(data.SLASH, data.STAR) {
		operator := v.previous()
		right := v.unary()
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) unary() data.Expression {
	if v.match(data.BANG, data.MINUS) {
		operator := v.previous()
		right := v.unary()
		return data.UnaryExpr{operator, right}
	}

	return v.primary()
}

func (v *Parser) primary() data.Expression {
	if v.match(data.FALSE) {
		return data.LiteralExpr{data.BooleanLiteral{false}}
	}
	if v.match(data.TRUE) {
		return data.LiteralExpr{data.BooleanLiteral{true}}
	}
	if v.match(data.NUMBER, data.STRING) {
		return data.LiteralExpr{v.previous().Literal}
	}
	if v.match(data.NIL) {
		return data.LiteralExpr{data.NullLiteral{}}
	}
	if v.match(data.LEFT_PAREN) { // todo: left to error part
		expr := v.expression()
		v.consume(data.RIGHT_PAREN, "message")
		return data.GroupingExpr{expr}
	}
	return data.BinaryExpr{} // todo: what to return by default?
}
