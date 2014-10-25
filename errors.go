package main

import (
  "fmt"
  "errors"
  "runtime"
)

const maxStackLength = 100000

// a simple wrapper around native errors that also provides a stack trace
type Error struct {
  e error
  context string
  Stack string
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

// create an error with a basic message
func NewError(s string) error {
  err := errors.New(s)

  stackBuf := make([]byte, maxStackLength, maxStackLength)
  bytesRead := runtime.Stack(stackBuf, false)
  stack := string(stackBuf[:bytesRead])

  return Error{err, "", stack}
}

// create an error with a formatted message
func NewErrorf(f string, a ...interface{}) error {
  err := fmt.Errorf(f, a...)

  stackBuf := make([]byte, maxStackLength, maxStackLength)
  bytesRead := runtime.Stack(stackBuf, false)
  stack := string(stackBuf[:bytesRead])

  return Error{err, "", stack}
}

// useful for wrapping errors received from core and third party libraries, providing them stack
// traces
func NewMaskedError(underlying error) error {
  stackBuf := make([]byte, maxStackLength, maxStackLength)
  bytesRead := runtime.Stack(stackBuf, false)
  stack := string(stackBuf[:bytesRead])

  return Error{underlying, "", stack}
}

// useful for wrapping errors received from core and third party libraries, providing them stack
// traces. additionally, this method allows for an additional formatted message to be saved for
// context.
func NewMaskedErrorf(underlying error, f string, a ...interface{}) error {
  stackBuf := make([]byte, maxStackLength, maxStackLength)
  bytesRead := runtime.Stack(stackBuf, false)
  stack := string(stackBuf[:bytesRead])

  return Error{underlying, fmt.Sprintf(f, a...), stack}
}