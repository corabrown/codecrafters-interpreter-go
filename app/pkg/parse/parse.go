package parse

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
)

type Parser struct {
	tokens  []data.Token
	current int
	errors  []*errors.Error
}

func NewParser(tokens []data.Token) Parser {
	return Parser{tokens: tokens}
}

func (v *Parser) Parse() (data.Expression, []*errors.Error) {
	return v.expression(), v.errors
}

func (v *Parser) ParseError() bool {
	return len(v.errors) > 0
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

func (v *Parser) consume(t data.TokenType, message string) (data.Token, *errors.Error) {
	if v.check(t) {
		return v.advance(), nil
	}

	err := v.error(v.peek(), message)
	return data.Token{}, err
}

func (v *Parser) expression() data.Expression {
	return v.equality()
}

func (v *Parser) equality() data.Expression {
	expr := v.comparison()

	for v.match(data.BANG_EQUAL, data.EQUAL_EQUAL) {
		operator := v.previous()
		right := v.comparison()
		if right == nil {
			v.errors = append(v.errors, v.error(v.peek(), "need non-null"))
		}
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) comparison() data.Expression {
	expr := v.term()

	for v.match(data.GREATER, data.GREATER_EQUAL, data.LESS, data.LESS_EQUAL) {
		operator := v.previous()
		right := v.term()
		if right == nil {
			v.errors = append(v.errors, v.error(v.peek(), "need non-null"))
		}
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) term() data.Expression {
	expr := v.factor()

	for v.match(data.MINUS, data.PLUS) {
		operator := v.previous()
		right := v.factor()
		if right == nil {
			v.errors = append(v.errors, v.error(v.peek(), "need non-null"))
		}
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) factor() data.Expression {
	expr := v.unary()

	for v.match(data.SLASH, data.STAR) {
		operator := v.previous()
		right := v.unary()
		if right == nil {
			v.errors = append(v.errors, v.error(v.peek(), "need non-null"))
		}
		expr = data.BinaryExpr{expr, operator, right}
	}

	return expr
}

func (v *Parser) unary() data.Expression {
	if v.match(data.BANG, data.MINUS) {
		operator := v.previous()
		right := v.unary()
		if right == nil {
			v.errors = append(v.errors, v.error(v.peek(), "need non-null"))
		}
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
		if _, err := v.consume(data.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			v.syncronize()
		}
		return data.GroupingExpr{expr}
	}
	return nil // todo: what to return by default?
}

func (v *Parser) error(token data.Token, message string) (err *errors.Error) {

	defer func() {
		if err != nil {
			v.errors = append(v.errors, err)
		}
	}()

	if token.TokenType == data.EOF {
		err = errors.NewError(token.Line, message, " at end")
		return
	}
	err = errors.NewError(token.Line, message, fmt.Sprintf(" at '%v'", token.Lexeme))
	return
}

func (v *Parser) syncronize() {
	v.advance()

	for !v.isAtEnd() {
		if v.previous().TokenType == data.SEMICOLON {
			return
		}
	}
	switch v.peek().TokenType {
	case data.CLASS, data.FUN, data.VAR, data.FOR, data.IF, data.WHILE, data.PRINT, data.RETURN:
		return
	}

	v.advance()
}
