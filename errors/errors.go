package errors

import (
    "errors"
    e "github.com/pkg/errors"
)

type Error interface {
    error
}

func New(text string) error {
    return errors.New(text)
}

func Wrap(err error, message string) error {
    return e.Wrap(err, message)
}

func Cause(err error) error {
    return e.Cause(err)
}

func Is(err, target error) bool {
    return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
    return errors.As(err, target)
}

func AsErrorObj(err error) (e *ErrorObj, ok bool) {
    ok = errors.As(err, &e)
    return
}
