package monitor

// const
const (
	XNameSpace       = xNameSpace
	xNameSpace       = "x"
	orderSubSystem   = "order"
	stakingSubSystem = "staking"
	streamSubSystem  = "stream"
	portSubSystem    = "port"
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
