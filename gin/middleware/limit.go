package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	capacity  int        // 桶的容量
	rate      float64    // 令牌放入速率，每秒存放的数量
	tokens    float64    // 当前令牌数量
	lastToken time.Time  // 上一次放令牌的时间
	mtx       sync.Mutex // 互斥锁
}

func (tb *TokenBucket) Allow() bool {
	tb.mtx.Lock()
	defer tb.mtx.Unlock()
	now := time.Now()
	// 计算需要放的令牌数量
	tb.tokens = tb.tokens + tb.rate*now.Sub(tb.lastToken).Seconds()
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}
	// 判断是否允许请求
	if tb.tokens >= 1 {
		tb.tokens--
		tb.lastToken = now
		return true
	} else {
		return false
	}
}

func NewTokenBucket(capacity int, rateBySec float64) *TokenBucket {
	return &TokenBucket{
		capacity:  capacity,
		rate:      rateBySec,
		tokens:    float64(capacity),
		lastToken: time.Now(),
	}
}

func (tb *TokenBucket) LimitHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !tb.Allow() {
			c.String(503, "Too many request")
			c.Abort()
			return
		}
		c.Next()
	}
}
