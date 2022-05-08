/*
 * @File   : prometheus
 * @Author : huangbin
 *
 * @Created on 2020/9/11 4:34 下午
 * @Project : microapp
 * @Software: GoLand
 * @Description  :
 */

package prometheus

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMonitor struct {
	metrics       map[string]*MetricNode
	promeRegister *prometheus.Registry
	router        string
	w             sync.WaitGroup
	stopCh        chan struct{}
	mtxRW         sync.RWMutex
}

/*
 * @Method   : NewPrometheusMonitor/新建prometheus实例
 *
 * @param    : router string 获取数据指标的接口,用于prometheus调用
 * @Return   : *PrometheusMonitor  实例指针
 *
 * @Description :
 */
func NewPrometheusMonitor(router string) *PrometheusMonitor {
	ins := &PrometheusMonitor{
		router: router,
	}

	if len(ins.router) < 1 {
		ins.router = CSTDefaultRouter
	}
	ins.stopCh = make(chan struct{})
	ins.promeRegister = prometheus.NewPedanticRegistry()
	ins.metrics = make(map[string]*MetricNode)

	ins.w.Add(1)
	go ins.timingClearMetricCache()

	return ins
}

func (mo *PrometheusMonitor) ResetCacheHttpMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// 数据被获取后，清除本地缓存
		mo.resetMetrics(false)
	})
}

func (mo *PrometheusMonitor) ResetCacheGINMiddleWare(ctx *gin.Context) {
	ctx.Next()
	// 数据被获取后，清除本地缓存
	mo.resetMetrics(false)
}

func (mo *PrometheusMonitor) resetMetrics(force bool) {
	mo.mtxRW.RLock()
	defer mo.mtxRW.RUnlock()

	now := time.Now()
	for _, v := range mo.metrics {
		switch v.MetricType {
		case CSTMetricCounterVec:
			if ((v.ClearCache.Enable && now.Sub(v.ClearCache.Point) >= v.ClearCache.Period) || force) && v.HaveRegister {
				v.Metric.(*prometheus.CounterVec).Reset()
				v.ClearCache.Point = now
			}
		case CSTMetricGaugeVec:
			if ((v.ClearCache.Enable && now.Sub(v.ClearCache.Point) >= v.ClearCache.Period) || force) && v.HaveRegister {
				v.Metric.(*prometheus.GaugeVec).Reset()
				v.ClearCache.Point = time.Now()
			}
		case CSTMetricHistogramVec:
			if ((v.ClearCache.Enable && now.Sub(v.ClearCache.Point) >= v.ClearCache.Period) || force) && v.HaveRegister {
				v.Metric.(*prometheus.HistogramVec).Reset()
				v.ClearCache.Point = now
			}
		case CSTMetricSummaryVec:
			if ((v.ClearCache.Enable && now.Sub(v.ClearCache.Point) >= v.ClearCache.Period) || force) && v.HaveRegister {
				v.Metric.(*prometheus.SummaryVec).Reset()
				v.ClearCache.Point = now
			}
		}
	}
}

func (mo *PrometheusMonitor) timingClearMetricCache() {
	defer func() {
		mo.w.Done()
	}()

	now := time.Now()
	next := now.Add(time.Minute * 1)
	next = time.Date(next.Year(), next.Month(), next.Day(), 23, 30, 0, 0, next.Location())
	t := time.NewTimer(next.Sub(now))

	for {
		select {
		case <-mo.stopCh:
			return
		case <-t.C:
			mo.resetMetrics(true)

			now := time.Now()
			next = now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 23, 30, 0, 0, next.Location())
			t = time.NewTimer(next.Sub(now))
			fmt.Printf("[%s] ResetAllMetrics, NextReset Point: %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				next.Format("2006-01-02 15:04:05"))
		}
	}
}
