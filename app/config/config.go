package config

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/spf13/viper"
)

type OecConfig struct {
	// mempool.recheck
	mempoolRecheck bool
	// mempool.force_recheck_gap
	mempoolForceRecheckGap int64
	// mempool.size
	mempoolSize int

	// log_level
	logLevel string

	// rpc.disable-api
	rpcDisableApi string
	// rpc.rate-limit-api
	rpcRateLimitApi string
	// rpc.rate-limit-burst
	rpcRateLimitBurst int
	// rpc.rate-limit-count
	rpcRateLimitCount int
	// gas-limit-buffer
	gasLimitBuffer uint64
	// enable-dynamic-gp
	enableDynamicGp bool
	// dynamic-gp-weight
	dynamicGpWeight int
}

func NewOecConfig() *OecConfig {
	c := &OecConfig{}
	loadFromConfig(c)
	fmt.Printf("%+v\n", c)
	loadFromApollo(c)
	fmt.Printf("%+v\n", c)

	return c
}

func loadFromConfig(c *OecConfig) {
	c.SetMempoolRecheck(viper.GetBool("mempool.recheck"))
	c.SetMempoolForceRecheckGap(viper.GetInt64("mempool.force_recheck_gap"))
	c.SetMempoolSize(viper.GetInt("mempool.size"))
	c.SetLogLevel(viper.GetString("log_level"))
	c.SetRpcDisableApi(viper.GetString("rpc.disable-api"))
	c.SetRpcRateLimitApi(viper.GetString("rpc.rate-limit-api"))
	c.SetRpcRateLimitBurst(viper.GetInt("rpc.rate-limit-burst"))
	c.SetRpcRateLimitCount(viper.GetInt("rpc.rate-limit-count"))
	c.SetGasLimitBuffer(viper.GetUint64("gas-limit-buffer"))
	c.SetEnableDynamicGp(viper.GetBool("enable-dynamic-gp"))
	c.SetDynamicGpWeight(viper.GetInt("dynamic-gp-weight"))
}

func loadFromApollo(c *OecConfig) {
	client := NewApollo(nil)
	cache := client.GetConfigCache("rpc-node")
	cache.Range(func(key, value interface{}) bool {
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
		case "log_level":
			c.SetLogLevel(v)
		case "rpc.disable-api":
			c.SetRpcDisableApi(v)
		case "rpc.rate-limit-api":
			c.SetRpcRateLimitApi(v)
		case "rpc.rate-limit-burst":
			r, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			c.SetRpcRateLimitBurst(r)
		case "rpc.rate-limit-count":
			r, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			c.SetRpcRateLimitCount(r)
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
		return true
	})
}

var oecConfig *OecConfig
var once sync.Once

func GetOecConfig() *OecConfig {
	once.Do(func() {
		oecConfig = NewOecConfig()
	})
	return oecConfig
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
	c.mempoolSize = value
}

func (c *OecConfig) GetLogLevel() string {
	return c.logLevel
}

func (c *OecConfig) SetLogLevel(value string) {
	c.logLevel = value
}

func (c *OecConfig) GetRpcDisableApi() string {
	return c.rpcDisableApi
}

func (c *OecConfig) SetRpcDisableApi(value string) {
	c.rpcDisableApi = value
}

func (c *OecConfig) GetRpcRateLimitApi() string {
	return c.rpcRateLimitApi
}

func (c *OecConfig) SetRpcRateLimitApi(value string) {
	c.rpcRateLimitApi = value
}

func (c *OecConfig) GetRpcRateLimitBurst() int {
	return c.rpcRateLimitBurst
}

func (c *OecConfig) SetRpcRateLimitBurst(value int) {
	c.rpcRateLimitBurst = value
}

func (c *OecConfig) GetRpcRateLimitCount() int {
	return c.rpcRateLimitCount
}

func (c *OecConfig) SetRpcRateLimitCount(value int) {
	c.rpcRateLimitCount = value
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
	c.dynamicGpWeight = value
}
