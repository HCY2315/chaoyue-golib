package middleware

import (
	"fmt"
	"net/http"
	"strings"

	gin2 "github.com/HCY2315/chaoyue-golib/gin"
	"github.com/HCY2315/chaoyue-golib/log"
	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/thirdparty"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ErrMiddlewareCallback func(code, status int, msg string, ctx *gin.Context)

func DingMsgErrMiddlewareCallback(ding thirdparty.DingService) ErrMiddlewareCallback {
	return func(code, status int, msg string, c *gin.Context) {
		if status < 300 {
			return
		}
		body := c.GetString(CtxKeyDebugReqBody)
		method := c.Request.Method
		fullURL := c.Request.Host + c.Request.URL.String()
		headers := c.Request.Header
		textToSend := fmt.Sprintf(`
错误:%s
status: %d
errCode: %d
method: %s
url:%s
req body:%s
headers:%+v
`, msg, status, code, method, fullURL, body, headers)
		if errSend := ding.SendPlain(textToSend); errSend != nil {
			log.Errorf("发送钉钉失败:%s\n%s", errSend.Error(), textToSend)
		}
	}
}

func BuildErrMiddleware(trans ut.Translator, cbs ...ErrMiddlewareCallback) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 合并Msg
		if len(c.Errors) == 0 {
			return
		}
		es := c.Errors
		var errMsgBuilder strings.Builder
		for _, err := range c.Errors {
			if valErr, ok := err.Err.(validator.ValidationErrors); ok {
				transErrMsgMap := valErr.Translate(trans)
				var errMsgSB strings.Builder
				errMsgSB.WriteString("请检查以下字段: ")
				for _, v := range transErrMsgMap {
					errMsgSB.WriteString(v)
					errMsgSB.WriteString(";")
				}
				c.JSON(http.StatusBadRequest, gin2.GeneralVO{
					Code: errors.ErrCodeBadRequest,
					Msg:  errMsgSB.String(),
				})
				return
			}
			errMsgBuilder.WriteString(err.Error())
			errMsgBuilder.WriteRune('\n')
		}
		var errVO gin2.GeneralVO
		errVO.Msg = errMsgBuilder.String()
		// 挑选主要HTTP status, code
		httpStatus, code := selectErrorCode(es)
		errVO.Code = code
		c.JSON(httpStatus, errVO)
		log.Errorf("%s-%v status:%d, code:%d, msg:%s", c.Request.Method, c.Request.URL, httpStatus, code, errVO.Msg)
		for _, cb := range cbs {
			cb(errVO.Code, httpStatus, errVO.Msg, c)
		}
	}
}

func selectErrorCode(es []*gin.Error) (retStatus int, retCode int) {
	if len(es) == 0 {
		return http.StatusOK, 0
	}
	// StatusErr > CodeErr > error
	// Status 500+ > 400+ > 300+ > 其他选大
	// Code选大
	// 其余都是500 3500
	var statusErr errors.ErrorWithCodeAndStatus
	var codeErr errors.ErrorWithCode
	for _, e := range es {
		err := e.Err
		switch err.(type) {
		case errors.ErrorWithCodeAndStatus:
			ecStatus := err.(errors.ErrorWithCodeAndStatus)
			statusErr = overlayECS(ecStatus, statusErr)
		case errors.ErrorWithCode:
			ec := err.(errors.ErrorWithCode)
			if codeErr == nil || ec.Code() > codeErr.Code() {
				codeErr = ec
			}
			continue
		default:
			continue
		}
	}
	if statusErr == nil {
		retStatus = http.StatusInternalServerError
		if codeErr == nil {
			retCode = errors.ErrCodeInternal
		} else {
			retCode = codeErr.Code()
		}
		return
	}
	return statusErr.HTTPStatus(), statusErr.Code()
}

func overlayECS(e2, e1 errors.ErrorWithCodeAndStatus) errors.ErrorWithCodeAndStatus {
	if e1 == nil {
		return e2
	}
	if e2 == nil {
		return e1
	}
	s1 := e1.HTTPStatus()
	s2 := e2.HTTPStatus()
	f1, f2 := isFocusedStatus(s1), isFocusedStatus(s2)
	if f2 {
		if f1 {
			if s2 > s1 {
				return e2
			}
			return e1
		}
		return e2
	}
	if f1 {
		return e1
	}
	if s2 > s1 {
		return e2
	}
	return e1
}

func isFocusedStatus(status int) bool {
	return status >= http.StatusOK && status < 600
}
