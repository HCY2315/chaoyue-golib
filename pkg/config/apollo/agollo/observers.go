package agollo

import (
	"sync"
)

const (
	namespaceToKeySep = ":"
)

// Observer 配置项更改回调函数
type Observer func(key string, oldValue interface{}, newValue interface{})

// observerMap stores the obersers registered for some keys
type observerMap struct {
	observersDelegate map[string]Observer // key->observer yaml版本的namespace的key支持多层级比如config.yaml:hystrix.ask_app_server，key的格式为${namespace}:${real_key}
	lock              sync.RWMutex
}

func newObservers() *observerMap {
	return &observerMap{
		lock:              sync.RWMutex{},
		observersDelegate: map[string]Observer{},
	}
}

func (o *observerMap) add(namespace, configKey string, cb Observer) {
	o.lock.Lock()
	defer o.lock.Unlock()

	key := namespace + namespaceToKeySep + configKey

	o.observersDelegate[key] = cb
}

func (o *observerMap) get(namespace, configKey string) Observer {
	o.lock.RLock()
	defer o.lock.RUnlock()

	key := namespace + namespaceToKeySep + configKey
	return o.observersDelegate[key]
}

func (o *observerMap) getObservers() map[string]Observer {
	o.lock.RLock()
	defer o.lock.RUnlock()

	res := map[string]Observer{}
	for k, v := range o.observersDelegate {
		res[k] = v
	}
	return res
}
