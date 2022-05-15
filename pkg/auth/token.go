package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HCY2315/chaoyue-golib/log"
	"github.com/HCY2315/chaoyue-golib/pkg/common"
	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/utils"
	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenInvalid = errors.ErrorWithCodeAndHTTPStatus(3401, "token 不正确", http.StatusUnauthorized)
)

type TokenValues common.TypedNamedValues

type TokenBuilder interface {
	Token(kvs TokenValues) (string, error)
	Parse(token string) (TokenValues, error)
}

type jwtTokenBuilder struct {
	// 全局添加的值
	global common.NamedValues
	// 有效期
	ttl time.Duration
	// 加密用
	filter     common.DoubleDirectionTextFilter
	alg        jwt.SigningMethod
	signSecret []byte
}

func (j *jwtTokenBuilder) Token(kvs TokenValues) (string, error) {
	// jwt generate
	cms := jwt.MapClaims{}
	if errCopy := namedValueToClaim(cms, j.global); errCopy != nil {
		return "", errCopy
	}
	if errCopy := namedValueToClaim(cms, kvs); errCopy != nil {
		return "", errCopy
	}
	token, errSign := jwt.NewWithClaims(j.alg, cms).SignedString(j.signSecret)
	if errSign != nil {
		return "", errors.Wrap(errSign, "sign")
	}
	if j.filter == nil {
		return token, nil
	}
	log.Debugf("token gen with values %+v", cms)
	// encrypt
	encryptToken, errEncrypt := j.filter.Forward(token)
	if errEncrypt != nil {
		return "", errors.Wrap(errEncrypt, "encrypt")
	}
	return encryptToken, nil
}

func namedValueToClaim(cms jwt.MapClaims, nvs common.NamedValues) error {
	globalKeys, errGetGlobalKeys := nvs.Keys()
	if errGetGlobalKeys != nil {
		return errors.Wrap(errGetGlobalKeys, "get keys")
	}
	for _, k := range globalKeys {
		value, errGetValue := nvs.Get(k)
		if errGetValue != nil {
			return errors.Wrap(errGetValue, "get value by %s", k)
		}
		cms[k] = value
	}
	return nil
}

func (j *jwtTokenBuilder) Parse(token string) (TokenValues, error) {
	plainToken := token
	if j.filter != nil {
		var errDecrypt error
		plainToken, errDecrypt = j.filter.Backward(token)
		if errDecrypt != nil {
			return nil, errors.Wrap(errDecrypt, "decrypt %s", token)
		}
	}
	parsed, errParse := jwt.Parse(plainToken, func(token *jwt.Token) (interface{}, error) {
		if alg := token.Method.Alg(); alg != j.alg.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", alg)
		}
		return j.signSecret, nil
	})
	if errParse != nil {
		return nil, errors.Wrap(errParse, "Parse token [%s]", plainToken)
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token %s", plainToken)
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Wrap(errors.ErrShouldNotHappen, "claim not map")
	}
	return claimsToToken(claims)
}

func claimsToToken(cms jwt.MapClaims) (TokenValues, error) {
	t := common.NewSimpleMemNamedValues()
	for k, v := range cms {
		//fixme put int but float64 as result
		var value interface{} = v
		fl, isFloat := v.(float64)
		if isFloat && utils.IsThisFloatAInt(fl) {
			value = int(fl)
		}
		if errSet := t.Set(k, value); errSet != nil {
			return nil, errors.Wrap(errSet, "set %s -> %v to NamedValues", k, v)
		}
	}
	return t, nil
}

type JwtTokenBuildOptions func(b *jwtTokenBuilder)

func WithSetAlgOption(alg jwt.SigningMethod) JwtTokenBuildOptions {
	return func(b *jwtTokenBuilder) {
		b.alg = alg
	}
}

//WithValuesGetter 相同key后面会覆盖前面的
func WithValuesGetter(globalNamedValues common.NamedValues) JwtTokenBuildOptions {
	return func(b *jwtTokenBuilder) {
		b.global = globalNamedValues
	}
}

func WithFilter(f common.DoubleDirectionTextFilter) JwtTokenBuildOptions {
	return func(b *jwtTokenBuilder) {
		b.filter = f
	}
}

const (
	JWTClaimExpireKey = "exp"
)

func NewJwtTokenBuilder(ttl time.Duration, signSecret []byte, ops ...JwtTokenBuildOptions) *jwtTokenBuilder {
	globalNV := common.NewSimpleMemNamedValuesWithHook(common.BuildExpireTimeGetHookFunc(ttl, JWTClaimExpireKey))
	if errSet := globalNV.Set(JWTClaimExpireKey, nil); errSet != nil {
		panic(errSet)
	}
	r := &jwtTokenBuilder{
		ttl:        ttl,
		signSecret: signSecret,
		global:     globalNV,
		alg:        jwt.SigningMethodHS256,
	}
	for _, op := range ops {
		op(r)
	}

	return r
}
