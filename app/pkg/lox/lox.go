package lox

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/data"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/errors"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/scan"
)

type Lox struct {
	fileContents []byte
	scan.Scanner
	parse.Parser
	tokens     []data.Token
	expression data.Expression
	errors     []*errors.Error
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
	if l.fileContents == nil {
		return
	}

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

func (v *Lox) HadError() bool {
	return len(v.errors) > 0
}

func (v *Lox) ReportErrors() {
	for _, e := range v.errors {
		e.Report()
	}
}

func (v *Lox) GetExpression() data.Expression {
	return v.expression
}
