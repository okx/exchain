package config

type IDynamicConfig interface {
	GetIavlCacheSize() int
	GetIavlFSCacheSize() int64
}

var DynamicConfig IDynamicConfig = MockDynamicConfig{}

func SetDynamicConfig(c IDynamicConfig) {
	DynamicConfig = c
}

type MockDynamicConfig struct {
}

func (d MockDynamicConfig) GetIavlCacheSize() int {
	return 10000
}

func (d MockDynamicConfig) GetIavlFSCacheSize() int64 {
	return 10000
}
