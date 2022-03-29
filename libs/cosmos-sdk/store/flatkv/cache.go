package flatkv

import (
	"sync"

	dbm "github.com/okex/exchain/libs/tm-db"
)

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

// write cache to db and clear cache
func (c *Cache) write(db dbm.DB, version int64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	batch := db.NewBatch()
	defer batch.Close()
	for key, cValue := range c.data {
		if cValue.deleted {
			batch.Delete([]byte(key))
		} else if cValue.dirty {
			batch.Set([]byte(key), cValue.value)
		}
		delete(c.data, key)
	}
	setLatestVersion(batch, version)
	batch.Write()
}
