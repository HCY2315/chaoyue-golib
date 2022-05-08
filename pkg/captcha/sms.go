package captcha

import (
	"context"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/log"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/thirdparty"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ErrVerifyCodeNotMatch = errors.ErrorWithCodeAndHTTPStatus(3400, "验证码不匹配", 400)
)

type SMSVerifyService interface {
	SendCodeToPhone(phone string) error
	VerifyMatch(phone, code string) (bool, error)
}

type DebugSMSVerify struct {
}

func (d DebugSMSVerify) SendCodeToPhone(phone string) error {
	log.Debugf("[SMS] send code to phone %s", phone)
	return nil
}

func (d DebugSMSVerify) VerifyMatch(phone, code string) (bool, error) {
	log.Debugf("[SMS] check phone %s match code %s", phone, code)
	return true, nil
}

func NewDebugSMSVerify() *DebugSMSVerify {
	return &DebugSMSVerify{}
}

// refactor factory
type SMSVerifyWithSuperCode struct {
	superCode string
	impl      SMSVerifyService
}

func (s *SMSVerifyWithSuperCode) SendCodeToPhone(phone string) error {
	return s.impl.SendCodeToPhone(phone)
}

func (s *SMSVerifyWithSuperCode) VerifyMatch(phone, code string) (bool, error) {
	if code == s.superCode {
		log.Warnf("SMS verify %s using SUPER-CODE", phone)
		return true, nil
	}
	return s.impl.VerifyMatch(phone, code)
}

func NewSMSVerifyWithSuperCode(superCode string, impl SMSVerifyService) *SMSVerifyWithSuperCode {
	return &SMSVerifyWithSuperCode{superCode: superCode, impl: impl}
}

// refactor kv store
type SMSVerifyImpl struct {
	prefix       string
	ttl          time.Duration
	aliSMSSender thirdparty.SmsSendService
	redisCli     *redis.Client
}

func NewSmsVerifyImpl(prefix string, ttl time.Duration, aliSMSSender thirdparty.SmsSendService, redisCli *redis.Client) *SMSVerifyImpl {
	return &SMSVerifyImpl{prefix: prefix, ttl: ttl, aliSMSSender: aliSMSSender, redisCli: redisCli}
}

func (a *SMSVerifyImpl) SendCodeToPhone(phone string) error {
	code := utils.RandDigitString(6)
	if err := a.aliSMSSender.SendOne(phone, code); err != nil {
		return errors.Wrap(err, "call ali sms service")
	}
	k := a.persistKey(phone)
	ctx := context.Background()
	if err := a.redisCli.Set(ctx, k, code, a.ttl).Err(); err != nil {
		return errors.Wrap(err, "set store:%s -> %s", k, code)
	}
	return nil
}

func (a *SMSVerifyImpl) persistKey(id string) string {
	return a.prefix + ":" + id
}

func (a *SMSVerifyImpl) VerifyMatch(phone, code string) (bool, error) {
	k := a.persistKey(phone)
	ctx := context.Background()
	s, err := a.redisCli.Get(ctx, k).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, errors.Wrap(err, "get from store %s", k)
	}
	return s == code, nil
}
