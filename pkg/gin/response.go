package gin

import "git.cestong.com.cn/cecf/cecf-golib/pkg/errors"

type GeneralVO struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewGeneralVO(code int, msg string) *GeneralVO {
	return &GeneralVO{
		Code: code,
		Msg:  msg,
	}
}

const (
	OKCode = 0
	OKMsg  = "OK"
)

func SuccessGeneralVO() GeneralVO {
	return GeneralVO{
		Code: OKCode,
		Msg:  OKMsg,
	}
}

func NewGeneralVOFromErr(err error) GeneralVO {
	var v GeneralVO
	if ec, ok := err.(errors.ErrorWithCode); ok {
		v.Code = ec.Code()
	}
	v.Msg = err.Error()
	return v
}

type EnhancedRespErr struct {
	Code     int      `json:"code"`
	Msg      string   `json:"msg"`
	Tags     []string `json:"tags"`
	UserHint string   `json:"userHint"`
	HelpURL  string   `json:"helpURl"`
}

func (e *EnhancedRespErr) Error() string {
	return e.Msg
}

type GeneralResponseVO struct {
	Error *EnhancedRespErr `json:"error"`
	Data  interface{}      `json:"data"`
}
