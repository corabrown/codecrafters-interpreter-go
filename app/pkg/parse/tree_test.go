package parse

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
)

func Test_AstPrinter(t *testing.T) {
	expr := data.BinaryExpr{
		Left: data.UnaryExpr{
			Operator: data.Token{TokenType: data.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    data.LiteralExpr{Value: data.NumberLiteral{Val: 123}},
		},
		Operator: data.Token{TokenType: data.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right:    data.GroupingExpr{Expression: data.LiteralExpr{Value: data.NumberLiteral{Val: 45.67}}},
	}

	printer := &AstPrinter{}
	output := printer.Print(expr)

	if output != "(* (- 123.0) (group 45.67))" {
		t.Errorf("output: %v does not match expected: %v", output, "(* (- 123) (group 45.67))")
	}
}
