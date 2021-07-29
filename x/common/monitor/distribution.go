package monitor

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// DistrMetric is the struct of metric in order module
type DistrMetric struct {
	TotalFee            metrics.Gauge
	FeeToControlledVals metrics.Gauge
	FeeToOtherVals      metrics.Gauge
	FeeToCommunityPool  metrics.Gauge
}

// DefaultDistrMetric returns Metrics build using Prometheus client library if Prometheus is enabled
// Otherwise, it returns no-op Metrics
func DefaultDistrMetric(config *prometheusConfig) *DistrMetric {
	if config.Prometheus {
		return NewDistrMetric()
	}
	return NopDistrMetric()
}

// NewDistrMetric returns a pointer of a new OrderMetric object
func NewDistrMetric(labelsAndValues ...string) *DistrMetric {
	var labels []string
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &DistrMetric{
		TotalFee: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: distrSubSystem,
			Name:      "total_fee",
			Help:      "the total fee that is accumulated every block",
		}, labels).With(labelsAndValues...),
		FeeToControlledVals: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: distrSubSystem,
			Name:      "fee_to_controlled_vals",
			Help:      "the fee that is distributed to controlled vals",
		}, labels).With(labelsAndValues...),
		FeeToOtherVals: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: distrSubSystem,
			Name:      "fee_to_other_vals",
			Help:      "the fee that is distributed to other vals",
		}, labels).With(labelsAndValues...),
		FeeToCommunityPool: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: distrSubSystem,
			Name:      "fee_to_community_pool",
			Help:      "the fee that is distributed to community pool",
		}, labels).With(labelsAndValues...),
	}
}

// NopDistrMetric returns a pointer of a no-op Metrics
func NopDistrMetric() *DistrMetric {
	return &DistrMetric{
		TotalFee:            discard.NewGauge(),
		FeeToControlledVals: discard.NewGauge(),
		FeeToOtherVals:      discard.NewGauge(),
		FeeToCommunityPool:  discard.NewGauge(),
	}
}