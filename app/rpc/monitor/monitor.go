package monitor

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	FlagEnableMonitor = "rpc.enable-monitor"
	MetricsNamespace  = "x"
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this package.
	MetricsSubsystem = "rpc"

	MetricsFieldName   = "apis"
	MetricsMethodLabel = "method"

	MetricsCounterNamePattern   = "%s_%s_count"
	MetricsHistogramNamePattern = "%s_%s_duration"
)

// RpcMetrics ...
type RpcMetrics struct {
	Counter   metrics.Counter
	Histogram metrics.Histogram
}

type Monitor struct {
	method   string
	logger   log.Logger
	lastTime time.Time
	metrics  *RpcMetrics
}

func MakeMonitorMetrics(namespace string) *RpcMetrics {

	return &RpcMetrics{
		Counter: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      fmt.Sprintf(MetricsCounterNamePattern, namespace, MetricsFieldName),
			Help:      fmt.Sprintf("Total request number of %s/%s method.", namespace, MetricsFieldName),
		}, []string{MetricsMethodLabel}),
		Histogram: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      fmt.Sprintf(MetricsHistogramNamePattern, namespace, MetricsFieldName),
			Help:      fmt.Sprintf("Request duration of %s/%s method.", namespace, MetricsFieldName),
			Buckets:   []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.8, 1, 3, 5, 8, 10},
		}, []string{MetricsMethodLabel}),
	}
}

func GetMonitor(method string, logger log.Logger, metrics *RpcMetrics) *Monitor {
	return &Monitor{
		method:  method,
		logger:  logger,
		metrics: metrics,
	}
}

func (m *Monitor) OnBegin() *Monitor {
	m.lastTime = time.Now()

	if m.metrics == nil {
		return m
	}
	m.metrics.Counter.With(MetricsMethodLabel, m.method).Add(1)

	return m
}

func (m *Monitor) OnEnd(args ...interface{}) {
	elapsed := time.Since(m.lastTime).Seconds()
	m.logger.Debug("RPC", MetricsMethodLabel, m.method, "Elapsed", elapsed*1e3, "Params", args)

	if m.metrics == nil {
		return
	}
	m.metrics.Histogram.With(MetricsMethodLabel, m.method).Observe(elapsed)
}
