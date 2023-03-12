package main

import (
	"errors"
	"fmt"
)

type Error struct {
	Message string
}

func (e *Error) Raw() error {
	return errors.New(e.Message)
}

func (e *Error) Format(v ...interface{}) error {
	return fmt.Errorf(e.Message, v...)
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}

var (
	ValuesTooManyError   = NewError("too many values provided to '%s'.")
	ValuesTooLittleError = NewError("too little values provided to '%s'.")

	ConversionError = NewError("failed '%s' conversion.")

	SyntaxError = NewError("syntax error '%s'.")

	InvalidObjectError = NewError("invalid object '%s'.")
	NoObjectError      = NewError("no object provided to '%s'.")

	NotAssignableError = NewError("object '%s' is not assignable (P.S. '@' is forbidden - builtins only!).")
	NotCallableError   = NewError("object '%s' is not callable.")

	NumericsError = NewError("numerics error '%s'.")
)
