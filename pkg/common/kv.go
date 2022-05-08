package common

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"time"
)

type NamedValues interface {
	Keys() ([]string, error)
	Get(k string) (interface{}, error)
	Set(k string, value interface{}) error
}

type TypedNamedValues interface {
	NamedValues
	GetInt(k string) (int, error)
	GetBool(k string) (bool, error)
	GetString(k string) (string, error)
}

//NewSimpleMemNamedValues 少量kv存储，类似url.Values，非线程安全
func NewSimpleMemNamedValues() *memNamedValues {
	return &memNamedValues{
		m: make(map[string]interface{}, 0),
	}
}

func BuildExpireTimeGetHookFunc(ttl time.Duration, key string) GetHookFunc {
	return func(values NamedValues, k string) (valueOverride interface{}) {
		if k == key {
			return time.Now().Add(ttl).UnixMilli()
		}
		return nil
	}
}
func NewSimpleMemNamedValuesWithHook(hk GetHookFunc) *memNamedValues {
	m := NewSimpleMemNamedValues()
	m.getHook = hk
	return m
}

type GetHookFunc func(values NamedValues, k string) (valueOverride interface{})

type memNamedValues struct {
	m       map[string]interface{}
	getHook GetHookFunc
}

func (m memNamedValues) GetInt(k string) (int, error) {
	v, e := m.Get(k)
	if e != nil {
		return 0, e
	}
	i, ok := v.(int)
	if !ok {
		return 0, errors.Wrap(errors.ErrType, "%+v not int", v)
	}
	return i, nil
}

func (m memNamedValues) GetBool(k string) (bool, error) {
	v, e := m.Get(k)
	if e != nil {
		return false, e
	}
	i, ok := v.(bool)
	if !ok {
		return false, errors.Wrap(errors.ErrType, "%+v not bool", v)
	}
	return i, nil
}

func (m memNamedValues) GetString(k string) (string, error) {
	v, e := m.Get(k)
	if e != nil {
		return "", e
	}
	i, ok := v.(string)
	if !ok {
		return "", errors.Wrap(errors.ErrType, "%+v not string", v)
	}
	return i, nil
}

func (m memNamedValues) Keys() ([]string, error) {
	ks := make([]string, 0, len(m.m))
	for k := range m.m {
		ks = append(ks, k)
	}
	return ks, nil
}

func (m memNamedValues) Get(k string) (interface{}, error) {
	if m.getHook != nil {
		ov := m.getHook(&m, k)
		if ov != nil {
			return ov, nil
		}
	}
	v, find := m.m[k]
	if !find {
		return nil, errors.ErrEntryNotFound
	}
	return v, nil
}

func (m memNamedValues) Set(k string, value interface{}) error {
	m.m[k] = value
	return nil
}

type sortedMemNamedValues struct {
	memNamedValues
	keys []string
}

//NewSortedMemNamedValues Keys返回时按Set FIFO
func NewSortedMemNamedValues() NamedValues {
	return &sortedMemNamedValues{
		memNamedValues: memNamedValues{m: make(map[string]interface{}, 0)},
		keys:           nil,
	}
}

func (s *sortedMemNamedValues) Set(k string, value interface{}) error {
	s.keys = append(s.keys, k)
	for i := 0; i < len(s.keys)-1; i++ {
		if s.keys[i] == k {
			s.keys = append(s.keys[:i], s.keys[i+1:]...)
			break
		}
	}
	return s.memNamedValues.Set(k, value)
}

func (s *sortedMemNamedValues) Keys() ([]string, error) {
	return s.keys, nil
}

type OverrideNamedValues struct {
	nvs []NamedValues
}

func (o OverrideNamedValues) Keys() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (o OverrideNamedValues) Get(k string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (o OverrideNamedValues) Set(k string, value interface{}) error {
	//TODO implement me
	panic("implement me")
}
