package config

import "time"

type IDynamicConfig interface {
	GetMempoolRecheck() bool
	GetMempoolForceRecheckGap() int64
	GetMempoolSize() int
	GetMaxTxNumPerBlock() int64
	GetMaxGasUsedPerBlock() int64
	GetMempoolFlush() bool
	GetNodeKeyWhitelist() []string
	GetCsTimeoutPropose() time.Duration
	GetCsTimeoutProposeDelta() time.Duration
	GetCsTimeoutPrevote() time.Duration
	GetCsTimeoutPrevoteDelta() time.Duration
	GetCsTimeoutPrecommit() time.Duration
	GetCsTimeoutPrecommitDelta() time.Duration
	GetEnableWtx() bool
}

var DynamicConfig IDynamicConfig = MockDynamicConfig{}

func SetDynamicConfig(c IDynamicConfig) {
	DynamicConfig = c
}

type MockDynamicConfig struct {
}

func (d MockDynamicConfig) GetMempoolRecheck() bool {
	return DefaultMempoolConfig().Recheck
}

func (d MockDynamicConfig) GetMempoolForceRecheckGap() int64 {
	return DefaultMempoolConfig().ForceRecheckGap
}

func (d MockDynamicConfig) GetMempoolSize() int {
	return DefaultMempoolConfig().Size
}

func (d MockDynamicConfig) GetMaxTxNumPerBlock() int64 {
	return DefaultMempoolConfig().MaxTxNumPerBlock
}

func (d MockDynamicConfig) GetMaxGasUsedPerBlock() int64 {
	return DefaultMempoolConfig().MaxGasUsedPerBlock
}

func (d MockDynamicConfig) GetMempoolFlush() bool {
	return false
}

func (d MockDynamicConfig) GetNodeKeyWhitelist() []string {
	return []string{}
}

func (d MockDynamicConfig) GetCsTimeoutPropose() time.Duration {
	return DefaultConsensusConfig().TimeoutPropose
}
func (d MockDynamicConfig) GetCsTimeoutProposeDelta() time.Duration {
	return DefaultConsensusConfig().TimeoutProposeDelta
}
func (d MockDynamicConfig) GetCsTimeoutPrevote() time.Duration {
	return DefaultConsensusConfig().TimeoutPrevote
}
func (d MockDynamicConfig) GetCsTimeoutPrevoteDelta() time.Duration {
	return DefaultConsensusConfig().TimeoutPrecommitDelta
}
func (d MockDynamicConfig) GetCsTimeoutPrecommit() time.Duration {
	return DefaultConsensusConfig().TimeoutPrecommit
}
func (d MockDynamicConfig) GetCsTimeoutPrecommitDelta() time.Duration {
	return DefaultConsensusConfig().TimeoutPrecommitDelta
}

func (d MockDynamicConfig) GetEnableWtx() bool {
	return false
}
