package monitor

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// StreamMetrics is the struct of metric in stream module
type StreamMetrics struct {
	CacheSize metrics.Gauge
}

// DefaultStreamMetrics returns Metrics build using Prometheus client library if Prometheus is enabled
// Otherwise, it returns no-op Metrics
func DefaultStreamMetrics(config *prometheusConfig) *StreamMetrics {
	if config.Prometheus {
		return NewStreamMetrics()
	}
	return NopStreamMetrics()
}

// NewStreamMetrics returns a pointer of a new StreamMetrics object
func NewStreamMetrics(labelsAndValues ...string) *StreamMetrics {
	var labels []string
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &StreamMetrics{
		CacheSize: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: streamSubSystem,
			Name:      "cache_size",
			Help:      "the excuting cache111 queue size in stream module.",
		}, labels).With(labelsAndValues...),
	}
}

// NopStreamMetrics returns a pointer of no-op Metrics
func NopStreamMetrics() *StreamMetrics {
	return &StreamMetrics{
		//PulsarSendNum: discard.NewGauge(),
		CacheSize: discard.NewGauge(),
	}
}
