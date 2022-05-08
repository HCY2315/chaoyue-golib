/*
 * @File   : handle
 * @Author : huangbin
 *
 * @Created on 2020/9/7 10:47 上午
 * @Project : prometheus
 * @Software: GoLand
 * @Description  : 该文件为对外提供的使用接口
 */

package prometheus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func (mo *PrometheusMonitor) GetRegistry() *prometheus.Registry {
	return mo.promeRegister
}

/*
 * @Method   : Start/启动prometheus
 *
 * @param    : router *gin.Engine  gin句柄/如果传空则使用原生http注册接口
 * @Return   : error  错误
 *
 * @Description : 1、利用原生http或gin注册consul和prometheus接口，具体取决于使用方
 *                2、注册服务节点到consul用于prometheus服务发现
 */
func (mo *PrometheusMonitor) Start(router *gin.Engine) error {
	if router == nil {
		http.Handle(mo.router, mo.ResetCacheHttpMiddleWare(promhttp.HandlerFor(mo.promeRegister, promhttp.HandlerOpts{})))
	} else {
		router.GET(mo.router, mo.ResetCacheGINMiddleWare, gin.WrapH(promhttp.HandlerFor(mo.promeRegister, promhttp.HandlerOpts{})))
	}

	for _, v := range mo.metrics {
		mo.promeRegister.MustRegister(v.Metric)
	}

	return nil
}

func (mo *PrometheusMonitor) Stop() error {
	close(mo.stopCh)
	//fmt.Println("Consul DelRegisterService: " + string(rep))
	return nil
}

func (mo *PrometheusMonitor) RegisterHttpHandler(r *mux.Router) {
	r.Handle(mo.router, mo.ResetCacheHttpMiddleWare(promhttp.HandlerFor(mo.promeRegister, promhttp.HandlerOpts{})))
}

func (mo *PrometheusMonitor) RegisterGinHandler(router *gin.Engine) {
	router.GET(mo.router, mo.ResetCacheGINMiddleWare, gin.WrapH(promhttp.HandlerFor(mo.promeRegister, promhttp.HandlerOpts{})))
}

/*
 * @Method   : RegisterMetric/注册具体的数据指标(通用)
 *
 * @param    : metricType MetricConstType 指标类型/使用预定义的常量类型
 * @param    : metricName string 指标名称/与注册指标时标签名称一致
 * @param    : labels []string   标签名列表
 * @param    : timingClearCache bool   是否开启定时清除本地指标缓存
 * @param    : period time.Duration    开启定时清理后，多长时间清理一次
 * @Return   : error  错误
 *
 * @Description : 接口会在每天的23：30分总体清理一次缓存；如果有数据的label的value值
 *                非常的多元化，且不可预期，则可以设置一个以小时为单位的清理周期
 */
func (mo *PrometheusMonitor) RegisterMetric(metricType MetricConstType, metricName string, labels []string, timingClearCache bool, period time.Duration) error {
	mo.mtxRW.Lock()
	defer mo.mtxRW.Unlock()

	node := MetricNode{}
	switch metricType {
	case CSTMetricCounter:
		col := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: metricName,
			Help: "help " + metricName,
		},
			labels)
		node.MetricType = metricType
		node.MetricName = metricName
		node.Metric = col
		node.ClearCache.Enable = timingClearCache
		node.ClearCache.Period = period
		node.ClearCache.Point = time.Now()
		mo.metrics[node.MetricName] = &node
	case CSTMetricGauge:
		col := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricName,
			Help: "help " + metricName,
		},
			labels)
		node.MetricType = metricType
		node.MetricName = metricName
		node.Metric = col
		node.ClearCache.Enable = timingClearCache
		node.ClearCache.Period = period
		node.ClearCache.Point = time.Now()
		mo.metrics[node.MetricName] = &node
	case CSTMetricHistogram:
		col := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: metricName,
			Help: "help " + metricName,
		},
			labels)
		node.MetricType = metricType
		node.MetricName = metricName
		node.Metric = col
		node.ClearCache.Enable = timingClearCache
		node.ClearCache.Period = period
		node.ClearCache.Point = time.Now()
		mo.metrics[node.MetricName] = &node
	case CSTMetricSummary:
		col := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: metricName,
			Help: "help " + metricName,
		},
			labels)
		node.MetricType = metricType
		node.MetricName = metricName
		node.Metric = col
		node.ClearCache.Enable = timingClearCache
		node.ClearCache.Period = period
		node.ClearCache.Point = time.Now()
		mo.metrics[node.MetricName] = &node
	default:
		return fmt.Errorf("[prometheus] UnKnow Metric Type For: %d", metricType)
	}

	//fmt.Println("[prometheus] Register Metric Type : ", metricType)

	return nil
}

func (mo *PrometheusMonitor) RegisterMetricNode(node *MetricNode) error {
	mo.mtxRW.Lock()
	defer mo.mtxRW.Unlock()

	if m, ok := mo.metrics[node.MetricName]; ok && m.HaveRegister {
		return fmt.Errorf("%s have registered .", node.MetricName)
	}

	node.HaveRegister = true
	mo.metrics[node.MetricName] = node
	mo.promeRegister.MustRegister(node.Metric)

	return nil
}

func (mo *PrometheusMonitor) AddMetricNode(node *MetricNode) error {
	mo.mtxRW.Lock()
	defer mo.mtxRW.Unlock()

	if _, ok := mo.metrics[node.MetricName]; ok {
		return fmt.Errorf("%s have registered .", node.MetricName)
	}

	node.HaveRegister = false
	mo.metrics[node.MetricName] = node

	return nil
}

/*
 * @Method   : RegisterMetricSource/注册具体的数据指标
 *
 * @param    : metricType MetricConstType 指标类型/使用预定义的常量类型
 * @param    : metricName string 指标名称/与注册指标时标签名称一致
 * @param    : metricCollector prometheus.Collector   指标实例，可自定义指标
 * @Return   : error  错误
 *
 * @Description :
 */
func (mo *PrometheusMonitor) RegisterMetricSource(metricType MetricConstType, metricName string, metricCollector prometheus.Collector) error {
	mo.mtxRW.Lock()
	defer mo.mtxRW.Unlock()

	if metricType >= CSTMetricMax || metricType == CSTMetricNone {
		return fmt.Errorf("[prometheus] UnKnow Metric Type For: %d", metricType)
	}

	node := MetricNode{
		MetricType: metricType,
		MetricName: metricName,
		Metric:     metricCollector,
	}

	mo.metrics[node.MetricName] = &node
	//mo.log.Debug("[prometheus] Register Metric Type : ", metricType)

	return nil
}

/*
 * @Method   : PushMetricValue/上报指标数据
 *
 * @param    : metricName string 指标名称/与注册指标时标签名称一致
 * @param    : val float64       上报数据
 * @param    : labels map[string]string  标签/与注册指标时标签一一匹配(例："host": "127.0.0.1")
 * @Return   : error  错误
 *
 * @Description :
 */
func (mo *PrometheusMonitor) PushMetricValue(metricName string, val float64, labels map[string]string) error {
	mo.mtxRW.RLock()
	defer mo.mtxRW.RUnlock()

	m, ok := mo.metrics[metricName]
	if !ok {
		return fmt.Errorf("[prometheus] Unknow Metric Name[%s] In The Metrics List", metricName)
	}

	switch m.MetricType {
	case CSTMetricCounter:
		m.Metric.(*prometheus.CounterVec).With(labels).Add(val)
	case CSTMetricGauge:
		m.Metric.(*prometheus.GaugeVec).With(labels).Set(val)
	case CSTMetricHistogram:
		m.Metric.(*prometheus.HistogramVec).With(labels).Observe(val)
	case CSTMetricSummary:
		m.Metric.(*prometheus.SummaryVec).With(labels).Observe(val)
	}

	return nil
}

/*
 * @Method   : ResetSingleMetrics/重置metric本地缓存
 *
 * @param    : metricName string 指标名称/与注册指标时标签名称一致
 * @Return   : error  错误
 *
 * @Description :
 */
func (mo *PrometheusMonitor) ResetSingleMetrics(metricName string) error {
	mo.mtxRW.RLock()
	defer mo.mtxRW.RUnlock()

	m, ok := mo.metrics[metricName]
	if !ok {
		return fmt.Errorf("[prometheus] Unknow Metric Name[%s] In The Metrics List", metricName)
	}

	switch m.MetricType {
	case CSTMetricCounter:
		m.Metric.(*prometheus.CounterVec).Reset()
	case CSTMetricGauge:
		m.Metric.(*prometheus.GaugeVec).Reset()
	case CSTMetricHistogram:
		m.Metric.(*prometheus.HistogramVec).Reset()
	case CSTMetricSummary:
		m.Metric.(*prometheus.SummaryVec).Reset()
	}

	return nil
}

func (mo *PrometheusMonitor) FindMetric(metricName string) (*MetricNode, error) {
	mo.mtxRW.RLock()
	defer mo.mtxRW.RUnlock()

	m, ok := mo.metrics[metricName]
	if !ok {
		return nil, fmt.Errorf("[prometheus] Unknow Metric Name[%s] In The Metrics List", metricName)
	}

	return m, nil
}
