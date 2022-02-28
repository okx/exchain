package p2p

import (
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type LabelValues []string

// With validates the input, and returns a new aggregate labelValues.
func (lvs LabelValues) With(labelValues ...string) LabelValues {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}
	return append(lvs, labelValues...)
}

// Counter implements Counter, via a Prometheus CounterVec.
type Counter struct {
	cv     *prometheus.CounterVec
	lvs    LabelValues
	labels prometheus.Labels
	c      prometheus.Counter
}

// NewCounterFrom constructs and registers a Prometheus CounterVec,
// and returns a usable Counter object.
func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *Counter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return NewCounter(cv)
}

// NewCounter wraps the CounterVec and returns a usable Counter object.
func NewCounter(cv *prometheus.CounterVec) *Counter {
	counter := &Counter{
		cv: cv,
	}
	counter.labels = makeLabels(counter.lvs...)
	counter.c, _ = counter.cv.GetMetricWith(counter.labels)
	return counter
}

// With implements Counter.
func (c *Counter) With(labelValues ...string) metrics.Counter {
	counter := &Counter{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
	counter.labels = makeLabels(counter.lvs...)
	counter.c, _ = counter.cv.GetMetricWith(counter.labels)
	return counter
}

// Add implements Counter.
func (c *Counter) Add(delta float64) {
	if c.c != nil {
		c.c.Add(delta)
	} else {
		c.cv.With(c.labels).Add(delta)
	}
}

func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
