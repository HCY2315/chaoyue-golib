package middleware

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/log"
	"github.com/gin-gonic/gin"
	"time"
)

func BuildAccessLogMiddleware(logTag string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Now().Sub(start) / time.Millisecond
		log.Infof("%s %d %d %s %v", logTag, c.Writer.Status(), latency, c.Request.Method, c.Request.URL)
	}
}