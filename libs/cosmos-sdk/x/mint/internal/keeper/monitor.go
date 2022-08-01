package keeper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// MintMetric is the struct of metric in mint module
type MintMetric struct {
	ToTreasury metrics.Gauge
}

// DefaultMintMetric returns Metrics build using Prometheus client library
func DefaultMintMetric(labelsAndValues ...string) *MintMetric {
	var labels []string
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &MintMetric{
		ToTreasury: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: "x",
			Subsystem: "mint",
			Name:      "to_treasury",
			Help:      "the minted coins allocated to treasury that are accumulated every block",
		}, labels).With(labelsAndValues...),
	}
}
