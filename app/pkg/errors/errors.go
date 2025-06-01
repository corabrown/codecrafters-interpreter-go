package errors

import "fmt"

type Error struct {
	line    int
	message string
}

func (v Error) Report(where string) {
	fmt.Printf("[line %v] Error: %v: %v\n", v.line, where, v.message)
}

func NewError(line int, message string) Error {
	return Error{line: line, message: message}
}
