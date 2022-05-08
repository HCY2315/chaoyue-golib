package source

import (
	"fmt"
	"os"

	"github.com/HCY2315/chaoyue-golib/pkg/monitor/metrics"
	"github.com/shirou/gopsutil/process"
)

type ProcessMetrics struct {
	proc                *process.Process
	metricNamePrefix    string
	metricDefaultLabels []metrics.Label
}

func NewProcessMetrics(pid int, metricNamePrefix string, metricDefaultLabels []metrics.Label) (*ProcessMetrics, error) {
	proc, errNewProc := process.NewProcess(int32(pid))
	if errNewProc != nil {
		return nil, fmt.Errorf("new process failed:%w", errNewProc)
	}
	return &ProcessMetrics{
		proc:                proc,
		metricNamePrefix:    metricNamePrefix,
		metricDefaultLabels: metricDefaultLabels,
	}, nil
}

func NewProcessMetricsForSelf(namePrefix string, labels ...metrics.Label) (*ProcessMetrics, error) {
	return NewProcessMetrics(os.Getpid(), namePrefix, labels)
}

func (p *ProcessMetrics) GetMetrics() []MetricRecords {
	mi, _ := p.proc.MemoryInfo()
	memUsedPercent, _ := p.proc.MemoryPercent()
	memRSS := mi.RSS
	memStack := mi.Stack

	memRecords := NewMetricRecordsWithPrefix(p.metricNamePrefix, MetricTypeGauge, p.metricDefaultLabels...)
	memRecords.AddMetric(float64(memUsedPercent), "mem", "percent").
		AddMetric(float64(memRSS), "mem", "rss").
		AddMetric(float64(memStack), "mem", "stack")

	ti, _ := p.proc.Times()
	cpuUser := ti.User
	cpuSystem := ti.System
	cpuIOWait := ti.Iowait
	cpuPercent, _ := p.proc.Percent(0)
	cpuRecords := NewMetricRecordsWithPrefix(p.metricNamePrefix, MetricTypeGauge, p.metricDefaultLabels...)
	cpuRecords.AddMetric(cpuPercent, "cpu", "percent").
		AddMetric(cpuUser, "cpu", "user").
		AddMetric(cpuSystem, "cpu", "sys").
		AddMetric(cpuIOWait, "cpu", "iowait").
		AddMetric(cpuPercent, "cpu", "percent")

	ioCounter, _ := p.proc.IOCounters()
	ioRecords := NewMetricRecordsWithPrefix(p.metricNamePrefix, MetricTypeGauge, p.metricDefaultLabels...)
	ioRecords.AddMetric(float64(ioCounter.ReadBytes), "io", "read", "bytes").
		AddMetric(float64(ioCounter.WriteBytes), "io", "write", "bytes").
		AddMetric(float64(ioCounter.ReadCount), "io", "read", "cnt").
		AddMetric(float64(ioCounter.WriteCount), "io", "write", "cnt")

	otherRecord := NewMetricRecordsWithPrefix(p.metricNamePrefix, MetricTypeGauge, p.metricDefaultLabels...)
	dfCnt, _ := p.proc.NumFDs()
	threadCnt, _ := p.proc.NumThreads()
	otherRecord.AddMetric(float64(dfCnt), "fd", "cnt").
		AddMetric(float64(threadCnt), "thread", "cnt")
	return []MetricRecords{
		*cpuRecords,
		*memRecords,
		*ioRecords,
		*otherRecord,
	}
}
