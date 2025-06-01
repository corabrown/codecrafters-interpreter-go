package errors

import (
	"fmt"
	"os"
)

type Error struct {
	line    int
	message string
}

func (v Error) Report(where string) {
	fmt.Fprint(os.Stderr, fmt.Sprintf("[line %d] Error: %v: %v\n", v.line, where, v.message))
}

func NewError(line int, message string) Error {
	return Error{line: line, message: message}
}
