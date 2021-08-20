package config

import (
	"strconv"
	"sync"

	cmconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
)

type OecConfig struct {
	// mempool.recheck
	mempoolRecheck bool
	// mempool.force_recheck_gap
	mempoolForceRecheckGap int64
	// mempool.size
	mempoolSize int

	// gas-limit-buffer
	gasLimitBuffer uint64
	// enable-dynamic-gp
	enableDynamicGp bool
	// dynamic-gp-weight
	dynamicGpWeight int
}

const FlagEnableDynamic = "config.enable-dynamic"

var oecConfig *OecConfig
var once sync.Once

func GetOecConfig() *OecConfig {
	once.Do(func() {
		oecConfig = NewOecConfig()
	})
	return oecConfig
}

func NewOecConfig() *OecConfig {
	c := &OecConfig{}
	c.loadFromConfig()

	if viper.GetBool(FlagEnableDynamic) {
		loaded := c.loadFromApollo()
		if !loaded {
			panic("failed to connect apollo or no config items in apollo")
		}
	}

	return c
}

func RegisterDynamicConfig() {
	// set the dynamic config
	oecConfig := GetOecConfig()
	cmconfig.SetDynamicConfig(oecConfig)
	tmconfig.SetDynamicConfig(oecConfig)
}

func (c *OecConfig) loadFromConfig() {
	c.SetMempoolRecheck(viper.GetBool("mempool.recheck"))
	c.SetMempoolForceRecheckGap(viper.GetInt64("mempool.force_recheck_gap"))
	c.SetMempoolSize(viper.GetInt("mempool.size"))
	c.SetGasLimitBuffer(viper.GetUint64("gas-limit-buffer"))
	c.SetEnableDynamicGp(viper.GetBool("enable-dynamic-gp"))
	c.SetDynamicGpWeight(viper.GetInt("dynamic-gp-weight"))
}

func (c *OecConfig) loadFromApollo() bool {
	client := NewApolloClient(c)
	return client.LoadConfig()
}

func (c *OecConfig) update(key, value interface{}) {
	k, v := key.(string), value.(string)
	switch k {
	case "mempool.recheck":
		r, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		c.SetMempoolRecheck(r)
	case "mempool.force_recheck_gap":
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		c.SetMempoolForceRecheckGap(r)
	case "mempool.size":
		r, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		c.SetMempoolSize(r)
	case "gas-limit-buffer":
		r, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err)
		}
		c.SetGasLimitBuffer(r)
	case "enable-dynamic-gp":
		r, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		c.SetEnableDynamicGp(r)
	case "dynamic-gp-weight":
		r, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		c.SetDynamicGpWeight(r)
	}
}

func (c *OecConfig) GetMempoolRecheck() bool {
	return c.mempoolRecheck
}

func (c *OecConfig) SetMempoolRecheck(value bool) {
	c.mempoolRecheck = value
}

func (c *OecConfig) GetMempoolForceRecheckGap() int64 {
	return c.mempoolForceRecheckGap
}

func (c *OecConfig) SetMempoolForceRecheckGap(value int64) {
	c.mempoolForceRecheckGap = value
}

func (c *OecConfig) GetMempoolSize() int {
	return c.mempoolSize
}

func (c *OecConfig) SetMempoolSize(value int) {
	if value < 0 {
		return
	}
	c.mempoolSize = value
}

func (c *OecConfig) GetGasLimitBuffer() uint64 {
	return c.gasLimitBuffer
}
func (c *OecConfig) SetGasLimitBuffer(value uint64) {
	c.gasLimitBuffer = value
}

func (c *OecConfig) GetEnableDynamicGp() bool {
	return c.enableDynamicGp
}
func (c *OecConfig) SetEnableDynamicGp(value bool) {
	c.enableDynamicGp = value
}

func (c *OecConfig) GetDynamicGpWeight() int {
	return c.dynamicGpWeight
}
func (c *OecConfig) SetDynamicGpWeight(value int) {
	if value <= 0 {
		value = 1
	} else if value > 100 {
		value = 100
	}
	c.dynamicGpWeight = value
}
