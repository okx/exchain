package keeper

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

type evmCache struct {
	use      bool
	param    types.Params
	paramGas uint64

	chainConfig    types.ChainConfig
	chainConfigGas uint64

	blackList map[ethcmn.Address]bool
}

func newConfigCache() *evmCache {
	return &evmCache{
		blackList: make(map[ethcmn.Address]bool),
	}
}

func (c *evmCache) SetSkipFlag(skip bool) *evmCache {
	c.use = !skip
	return c
}

func (c *evmCache) useCache() bool {
	return c.use
}

func (c *evmCache) SetBlackList(blackList []sdk.AccAddress) {
	if !c.useCache() {
		return
	}
	for _, v := range blackList {
		c.blackList[ethcmn.BytesToAddress(v)] = true
	}
}

func (c *evmCache) IsBlackList(addr sdk.AccAddress) (bool, bool) {
	if !c.useCache() {
		return false, false
	}
	if len(c.blackList) == 0 {
		return false, false
	}
	return c.blackList[ethcmn.BytesToAddress(addr)], true
}

func (c *evmCache) BlackListLen() (int, bool) {
	if !c.useCache() {
		return 0, false
	}
	return len(c.blackList), true
}

func (c *evmCache) CleanBlackList() {
	if !c.useCache() {
		return
	}
	c.blackList = make(map[ethcmn.Address]bool)
}

func (c *evmCache) GetParams() (types.Params, uint64) {
	if !c.useCache() {
		return types.Params{}, 0
	}
	return c.param, c.paramGas
}

func (c *evmCache) setParams(data types.Params, gasConsumed uint64) {
	if !c.useCache() {
		return
	}
	if c.paramGas != 0 {
		return
	}
	c.param = data
	c.paramGas = gasConsumed
}

func (c *evmCache) GetChainConfig() (types.ChainConfig, uint64) {
	if !c.useCache() {
		return types.ChainConfig{}, 0
	}
	return c.chainConfig, c.chainConfigGas
}

func (c *evmCache) setChainConfig(data types.ChainConfig, gasConsumed uint64) {
	if !c.useCache() {
		return
	}
	if c.chainConfigGas != 0 {
		return
	}
	c.chainConfig = data
	c.chainConfigGas = gasConsumed
}

func (c *evmCache) Clean() {
	c.use = false
	c.param = types.Params{}
	c.paramGas = 0
}
