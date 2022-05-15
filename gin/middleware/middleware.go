// Package middleware 包含gin框架常用的中间件
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alibaba/sentinel-golang/core/system"
	sentinelPlugin "github.com/alibaba/sentinel-golang/pkg/adapters/gin"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			// 允许访问所有域
			c.Header("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 允许跨域设置,可以返回其他子段
			// 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			// 缓存请求信息 单位为秒
			c.Header("Access-Control-Max-Age", "172800")
			// 跨域请求是否需要带cookie信息 默认设置为true
			c.Header("Access-Control-Allow-Credentials", "true")
			// 设置返回格式是json
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

// Sentinel 限流
func Sentinel() gin.HandlerFunc {
	if _, err := system.LoadRules([]*system.Rule{
		{
			MetricType:   system.InboundQPS,
			TriggerCount: 200,
			Strategy:     system.BBR,
		},
	}); err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	return sentinelPlugin.SentinelMiddleware()
}

// Trace 链路追踪
func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var sp opentracing.Span
		opName := ctx.Request.URL.Path
		// Attempt to join a trace by getting trace context from the headers.
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			// If for whatever reason we can't join, go ahead an start a new root span.
			sp = opentracing.StartSpan(opName)
		} else {
			sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
		}
		ctx.Set("traceSpan", sp)
		ctx.Next()
		sp.Finish()
	}
}
