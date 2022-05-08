package metrics

// The MetricSink interface is used to transmit metrics information
// to an external system
type MetricSink interface {
	// A Gauge should retain the last value it is set to
	SetGauge(key []string, val float32)
	SetGaugeWithLabels(key []string, val float32, labels []Label)

	// Should emit a Key/Value pair for each call
	EmitKey(key []string, val float32)

	// Counters should accumulate values
	IncrCounter(key []string, val float32)
	IncrCounterWithLabels(key []string, val float32, labels []Label)

	// Samples are for timing information, where quantiles are used
	AddSample(key []string, val float32)
	AddSampleWithLabels(key []string, val float32, labels []Label)

	AddSet(key []string, val string)
	AddSetWithLabels(key []string, val string, labels []Label)

	Shutdown()
}

// BlackholeSink is used to just blackhole messages
type BlackholeSink struct{}

func (*BlackholeSink) SetGauge(key []string, val float32)                              {}
func (*BlackholeSink) SetGaugeWithLabels(key []string, val float32, labels []Label)    {}
func (*BlackholeSink) EmitKey(key []string, val float32)                               {}
func (*BlackholeSink) IncrCounter(key []string, val float32)                           {}
func (*BlackholeSink) IncrCounterWithLabels(key []string, val float32, labels []Label) {}
func (*BlackholeSink) AddSample(key []string, val float32)                             {}
func (*BlackholeSink) AddSampleWithLabels(key []string, val float32, labels []Label)   {}
func (*BlackholeSink) AddSet(key []string, val string)                                 {}
func (*BlackholeSink) AddSetWithLabels(key []string, val string, labels []Label)       {}
func (*BlackholeSink) Shutdown()                                                       {}
