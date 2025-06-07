package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/lox"
	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/parse"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	if _, ok := acceptedCommands[command]; !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lox.NewLox(fileContents)

	defer func() {
		if l.HadError() {
			l.ReportErrors()
			os.Exit(65)
		}
	}()

	switch command {
	case "tokenize":
		l.ScanFile()
		for _, t := range l.GetTokens() {
			fmt.Fprint(os.Stdout, t.ToString())
		}
	case "parse":
		l.ParseFile()
		if l.ParseError() {
			return
		}
		printer := &parse.AstPrinter{}
		fmt.Fprint(os.Stdout, printer.Print(l.GetExpression()))
	}
}

var acceptedCommands = map[string]struct{}{
	"tokenize": {},
	"parse":    {},
}
