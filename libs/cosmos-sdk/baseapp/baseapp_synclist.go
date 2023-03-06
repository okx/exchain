package baseapp

import (
	"container/list"
	"sync"

	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
)

type cacheRWSetList struct {
	mtx sync.Mutex
	mps *list.List
}

func newCacheRWSetList() *cacheRWSetList {
	return &cacheRWSetList{
		mps: list.New(),
	}
}

func (c *cacheRWSetList) Len() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.mps.Len()
}

func (c *cacheRWSetList) PutRwSet(rw types.MsRWSet) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.mps.PushBack(rw)
}

func (c *cacheRWSetList) GetRWSet() types.MsRWSet {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.mps.Len() > 0 {
		front := c.mps.Remove(c.mps.Front())
		return front.(types.MsRWSet)
	}

	return make(types.MsRWSet)
}

func (c *cacheRWSetList) Range(cb func(c types.MsRWSet)) {
	c.mtx.Lock()
	for i := c.mps.Front(); i != nil; i = i.Next() {
		cb(i.Value.(types.MsRWSet))
	}
	c.mtx.Unlock()
}

func (c *cacheRWSetList) Clear() {
	c.mtx.Lock()
	c.mps.Init()
	c.mtx.Unlock()
}

type cacheMultiStoreList struct {
	mtx    sync.Mutex
	stores *list.List
}

func newCacheMultiStoreList() *cacheMultiStoreList {
	return &cacheMultiStoreList{
		stores: list.New(),
	}
}

func (c *cacheMultiStoreList) Len() int {
	c.mtx.Lock()

	defer c.mtx.Unlock()
	return c.stores.Len()

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
	c.mtx.Lock()
	for i := c.stores.Front(); i != nil; i = i.Next() {
		cb(i.Value.(types.CacheMultiStore))
	}
	c.mtx.Unlock()
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
