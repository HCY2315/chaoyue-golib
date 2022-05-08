package captcha

import (
	"context"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/log"
	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
)

// 图形验证码

type DrawCaptchaService interface {
	GenerateBase64Img() (string, string, error)
	VerifyMatch(id, code string) (bool, error)
}

type SimpleDrawCaptchaService struct {
	driver base64Captcha.Driver
	store  DrawCaptchaStore
}

type SetDrawCaptchaServiceOptions func(s *SimpleDrawCaptchaService)

// todo 支持其他类型driver
func SetDriverType(typ string) SetDrawCaptchaServiceOptions {
	return func(s *SimpleDrawCaptchaService) {
		s.driver = base64Captcha.NewDriverString(128, 256, 5, base64Captcha.OptionShowHollowLine,
			4, "aa04597427053978fc168df1129ae74a78739e1b", nil, nil, nil)
	}
}

func NewSimpleDrawCaptchaService(store DrawCaptchaStore, opts ...SetDrawCaptchaServiceOptions) *SimpleDrawCaptchaService {
	ret := &SimpleDrawCaptchaService{store: store}
	for _, opt := range opts {
		opt(ret)
	}
	if ret.driver == nil {
		SetDriverType("")(ret)
	}
	return ret
}

func (s *SimpleDrawCaptchaService) GenerateBase64Img() (string, string, error) {
	c := base64Captcha.NewCaptcha(s.driver, s.store)
	return c.Generate()
}

func (s *SimpleDrawCaptchaService) VerifyMatch(id, code string) (bool, error) {
	return s.store.Verify(id, code, true), nil
}

type DrawCaptchaStore base64Captcha.Store

type RedisDrawCaptchaStore struct {
	prefix   string
	ttl      time.Duration
	redisCli *redis.Client
}

func NewRedisDrawCaptchaStore(prefix string, ttl time.Duration, redisCli *redis.Client) *RedisDrawCaptchaStore {
	return &RedisDrawCaptchaStore{prefix: prefix, ttl: ttl, redisCli: redisCli}
}

func (r *RedisDrawCaptchaStore) Set(id string, value string) {
	k := idToKey(r.prefix, id)
	log.Debugf("[drawCode] set %s->%s", k, value)
	err := r.redisCli.Set(context.Background(), k, value, r.ttl).Err()
	if err != nil {
		log.Errorf("set draw captcha %s", err.Error())
	}
	return
}

func idToKey(prefix, id string) string {
	return prefix + ":" + id
}

func (r *RedisDrawCaptchaStore) Get(id string, clear bool) string {
	k := idToKey(r.prefix, id)
	ctx := context.Background()
	log.Debugf("[drawCode] get %s", k)
	value, err := r.redisCli.Get(ctx, k).Result()
	if err != nil {
		if err == redis.Nil {
			return ""
		}
		log.Errorf("redis get %s failed:%s", k, err.Error())
		return ""
	}
	log.Debugf("[drawCode] get %s -> %s", k, value)
	if clear {
		_, errDel := r.redisCli.Del(ctx, k).Result()
		if errDel != nil {
			log.Errorf("redis del %s failed:%s", k, errDel.Error())
		}
	}
	return value
}

func (r *RedisDrawCaptchaStore) Verify(id, answer string, clear bool) bool {
	valueInDB := r.Get(id, clear)
	return valueInDB == answer
}
