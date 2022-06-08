package captcha

import (
	"context"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/thirdparty"
	"github.com/HCY2315/chaoyue-golib/pkg/utils"
	"github.com/go-redis/redis/v8"
)

var (
	ErrVerifyCodeNotMatch = errors.ErrorWithCodeAndHTTPStatus(3400, "验证码不匹配", 400)
)

type SMSVerifyService interface {
	SendCodeToPhone(phone string) error
	VerifyMatch(phone, code string) (bool, error)
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
