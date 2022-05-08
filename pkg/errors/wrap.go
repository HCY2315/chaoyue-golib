package errors

import "github.com/pkg/errors"

func Wrap(err error, format string, values ...interface{}) error {
	switch err.(type) {
	case ErrorWithCodeAndStatus:
		return ExtendErrWithCodeAndStatusContextMsg(err.(ErrorWithCodeAndStatus), format, values...)
	case ErrorWithCode:
		ec := err.(ErrorWithCode)
		return ExtendErrWithCodeContextMsg(ec, format, values...)
	default:
		return errors.Wrapf(err, format, values...)
	}
}
