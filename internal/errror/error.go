package errror

import (
	"fmt"
)

type ErrorWithMessage struct {
	err error
	msg string
}

func (e *ErrorWithMessage) Unwrap() []error {
	return []error{e.err}
}

func (e *ErrorWithMessage) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func WrapError(err error, msg string) error {
	return &ErrorWithMessage{err: err, msg: msg}
}
