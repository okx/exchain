package oec

import "sync"


type OecConfig struct {
	tpb uint64
	maxOpen uint64
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


func (c *OecConfig) GetTpb() uint64 {
	return c.tpb
}

func (c *OecConfig) SetTpb(tpb uint64) {
	c.tpb = tpb
}

func (c *OecConfig) GetMaxOpen() uint64 {
	return c.maxOpen
}

func (c *OecConfig) SetMaxOpen(m uint64) {
	c.maxOpen = m
}

