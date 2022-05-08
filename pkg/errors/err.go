package errors

import (
	"fmt"
	"net/http"
)

type ErrorWithCode interface {
	error
	Code() int
}

type ErrorWithHTTPStatus interface {
	HTTPStatus() int
}

type ErrorWithCodeAndStatus interface {
	ErrorWithHTTPStatus
	ErrorWithCode
}

type simpleErrorWithCode struct {
	code int
	msg  string
}

func (s simpleErrorWithCode) Error() string {
	return s.msg
}

func (s simpleErrorWithCode) Code() int {
	return s.code
}

func NewErrorWithCode(code int, msg string) ErrorWithCode {
	return &simpleErrorWithCode{
		code: code,
		msg:  msg,
	}
}

func ErrorWithCodeFromErr(code int, err error) ErrorWithCode {
	return NewErrorWithCode(code, err.Error())
}

func ExtendErrWithCodeContextMsg(ec ErrorWithCode, format string, values ...interface{}) ErrorWithCode {
	msg := fmt.Sprintf("%s <- (%s)", fmt.Sprintf(format, values...), ec.Error())
	return NewErrorWithCode(ec.Code(), msg)
}

func ExtendErrWithCodeAndStatusContextMsg(ecs ErrorWithCodeAndStatus, format string, values ...interface{}) ErrorWithCodeAndStatus {
	msg := fmt.Sprintf("%s <- (%s)", fmt.Sprintf(format, values...), ecs.Error())
	return ErrorWithCodeAndHTTPStatus(ecs.Code(), msg, ecs.HTTPStatus())
}

type simpleErrorWithCodeAndHTTPStatus struct {
	simpleErrorWithCode
	httpStatus int
}

func ErrorWithCodeAndHTTPStatus(code int, msg string, status int) ErrorWithCodeAndStatus {
	return simpleErrorWithCodeAndHTTPStatus{
		simpleErrorWithCode: simpleErrorWithCode{
			code: code,
			msg:  msg,
		},
		httpStatus: status,
	}
}

const (
	CodeOffsetToStatus = 3000
)

func ErrorWithHTTPStatusOffsetCode(msg string, status int) ErrorWithCodeAndStatus {
	return ErrorWithCodeAndHTTPStatus(status+CodeOffsetToStatus, msg, status)
}

func (s simpleErrorWithCodeAndHTTPStatus) HTTPStatus() int {
	return s.httpStatus
}

const (
	ErrCodeInternal        = 3500
	ErrCodeBadRequest      = 3400
	ErrCodeEntryNotFound   = 3404
	ErrCodeShouldNotHappen = 3599
)

var (
	ErrEntryNotFound   = ErrorWithCodeAndHTTPStatus(ErrCodeEntryNotFound, "找不到对象", http.StatusNotFound)
	ErrType            = ErrorWithCodeAndHTTPStatus(3451, "类型错误", 451)
	ErrShouldNotHappen = ErrorWithCodeAndHTTPStatus(ErrCodeShouldNotHappen, "代码错误", 599)
	ErrInternal        = ErrorWithCodeAndHTTPStatus(ErrCodeInternal, "服务错误", http.StatusInternalServerError)
	ErrBadRequest      = ErrorWithCodeAndHTTPStatus(ErrCodeBadRequest, "错误请求", http.StatusBadRequest)
	ErrUnauthorized    = ErrorWithCodeAndHTTPStatus(3401, "未登录或凭证已过期", http.StatusUnauthorized)
)
