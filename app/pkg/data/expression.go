package data

type Expression interface {
	Accept(visitor Visitor)
}

type BinaryExpr struct {
	Left     Expression
	Operator Token
	Right    Expression
}

func (v BinaryExpr) Accept(visitor Visitor) {
	visitor.VisitBinary(v)
}

type GroupingExpr struct {
	Expression Expression
}

func (v GroupingExpr) Accept(visitor Visitor) {
	visitor.VisitGrouping(v)
}

type LiteralExpr struct {
	Value Literal
}

func (v LiteralExpr) Accept(visitor Visitor) {
	visitor.VisitLiteral(v)
}

type UnaryExpr struct {
	Operator Token
	Right    Expression
}

func (v UnaryExpr) Accept(visitor Visitor) {
	visitor.VisitUnary(v)
}
