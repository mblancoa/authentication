package errors

import (
	"fmt"
)

type ErrorCode string

const (
	Error               ErrorCode = "Error"
	RuntimeError                  = "Runtime"
	NotFoundError                 = "Not Found"
	AuthenticationError           = "Authentication error"
)

type basicError struct {
	Code    ErrorCode
	Message string
	Cause   string
}

func NewError(code ErrorCode, message string) error {
	return basicError{
		Code:    code,
		Message: message,
	}
}

func NewErrorf(code ErrorCode, format string, a ...any) error {
	return basicError{
		Code:    code,
		Message: fmt.Sprintf(format, a),
	}
}

func NewErrorByCause(code ErrorCode, message string, cause error) error {
	msg := ""
	if message != "" {
		if cause.Error() != "" {
			msg = fmt.Sprintf("%s\nCaused by %s", message, cause)
		} else {
			msg = message
		}
	} else {
		msg = cause.Error()
	}
	return basicError{
		Code:    code,
		Message: msg,
	}
}

func NewGenericError(message string) error {
	return NewError(Error, message)
}

func NewGenericErrorf(format string, a ...any) error {
	return NewErrorf(Error, format, a)
}

func NewGenericErrorByCause(message string, cause error) error {
	return NewErrorByCause(Error, message, cause)
}

func NewRuntimeError(message string) error {
	return NewError(RuntimeError, message)
}

func NewRuntimeErrorf(format string, a ...any) error {
	return NewErrorf(RuntimeError, format, a)
}

func NewRuntimeErrorByCause(message string, cause error) error {
	return NewErrorByCause(RuntimeError, message, cause)
}

func NewNotFoundError(message string) error {
	return NewError(NotFoundError, message)
}

func NewAuthenticationError(message string) error {
	return NewErrorf(AuthenticationError, message)
}

func (error basicError) Error() string {
	return error.Message
}
