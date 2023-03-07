package types

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

const cacheSize = 1024

var paramsCache = NewCache()

type Cache struct {
	params           Params
	needParamsUpdate bool
	paramsMutex      sync.RWMutex

	feeSplits *lru.Cache
	shares    *lru.Cache
}

func NewCache() *Cache {
	c := &Cache{
		params:           DefaultParams(),
		needParamsUpdate: true,
	}

	c.feeSplits, _ = lru.New(cacheSize)
	c.shares, _ = lru.New(cacheSize)
	return c
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
	c.feeSplits.Add(contract, feeSplit)
}

// DeleteFeeSplit The change in feeSplit is only related to the user tx(register,update,cancel)
func (c *Cache) DeleteFeeSplit(contract common.Address, isCheckTx bool) {
	if isCheckTx {
		return
	}
	c.feeSplits.Remove(contract)
}

func (c *Cache) GetFeeSplit(contract common.Address) (FeeSplit, bool) {
	feeSplit, found := c.feeSplits.Get(contract)
	if found {
		return feeSplit.(FeeSplit), true
	}
	return FeeSplit{}, false
}

// UpdateShare The change in share is only related to the proposal
func (c *Cache) UpdateShare(contract common.Address, share sdk.Dec, isCheckTx bool) {
	if isCheckTx {
		return
	}
	c.shares.Add(contract, share)
}

func (c *Cache) GetShare(contract common.Address) (sdk.Dec, bool) {
	share, found := c.shares.Get(contract)
	if found {
		return share.(sdk.Dec), true
	}
	return sdk.Dec{}, false
}

func SetParamsNeedUpdate() {
	paramsCache.SetNeedParamsUpdate()
}

func GetParamsCache() *Cache {
	return paramsCache
}
