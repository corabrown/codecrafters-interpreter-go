package parse

import "github.com/codecrafters-io/interpreter-starter-go/app/pkg/scan"

type Expression interface {
	accept(visitor Visitor)
}

type BinaryExpr struct {
	left     Expression
	operator scan.Token
	right    Expression
}

func (v BinaryExpr) accept(visitor Visitor) {
	visitor.visitBinary(v)
}

type GroupingExpr struct {
	expression Expression
}

func (v GroupingExpr) accept(visitor Visitor) {
	visitor.visitGrouping(v)
}

type LiteralExpr struct {
	value scan.Literal
}

func (v LiteralExpr) accept(visitor Visitor) {
	visitor.visitLiteral(v)
}

type UnaryExpr struct {
	operator scan.Token
	right Expression 
}

func (v UnaryExpr) accept(visitor Visitor) {
	visitor.visitUnary(v)
}
