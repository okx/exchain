package baseapp

import (
	"container/list"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
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

func (c *cacheMultiStoreList) PushStore(store types.CacheMultiStore) {
	c.mtx.Lock()
	c.stores.PushBack(store)
	c.mtx.Unlock()
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
