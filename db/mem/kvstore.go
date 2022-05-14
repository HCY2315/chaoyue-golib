package db

type KVStore interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
