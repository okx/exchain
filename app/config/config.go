package config

import "sync"


type OecConfig struct {
	tpb int64
	maxOpen int64
}

func NewOecConfig() *OecConfig {
	return &OecConfig{
		tpb: 300,
		maxOpen: 5000,
	}
}

var oecConfig *OecConfig
var once sync.Once

func GetOecConfig() *OecConfig {
	once.Do(func() {
		oecConfig = NewOecConfig()
	})
	return oecConfig
}


func (c *OecConfig) GetTpb() int64 {
	return c.tpb
}

func (c *OecConfig) SetTpb(tpb int64) {
	c.tpb = tpb
}

func (c *OecConfig) GetMaxOpen() int64 {
	return c.maxOpen
}

func (c *OecConfig) SetMaxOpen(m int64) {
	c.maxOpen = m
}

