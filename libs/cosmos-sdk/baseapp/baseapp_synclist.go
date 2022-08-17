package baseapp

import (
	"container/list"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"sync"
)

type cacheMultiStoreList struct {
	mtx    sync.Mutex
	stores *list.List
}

func newCacheMultiStoreList() *cacheMultiStoreList {
	return &cacheMultiStoreList{
		stores: list.New(),
	}
}

func (c *cacheMultiStoreList) PushStores(stores map[int]types.CacheMultiStore) {
	c.mtx.Lock()
	for _, v := range stores {
		c.stores.PushBack(v)
	}
	c.mtx.Unlock()
}

func (c *cacheMultiStoreList) PushStore(store types.CacheMultiStore) {
	c.mtx.Lock()
	c.stores.PushBack(store)
	c.mtx.Unlock()
}

func (c *cacheMultiStoreList) Range(cb func(c types.CacheMultiStore)) {
	for i := c.stores.Front(); i != nil; i = i.Next() {
		cb(i.Value.(types.CacheMultiStore))
	}
}

func (c *cacheMultiStoreList) GetStoreWithParent(parent types.CacheMultiStore) types.CacheMultiStore {
	c.mtx.Lock()
	if c.stores.Len() > 0 {
		front := c.stores.Remove(c.stores.Front()).(types.CacheMultiStore)
		c.mtx.Unlock()
		front.(types.CacheMultiStoreResetter).Reset(parent)
		return front

	}
	c.mtx.Unlock()
	return parent.CacheMultiStore()
}

func (c *cacheMultiStoreList) GetStore() types.CacheMultiStore {
	c.mtx.Lock()
	if c.stores.Len() > 0 {
		front := c.stores.Remove(c.stores.Front())
		c.mtx.Unlock()
		return front.(types.CacheMultiStore)
	}
	c.mtx.Unlock()
	return nil
}

func (c *cacheMultiStoreList) Clear() {
	c.mtx.Lock()
	c.stores.Init()
	c.mtx.Unlock()
}
