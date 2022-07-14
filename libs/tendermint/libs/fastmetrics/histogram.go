package fastmetrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Histogram implements Histogram via a Prometheus HistogramVec. The difference
// between a Histogram and a Summary is that Histograms require predefined
// quantile buckets, and can be statistically aggregated.
type Histogram struct {
	hv  *prometheus.HistogramVec
	lvs LabelValues

	labels prometheus.Labels
	o      prometheus.Observer
}

// NewHistogramFrom constructs and registers a Prometheus HistogramVec,
// and returns a usable Histogram object.
func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *Histogram {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(hv)
	return NewHistogram(hv)
}

// NewHistogram wraps the HistogramVec and returns a usable Histogram object.
func NewHistogram(hv *prometheus.HistogramVec) *Histogram {
	his := &Histogram{
		hv: hv,
	}
	his.labels = makeLabels(his.lvs...)
	his.o, _ = his.hv.GetMetricWith(his.labels)
	return his
}

// With implements Histogram.
func (h *Histogram) With(labelValues ...string) metrics.Histogram {
	his := &Histogram{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
	his.labels = makeLabels(his.lvs...)
	his.o, _ = his.hv.GetMetricWith(his.labels)
	return his
}

// Observe implements Histogram.
func (h *Histogram) Observe(value float64) {
	if h.o != nil {
		h.o.Observe(value)
	} else {
		h.hv.With(h.labels).Observe(value)
	}
}
