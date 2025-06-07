package data

type Visitor interface {
	VisitBinary(v BinaryExpr)
	VisitGrouping(v GroupingExpr)
	VisitLiteral(v LiteralExpr)
	VisitUnary(v UnaryExpr)
}
