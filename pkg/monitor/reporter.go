package monitor

import (
	"context"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/task"

	"github.com/HCY2315/chaoyue-golib/pkg/monitor/metrics"

	"github.com/HCY2315/chaoyue-golib/pkg/monitor/source"

	"github.com/google/uuid"
)

type Reporter struct {
	sink       metrics.MetricSink
	sourceList []source.MetricSource
	name       string
}

func NewReporter(sink metrics.MetricSink, name string, sourceList ...source.MetricSource) *Reporter {
	return &Reporter{sink: sink, sourceList: sourceList, name: name}
}

func (r *Reporter) Name() string {
	if r.name != "" {
		return r.name
	}
	return "metric-reporter"
}

func (r *Reporter) Execute(ctx context.Context) error {
	for _, metricSource := range r.sourceList {
		records := metricSource.GetMetrics()
		for _, record := range records {
			r.reportRecord(record)
		}
	}
	return nil
}

func (r *Reporter) reportRecord(record source.MetricRecords) {
	record.Foreach(func(name string, metricType source.MetricType, value float32, labels []metrics.Label) {
		switch metricType {
		case source.MetricTypeGauge:
			r.sink.SetGaugeWithLabels([]string{name}, value, labels)
		case source.MetricTypeCounter:
			r.sink.IncrCounterWithLabels([]string{name}, value, labels)
		case source.MetricTypeSample:
			r.sink.AddSampleWithLabels([]string{name}, value, labels)
		}
	})
}

func NewIntervalMetricReporterWithSourceList(sink metrics.MetricSink,
	reportInterval time.Duration, taskName string, sourceList ...source.MetricSource) (*task.ScheduleItem, error) {
	reporter := NewReporter(sink, taskName, sourceList...)
	taskId := uuid.New().String()
	return &task.ScheduleItem{
		ScheduleItemDesc: task.NewIntervalScheduleItem(reportInterval),
		Task:             reporter,
		TaskId:           taskId,
	}, nil
}

func BuildGrPoolMetricSourceFromMap(metricPrefix string, grPoolMap map[string]source.RunningCapable,
	defaultLabels ...metrics.Label) []source.MetricSource {
	grPoolSourceList := make([]source.MetricSource, 0, len(grPoolMap))
	for name, grPool := range grPoolMap {
		grPoolSourceList = append(grPoolSourceList, source.NewGrPoolMetrics(metricPrefix+"_"+name, grPool, defaultLabels...))
	}
	return grPoolSourceList
}

func BuildProcessMetricSource(metricPrefix string, defaultLabels []metrics.Label) ([]source.MetricSource, error) {
	processMetricSource, err := source.NewProcessMetricsForSelf(metricPrefix, defaultLabels...)
	if err != nil {
		return nil, err
	}
	golangMetric := source.NewGolangMetricSource(metricPrefix, defaultLabels...)
	return []source.MetricSource{
		processMetricSource,
		golangMetric,
	}, nil
}
