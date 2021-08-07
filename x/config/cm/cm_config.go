package cm

type CmConfig interface {
	GetMaxOpen() uint64
}

var Config CmConfig


func SetConfig(c CmConfig) {
	Config = c
}