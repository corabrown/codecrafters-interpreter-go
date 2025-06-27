package lox

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/evaluate"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/scan"
)

type Lox struct {
	fileContents []byte
	scan.Scanner
	parse.Parser
	evaluate.Interpreter
	tokens        []data.Token
	expression    data.Expression
	Value         interface{}
	errors        []*errors.Error
	runtimeErrors []*errors.RuntimeError
}

func NewLox(fileContents []byte) Lox {
	if len(fileContents) == 0 {
		fmt.Println("EOF  null")
		return Lox{}
	}

	return Lox{fileContents: fileContents}
}

func (l *Lox) ScanFile() {
	if l.fileContents == nil {
		return
	}
	l.Scanner = scan.NewScanner(string(l.fileContents))
	tokens, scanErrors := l.Scan()
	l.tokens = tokens
	l.errors = append(l.errors, scanErrors...)
}

func (l *Lox) ParseFile() {
	if !l.Scanned() {
		l.ScanFile()
		if l.HadError() {
			return
		}
	}

	l.Parser = parse.NewParser(l.GetTokens())
	expr, parseErrors := l.Parse()
	l.expression = expr
	l.errors = append(l.errors, parseErrors...)
}

func (l *Lox) HadError() bool {
	return len(l.errors) > 0
}

func (l *Lox) ReportErrors() {
	for _, e := range l.errors {
		e.Report()
	}
}

func (l *Lox) HadRuntimeError() bool {
	return len(l.runtimeErrors) > 0
}

func (l *Lox) ReportRuntimeErrors() {
	for _, e := range l.runtimeErrors {
		e.Report()
	}
}

func (l *Lox) GetExpression() data.Expression {
	return l.expression
}

func (l *Lox) EvaluateFile() {
	if !l.Parsed() {
		l.ParseFile()
		if l.HadError() {
			return
		}
	}

	l.Interpreter = evaluate.NewInterpreter(l.expression)
	errors := l.Evaluate()
	l.Value = l.Interpreter.GetValue()
	l.runtimeErrors = append(l.runtimeErrors, errors...)
}
