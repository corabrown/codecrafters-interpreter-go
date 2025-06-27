package errors

import (
	"fmt"
	"os"
)

type Error struct {
	line    int
	message string
	where   string
}

func (v *Error) Report() {
	fmt.Fprintf(os.Stderr, "[line %d] Error%v: %v\n", v.line, v.where, v.message)
}

func NewError(line int, message string, where string) *Error {
	return &Error{line: line, message: message, where: where}
}

type RuntimeError struct {
	line    int
	message string
}

func NewRuntimeError(line int, message string) *RuntimeError {
	return &RuntimeError{line: line, message: message}
}

func (v *RuntimeError) Report() {
	fmt.Fprintf(os.Stderr, "%v\n[line %v]", v.message, v.line)
}
