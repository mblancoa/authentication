package errors

import "fmt"

type basicError struct {
	Message string
	Cause   string
}

func NewError(message string) error {
	return basicError{
		Message: message,
	}
}

func NewErrorf(format string, a ...any) error {
	return basicError{
		Message: fmt.Sprintf(format, a),
	}
}

func NewErrorByCause(message string, cause error) error {
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
		Message: msg,
	}
}
func (error basicError) Error() string {
	return error.Message
}
