package evaluate

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
)

type Interpreter struct {
	expression data.Expression
	value      interface{}
	errors     []*errors.RuntimeError
}

const (
	runtimeUnaryErrorMessage  = "Operand must be number."
	runtimeBinaryErrorMessage = "Operands must be numbers."
)

func NewInterpreter(expr data.Expression) Interpreter {
	return Interpreter{expression: expr}
}

func (i *Interpreter) addError(token data.Token, message string) {
	i.errors = append(i.errors, errors.NewRuntimeError(token.Line, message))
}

func (i *Interpreter) GetValue() interface{} {
	return i.value
}

func (i *Interpreter) Evaluate() []*errors.RuntimeError {
	switch expr := i.expression.(type) {
	case data.LiteralExpr:
		i.VisitLiteral(expr)
	case data.GroupingExpr:
		i.VisitGrouping(expr)
	case data.UnaryExpr:
		i.VisitUnary(expr)
	case data.BinaryExpr:
		i.VisitBinary(expr)
	}

	return i.errors
}

func (i *Interpreter) VisitBinary(v data.BinaryExpr) {

	left := NewInterpreter(v.Left)
	left.Evaluate()

	right := NewInterpreter(v.Right)
	right.Evaluate()

	l, r := left.value, right.value

	switch v.Operator.TokenType {
	case data.MINUS:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf - rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.SLASH:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf / rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.STAR:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf * rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.PLUS:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf + rf
		} else if ls, rs, ok := getStrings(l, r); ok {
			i.value = ls + rs
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.GREATER:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf > rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.GREATER_EQUAL:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf >= rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.LESS:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf < rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.LESS_EQUAL:
		if lf, rf, ok := getFloats(l, r); ok {
			i.value = lf <= rf
		} else {
			i.addError(v.Operator, runtimeBinaryErrorMessage)
		}
	case data.BANG_EQUAL:
		i.value = !isEqual(l, r)
	case data.EQUAL_EQUAL:
		i.value = isEqual(l, r)
	}
}

func getFloats(left, right any) (float64, float64, bool) {
	if lf, okLeft := left.(float64); okLeft {
		if rf, okRight := right.(float64); okRight {
			return lf, rf, true
		}
	}
	return 0, 0, false
}

func getStrings(left, right any) (string, string, bool) {
	if ls, okLeft := left.(string); okLeft {
		if rs, okRight := right.(string); okRight {
			return ls, rs, true
		}
	}
	return "", "", false
}

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
			i.addError(v.Operator, runtimeUnaryErrorMessage)
		}
	case data.BANG:
		i.value = !isTruthy(right.value)

	}
}

func isTruthy(val any) bool {
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

func isEqual(left, right any) bool {
	if (left == nil) && (right == nil) {
		return true
	}
	if (left == nil) || (right == nil) {
		return false
	}

	if l, ok := left.(float64); ok {
		if r, ok := right.(float64); ok {
			return l == r
		}
	}

	if l, ok := left.(string); ok {
		if r, ok := right.(string); ok {
			return l == r
		}
	}

	if l, ok := left.(bool); ok {
		if r, ok := right.(bool); ok {
			return l == r
		}
	}

	return false
}
