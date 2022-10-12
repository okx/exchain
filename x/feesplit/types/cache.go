package types

import (
	"sync"
)

var paramsCache = NewCache()

type Cache struct {
	params           Params
	needParamsUpdate bool
	paramsMutex      sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		params:           DefaultParams(),
		needParamsUpdate: true,
	}
}

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

func SetParamsNeedUpdate() {
	paramsCache.SetNeedParamsUpdate()
}

func GetParamsCache() *Cache {
	return paramsCache
}
