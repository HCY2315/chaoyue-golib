package thirdparty

import (
	"encoding/json"

	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/utils"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dys "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

type SmsSendService interface {
	SendOne(phone, content string) error
}

type AliSMSService struct {
	domain       string
	accessKeyID  string
	accessSecret string
	sigName      string
	tempCode     string
	cli          *dys.Client
}

func NewAliSMSService(domain, accessKeyID, accessSecret, sigName, tempCode string) (*AliSMSService, error) {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyID,
		AccessKeySecret: &accessSecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(domain)
	cli, err := dys.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &AliSMSService{
		domain:       domain,
		accessKeyID:  accessKeyID,
		accessSecret: accessSecret,
		sigName:      sigName,
		tempCode:     tempCode,
		cli:          cli,
	}, nil
}

type aliyunSMSReqData struct {
	Code string `json:"code"`
}

// refactor: template
func (a *AliSMSService) SendOne(phone, content string) error {
	var tempData aliyunSMSReqData
	tempData.Code = content
	dataStr, err := json.Marshal(tempData)
	if err != nil {
		return err
	}
	sendSmsRequest := &dys.SendSmsRequest{
		SignName:      tea.String(a.sigName),
		TemplateCode:  tea.String(a.tempCode),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(utils.BytesToString(dataStr)),
	}
	resp, err := a.cli.SendSms(sendSmsRequest)
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return errors.Wrap(ErrThirdParty, "no resp body")
	}
	if body := resp.Body; *body.Code != "OK" {
		return errors.Wrap(ErrThirdParty, "resp err:%s-%s", *body.Code, *body.Message)
	}
	return nil
}
