package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	pcg "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type GinHandlerFilterFunc func(c *gin.Context) bool

func ExcludeByPath(path string) GinHandlerFilterFunc {
	return func(c *gin.Context) bool {
		return c.Request.URL.Path != path
	}
}

// req_count,count(server_id,status,host,ip,method, path)
// http_latency_histogram, latency(server_id, status, method, path, ip, host)
type PrometheusExporter struct {
	latencyHistogram pcg.ExemplarObserver
	serverID         string
	host             string
	ip               string
	filterFunc       func(c *gin.Context) bool
}

func NewPrometheusExporter(serverID, host, hostIP string, filter func(c *gin.Context) bool) *PrometheusExporter {
	constantLabels := pcg.Labels{
		"host":     host,
		"ip":       hostIP,
		"serverID": serverID,
	}
	h := promauto.NewHistogram(pcg.HistogramOpts{
		Name:        "http_access_histogram",
		Help:        "http请求延迟、method、status",
		ConstLabels: constantLabels,
		Buckets:     pcg.LinearBuckets(0, 50, 20),
	}).(pcg.ExemplarObserver)
	return &PrometheusExporter{
		serverID:         serverID,
		host:             host,
		ip:               host,
		latencyHistogram: h,
		filterFunc:       filter,
	}
}

func (p *PrometheusExporter) HandleFunc(c *gin.Context) {
	if !p.filterFunc(c) {
		c.Next()
		return
	}
	method := c.Request.Method
	begin := time.Now()
	c.Next()
	status := c.Writer.Status()
	latency := time.Since(begin).Milliseconds()
	labels := pcg.Labels{
		"method": method,
		"status": strconv.Itoa(status),
		"path":   c.Request.URL.Path,
	}
	p.latencyHistogram.ObserveWithExemplar(float64(latency), labels)
}

func (p *PrometheusExporter) ExportMetricsHandler() gin.HandlerFunc {
	hh := promhttp.Handler()
	return func(c *gin.Context) {
		hh.ServeHTTP(c.Writer, c.Request)
	}
}
