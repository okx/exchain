package monitor

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// StakingMetric is the struct of metric in order module
type StakingMetric struct {
	AllValidatorsShare                         metrics.Gauge
	ControlledValidatorsShare                  metrics.Gauge
	ControlledValidatorsShareRatio             metrics.Gauge
	AllValidatorsAndCandidateShare             metrics.Gauge
	ControlledValidatorsAndCandidateShare      metrics.Gauge
	ControlledValidatorsAndCandidateShareRatio metrics.Gauge
	OfficialValidatorStakingOKT                metrics.Gauge
	OfficialDelegatorStakingOKT                metrics.Gauge
	OfficialValidatorOutstandingOKT            metrics.Gauge
	CommunityValidatorStakingOKT               metrics.Gauge
	CommunityDelegatorStakingOKT               metrics.Gauge
	CommunityValidatorOutstandingOKT           metrics.Gauge
	TotalStakingOKT                            metrics.Gauge
	TotalSupplyOKT                             metrics.Gauge
}

// DefaultOrderMetrics returns Metrics build using Prometheus client library if Prometheus is enabled
// Otherwise, it returns no-op Metrics
func DefaultStakingMetric(config *prometheusConfig) *StakingMetric {
	if config.Prometheus {
		return NewStakingMetric()
	}
	return NopStakingMetric()
}

// NewOrderMetrics returns a pointer of a new OrderMetric object
func NewStakingMetric(labelsAndValues ...string) *StakingMetric {
	var labels []string
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &StakingMetric{
		AllValidatorsShare: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "all_validators_share",
			Help:      "the total share of all validators",
		}, labels).With(labelsAndValues...),
		ControlledValidatorsShare: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "controlled_validators_share",
			Help:      "the total share of all contraolled validators",
		}, labels).With(labelsAndValues...),
		ControlledValidatorsShareRatio: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "controlled_validators_share_ratio",
			Help:      "the ratio of all contraolled validators share to all validators share",
		}, labels).With(labelsAndValues...),
		AllValidatorsAndCandidateShare: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "all_validators_and_candidate_share",
			Help:      "the total share of all validators and candidate",
		}, labels).With(labelsAndValues...),
		ControlledValidatorsAndCandidateShare: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "controlled_validators_and_candidate_share",
			Help:      "the total share of all contraolled validators and candidate",
		}, labels).With(labelsAndValues...),
		ControlledValidatorsAndCandidateShareRatio: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "controlled_validators_and_candidate_share_ratio",
			Help:      "the ratio of all contraolled validators share to all validators share and candidate",
		}, labels).With(labelsAndValues...),
		OfficialValidatorStakingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "official_validator_staking_okt",
			Help:      "amount of staking okt to create validators official",
		}, labels).With(labelsAndValues...),
		OfficialDelegatorStakingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "official_delegator_staking_okt",
			Help:      "amount of staking okt for delegator official",
		}, labels).With(labelsAndValues...),
		OfficialValidatorOutstandingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "official_validator_outstanding_okt",
			Help:      "Not taken out rewards for validator official",
		}, labels).With(labelsAndValues...),
		CommunityValidatorStakingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "community_validator_staking_okt",
			Help:      "amount of staking okt to create validators community",
		}, labels).With(labelsAndValues...),
		CommunityDelegatorStakingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "community_delegator_staking_okt",
			Help:      "amount of staking okt for delegator community",
		}, labels).With(labelsAndValues...),
		TotalStakingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "total_staking_okt",
			Help:      "total amount of staking okt",
		}, labels).With(labelsAndValues...),
		TotalSupplyOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "total_supply_okt",
			Help:      "total supply okt",
		}, labels).With(labelsAndValues...),
		CommunityValidatorOutstandingOKT: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: xNameSpace,
			Subsystem: stakingSubSystem,
			Name:      "community_validator_outstanding_okt",
			Help:      "Not taken out rewards for validator community",
		}, labels).With(labelsAndValues...),
	}
}

// NopStakingMetric returns a pointer of a no-op Metrics
func NopStakingMetric() *StakingMetric {
	return &StakingMetric{
		AllValidatorsShare:                         discard.NewGauge(),
		ControlledValidatorsShare:                  discard.NewGauge(),
		ControlledValidatorsShareRatio:             discard.NewGauge(),
		AllValidatorsAndCandidateShare:             discard.NewGauge(),
		ControlledValidatorsAndCandidateShare:      discard.NewGauge(),
		ControlledValidatorsAndCandidateShareRatio: discard.NewGauge(),
		OfficialValidatorStakingOKT:                discard.NewGauge(),
		OfficialDelegatorStakingOKT:                discard.NewGauge(),
		OfficialValidatorOutstandingOKT:            discard.NewGauge(),
		CommunityValidatorStakingOKT:               discard.NewGauge(),
		CommunityDelegatorStakingOKT:               discard.NewGauge(),
		CommunityValidatorOutstandingOKT:           discard.NewGauge(),
		TotalStakingOKT:                            discard.NewGauge(),
		TotalSupplyOKT:                             discard.NewGauge(),
	}
}
