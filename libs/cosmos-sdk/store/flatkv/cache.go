package flatkv

import "sync"

type Cache struct {
	mtx  sync.RWMutex
	data map[string][]byte
}

func newCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) get(key []byte) (value []byte, ok bool) {
	strKey := string(key)
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	value, ok = c.data[strKey]
	return
}

func (c *Cache) add(key, value []byte) {
	strKey := string(key)
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[strKey] = value
}

func (c *Cache) delete(key []byte) {
	strKey := string(key)
	c.mtx.Lock()
	defer c.mtx.Unlock()
	delete(c.data, strKey)
}

func (c *Cache) copy() map[string][]byte {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	copyMap := make(map[string][]byte, len(c.data))
	for k, v := range c.data {
		copyMap[k] = v
	}
	return copyMap
}

func (c *Cache) reset() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data = make(map[string][]byte)
}
