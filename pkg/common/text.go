package common

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
	"strings"
)

type DoubleDirectionTextFilter interface {
	Forward(string) (string, error)
	Backward(string) (string, error)
}

type DoubleDirectionStringifyFilter interface {
	EncodeToString([]byte) string
	DecodeFromString(string) ([]byte, error)
}

type GCMCryptTextFilter struct {
	gcm             *GCMCipher
	stringifyFilter DoubleDirectionStringifyFilter
}

type Base64URLFilter struct {
}

func NewBase64URLFilter() *Base64URLFilter {
	return &Base64URLFilter{}
}

func (b Base64URLFilter) EncodeToString(bytes []byte) string {
	return base64.URLEncoding.EncodeToString(bytes)
}

func (b Base64URLFilter) DecodeFromString(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

type HexFilter struct {
}

func NewHexFilter() *HexFilter {
	return &HexFilter{}
}

func (h HexFilter) EncodeToString(i []byte) string {
	return hex.EncodeToString(i)
}

func (h HexFilter) DecodeFromString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func NewGCMCryptTextFilterByPasswd(passwd []byte, stringify DoubleDirectionStringifyFilter) (*GCMCryptTextFilter, error) {
	if stringify == nil {
		stringify = &Base64URLFilter{}
	}
	gcm, err := NewGCMCipher(passwd)
	if err != nil {
		return nil, err
	}
	return &GCMCryptTextFilter{
		gcm:             gcm,
		stringifyFilter: stringify,
	}, nil
}

func (g *GCMCryptTextFilter) Forward(s string) (string, error) {
	b, e := g.gcm.Encrypt(utils.String2Bytes(s))
	if e != nil {
		return "", e
	}
	return g.stringifyFilter.EncodeToString(b), nil
}

func (g *GCMCryptTextFilter) Backward(s string) (string, error) {
	bytes, errBytes := g.stringifyFilter.DecodeFromString(s)
	if errBytes != nil {
		return "", errors.Wrap(errBytes, "decode bytes from %s", s)
	}
	b, e := g.gcm.Decrypt(bytes)
	if e != nil {
		return "", e
	}
	return utils.Bytes2String(b), nil
}

func ConcatMsg(ss ...string) string {
	if len(ss) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteString(s)
	}
	return sb.String()
}

func ConcatErrMsg(es ...error) string {
	if len(es) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, s := range es {
		sb.WriteString(s.Error())
	}
	return sb.String()
}

func ConcatStringer(ss ...fmt.Stringer) string {
	if len(ss) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteString(s.String())
	}
	return sb.String()
}
