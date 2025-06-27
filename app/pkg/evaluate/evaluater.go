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

func NewInterpreter(expr data.Expression) Interpreter {
	return Interpreter{expression: expr}
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

	in, exprType := getLeftAndRightVals(v)

	if v.Operator.TokenType == data.BANG_EQUAL {
		i.value = !isEqual(in.left, in.right)
	}
	if v.Operator.TokenType == data.EQUAL_EQUAL {
		i.value = isEqual(in.left, in.right)
	}

	switch exprType {
	case floatType:
		l, r := in.leftFloat, in.rightFloat
		switch v.Operator.TokenType {
		case data.MINUS:
			i.value = l - r
		case data.SLASH:
			i.value = l / r
		case data.STAR:
			i.value = l * r
		case data.PLUS:
			i.value = l + r
		case data.GREATER:
			i.value = l > r
		case data.GREATER_EQUAL:
			i.value = l >= r
		case data.LESS:
			i.value = l < r
		case data.LESS_EQUAL:
			i.value = l <= r
		}
	case stringType:
		if v.Operator.TokenType == data.PLUS {
			i.value = in.leftString + in.rightString
		}
	}
}

type binaryInput struct {
	leftFloat   float64
	rightFloat  float64
	leftString  string
	rightString string
	left        any
	right       any
}

type binaryExpressionType string

const (
	stringType  binaryExpressionType = "string"
	floatType   binaryExpressionType = "float"
	invalidType binaryExpressionType = "invalid"
)

func getLeftAndRightVals(v data.BinaryExpr) (b binaryInput, exprType binaryExpressionType) {
	exprType = invalidType

	left := NewInterpreter(v.Left)
	left.Evaluate()

	right := NewInterpreter(v.Right)
	right.Evaluate()

	b.left = left.value
	b.right = right.value

	if l, ok := left.value.(float64); ok {
		if r, ok := right.value.(float64); ok {
			b.leftFloat = l
			b.rightFloat = r
			exprType = floatType
		}
	} else if l, ok := left.value.(string); ok {
		if r, ok := right.value.(string); ok {
			b.leftString = l
			b.rightString = r
			exprType = stringType
		}
	}

	return
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
			i.errors = append(i.errors, errors.NewRuntimeError("Operand must be a number."))
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
