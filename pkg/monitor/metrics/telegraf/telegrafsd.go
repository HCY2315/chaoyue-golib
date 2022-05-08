package telegraf

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"ccp-server/framework/monitor/metrics"
)

const (
	// We force flush the statsite metrics after this period of
	// inactivity. Prevents stats from getting stuck in a buffer
	// forever.
	flushInterval = 100 * time.Millisecond

	// statsdMaxLen is the maximum size of a packet
	// to send to statsd
	statsdMaxLen = 1400
)

type TelegrafStatsdSink struct {
	addr        string
	metricQueue chan string
	stopped     chan struct{}
}

// NewStatsdSinkFromURL creates an StatsdSink from a URL. It is used
// (and tested) from NewMetricSinkFromURL.
func NewTelegrafStatsdSinkFromURL(u *url.URL) (metrics.MetricSink, error) {
	return NewTelegrafStatsdSink(u.Host)
}

// NewStatsdSink is used to create a new StatsdSink
func NewTelegrafStatsdSink(addr string) (metrics.MetricSink, error) {
	s := &TelegrafStatsdSink{
		addr:        addr,
		metricQueue: make(chan string, 4096),
		stopped:     make(chan struct{}),
	}
	go s.flushMetrics()
	return s, nil
}

// Close is used to stop flushing to statsd
func (s *TelegrafStatsdSink) Shutdown() {
	close(s.metricQueue)
	<-s.stopped
}

func (s *TelegrafStatsdSink) SetGauge(key []string, val float32) {
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|g\n", flatKey, val))
}

func (s *TelegrafStatsdSink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	flatKey, flatLabel := s.flattenKeyLabels(key, labels)
	s.pushMetric(fmt.Sprintf("%s,%s:%f|g\n", flatKey, flatLabel, val))
}

func (s *TelegrafStatsdSink) EmitKey(key []string, val float32) {
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|kv\n", flatKey, val))
}

func (s *TelegrafStatsdSink) IncrCounter(key []string, val float32) {
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|c\n", flatKey, val))
}

func (s *TelegrafStatsdSink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	flatKey, flatLabel := s.flattenKeyLabels(key, labels)
	s.pushMetric(fmt.Sprintf("%s,%s:%f|c\n", flatKey, flatLabel, val))
}

func (s *TelegrafStatsdSink) AddSample(key []string, val float32) {
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%f|ms\n", flatKey, val))
}

func (s *TelegrafStatsdSink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	flatKey, flatLabel := s.flattenKeyLabels(key, labels)
	s.pushMetric(fmt.Sprintf("%s,%s:%f|ms\n", flatKey, flatLabel, val))
}

func (s *TelegrafStatsdSink) AddSet(key []string, val string) {
	flatKey := s.flattenKey(key)
	s.pushMetric(fmt.Sprintf("%s:%s|s\n", flatKey, val))
}
func (s *TelegrafStatsdSink) AddSetWithLabels(key []string, val string, labels []metrics.Label) {
	flatKey, flatLabel := s.flattenKeyLabels(key, labels)
	s.pushMetric(fmt.Sprintf("%s,%s:%s|s\n", flatKey, flatLabel, val))
}

// Flattens the key for formatting, removes spaces
func (s *TelegrafStatsdSink) flattenKey(parts []string) string {
	joined := strings.Join(parts, ".")
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

func (s *TelegrafStatsdSink) flattenLabel(labels []metrics.Label) string {
	var r string
	for _, label := range labels {
		r = r + label.Name
		r = r + "="
		r = r + label.Value
		r = r + ","
	}
	return strings.Trim(r, ",")
}

// Flattens the key along with labels for formatting, removes spaces
func (s *TelegrafStatsdSink) flattenKeyLabels(parts []string, labels []metrics.Label) (string, string) {
	return s.flattenKey(parts), s.flattenLabel(labels)
}

// Does a non-blocking push to the metrics queue
func (s *TelegrafStatsdSink) pushMetric(m string) {
	select {
	case s.metricQueue <- m:
	default:
	}
}

// Flushes metrics
func (s *TelegrafStatsdSink) flushMetrics() {
	var sock net.Conn
	var err error
	var wait <-chan time.Time
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

CONNECT:
	// Create a buffer
	buf := bytes.NewBuffer(nil)

	// Attempt to connect
	sock, err = net.Dial("udp", s.addr)
	if err != nil {
		log.Printf("[ERR] Error connecting to statsd! Err: %s", err)
		goto WAIT
	}

	for {
		select {
		case metric, ok := <-s.metricQueue:
			// Get a metric from the queue
			if !ok {
				goto QUIT
			}

			// Check if this would overflow the packet size
			if len(metric)+buf.Len() > statsdMaxLen {
				_, err := sock.Write(buf.Bytes())
				buf.Reset()
				if err != nil {
					log.Printf("[ERR] Error writing to statsd! Err: %s", err)
					goto WAIT
				}
			}

			// Append to the buffer
			buf.WriteString(metric)

		case <-ticker.C:
			if buf.Len() == 0 {
				continue
			}

			_, err := sock.Write(buf.Bytes())
			buf.Reset()
			if err != nil {
				log.Printf("[ERR] Error flushing to statsd! Err: %s", err)
				goto WAIT
			}
		}
	}

WAIT:
	// Wait for a while
	wait = time.After(time.Duration(5) * time.Second)
	for {
		select {
		// Dequeue the messages to avoid backlog
		case _, ok := <-s.metricQueue:
			if !ok {
				goto QUIT
			}
		case <-wait:
			goto CONNECT
		}
	}
QUIT:
	s.metricQueue = nil
	close(s.stopped)
}
