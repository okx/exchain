package types

import (
	"container/list"
	"sync"
)

var (
	verifySigCache *Cache
	cacheOnce      sync.Once
)

const cacheSize = 1000000

func init() {
	cacheOnce.Do(func() {
		verifySigCache = newCache()
	})
}

type Cache struct {
	mtx   sync.RWMutex
	data  map[string]*list.Element
	queue *list.List
}

type cacheNode struct {
	key   string
	value *ethSigCache
}

func newCache() *Cache {
	return &Cache{
		data:  make(map[string]*list.Element, cacheSize),
		queue: list.New(),
	}
}

func (c *Cache) Get(key string) (*ethSigCache, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	// validate key
	if !validateKey(key) {
		return nil, false
	}
	// get cache
	if elem, ok := c.data[key]; ok {
		return elem.Value.(*cacheNode).value, true
	}
	return nil, false
}

func (c *Cache) Add(key string, value *ethSigCache) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	// validate key
	if !validateKey(key) {
		return
	}
	// add cache
	node := &cacheNode{key, value}
	elem := c.queue.PushBack(node)
	c.data[key] = elem

	for c.queue.Len() > cacheSize {
		oldest := c.queue.Front()
		oldKey := c.queue.Remove(oldest).(*cacheNode).key
		delete(c.data, oldKey)
	}
}

func validateKey(key string) bool {
	if key == "" {
		return false
	}
	return true
}
