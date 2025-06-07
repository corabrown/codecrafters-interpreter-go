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
	fmt.Fprint(os.Stderr, fmt.Sprintf("[line %d] Error%v: %v\n", v.line, v.where, v.message))
}

func NewError(line int, message string, where string) *Error {
	return &Error{line: line, message: message, where: where}
}
