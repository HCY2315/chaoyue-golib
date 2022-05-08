/*
 * @File   : metrics
 * @Author : huangbin
 *
 * @Created on 2020/11/3 3:10 下午
 * @Project : framework
 * @Software: GoLand
 * @Description  :
 */

package prometheus

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"ccp-server/framework/common"
	"ccp-server/framework/monitor/metrics"
	"ccp-server/framework/utils"
)

type PrometheusSink struct {
	router    string
	addr      string //127.0.0.1:1250
	prome     *PrometheusMonitor
	mtxLock   [CSTMetricMax]sync.Mutex
	mtxCCache sync.Mutex
	r         *mux.Router
	logger    common.ILogger

	stopMonitorChan chan struct{}
}

func NewPrometheusSink(log common.ILogger, opts ...Option) metrics.MetricSink {
	options := newOptions(opts...)
	if options.Id == "" {
		options.Id = options.Name + "-" + uuid.New().String()
	}
	r := mux.NewRouter()
	ins := &PrometheusSink{
		router:          options.PromePath,
		addr:            fmt.Sprintf(":%d", options.Port),
		prome:           NewPrometheusMonitor(options.PromePath),
		r:               r,
		logger:          log,
		stopMonitorChan: make(chan struct{}),
	}
	prometheus.DefaultGatherer = ins.prome.GetRegistry()
	prometheus.DefaultRegisterer = ins.prome.GetRegistry()
	// 设置服务发现
	switch options.Registry {
	case "consul":
		localIPList, _ := utils.LocalIPv4s()
		localIP := localIPList[0]
		nodeInfo := Node{
			Id:      options.Id,
			Name:    options.Name,
			Address: localIP,
			Port:    options.Port,
			Checks: HealthCheckInfo{
				DeregisterCriticalServiceAfter: "30s",
				Http:                           fmt.Sprintf("http://%s:%d/%sconsuld", localIP, options.Port, options.PromePath),
				Interval:                       "5s",
			},
		}
		ins.logger.Debugf("register node info %v in consul for prometheus", nodeInfo)
		consulAgentLen := len(options.RegistryAddress)
		consulAgentIdx := rand.Intn(consulAgentLen)
		consulClient := NewConsulRegistry(nodeInfo, options.RegistryAddress[consulAgentIdx])
		r.HandleFunc(options.PromePath+"consuld", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
		})
		_, _ = consulClient.RegisterService()
	}

	_ = ins.Start()
	return ins
}

func (p *PrometheusSink) SetGauge(key []string, val float32) {
	pk := p.flattenKey(key)
	var err error = nil
	var metric *MetricNode = nil

	fnRegister := func() {
		p.mtxLock[CSTMetricGauge].Lock()
		defer p.mtxLock[CSTMetricGauge].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
			fmt.Println("SetGauge: ", pk)
		}

		col := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: pk,
			Help: "Gauge " + pk,
		})
		metric.MetricType = CSTMetricGauge
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}
	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(prometheus.Gauge).Set(float64(val))
}

func (p *PrometheusSink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	pk := p.flattenKey(key)
	labs := p.flattenLabelCell(labels)
	var err error = nil
	var metric *MetricNode = nil

	fnRegister := func() {
		p.mtxLock[CSTMetricGaugeVec].Lock()
		defer p.mtxLock[CSTMetricGaugeVec].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
		}

		names := p.flattenLabelKeys(labels)
		col := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: pk,
			Help: "GaugeVec " + pk,
		},
			names)
		metric.MetricType = CSTMetricGaugeVec
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}

	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(*prometheus.GaugeVec).With(labs).Set(float64(val))
}

func (p *PrometheusSink) IncrCounter(key []string, val float32) {
	pk := p.flattenKey(key)
	var err error = nil
	var metric *MetricNode = nil

	fnRegister := func() {
		p.mtxLock[CSTMetricCounter].Lock()
		defer p.mtxLock[CSTMetricCounter].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
			fmt.Println("IncrCounter: ", pk)
		}

		col := prometheus.NewCounter(prometheus.CounterOpts{
			Name: pk,
			Help: "Counter " + pk,
		})
		metric.MetricType = CSTMetricCounter
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}

	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(prometheus.Counter).Add(float64(val))
}

func (p *PrometheusSink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	pk := p.flattenKey(key)
	labs := p.flattenLabelCell(labels)
	var err error = nil
	var metric *MetricNode = nil

	fnRegister := func() {
		p.mtxLock[CSTMetricCounterVec].Lock()
		defer p.mtxLock[CSTMetricCounterVec].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
		}

		names := p.flattenLabelKeys(labels)
		col := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: pk,
			Help: "CounterVec " + pk,
		},
			names)
		metric.MetricType = CSTMetricCounterVec
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}

	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(*prometheus.CounterVec).With(labs).Add(float64(val))
}

func (p *PrometheusSink) AddSummary(key []string, val float32) {
	pk := p.flattenKey(key)
	var err error = nil
	var metric *MetricNode = nil

	fnRegister := func() {
		p.mtxLock[CSTMetricSummary].Lock()
		defer p.mtxLock[CSTMetricSummary].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
			fmt.Println("AddSummary: ", pk)
		}

		col := prometheus.NewSummary(prometheus.SummaryOpts{
			Name: pk,
			Help: "Summary " + pk,
		})
		metric.MetricType = CSTMetricSummary
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}

	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(prometheus.Summary).Observe(float64(val))
}

func (p *PrometheusSink) AddSummaryWithLabels(key []string, val float32, labels []metrics.Label) {
	pk := p.flattenKey(key)
	labs := p.flattenLabelCell(labels)
	var metric *MetricNode = nil
	var err error

	fnRegister := func() {
		p.mtxLock[CSTMetricSummaryVec].Lock()
		defer p.mtxLock[CSTMetricSummaryVec].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
		}

		names := p.flattenLabelKeys(labels)
		col := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: pk,
			Help: "SummaryVec " + pk,
		},
			names)
		metric.MetricType = CSTMetricSummaryVec
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}

	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(*prometheus.SummaryVec).With(labs).Observe(float64(val))
}

func (p *PrometheusSink) AddSample(key []string, val float32) {
	p.AddHistogram(key, val)
}
func (p *PrometheusSink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	p.AddHistogramWithLabels(key, val, labels)
}

func (p *PrometheusSink) AddHistogram(key []string, val float32) {
	pk := p.flattenKey(key)
	var metric *MetricNode = nil
	var err error

	fnRegister := func() {
		p.mtxLock[CSTMetricHistogram].Lock()
		defer p.mtxLock[CSTMetricHistogram].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
			fmt.Println("AddHistogram: ", pk)
		}

		col := prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: pk,
			Help: "Histogram " + pk,
		})

		metric.MetricType = CSTMetricHistogram
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}
	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(prometheus.Histogram).Observe(float64(val))
}

func (p *PrometheusSink) AddHistogramWithLabels(key []string, val float32, labels []metrics.Label) {
	pk := p.flattenKey(key)
	labs := p.flattenLabelCell(labels)

	var metric *MetricNode = nil
	var err error
	fnRegister := func() {
		p.mtxLock[CSTMetricHistogramVec].Lock()
		defer p.mtxLock[CSTMetricHistogramVec].Unlock()

		if metric != nil {
			if metric.HaveRegister {
				return
			}
		} else {
			metric = &MetricNode{}
		}

		names := p.flattenLabelKeys(labels)
		col := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: pk,
			Help: "HistogramVec " + pk,
		},
			names)
		metric.MetricType = CSTMetricHistogramVec
		metric.MetricName = pk
		metric.Metric = col

		_ = p.prome.RegisterMetricNode(metric)
	}
	metric, err = p.prome.FindMetric(pk)
	if err != nil {
		fnRegister()
	}

	if !metric.HaveRegister {
		fnRegister()
	}

	metric.Metric.(*prometheus.HistogramVec).With(labs).Observe(float64(val))
}

/*
 * @Method   : SetMetricCacheClearPeriod/设置指标缓存清楚策略
 *
 * @param    : key []string  指标名字
 * @param    : period time.Duration 指标缓存清除周期/time.second*1
 * @Return   : error  错误
 *
 * @Description : 接口会在每天的23：30分总体清理一次缓存；如果有数据的label的value值
 *                非常的多元化，且不可预期，则可以设置一个以小时为单位的清理周期
 */
func (p *PrometheusSink) SetMetricCacheClearPeriod(key []string, period time.Duration) error {
	pk := p.flattenKey(key)
	metric, err := p.prome.FindMetric(pk)
	if err != nil {
		func() {
			p.mtxCCache.Lock()
			defer p.mtxCCache.Unlock()

			node := &MetricNode{}
			node.MetricName = pk
			node.ClearCache.Enable = true
			node.ClearCache.Period = period
			node.ClearCache.Point = time.Now()

			metric = node
			p.prome.AddMetricNode(metric)
		}()
		return nil
	}

	p.mtxCCache.Lock()
	defer p.mtxCCache.Unlock()

	if !metric.ClearCache.Enable {
		metric.ClearCache.Enable = true
		metric.ClearCache.Period = period
		metric.ClearCache.Point = time.Now()

		return nil
	}

	return fmt.Errorf("%s have setted cache clear period .", pk)
}

func (p *PrometheusSink) Run() error {
	p.prome.RegisterHttpHandler(p.r)
	err := http.ListenAndServe(p.addr, nil)
	if err != nil {
		return fmt.Errorf("http.ListenAndServ: %s", err.Error())
	}

	return nil
}

func (p *PrometheusSink) Start() error {
	p.prome.RegisterHttpHandler(p.r)
	go func() {
		for {
			select {
			case <-p.stopMonitorChan:
				return
			default:
				err := http.ListenAndServe(p.addr, p.r)
				if err != nil {
					p.logger.Warnf("[Prometheus] init http server failed: %v", err)
					time.Sleep(10 * time.Millisecond)
				}
			}
		}
	}()
	return nil
}

// Flattens the key for formatting, removes spaces
func (p *PrometheusSink) flattenKey(parts []string) string {
	joined := strings.Join(parts, "_")
	return strings.Map(func(r rune) rune {
		switch r {
		case ':':
			fallthrough
		case ' ':
			return '_'
		default:
			return r
		}
	}, joined)
}

func (p *PrometheusSink) flattenLabelKeys(labels []metrics.Label) []string {
	var keys []string

	for _, v := range labels {
		keys = append(keys, v.Name)
	}

	return keys
}

func (p *PrometheusSink) flattenLabelCell(labels []metrics.Label) map[string]string {
	ls := make(map[string]string)

	for _, v := range labels {
		ls[v.Name] = v.Value
	}

	return ls
}

func (p *PrometheusSink) EmitKey(key []string, val float32)                                 {}
func (p *PrometheusSink) AddSet(key []string, val string)                                   {}
func (p *PrometheusSink) AddSetWithLabels(key []string, val string, labels []metrics.Label) {}
func (p *PrometheusSink) Shutdown() {
	close(p.stopMonitorChan)
}
