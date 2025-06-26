package evaluate

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
)

type Interpreter struct {
	expression data.Expression
	value      interface{}
	errors     []*errors.Error
}

func NewInterpreter(expr data.Expression) Interpreter {
	return Interpreter{expression: expr}
}

func (i *Interpreter) GetValue() interface{} {
	return i.value
}

func (i *Interpreter) Evaluate() {
	switch expr := i.expression.(type) {
	case data.LiteralExpr:
		i.VisitLiteral(expr)
	case data.GroupingExpr:
		i.VisitGrouping(expr)
	case data.UnaryExpr:
		i.VisitUnary(expr)
	}
}

func (i *Interpreter) VisitBinary(v data.BinaryExpr) {}

func (i *Interpreter) VisitGrouping(v data.GroupingExpr) {
	v.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteral(v data.LiteralExpr) {
	switch lit := v.Value.(type) {
	case data.BooleanLiteral:
		i.value = lit.Val
	case data.NullLiteral:
		i.value = nil
	case data.NumberLiteral:
		i.value = lit.Val
	case data.StringLiteral:
		i.value = lit.Val
	}
}

func (i *Interpreter) VisitUnary(v data.UnaryExpr) {
	right := NewInterpreter(v.Right)
	right.Evaluate()

	switch v.Operator.TokenType {
	case data.MINUS:
		switch r := right.value.(type) {
		case float64:
			i.value = -r
		default:
			i.errors = append(i.errors, errors.NewError(0, "incorrect use of unary -", ""))
		}
	case data.BANG:
		i.value = !isTruthy(right.value)

	}
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	default:
		return true
	}
}
