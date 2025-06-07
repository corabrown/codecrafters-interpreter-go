package parse

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/scan"
)

func Test_AstPrinter(t *testing.T) {
	expr := BinaryExpr{
		left: UnaryExpr{
			operator: scan.Token{scan.MINUS, "-", nil, 1},
			right:    LiteralExpr{value: scan.NumberLiteral{123}},
		},
		operater: scan.Token{scan.STAR, "*", nil, 1},
		right:    GroupingExpr{LiteralExpr{scan.NumberLiteral{45.67}}},
	}

	printer := &AstPrinter{}
	output := printer.Print(expr)

	if output != "(* (- 123.0) (group 45.67))" {
		t.Errorf("output: %v does not match expected: %v", output, "(* (- 123) (group 45.67))")
	}
}
