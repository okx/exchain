package app

import (
	"github.com/okex/exchain/x/common/monitor"
)

var (
	// init monitor prometheus metrics
	orderMetrics   = monitor.DefaultOrderMetrics(monitor.DefaultPrometheusConfig())
	streamMetrics  = monitor.DefaultStreamMetrics(monitor.DefaultPrometheusConfig())
	stakingMetrics = monitor.DefaultStakingMetric(monitor.DefaultPrometheusConfig())
	distrMetrics   = monitor.DefaultDistrMetric(monitor.DefaultPrometheusConfig())
)
