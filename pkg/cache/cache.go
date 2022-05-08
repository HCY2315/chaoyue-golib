// Package cache 实现缓存功能
package cache

type ICache interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
}
