package source

import (
	"runtime"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/monitor/metrics"
)

type GolangMetricSource struct {
	prefix        string
	defaultLabels []metrics.Label
}

func NewGolangMetricSource(prefix string, defaultLabels ...metrics.Label) *GolangMetricSource {
	return &GolangMetricSource{prefix: prefix, defaultLabels: defaultLabels}
}

func (g *GolangMetricSource) GetMetrics() []MetricRecords {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	goRoutineNum := runtime.NumGoroutine()
	//TODO: block & mutex(need runtime.SetBlockProfileRate(1) runtime.SetMutexProfileFraction(1))
	//mutex,block,thread
	record := NewMetricRecordsWithPrefix(g.prefix, MetricTypeGauge)
	record.AddMetric(float64(memStat.Alloc), "runtime", "heap", "alloc").
		AddMetric(float64(memStat.HeapObjects), "runtime", "heap", "objects").
		AddMetric(float64(memStat.HeapInuse), "runtime", "heap", "inuse").
		AddMetric(float64(memStat.NumGC), "runtime", "gc", "cnt").
		AddMetric(float64(memStat.StackInuse), "runtime", "stack", "inuse").
		AddMetric(float64(goRoutineNum), "runtime", "goroutine", "cnt")
	return []MetricRecords{*record}
}
