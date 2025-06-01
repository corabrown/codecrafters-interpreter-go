package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/pkg/scan"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	commands := os.Args
	_ = commands
	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var s scan.Scanner
	if len(fileContents) > 0 {
		s = scan.Scan(string(fileContents))
	} else {
		fmt.Println("EOF  null")
	}

	lox := Lox{hadError: s.ScanError()}
	if lox.hadError {
		os.Exit(65)
	}
}

type Lox struct {
	hadError bool
}
