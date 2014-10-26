package chewcrew

import (
	"errors"
	"fmt"
	"runtime"
)

const maxStackLength = 100000

// A simple wrapper around native errors that also provides a stack trace and an optional context
// message.
type Error struct {
	e       error
	context string
	Stack   string
}

func (e Error) Error() string {
	s := ""

	// prefix with context if there is one
	if e.context != "" {
		s = e.context + " : "
	}

	// provide error message and stack
	s += fmt.Sprintf("%s\n%s", e.e.Error(), e.Stack)

	return s
}

// Create an error with a basic message.
func NewError(s string) error {
	err := errors.New(s)

	stackBuf := make([]byte, maxStackLength, maxStackLength)
	bytesRead := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:bytesRead])

	return Error{err, "", stack}
}

// Create an error with a formatted message.
func NewErrorf(f string, a ...interface{}) error {
	err := fmt.Errorf(f, a...)

	stackBuf := make([]byte, maxStackLength, maxStackLength)
	bytesRead := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:bytesRead])

	return Error{err, "", stack}
}

// Useful for wrapping errors received from core and third party libraries, providing them a stack
// trace.
func NewMaskedError(underlying error) error {
	stackBuf := make([]byte, maxStackLength, maxStackLength)
	bytesRead := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:bytesRead])

	return Error{underlying, "", stack}
}

// Useful for wrapping errors received from core and third party libraries, providing them a stack
// trace. Additionally, this method allows for a message to be saved for context.
func NewMaskedErrorWithContext(underlying error, context string) error {
	stackBuf := make([]byte, maxStackLength, maxStackLength)
	bytesRead := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:bytesRead])

	return Error{underlying, context, stack}
}

// Useful for wrapping errors received from core and third party libraries, providing them a stack
// trace. Additionally, this method allows for a formatted message to be saved for context.
func NewMaskedErrorWithContextf(underlying error, f string, a ...interface{}) error {
	stackBuf := make([]byte, maxStackLength, maxStackLength)
	bytesRead := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:bytesRead])

	return Error{underlying, fmt.Sprintf(f, a...), stack}
}
