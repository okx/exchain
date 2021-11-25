package config

type IDynamicConfig interface {
	GetIavlCacheSize() int
}

var DynamicConfig IDynamicConfig

func SetDynamicConfig(c IDynamicConfig) {
	DynamicConfig = c
}
