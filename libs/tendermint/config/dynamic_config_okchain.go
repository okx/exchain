package config

type IDynamicConfig interface {
	GetMempoolRecheck() bool
	GetMempoolForceRecheckGap() int64
	GetMempoolSize() int
	GetMaxTxNumPerBlock() int64
	GetMaxGasUsedPerBlock() int64
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
