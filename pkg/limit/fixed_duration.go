package limit

import (
	"context"
	"strconv"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/go-redis/redis/v8"
)

//FixedDurationLimiter 固定时间段内访问限制，超过时间段后重置
type FixedDurationLimiter interface {
	TryAccess(id string) (allow bool, nextTry time.Time, err error)
}

type FixedDurationResolveFunc func(id string, accessTime time.Time) (right time.Time, durationID string)

func BuildFixedMinuteDurationResolve(minutes int64) FixedDurationResolveFunc {
	return func(id string, accessTime time.Time) (time.Time, string) {
		nano := accessTime.UnixNano()
		step := minutes * int64(time.Minute)
		minuteIndex := nano / step
		rightNano := (minuteIndex + 1) * step
		rightSeconds := rightNano / int64(time.Second)
		right := time.Unix(rightSeconds, rightNano-rightSeconds*int64(time.Second))
		return right.In(accessTime.Location()), strconv.FormatInt(minuteIndex, 10)
	}
}

type RedisFixedDurationLimiter struct {
	redisCli        *redis.Client
	durationResolve FixedDurationResolveFunc
	limitUpper      int64
}

func NewRedisFixedDurationLimiter(redisCli *redis.Client, durationResolve FixedDurationResolveFunc,
	limitUpper int64) *RedisFixedDurationLimiter {
	return &RedisFixedDurationLimiter{
		redisCli:        redisCli,
		durationResolve: durationResolve,
		limitUpper:      limitUpper,
	}
}

func (r *RedisFixedDurationLimiter) TryAccess(id string) (bool, time.Time, error) {
	accessTime := time.Now()
	rightRange, durationID := r.durationResolve(id, time.Now())
	nextTry := rightRange.Add(time.Second * 1)
	key := r.accessKey(id, durationID)
	ctx := context.Background()
	nowCnt, errIncr := r.redisCli.Incr(ctx, key).Result()
	if errIncr != nil {
		return false, nextTry, errors.Wrap(errIncr, "incr %s", key)
	}
	if nowCnt == 1 {
		ttl := rightRange.Sub(accessTime)
		if errExpire := r.redisCli.Expire(ctx, key, ttl).Err(); errExpire != nil {
			return false, nextTry, errors.Wrap(errExpire, "expire %s", key)
		}
	}
	if nowCnt > r.limitUpper {
		return false, nextTry, nil
	}
	return true, nextTry, nil
}

const (
	FixDurationLimiterRedisPrefix = "lim:fixed"
)

func (r *RedisFixedDurationLimiter) accessKey(id string, durationID string) string {
	return FixDurationLimiterRedisPrefix + ":" + id + ":" + durationID
}
