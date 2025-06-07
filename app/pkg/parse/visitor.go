package parse

type Visitor interface {
	visitBinary(v BinaryExpr)
	visitGrouping(v GroupingExpr)
	visitLiteral(v LiteralExpr)
	visitUnary(v UnaryExpr)
}
