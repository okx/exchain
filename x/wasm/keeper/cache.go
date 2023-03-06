package keeper

import (
	"sync"

	"github.com/okx/okbchain/x/wasm/types"
)

var wasmParamsCache = NewCache()

type Cache struct {
	paramsCache      types.Params
	needParamsUpdate bool
	paramsMutex      sync.RWMutex

	blockedContractMethodsCache map[string]*types.ContractMethods
	needBlockedUpdate           bool
	blockedMutex                sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		paramsCache:                 types.DefaultParams(),
		blockedContractMethodsCache: make(map[string]*types.ContractMethods, 0),
		needParamsUpdate:            true,
		needBlockedUpdate:           true,
	}
}

func (c *Cache) UpdateParams(params types.Params) {
	c.paramsMutex.Lock()
	defer c.paramsMutex.Unlock()
	c.paramsCache = params
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

func (c *Cache) GetParams() types.Params {
	c.paramsMutex.RLock()
	defer c.paramsMutex.RUnlock()
	return types.Params{
		CodeUploadAccess:             c.paramsCache.CodeUploadAccess,
		InstantiateDefaultPermission: c.paramsCache.InstantiateDefaultPermission,
		UseContractBlockedList:       c.paramsCache.UseContractBlockedList,
		VmbridgeEnable:               c.paramsCache.VmbridgeEnable,
	}
}

func (c *Cache) SetNeedBlockedUpdate() {
	c.blockedMutex.Lock()
	defer c.blockedMutex.Unlock()
	c.needBlockedUpdate = true
}

func (c *Cache) IsNeedBlockedUpdate() bool {
	c.blockedMutex.RLock()
	defer c.blockedMutex.RUnlock()
	return c.needBlockedUpdate
}

func (c *Cache) GetBlockedContractMethod(addr string) (contract *types.ContractMethods) {
	c.blockedMutex.RLock()
	bc := c.blockedContractMethodsCache[addr]
	c.blockedMutex.RUnlock()
	return bc
}

func (c *Cache) UpdateBlockedContractMethod(cms []*types.ContractMethods) {
	c.blockedMutex.Lock()
	c.blockedContractMethodsCache = make(map[string]*types.ContractMethods, len(cms))
	for i, _ := range cms {
		c.blockedContractMethodsCache[cms[i].ContractAddr] = cms[i]
	}
	c.needBlockedUpdate = false
	c.blockedMutex.Unlock()
}

func GetWasmParamsCache() *Cache {
	return wasmParamsCache
}

func SetNeedParamsUpdate() {
	GetWasmParamsCache().SetNeedParamsUpdate()
}
