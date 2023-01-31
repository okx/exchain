package config

import "time"

type IDynamicConfig interface {
	GetMempoolRecheck() bool
	GetMempoolForceRecheckGap() int64
	GetMempoolSize() int
	GetMempoolCacheSize() int
	GetMaxTxNumPerBlock() int64
	GetEnableDeleteMinGPTx() bool
	GetMaxGasUsedPerBlock() int64
	GetEnablePGU() bool
	GetPGUAdjustment() float64
	GetMempoolFlush() bool
	GetNodeKeyWhitelist() []string
	GetMempoolCheckTxCost() bool
	GetSentryAddrs() []string
	GetCsTimeoutPropose() time.Duration
	GetCsTimeoutProposeDelta() time.Duration
	GetCsTimeoutPrevote() time.Duration
	GetCsTimeoutPrevoteDelta() time.Duration
	GetCsTimeoutPrecommit() time.Duration
	GetCsTimeoutPrecommitDelta() time.Duration
	GetCsTimeoutCommit() time.Duration
	GetEnableWtx() bool
	GetDeliverTxsExecuteMode() int
	GetEnableHasBlockPartMsg() bool
	GetCommitGapOffset() int64
	GetIavlAcNoBatch() bool
}

var DynamicConfig IDynamicConfig = MockDynamicConfig{}

func SetDynamicConfig(c IDynamicConfig) {
	DynamicConfig = c
}

type MockDynamicConfig struct {
	enableDeleteMinGPTx bool
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

func (d MockDynamicConfig) GetMempoolCacheSize() int {
	return DefaultMempoolConfig().CacheSize
}

func (d MockDynamicConfig) GetMaxTxNumPerBlock() int64 {
	return DefaultMempoolConfig().MaxTxNumPerBlock
}

func (d MockDynamicConfig) GetMaxGasUsedPerBlock() int64 {
	return DefaultMempoolConfig().MaxGasUsedPerBlock
}

func (d MockDynamicConfig) GetEnablePGU() bool {
	return false
}

func (d MockDynamicConfig) GetPGUAdjustment() float64 {
	return 1
}

func (d MockDynamicConfig) GetMempoolFlush() bool {
	return false
}

func (d MockDynamicConfig) GetNodeKeyWhitelist() []string {
	return []string{}
}

func (d MockDynamicConfig) GetMempoolCheckTxCost() bool {
	return false
}

func (d MockDynamicConfig) GetSentryAddrs() []string {
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
func (d MockDynamicConfig) GetCsTimeoutCommit() time.Duration {
	return DefaultConsensusConfig().TimeoutCommit
}

func (d MockDynamicConfig) GetEnableWtx() bool {
	return false
}
func (d MockDynamicConfig) GetDeliverTxsExecuteMode() int {
	return 0
}

func (d MockDynamicConfig) GetEnableHasBlockPartMsg() bool {
	return false
}

func (d MockDynamicConfig) GetEnableDeleteMinGPTx() bool {
	return d.enableDeleteMinGPTx
}

func (d *MockDynamicConfig) SetEnableDeleteMinGPTx(enable bool) {
	d.enableDeleteMinGPTx = enable
}

func (d MockDynamicConfig) GetCommitGapOffset() int64 {
	return 0
}

func (d MockDynamicConfig) GetIavlAcNoBatch() bool {
	return false
}
