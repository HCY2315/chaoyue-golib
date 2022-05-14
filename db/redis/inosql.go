package db

import "time"

type RedisZ struct {
	Score  float64
	Member interface{}
}

type INoSqlClient interface {
	Enable() bool
	Stop()
	Do(args ...interface{}) (interface{}, error)
	Ping() (string, error)
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	GetSet(key string, value interface{}) (string, error)
	Del(keys ...string) error
	Expire(key string, expiration time.Duration) error
	Exists(keys ...string) (int64, error)
	Incr(key string) (int64, error)
	IncrBy(key string, increment int64) (int64, error)
	Decr(key string) (int64, error)
	DecrBy(key string, decrement int64) (int64, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	MGet(keys ...string) ([]interface{}, error)
	MSet(pairs ...interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HSet(key, field string, value interface{}) error
	HSetNX(key, field string, value interface{}) (bool, error)
	HDel(key string, fields ...string) error
	HKeys(key string) ([]string, error)
	HVals(key string) ([]string, error)
	HLen(key string) (int64, error)
	HMGet(key string, fields ...string) ([]interface{}, error)
	HMSet(key string, fields map[string]interface{}) error
	LPush(key string, values ...interface{}) error
	LPop(key string) (string, error)
	LRange(key string, start, stop int64) ([]string, error)
	LLen(key string) (int64, error)
	RPush(key string, values ...interface{}) error
	RPop(key string) (string, error)
	SAdd(key string, members ...interface{}) (int64, error)
	SRem(key string, members ...interface{}) error
	SIsMember(key string, member interface{}) (bool, error)
	SMembers(key string) ([]string, error)
	SCard(key string) (int64, error)
	ZAdd(key string, score float64, value interface{}) error
	ZRem(key string, members ...interface{}) error
	ZCard(key string) (int64, error)
	ZScore(key, member string) (float64, error)
	ZCount(key, min, max string) (int64, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScoresEx(key string, start, stop int64) ([]RedisZ, error)
	ZRangeByScoreEx(key string, min, max string, offset, count int64) ([]string, error)
	ZRangeByScoreWithScoresEx(key string, min, max string, offset, count int64) ([]RedisZ, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeByScoreWithScoresEx(key string, min, max string, offset, count int64) ([]RedisZ, error)
	ZRank(key, member string) (int64, error)
	ZRevRank(key, member string) (int64, error)
	Scan(cursor uint64, match string, count int64) (keys []string, cursor2 uint64, err error)
	SScan(key string, cursor uint64, match string, count int64) (keys []string, cursor2 uint64, err error)
	HScan(key string, cursor uint64, match string, count int64) (keys []string, cursor2 uint64, err error)
	ZScan(key string, cursor uint64, match string, count int64) (keys []string, cursor2 uint64, err error)
	Pipeline() INoSqlPipeline
	Publish(channel string, msg interface{}) error
	Subscribe(callback func(channel, msg string, err error), channels ...string) error
	HIncrBy(h, k string, value int64) (int64, error)
	TTL(key string) (time.Duration, error)
	Info(section ...string) (string, error)
}

type INoSqlPipeline interface {
	Close() error
	Exec() error
	Do(args ...interface{}) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) error
	Expire(key string, expiration time.Duration) error
	Incr(key string) (int64, error)
	IncrBy(key string, ememnt int64) (int64, error)
	Decr(key string) (int64, error)
	DecrBy(key string, ememnt int64) (int64, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	MSet(pairs ...interface{}) error
	HSet(key, field string, value interface{}) error
	HSetNX(key, field string, value interface{}) (bool, error)
	HDel(key string, fields ...string) error
	HIncrBy(key, field string, incr int64) (int64, error)
	HMSet(key string, fields map[string]interface{}) error
	LPush(key string, values ...interface{}) error
	LPop(key string) (string, error)
	RPush(key string, values ...interface{}) error
	RPop(key string) (string, error)
	SAdd(key string, members ...interface{}) (int64, error)
	SRem(key string, members ...interface{}) error
	ZAdd(key string, score float64, value interface{}) error
	ZRem(key string, members ...interface{}) error
	Publish(channel string, msg interface{}) (err error)
}
