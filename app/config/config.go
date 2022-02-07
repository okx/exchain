package config

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	iavlconfig "github.com/okex/exchain/libs/iavl/config"
	tmconfig "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/spf13/viper"
)

type OecConfig struct {
	// mempool.recheck
	mempoolRecheck bool
	// mempool.force_recheck_gap
	mempoolForceRecheckGap int64
	// mempool.size
	mempoolSize int
	// mempool.flush
	mempoolFlush bool
	// mempool.max_tx_num_per_block
	maxTxNumPerBlock int64
	// mempool.max_gas_used_per_block
	maxGasUsedPerBlock int64

	// gas-limit-buffer
	gasLimitBuffer uint64
	// enable-dynamic-gp
	enableDynamicGp bool
	// dynamic-gp-weight
	dynamicGpWeight int

	// consensus.timeout_propose
	csTimeoutPropose time.Duration
	// consensus.timeout_propose_delta
	csTimeoutProposeDelta time.Duration
	// consensus.timeout_prevote
	csTimeoutPrevote time.Duration
	// consensus.timeout_prevote_delta
	csTimeoutPrevoteDelta time.Duration
	// consensus.timeout_precommit
	csTimeoutPrecommit time.Duration
	// consensus.timeout_precommit_delta
	csTimeoutPrecommitDelta time.Duration

	// iavl-cache-size
	iavlCacheSize  int
}

const (
	FlagEnableDynamic = "config.enable-dynamic"

	FlagMempoolRecheck         = "mempool.recheck"
	FlagMempoolForceRecheckGap = "mempool.force_recheck_gap"
	FlagMempoolSize            = "mempool.size"
	FlagMempoolFlush           = "mempool.flush"
	FlagMaxTxNumPerBlock       = "mempool.max_tx_num_per_block"
	FlagMaxGasUsedPerBlock     = "mempool.max_gas_used_per_block"
	FlagGasLimitBuffer         = "gas-limit-buffer"
	FlagEnableDynamicGp        = "enable-dynamic-gp"
	FlagDynamicGpWeight        = "dynamic-gp-weight"

	FlagCsTimeoutPropose        = "consensus.timeout_propose"
	FlagCsTimeoutProposeDelta   = "consensus.timeout_propose_delta"
	FlagCsTimeoutPrevote        = "consensus.timeout_prevote"
	FlagCsTimeoutPrevoteDelta   = "consensus.timeout_prevote_delta"
	FlagCsTimeoutPrecommit      = "consensus.timeout_precommit"
	FlagCsTimeoutPrecommitDelta = "consensus.timeout_precommit_delta"
)

<<<<<<< HEAD
var oecConfig *OecConfig
var once sync.Once
var confLogger log.Logger
=======
var (
	testnetNodeIdWhitelist = []string{}

	mainnetNodeIdWhitelist = []string{}

	oecConfig  *OecConfig
	once       sync.Once
	confLogger log.Logger
)
>>>>>>> 0076fe56 (log end)

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

func RegisterDynamicConfig(logger log.Logger) {
	confLogger = logger
	// set the dynamic config
	oecConfig := GetOecConfig()
	tmconfig.SetDynamicConfig(oecConfig)
	iavlconfig.SetDynamicConfig(oecConfig)
}

func (c *OecConfig) loadFromConfig() {
	c.SetMempoolRecheck(viper.GetBool(FlagMempoolRecheck))
	c.SetMempoolForceRecheckGap(viper.GetInt64(FlagMempoolForceRecheckGap))
	c.SetMempoolSize(viper.GetInt(FlagMempoolSize))
	c.SetMempoolFlush(viper.GetBool(FlagMempoolFlush))
	c.SetMaxTxNumPerBlock(viper.GetInt64(FlagMaxTxNumPerBlock))
	c.SetMaxGasUsedPerBlock(viper.GetInt64(FlagMaxGasUsedPerBlock))
	c.SetGasLimitBuffer(viper.GetUint64(FlagGasLimitBuffer))
	c.SetEnableDynamicGp(viper.GetBool(FlagEnableDynamicGp))
	c.SetDynamicGpWeight(viper.GetInt(FlagDynamicGpWeight))
	c.SetCsTimeoutPropose(viper.GetDuration(FlagCsTimeoutPropose))
	c.SetCsTimeoutProposeDelta(viper.GetDuration(FlagCsTimeoutProposeDelta))
	c.SetCsTimeoutPrevote(viper.GetDuration(FlagCsTimeoutPrevote))
	c.SetCsTimeoutPrevoteDelta(viper.GetDuration(FlagCsTimeoutPrevoteDelta))
	c.SetCsTimeoutPrecommit(viper.GetDuration(FlagCsTimeoutPrecommit))
	c.SetCsTimeoutPrecommitDelta(viper.GetDuration(FlagCsTimeoutPrecommitDelta))
	c.SetIavlCacheSize(viper.GetInt(iavl.FlagIavlCacheSize))
}

func (c *OecConfig) loadFromApollo() bool {
	client := NewApolloClient(c)
	return client.LoadConfig()
}

func (c *OecConfig) format() string {
	return fmt.Sprintf(`OEC config:
	mempool.recheck: %v
	mempool.force_recheck_gap: %d
	mempool.size: %d
	mempool.flush: %v
	mempool.max_tx_num_per_block: %d
	mempool.max_gas_used_per_block: %d

	gas-limit-buffer: %d
	enable-dynamic-gp: %v
	dynamic-gp-weight: %d

	consensus.timeout_propose: %s
	consensus.timeout_propose_delta: %s
	consensus.timeout_prevote: %s
	consensus.timeout_prevote_delta: %s
	consensus.timeout_precommit: %s
	consensus.timeout_precommit_delta: %s
	
	iavl-cache-size: %d`,
		c.GetMempoolRecheck(),
		c.GetMempoolForceRecheckGap(),
		c.GetMempoolSize(),
		c.GetMempoolFlush(),
		c.GetMaxTxNumPerBlock(),
		c.GetMaxGasUsedPerBlock(),
		c.GetGasLimitBuffer(),
		c.GetEnableDynamicGp(),
		c.GetDynamicGpWeight(),
		c.GetCsTimeoutPropose(),
		c.GetCsTimeoutProposeDelta(),
		c.GetCsTimeoutPrevote(),
		c.GetCsTimeoutPrevoteDelta(),
		c.GetCsTimeoutPrecommit(),
		c.GetCsTimeoutPrecommitDelta(),
		c.GetIavlCacheSize(),
	)
}

func (c *OecConfig) update(key, value interface{}) {
	k, v := key.(string), value.(string)
	switch k {
	case FlagMempoolRecheck:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetMempoolRecheck(r)
	case FlagMempoolForceRecheckGap:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetMempoolForceRecheckGap(r)
	case FlagMempoolSize:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetMempoolSize(r)
	case FlagMempoolFlush:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetMempoolFlush(r)
	case FlagMaxTxNumPerBlock:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetMaxTxNumPerBlock(r)
	case FlagMaxGasUsedPerBlock:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetMaxGasUsedPerBlock(r)
	case FlagGasLimitBuffer:
		r, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return
		}
		c.SetGasLimitBuffer(r)
	case FlagEnableDynamicGp:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetEnableDynamicGp(r)
	case FlagDynamicGpWeight:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetDynamicGpWeight(r)
	case FlagCsTimeoutPropose:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPropose(r)
	case FlagCsTimeoutProposeDelta:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutProposeDelta(r)
	case FlagCsTimeoutPrevote:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrevote(r)
	case FlagCsTimeoutPrevoteDelta:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrevoteDelta(r)
	case FlagCsTimeoutPrecommit:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrecommit(r)
	case FlagCsTimeoutPrecommitDelta:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrecommitDelta(r)
	case iavl.FlagIavlCacheSize:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetIavlCacheSize(r)
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
	if value <= 0 {
		return
	}
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

func (c *OecConfig) GetMempoolFlush() bool {
	return c.mempoolFlush
}
func (c *OecConfig) SetMempoolFlush(value bool) {
	c.mempoolFlush = value
}

func (c *OecConfig) GetMaxTxNumPerBlock() int64 {
	return c.maxTxNumPerBlock
}
func (c *OecConfig) SetMaxTxNumPerBlock(value int64) {
	if value < 0 {
		return
	}
	c.maxTxNumPerBlock = value
}

func (c *OecConfig) GetMaxGasUsedPerBlock() int64 {
	return c.maxGasUsedPerBlock
}
func (c *OecConfig) SetMaxGasUsedPerBlock(value int64) {
	if value < -1 {
		return
	}
	c.maxGasUsedPerBlock = value
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

func (c *OecConfig) GetCsTimeoutPropose() time.Duration {
	return c.csTimeoutPropose
}
func (c *OecConfig) SetCsTimeoutPropose(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPropose = value
}

func (c *OecConfig) GetCsTimeoutProposeDelta() time.Duration {
	return c.csTimeoutProposeDelta
}
func (c *OecConfig) SetCsTimeoutProposeDelta(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutProposeDelta = value
}

func (c *OecConfig) GetCsTimeoutPrevote() time.Duration {
	return c.csTimeoutPrevote
}
func (c *OecConfig) SetCsTimeoutPrevote(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrevote = value
}

func (c *OecConfig) GetCsTimeoutPrevoteDelta() time.Duration {
	return c.csTimeoutPrevoteDelta
}
func (c *OecConfig) SetCsTimeoutPrevoteDelta(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrevoteDelta = value
}

func (c *OecConfig) GetCsTimeoutPrecommit() time.Duration {
	return c.csTimeoutPrecommit
}
func (c *OecConfig) SetCsTimeoutPrecommit(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrecommit = value
}

func (c *OecConfig) GetCsTimeoutPrecommitDelta() time.Duration {
	return c.csTimeoutPrecommitDelta
}
func (c *OecConfig) SetCsTimeoutPrecommitDelta(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrecommitDelta = value
}

func (c *OecConfig) GetIavlCacheSize() int {
	return c.iavlCacheSize
}
func (c *OecConfig) SetIavlCacheSize(value int) {
	c.iavlCacheSize = value
	iavl.IavlCacheSize = value
}
