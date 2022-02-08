package config

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	iavlconfig "github.com/okex/exchain/libs/iavl/config"
	tmconfig "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/spf13/viper"
)

var _ tmconfig.IDynamicConfig = &OecConfig{}
var _ iavlconfig.IDynamicConfig = &OecConfig{}

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
	// nodeKeyWhitelist
	nodeKeyWhitelist []string

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
	iavlCacheSize int
}

const (
	FlagEnableDynamic = "config.enable-dynamic"

	FlagMempoolRecheck         = "mempool.recheck"
	FlagMempoolForceRecheckGap = "mempool.force_recheck_gap"
	FlagMempoolSize            = "mempool.size"
	FlagMempoolFlush           = "mempool.flush"
	FlagMaxTxNumPerBlock       = "mempool.max_tx_num_per_block"
	FlagMaxGasUsedPerBlock     = "mempool.max_gas_used_per_block"
	FlagNodeKeyWhitelist       = "mempool.node_key_whitelist"
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

var (
	testnetNodeIdWhitelist = []string{
		// RPC nodes for users
		"3a339568305c5aff58a1f134437b608490e2ec6d",
		"b9e7bf85886f1d11ee5079726a268401bf7b6254",
		"54c5ffc54e10a311660d16a96d54ddc59edb5555",
		"d77e385de87acdd042973c5d3029b02db8d767ff",
		"5cfdc51d1502fbe44d1b2a7f1f37e1016ad5ee97",
		"704be3bf19866f2aa5c77b09f003e2b69c552927",
		"1767342f12cb0e1e393a42c56d63d7486b2c54cd",
		"d33084a8c7bab8c9b6f286378b5e3ac197caa41a",
		"6a96b0a094ec9aaff2b7148b0c5811618b41c101",
		// RPC nodes for developers
		"3a35faa50649164d59f07f31d78946ca07464e9c",
		"cee36e7fbc99eaa02bd9af692dae367a867c43f4",
		"fbcae686695cd17ee8319bbd6b9b0aaf0f10d8c4",
		"2c34f93a8665d694e56319ccdc6738b203c33848",
		"f689ab031c0758367af229aa8df65ac69762327d",
		"58c495e040a1576ebc1f386a7dc04c4e60ee63d7",
		"a2f685db92a88c18780d8d9cb1162ab61517ae64",
		"8d4b0539b95b60e1691eac77be4aa7295645d9d9",
		"6f328902a0bf5e7b922d6a5980dd6888097db984",
		"12503ae035dd7ff04e19b0ca2c9c8b54a0a56b22",
		// validator nodes
		"c39ca38c650b920f9b6c5a9aed7ff904124ec3ad",
		"d937e21fd489809add23dc3e55ed78d947217aa8",
		"a3eb3c129e49137d5e1665bbf87b6f2be70a0b85",
		"b171a9ef83b95c28182bc7aa7ea8639d04e572e7",
		"3a700a3849c401396b1c51eb65b1cfc1a8c4394b",
		"0208e66d4ca746ec535a0bf05409dc87df408b15",
		"ed1819fa1eae52ddec4c0f8cddd80b9cb7c68a22",
		"0b3ab9597a66f2f94c8efa4ccb6ed2a1f44d4184",
		"67b29551c7c3839ad6c93379991344266aec3829",
		"cd07b20b596aac923a1d5bb022581e279755aff1",
		"6ce06a89a968a4204d9dcb470f2275767c8dfa68",
		"6dd38d96df3ccbca95769ee15bdfdd952ad007c5",
		"fcc95bfee6ea74bdf385be3a29072329603676e5",
		"7b5b3041d2b3546a236b6df7ff7e06a19a5cae46",
		"c098585e299ff7afe6f354c4431550d6919bdd0d",
		"5b44fb4af4cfb72286162cb49a3bc04cb8187775",
		"358e3399b68fb67787f1386c685db2e75352d9eb",
		"96d9cb96041c053e63ff7d0c7d81dfab706136e4",
		"0de948586fb30293d1dd14a99ebc3f719deb7c6f",
		"284e87518752c8f655fe217113fa86ba7d6ca72f",
		"7f2b8a6b9b8b12247e6992aeb32d69e169c2f5ac",
	}

	mainnetNodeIdWhitelist = []string{
	}

	oecConfig  *OecConfig
	once       sync.Once
	confLogger log.Logger
)

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
	c.SetNodeKeyWhitelist(viper.GetString(FlagNodeKeyWhitelist))
}

func resolveNodeKeyWhitelist(plain string) []string {
	if len(plain) == 0 {
		return []string{}
	}
	return strings.Split(plain, ",")
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
	case FlagNodeKeyWhitelist:
		r, ok := value.(string)
		if !ok {
			return
		}
		c.SetNodeKeyWhitelist(r)
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

func (c *OecConfig) GetNodeKeyWhitelist() []string {
	return c.nodeKeyWhitelist
}
func (c *OecConfig) SetNodeKeyWhitelist(value string) {
	idList := resolveNodeKeyWhitelist(value)

	for _, id := range idList {
		if id == "testnet-node-ids" {
			c.nodeKeyWhitelist = append(c.nodeKeyWhitelist, testnetNodeIdWhitelist...)
		} else if id == "mainnet-node-ids" {
			c.nodeKeyWhitelist = append(c.nodeKeyWhitelist, mainnetNodeIdWhitelist...)
		} else {
			c.nodeKeyWhitelist = append(c.nodeKeyWhitelist, id)
		}
	}
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
