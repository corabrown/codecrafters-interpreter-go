package parse

import (
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
)

type AstPrinter struct {
	currentString string
}

func (v *AstPrinter) VisitBinary(b data.BinaryExpr) {
	v.currentString = v.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (v *AstPrinter) VisitGrouping(g data.GroupingExpr) {
	v.currentString = v.parenthesize("group", g.Expression)
}

func (v *AstPrinter) VisitLiteral(l data.LiteralExpr) {
	if l.Value == nil {
		v.currentString = "nil"
	}
	v.currentString = l.Value.ToString()
}

func (v *AstPrinter) VisitUnary(u data.UnaryExpr) {
	v.currentString = v.parenthesize(u.Operator.Lexeme, u.Right)
}

func (v *AstPrinter) parenthesize(name string, exprs ...data.Expression) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")

		if expr != nil {
			expr.Accept(v)
		}
		sb.WriteString(v.currentString)
		v.currentString = ""
	}
	sb.WriteString(")")

	return sb.String()
}

func (v *AstPrinter) Print(expr data.Expression) string {
	expr.Accept(v)
	return v.currentString
}
