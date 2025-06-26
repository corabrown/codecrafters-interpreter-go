package evaluate

import "github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"

type Interpreter struct {
	expression data.Expression
	value      interface{}
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
	}
}

func (i *Interpreter) VisitBinary(v data.BinaryExpr) {}

func (i *Interpreter) VisitGrouping(v data.GroupingExpr) {}

func (i *Interpreter) VisitLiteral(v data.LiteralExpr) {
	switch lit := v.Value.(type) {
	case data.BooleanLiteral:
		i.value = lit.Val
	case data.NullLiteral:
		i.value = nil
	}

}

func (i *Interpreter) VisitUnary(v data.UnaryExpr) {}
