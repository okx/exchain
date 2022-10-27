package fastmetrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge implements Gauge, via a Prometheus GaugeVec.
type Gauge struct {
	gv  *prometheus.GaugeVec
	lvs LabelValues

	labels prometheus.Labels
	g      prometheus.Gauge
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
	gauge := &Gauge{
		gv: gv,
	}

	gauge.labels = makeLabels(gauge.lvs...)
	gauge.g, _ = gauge.gv.GetMetricWith(gauge.labels)
	return gauge
}

// With implements Gauge.
func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	gauge := &Gauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
	gauge.labels = makeLabels(gauge.lvs...)
	gauge.g, _ = gauge.gv.GetMetricWith(gauge.labels)
	return gauge
}

// Set implements Gauge.
func (g *Gauge) Set(value float64) {
	if g.g != nil {
		g.g.Set(value)
	} else {
		g.gv.With(g.labels).Set(value)
	}
}

// Add is supported by Prometheus GaugeVecs.
func (g *Gauge) Add(delta float64) {
	if g.g != nil {
		g.g.Add(delta)
	} else {
		g.gv.With(g.labels).Add(delta)
	}
}
