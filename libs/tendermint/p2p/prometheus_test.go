package p2p

import (
	"testing"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func TestNewCounterFrom(t *testing.T) {
	labelsAndValues := []string{"chain_id", "exchain-66"}
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	counter := NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "namespace",
		Subsystem: MetricsSubsystem,
		Name:      "peer_send_bytes_total",
		Help:      "Number of bytes sent to a given peer.",
	}, append(labels, "peer_id", "chID")).With(labelsAndValues...)
	labels = []string{
		"peer_id", string("id"),
		"chID", getChIdStr(123),
	}
	counter.With(labels...).Add(123)
}
