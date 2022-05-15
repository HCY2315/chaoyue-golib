package thirdparty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/HCY2315/chaoyue-golib/log"
	"github.com/HCY2315/chaoyue-golib/pkg/errors"
)

type SNSInfoReqVO struct {
	Code string
}

type SNSInfoRespVO struct {
	OpenID string
}

type WeChatSNSService interface {
	GetUserSNSInfo(SNSInfoReqVO) (SNSInfoRespVO, error)
}

type WechatSnsResult struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	ErrorCode  int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

type WechatSNSImpl struct {
	wechatSNSUrlTemplate string
}

func NewWechatSNSImpl(wechatSNSUrlTemplate string) *WechatSNSImpl {
	return &WechatSNSImpl{wechatSNSUrlTemplate: wechatSNSUrlTemplate}
}

func (w *WechatSNSImpl) GetUserSNSInfo(req SNSInfoReqVO) (SNSInfoRespVO, error) {
	var resp SNSInfoRespVO
	snsURL := fmt.Sprintf(w.wechatSNSUrlTemplate, req.Code)
	wechatSNSResponse, errGet := http.Get(snsURL)
	if errGet != nil {
		return resp, errors.Wrap(errGet, "request wechat %s", snsURL)
	}
	defer wechatSNSResponse.Body.Close()
	body, errRead := ioutil.ReadAll(wechatSNSResponse.Body)
	if errRead != nil {
		return resp, errors.Wrap(errRead, "read wechat wechatSNSResponse body")
	}
	var snsResponse WechatSnsResult
	if errUnmarshal := json.Unmarshal(body, &snsResponse); errUnmarshal != nil {
		return resp, errors.Wrap(errUnmarshal, "marshal %s to json", body)
	}
	log.Infof("url: ", snsURL, "\nwechatSNSResponse:", string(body))
	if snsCode := snsResponse.ErrorCode; snsCode != 0 {
		if snsCode == 40029 {
			return resp, errors.Wrap(errors.ErrBadRequest, "微信code不正确或已过期，请重试")
		}
		return resp, errors.Wrap(errors.ErrBadRequest, "invalid response %+v for url:%s", snsResponse, snsURL)
	}
	resp.OpenID = snsResponse.OpenID
	return resp, nil
}

type wechatSNSDebugWrapper struct {
	serviceImpl        WeChatSNSService
	superCode          string
	openIDForSuperCode string
}

func NewWechatSNSDebugWrapper(serviceImpl WeChatSNSService, superCode string, openIDForSuperCode string) WeChatSNSService {
	return &wechatSNSDebugWrapper{serviceImpl: serviceImpl, superCode: superCode, openIDForSuperCode: openIDForSuperCode}
}

func (w *wechatSNSDebugWrapper) GetUserSNSInfo(req SNSInfoReqVO) (SNSInfoRespVO, error) {
	if req.Code == w.superCode {
		var resp SNSInfoRespVO
		resp.OpenID = w.openIDForSuperCode
		return resp, nil
	}
	return w.serviceImpl.GetUserSNSInfo(req)
}
