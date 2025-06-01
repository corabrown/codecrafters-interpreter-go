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

	if len(fileContents) > 0 {
		scan.Scan(string(fileContents))
	} else {
		fmt.Println("EOF  null")
	}

	lox := Lox{}
	if lox.hadError {
		os.Exit(65)
	}
}

type Error struct {
	line    int
	message string
}

func (v Error) report(where string) {
	fmt.Printf("[line %v] Error %v: %v", v.line, where, v.message)
}

type Lox struct {
	hadError bool
}
