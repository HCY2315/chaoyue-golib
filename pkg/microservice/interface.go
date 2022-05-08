// Package microservice 表示一个微服务实例。包括服务信息和创建方法
package microservice

//PersistObject 代表一个可以被持久化/反持久化的数据
type PersistObject interface {
	Encode() []byte
	Recover(persistBytes []byte) error
}
