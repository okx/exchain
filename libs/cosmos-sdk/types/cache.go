package types

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/spf13/viper"
	"time"
)

var (
	maxAccInMap         = 100000
	deleteAccCount      = 10000
	maxContractInMap    = 110000
	deleteContractCount = 10000

	FlagMultiCache          = "multi-cache"
	MaxAccInMultiCache      = "multi-cache-acc"
	MaxContractInMultiCache = "multi-cache-contract"
	UseCache                bool
)

type account interface {
	Copy() interface{}
	GetAddress() AccAddress
	SetAddress(AccAddress) error
	GetPubKey() crypto.PubKey
	SetPubKey(crypto.PubKey) error
	GetAccountNumber() uint64
	SetAccountNumber(uint64) error
	GetSequence() uint64
	SetSequence(uint64) error
	GetCoins() Coins
	SetCoins(Coins) error
	SpendableCoins(blockTime time.Time) Coins
	String() string
}

type storageWithCache struct {
	value []byte
	dirty bool
}

type accountWithCache struct {
	acc     account
	gas     uint64
	isDirty bool
}

type codeWithCache struct {
	code    []byte
	isDirty bool
}

type Cache struct {
	useCache  bool
	parent    *Cache
	gasConfig types.GasConfig

	storageMap map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache
	accMap     map[ethcmn.Address]*accountWithCache
	codeMap    map[ethcmn.Hash]*codeWithCache
}

func initCacheParam() {
	UseCache = viper.GetBool(FlagMultiCache)

	maxAccInMap = viper.GetInt(MaxAccInMultiCache)
	deleteAccCount = maxAccInMap / 10

	maxContractInMap = viper.GetInt(MaxContractInMultiCache)
	deleteContractCount = maxContractInMap / 10
}

func NewChainCache() *Cache {
	initCacheParam()
	return NewCache(nil, UseCache)
}

func NewCache(parent *Cache, useCache bool) *Cache {
	return &Cache{
		useCache: useCache,
		parent:   parent,

		storageMap: make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache, 0),
		accMap:     make(map[ethcmn.Address]*accountWithCache, 0),
		codeMap:    make(map[ethcmn.Hash]*codeWithCache),
		gasConfig:  types.KVGasConfig(),
	}

}

func (c *Cache) Skip() bool {
	if c == nil || !c.useCache {
		return true
	}
	return false
}

func (c *Cache) UpdateAccount(addr AccAddress, acc account, lenBytes int, isDirty bool) {
	if c.Skip() {
		return
	}
	ethAddr := ethcmn.BytesToAddress(addr.Bytes())
	c.accMap[ethAddr] = &accountWithCache{
		acc:     acc,
		isDirty: isDirty,
		gas:     types.Gas(lenBytes)*c.gasConfig.ReadCostPerByte + c.gasConfig.ReadCostFlat,
	}
}

func (c *Cache) UpdateStorage(addr ethcmn.Address, key ethcmn.Hash, value []byte, isDirty bool) {
	if c.Skip() {
		return
	}

	if _, ok := c.storageMap[addr]; !ok {
		c.storageMap[addr] = make(map[ethcmn.Hash]*storageWithCache, 0)
	}
	c.storageMap[addr][key] = &storageWithCache{
		value: value,
		dirty: isDirty,
	}
}

func (c *Cache) UpdateCode(key []byte, value []byte, isdirty bool) {
	if c.Skip() {
		return
	}
	hash := ethcmn.BytesToHash(key)
	c.codeMap[hash] = &codeWithCache{
		code:    value,
		isDirty: isdirty,
	}
}

func (c *Cache) GetAccount(addr ethcmn.Address) (account, uint64, bool, bool) {
	if c.Skip() {
		return nil, 0, false, false
	}

	if data, ok := c.accMap[addr]; ok {
		return data.acc, data.gas, ok, false
	}

	if c.parent != nil {
		acc, gas, ok, _ := c.parent.GetAccount(addr)
		return acc, gas, ok, true
	}
	return nil, 0, false, false
}

func (c *Cache) GetStorage(addr ethcmn.Address, key ethcmn.Hash) ([]byte, bool) {
	if c.Skip() {
		return nil, false
	}
	if _, hasAddr := c.storageMap[addr]; hasAddr {
		data, hasKey := c.storageMap[addr][key]
		if hasKey {
			return data.value, hasKey
		}
	}

	if c.parent != nil {
		return c.parent.GetStorage(addr, key)
	}
	return nil, false
}

func (c *Cache) GetCode(key []byte) ([]byte, bool) {
	if c.Skip() {
		return nil, false
	}

	hash := ethcmn.BytesToHash(key)
	if data, ok := c.codeMap[hash]; ok {
		return data.code, ok
	}

	if c.parent != nil {
		return c.parent.GetCode(hash.Bytes())
	}
	return nil, false
}

func (c *Cache) Write(updateDirty bool) {
	if c.Skip() {
		return
	}

	if !updateDirty {
		c.storageMap = make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache)
		c.accMap = make(map[ethcmn.Address]*accountWithCache)
		c.codeMap = make(map[ethcmn.Hash]*codeWithCache)
		return
	}

	if c.parent == nil {
		return
	}

	c.writeStorage()
	c.writeAcc()
	c.writeCode()
}

func (c *Cache) writeStorage() {
	for addr, storages := range c.storageMap {
		if _, ok := c.parent.storageMap[addr]; !ok {
			c.parent.storageMap[addr] = make(map[ethcmn.Hash]*storageWithCache, 0)
		}

		for key, v := range storages {
			if v.dirty {
				c.parent.storageMap[addr][key] = v
			}
		}
	}
	c.storageMap = make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache)
}

func (c *Cache) writeAcc() {
	for addr, v := range c.accMap {
		if v.isDirty {
			c.parent.accMap[addr] = v
		}
	}
	c.accMap = make(map[ethcmn.Address]*accountWithCache)
}

func (c *Cache) writeCode() {
	for hash, v := range c.codeMap {
		if v.isDirty {
			c.parent.codeMap[hash] = v
		}
	}
	c.codeMap = make(map[ethcmn.Hash]*codeWithCache)
}

func (c *Cache) TryDelete(logger log.Logger, height int64) {
	if c.Skip() {
		return
	}
	if height%1000 == 0 {
		c.logInfo(logger, "null")
	}

	if len(c.accMap) < maxAccInMap && len(c.storageMap) < maxContractInMap {
		return
	}

	deleteMsg := ""
	if len(c.accMap) >= maxAccInMap {
		deleteMsg += fmt.Sprintf("Acc:Deleted Before:%d", len(c.accMap))
		cnt := 0
		for key := range c.accMap {
			delete(c.accMap, key)
			cnt++
			if cnt > deleteAccCount {
				break
			}
		}
	}

	if len(c.storageMap) >= maxContractInMap {
		lenStorage := 0
		for _, v := range c.storageMap {
			lenStorage += len(v)
		}
		deleteMsg += fmt.Sprintf("Storage:Deleted Before:len(contract):%d, len(storage):%d", len(c.storageMap), lenStorage)
		cnt := 0
		for key := range c.storageMap {
			delete(c.storageMap, key)
			cnt++
			if cnt > deleteContractCount {
				break
			}
		}
	}
	if deleteMsg != "" {
		c.logInfo(logger, deleteMsg)
	}
}

func (c *Cache) logInfo(logger log.Logger, deleteMsg string) {
	lenStorage := 0
	for _, v := range c.storageMap {
		lenStorage += len(v)
	}
	nowStats := fmt.Sprintf("len(acc):%d len(contracts):%d len(storage):%d", len(c.accMap), len(c.storageMap), lenStorage)
	logger.Info("MultiCache", "deleteMsg", deleteMsg, "nowStats", nowStats)
}
