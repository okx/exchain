package config

const (
	DefaultCommitGapHeight = 100
)

type IDynamicConfig interface {
	GetIavlCacheSize() int
	GetIavlFSCacheSize() int64
	GetCommitGapHeight() int64
	SetCommitGapHeight(gap int64)
}

var DynamicConfig IDynamicConfig = MockDynamicConfig{commitGapHeight: DefaultCommitGapHeight}

func SetDynamicConfig(c IDynamicConfig) {
	DynamicConfig = c
}

type MockDynamicConfig struct {
	commitGapHeight int64
}

func (d MockDynamicConfig) GetIavlCacheSize() int {
	return 10000
}

func (d MockDynamicConfig) GetIavlFSCacheSize() int64 {
	return 10000
}

func (d MockDynamicConfig) GetCommitGapHeight() int64 {
	return d.commitGapHeight
}

func (d MockDynamicConfig) SetCommitGapHeight(gap int64) {
	d.commitGapHeight = gap
}
