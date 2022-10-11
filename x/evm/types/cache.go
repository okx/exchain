package types

import (
	"sync"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var evmParamsCache = NewCache()

type Cache struct {
	paramsCache      Params
	needParamsUpdate bool
	paramsMutex      sync.RWMutex

	blockedContractMethodsCache map[string]BlockedContract
	needBlockedUpdate           bool
	blockedMutex                sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		paramsCache:                 DefaultParams(),
		blockedContractMethodsCache: make(map[string]BlockedContract, 0),
		needParamsUpdate:            true,
		needBlockedUpdate:           true,
	}
}

func (c *Cache) SetNeedParamsUpdate() {
	c.paramsMutex.Lock()
	defer c.paramsMutex.Unlock()
	c.needParamsUpdate = true
}

func (c *Cache) GetParams(subspace Subspace, ctx sdk.Context) Params {
	c.paramsMutex.Lock()
	defer c.paramsMutex.Unlock()

	if c.needParamsUpdate {
		if !ctx.IsCheckTx() {
			var params Params
			subspace.GetParamSet(ctx, &params)
			c.paramsCache = params
		}

		c.needBlockedUpdate = false
	}

	return NewParams(c.paramsCache.EnableCreate,
		c.paramsCache.EnableCall,
		c.paramsCache.EnableContractDeploymentWhitelist,
		c.paramsCache.EnableContractBlockedList,
		c.paramsCache.MaxGasLimitPerTx,
		c.paramsCache.ExtraEIPs...)
}

func (c *Cache) SetNeedBlockedUpdate() {
	c.blockedMutex.Lock()
	defer c.blockedMutex.Unlock()
	c.needBlockedUpdate = true
}

func (c *Cache) GetBlockedContractMethod(addr string, csdb *CommitStateDB) (contract *BlockedContract) {
	c.blockedMutex.Lock()
	defer c.blockedMutex.Unlock()

	if c.needBlockedUpdate {
		if !csdb.ctx.IsCheckTx() {
			bcl := csdb.GetContractMethodBlockedList()
			c.blockedContractMethodsCache = make(map[string]BlockedContract, len(bcl))
			for i, _ := range bcl {
				c.blockedContractMethodsCache[string(bcl[i].Address)] = bcl[i]
			}
		}
		c.needBlockedUpdate = false
	}

	bc, ok := c.blockedContractMethodsCache[addr]
	if ok {
		return NewBlockContract(bc.Address, bc.BlockMethods)
	}
	return nil
}

func SetEvmParamsNeedUpdate() {
	GetEvmParamsCache().SetNeedParamsUpdate()
}

func GetEvmParamsCache() *Cache {
	return evmParamsCache
}
