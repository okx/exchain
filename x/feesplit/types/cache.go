package types

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var paramsCache = NewCache()

type Cache struct {
	params           Params
	needParamsUpdate bool
	paramsMutex      sync.RWMutex

	feeSplits     map[string]FeeSplit
	feeSplitMutex sync.RWMutex

	shares      map[string]sdk.Dec
	sharesMutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		params:           DefaultParams(),
		needParamsUpdate: true,

		feeSplits: make(map[string]FeeSplit, 0),
		shares:    make(map[string]sdk.Dec, 0),
	}
}

// UpdateParams  the update in params is relates to the proposal and initGenesis
func (c *Cache) UpdateParams(params Params, isCheckTx bool) {
	if isCheckTx {
		return
	}
	c.paramsMutex.Lock()
	defer c.paramsMutex.Unlock()
	c.params = params
	c.needParamsUpdate = false
}

func (c *Cache) SetNeedParamsUpdate() {
	c.paramsMutex.Lock()
	defer c.paramsMutex.Unlock()
	c.needParamsUpdate = true
}

func (c *Cache) IsNeedParamsUpdate() bool {
	c.paramsMutex.RLock()
	defer c.paramsMutex.RUnlock()
	return c.needParamsUpdate
}

func (c *Cache) GetParams() Params {
	c.paramsMutex.RLock()
	defer c.paramsMutex.RUnlock()
	return NewParams(c.params.EnableFeeSplit,
		c.params.DeveloperShares,
		c.params.AddrDerivationCostCreate,
	)
}

// UpdateFeeSplit The change in feeSplit is only related to the user tx(register,update,cancel)
func (c *Cache) UpdateFeeSplit(contract common.Address, feeSplit FeeSplit, isCheckTx bool) {
	if isCheckTx {
		return
	}
	c.feeSplitMutex.Lock()
	defer c.feeSplitMutex.Unlock()
	c.feeSplits[contract.String()] = feeSplit
}

// DeleteFeeSplit The change in feeSplit is only related to the user tx(register,update,cancel)
func (c *Cache) DeleteFeeSplit(contract common.Address, isCheckTx bool) {
	if isCheckTx {
		return
	}
	c.feeSplitMutex.Lock()
	defer c.feeSplitMutex.Unlock()
	delete(c.feeSplits, contract.String())
}

func (c *Cache) GetFeeSplit(contract common.Address) (FeeSplit, bool) {
	c.feeSplitMutex.RLock()
	defer c.feeSplitMutex.RUnlock()
	feeSplit, found := c.feeSplits[contract.String()]
	return feeSplit, found
}

// UpdateShare The change in share is only related to the proposal
func (c *Cache) UpdateShare(contract common.Address, share sdk.Dec, isCheckTx bool) {
	if isCheckTx {
		return
	}
	c.sharesMutex.Lock()
	defer c.sharesMutex.Unlock()
	c.shares[contract.String()] = share
}

func (c *Cache) GetShare(contract common.Address) (sdk.Dec, bool) {
	c.sharesMutex.RLock()
	defer c.sharesMutex.RUnlock()
	share, found := c.shares[contract.String()]
	return share, found
}

func SetParamsNeedUpdate() {
	paramsCache.SetNeedParamsUpdate()
}

func GetParamsCache() *Cache {
	return paramsCache
}
