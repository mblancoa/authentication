package errors

import (
	"fmt"
	"reflect"
)

type Code string

const (
	// Error is Default error
	Error               Code = "Error"
	NotFoundError       Code = "Not Found"
	AuthenticationError Code = "Authentication error"
)

type basicError struct {
	Code    Code   `json:"code"`
	Message string `json:"Message"`
	Cause   string `json:"-"`
}

func NewError(code Code, message string) error {
	return basicError{
		Code:    code,
		Message: message,
	}
}

func NewErrorf(code Code, format string, a ...any) error {
	return basicError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

func NewErrorByCause(code Code, message string, cause error) error {
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
	return NewErrorf(Error, format, a...)
}

func NewGenericErrorByCause(message string, cause error) error {
	return NewErrorByCause(Error, message, cause)
}

func NewNotFoundError(message string) error {
	return NewError(NotFoundError, message)
}

func NewAuthenticationError(message string) error {
	return NewErrorf(AuthenticationError, message)
}

func NewAuthenticationErrorByCause(message string, cause error) error {
	return NewErrorByCause(AuthenticationError, message, cause)
}

func (error basicError) Error() string {
	return error.Message
}

func GetCode(err error) Code {
	if "basicError" == reflect.TypeOf(err).Name() {
		return err.(basicError).Code
	}
	return ""
}

func GetCodeOrDefault(err error, def Code) Code {
	if "basicError" == reflect.TypeOf(err).Name() {
		return err.(basicError).Code
	}
	return def
}

func ManageErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
