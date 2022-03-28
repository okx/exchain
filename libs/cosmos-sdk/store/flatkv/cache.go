package flatkv

import "sync"

// cache value
type cValue struct {
	value   []byte
	deleted bool
	dirty   bool
}

// Cache defines flat kv cache
type Cache struct {
	mtx  sync.RWMutex
	data map[string]cValue
}

func newCache() *Cache {
	return &Cache{
		data: make(map[string]cValue),
	}
}

func (c *Cache) get(key []byte) ([]byte, bool) {
	strKey := string(key)
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	cacheValue, ok := c.data[strKey]
	if !ok {
		return nil, false
	}
	return cacheValue.value, true
}

func (c *Cache) add(key, value []byte, deleted bool, dirty bool) {
	strKey := string(key)
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[strKey] = cValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}

// return cache and clear cache
func (c *Cache) reset() map[string]cValue {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	copyMap := make(map[string]cValue, len(c.data))
	for k, v := range c.data {
		copyMap[k] = v
		delete(c.data, k)
	}
	return copyMap
}
