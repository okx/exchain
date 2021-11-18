package types

import (
	lru "github.com/hashicorp/golang-lru"
)

type Cache struct {
	cache *lru.Cache
}

func NewCache(size int) *Cache {
	c, err := lru.New(size)
	if err != nil {
		return nil
	}

	return &Cache{
		cache: c,
	}
}

func (c *Cache) Set(key, value interface{}) {
	c.cache.Add(key, value)
}

func (c *Cache) Get(key interface{}) interface{} {
	r, ok := c.cache.Get(key)
	if !ok {
		return nil
	}
	return r
}
