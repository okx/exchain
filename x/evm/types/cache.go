package types

import (
	"sync"
)

var EvmParamsCache = NewCache()

type Cache struct {
	paramsCache                 Params
	blockedContractMethodsCache map[string]BlockedContract
	needUpdate                  bool
	mutex                       sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		paramsCache:                 DefaultParams(),
		blockedContractMethodsCache: make(map[string]BlockedContract, 0),
		needUpdate:                  true,
	}
}

func (c *Cache) UpdateParams(params Params) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.paramsCache = params
	c.needUpdate = false
}

func (c *Cache) SetNeedUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.needUpdate = true
}

func (c *Cache) IsNeedUpdate() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.needUpdate
}

func (c Cache) GetParams() Params {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return NewParams(c.paramsCache.EnableCreate, c.paramsCache.EnableCall, c.paramsCache.EnableContractDeploymentWhitelist, c.paramsCache.EnableContractBlockedList, c.paramsCache.MaxGasLimitPerTx, c.paramsCache.ExtraEIPs...)
}

func (c *Cache) GetBlockedContractMethod(addr string) (contract *BlockedContract) {
	c.mutex.RLock()
	bc, ok := c.blockedContractMethodsCache[addr]
	c.mutex.RUnlock()
	if ok {
		NewBlockContract(bc.Address, bc.BlockMethods)
	}
	return nil
}

func (c *Cache) UpdateBlockedContractMethod(bcl BlockedContractList) {
	c.mutex.Lock()
	c.blockedContractMethodsCache = make(map[string]BlockedContract, 0)
	for i, _ := range bcl {
		c.blockedContractMethodsCache[bcl[i].String()] = bcl[i]
	}
	c.mutex.Unlock()
	c.needUpdate = false
}
