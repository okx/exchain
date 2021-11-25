package types

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"time"
)

type storageWithCache struct {
	value []byte
	dirty bool
}

type accountWithCache struct {
	acc     account
	gas     uint64
	isDirty bool
}

type Cache struct {
	useCache bool
	parent   *Cache

	storageMap map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache

	accMap map[ethcmn.Address]*accountWithCache
}

type account interface {
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

func NewCache(parent *Cache, useCache bool) *Cache {
	return &Cache{
		useCache: useCache,
		parent:   parent,

		storageMap: make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache, 0),
		accMap:     make(map[ethcmn.Address]*accountWithCache, 0),
	}

}

func (c *Cache) UpdateStorage(addr ethcmn.Address, key ethcmn.Hash, value []byte, isDirty bool) {
	ts := time.Now()
	defer func() {
		UpdateStorage += time.Now().Sub(ts)
	}()

	if !c.useCache {
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

func (c *Cache) UpdateAcc(addr AccAddress, acc account, lenBytes int, isDirty bool) {
	ts := time.Now()
	defer func() {
		UpdataAcc += time.Now().Sub(ts)
	}()

	if !c.useCache {
		return
	}
	ethAddr := ethcmn.BytesToAddress(addr.Bytes())
	//fmt.Println("update-----", ethAddr.String(), lenBytes*3+1000)
	c.accMap[ethAddr] = &accountWithCache{
		acc:     acc,
		isDirty: isDirty,
		gas:     uint64(lenBytes*3 + 1000),
	}
}

func (c *Cache) GetAcc(addr ethcmn.Address) (account, uint64, bool) {
	ts := time.Now()
	defer func() {
		AccountTs += time.Now().Sub(ts)
	}()

	if !c.useCache {
		return nil, 0, false
	}

	if data, ok := c.accMap[addr]; ok {
		return data.acc, data.gas, ok
	}

	if c.parent != nil {
		if data, ok := c.parent.accMap[addr]; ok {
			return data.acc, data.gas, ok
		}
	}
	return nil, 0, false

}

func (c *Cache) GetStorage(addr ethcmn.Address, key ethcmn.Hash) ([]byte, bool) {
	ts := time.Now()
	defer func() {
		StorageTs += time.Now().Sub(ts)
	}()
	if !c.useCache {
		return nil, false
	}
	if _, hasAddr := c.storageMap[addr]; hasAddr {
		data, hasKey := c.storageMap[addr][key]
		if hasKey {
			return data.value, hasKey
		}
	}

	if c.parent != nil {
		if _, hasAddr := c.parent.storageMap[addr]; hasAddr {
			data, hasKey := c.parent.storageMap[addr][key]
			if hasKey {
				return data.value, hasKey
			}
		}
	}
	return nil, false
}

//TODO delete
var (
	UpdateStorage = time.Duration(0)
	UpdataAcc     = time.Duration(0)
	WriteTs       = time.Duration(0)
	StorageTs     = time.Duration(0)
	AccountTs     = time.Duration(0)
)

func DisplayTs() {
	fmt.Println("Write", WriteTs.Seconds(), "UpdateStorage", UpdateStorage.Seconds(), "UpdateAcc", UpdataAcc.Seconds(), "GetStorage", StorageTs.Seconds(), "GetAcc", AccountTs.Seconds())
}

func (c *Cache) Write(updateDirty bool) {
	ts := time.Now()
	defer func() {
		WriteTs += time.Now().Sub(ts)
	}()
	if !c.useCache {
		return
	}
	c.writeStorage(updateDirty)
	c.writeAcc(updateDirty)
}

func (c *Cache) writeStorage(updateDirty bool) {
	for addr, storages := range c.storageMap {
		if _, ok := c.parent.storageMap[addr]; !ok {
			c.parent.storageMap[addr] = make(map[ethcmn.Hash]*storageWithCache, 0)
		}

		for key, v := range storages {
			if !v.dirty || (updateDirty && v.dirty) {
				c.parent.storageMap[addr][key] = v
			}
		}
	}
	c.storageMap = make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache)
}

func (c *Cache) writeAcc(updateDirty bool) {
	for addr, v := range c.accMap {
		if !v.isDirty || (updateDirty && v.isDirty) {
			c.parent.accMap[addr] = v
		}
	}
	c.accMap = make(map[ethcmn.Address]*accountWithCache)
}
