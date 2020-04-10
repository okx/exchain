package monitor

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// OrderMetric is the struct of metric in order module
type OrderMetric struct {
	FullFilledNum    metrics.Gauge
	PendingNum       metrics.Gauge
	CanceledNum      metrics.Gauge
	ExpiredNum       metrics.Gauge
	PartialFilledNum metrics.Gauge
}

// DefaultOrderMetrics returns Metrics build using Prometheus client library if Prometheus is enabled
// Otherwise, it returns no-op Metrics
func DefaultOrderMetrics(config *prometheusConfig) *OrderMetric {
	if config.Prometheus {
		return NewOrderMetrics()
	}
	return NopOrderMetrics()
}

// NewOrderMetrics returns a pointer of a new OrderMetric object
func NewOrderMetrics(labelsAndValues ...string) *OrderMetric {
	var labels []string
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &OrderMetric{
		FullFilledNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: orderSubSystem,
			Name:      "fullfilled",
			Help:      "the number of fullfilled order",
		}, labels).With(labelsAndValues...),
		PendingNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: orderSubSystem,
			Name:      "pending",
			Help:      "the number of pending order",
		}, labels).With(labelsAndValues...),
		CanceledNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: orderSubSystem,
			Name:      "canceled",
			Help:      "the number of canceled order",
		}, labels).With(labelsAndValues...),
		ExpiredNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: orderSubSystem,
			Name:      "expired",
			Help:      "the number of expired order",
		}, labels).With(labelsAndValues...),
		PartialFilledNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: orderSubSystem,
			Name:      "partial_filled",
			Help:      "the number of partial_filled order",
		}, labels).With(labelsAndValues...),
	}
}

// NopOrderMetrics returns a pointer of a no-op Metrics
func NopOrderMetrics() *OrderMetric {
	return &OrderMetric{
		FullFilledNum:    discard.NewGauge(),
		PendingNum:       discard.NewGauge(),
		CanceledNum:      discard.NewGauge(),
		ExpiredNum:       discard.NewGauge(),
		PartialFilledNum: discard.NewGauge(),
	}
}
