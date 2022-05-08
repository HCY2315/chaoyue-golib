package source

import (
	"strings"

	"git.cestong.com.cn/cecf/cecf-golib/pkg/monitor/metrics"
)

type MetricType uint8

const (
	_ MetricType = iota
	MetricTypeGauge
	MetricTypeCounter
	MetricTypeSample
)

type MetricRecords struct {
	prefix   string
	typ      MetricType
	valueMap map[string]float64
	labels   []metrics.Label
}

func NewMetricRecords(metricType MetricType, labels ...metrics.Label) *MetricRecords {
	return NewMetricRecordsWithPrefix("", metricType, labels...)
}

func NewMetricRecordsWithPrefix(prefix string, metricType MetricType, labels ...metrics.Label) *MetricRecords {
	return &MetricRecords{
		prefix:   prefix,
		typ:      metricType,
		valueMap: make(map[string]float64, 20),
		labels:   labels,
	}
}

const (
	metricNameDelimiter = "_"
)

func (mr *MetricRecords) AddMetric(value float64, names ...string) *MetricRecords {
	mr.valueMap[strings.Join(names, metricNameDelimiter)] = value
	return mr
}

func (mr *MetricRecords) Foreach(f func(name string, tye MetricType, value float32, labels []metrics.Label)) {
	for name, value := range mr.valueMap {
		if mr.prefix != "" {
			name = mr.prefix + metricNameDelimiter + name
		}
		f(name, mr.typ, float32(value), mr.labels)
	}
}

type MetricSource interface {
	GetMetrics() []MetricRecords
}
