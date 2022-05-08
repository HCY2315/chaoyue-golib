package apollo

import (
	"encoding/json"
	"reflect"
	"sync"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
)

type Map struct {
	m    map[string]interface{}
	lock sync.RWMutex
}

func (m *Map) Reset(vals map[string]interface{}) {
	m.lock.Lock()
	if nil == vals {
		m.m = map[string]interface{}{}
	} else {
		m.m = vals
	}
	m.lock.Unlock()
}

func (m *Map) GetMap() map[string]interface{} {
	return m.m
}

func (m *Map) Set(key string, value interface{}) {
	m.lock.Lock()
	m.m[key] = value
	m.lock.Unlock()
}

func (m *Map) Get(key string) (interface{}, bool) {
	m.lock.RLock()
	val, ok := m.m[key]
	m.lock.RUnlock()
	return val, ok
}

// json map对象转yaml raw
func jsonMapToYaml(m interface{}) (bytes string, err error) {

	var yamlByte, mByte []byte

	mByte, err = json.Marshal(&m)
	if err != nil {
		return
	}
	yamlByte, err = yaml.JSONToYAML(mByte)
	if err != nil {
		return
	}
	return string(yamlByte), err
}

// 校验yaml格式
func checkFormat(buf []byte) (err error) {
	c := make(map[string]interface{}, 0)
	err = yaml.Unmarshal(buf, &c)
	return errors.Wrap(err, "checkFormat")
}

// 校验json格式
func checkJsonValid(buf []byte) (err error) {
	m := make(map[string]interface{}, 0)
	err = utils.Unmarshal(buf, &m)
	return errors.Wrap(err, "checkJsonValid")
}

// 判断是否包含
func Contains(container interface{}, obj interface{}) bool {
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if containerValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if containerValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}
