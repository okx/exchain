package monitor

// const
const (
	xNameSpace       = "x"
	orderSubSystem   = "order"
	stakingSubSystem = "staking"
	streamSubSystem  = "stream"
)

type prometheusConfig struct {
	// when true, Prometheus metrics are served under /metrics on PrometheusListenAddr
	Prometheus bool
}

// DefaultPrometheusConfig returns a default PrometheusConfig pointer
func DefaultPrometheusConfig() *prometheusConfig {
	return &prometheusConfig{
		Prometheus: true,
	}
}
