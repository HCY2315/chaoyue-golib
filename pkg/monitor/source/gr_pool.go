package source

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/monitor/metrics"
)

type GrPoolMetrics struct {
	prefixName    string
	p             RunningCapable
	defaultLabels []metrics.Label
}

type RunningCapable interface {
	Cap() int
	Running() int
}

func NewGrPoolMetrics(prefixName string, p RunningCapable, defaultLabels ...metrics.Label) *GrPoolMetrics {
	return &GrPoolMetrics{prefixName: prefixName, p: p, defaultLabels: defaultLabels}
}

func (g *GrPoolMetrics) GetMetrics() []MetricRecords {
	capacity := g.p.Cap()
	running := g.p.Running()
	record := NewMetricRecordsWithPrefix(g.prefixName, MetricTypeSample, g.defaultLabels...)
	record.AddMetric(float64(capacity), "grpool", "cap").
		AddMetric(float64(running), "grpool", "running")
	return []MetricRecords{*record}
}
