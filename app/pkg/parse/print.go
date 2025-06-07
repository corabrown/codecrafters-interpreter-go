package parse

import "strings"

type AstPrinter struct {
	currentString string
}

func (v *AstPrinter) visitBinary(b BinaryExpr) {
	v.currentString = v.parenthesize(b.operator.Lexeme, b.left, b.right)
}

func (v *AstPrinter) visitGrouping(g GroupingExpr) {
	v.currentString = v.parenthesize("group", g.expression)
}

func (v *AstPrinter) visitLiteral(l LiteralExpr) {
	if l.value == nil {
		v.currentString = "nil"
	}
	v.currentString = l.value.ToString()
}

func (v *AstPrinter) visitUnary(u UnaryExpr) {
	v.currentString = v.parenthesize(u.operator.Lexeme, u.right)
}

func (v *AstPrinter) parenthesize(name string, exprs ...Expression) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		expr.accept(v)
		sb.WriteString(v.currentString)
		v.currentString = ""
	}
	sb.WriteString(")")

	return sb.String()
}

func (v *AstPrinter) Print(expr Expression) string {
	expr.accept(v)
	return v.currentString
}
