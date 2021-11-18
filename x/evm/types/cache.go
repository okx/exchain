package types

import (
	lru "github.com/hashicorp/golang-lru"
)

type Cache struct {
	enable bool
	cache  *lru.Cache
}

func NewCache(size int, enable bool) *Cache {
	c, err := lru.New(size)
	if err != nil {
		return nil
	}

	return &Cache{
		cache:  c,
		enable: enable,
	}
}

func (c *Cache) Set(key, value interface{}) {
	if !c.enable {
		return
	}
	c.cache.Add(key, value)
}

func (c *Cache) Get(key interface{}) interface{} {
	if !c.enable {
		return nil
	}
	r, ok := c.cache.Get(key)
	if !ok {
		return nil
	}
	return r
}
