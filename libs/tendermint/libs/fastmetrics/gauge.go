package fastmetrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge implements Gauge, via a Prometheus GaugeVec.
type Gauge struct {
	gv  *prometheus.GaugeVec
	lvs LabelValues
}

// NewGaugeFrom constructs and registers a Prometheus GaugeVec,
// and returns a usable Gauge object.
func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return NewGauge(gv)
}

// NewGauge wraps the GaugeVec and returns a usable Gauge object.
func NewGauge(gv *prometheus.GaugeVec) *Gauge {
	return &Gauge{
		gv: gv,
	}
}

// With implements Gauge.
func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}

// Set implements Gauge.
func (g *Gauge) Set(value float64) {
	g.gv.With(makeLabels(g.lvs...)).Set(value)
}

// Add is supported by Prometheus GaugeVecs.
func (g *Gauge) Add(delta float64) {
	g.gv.With(makeLabels(g.lvs...)).Add(delta)
}
