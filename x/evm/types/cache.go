package types

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
)

const (
	FlagEvmCodeCache = "evm-code-cache"
	CodeCacheSize    = 500000
)

var isEvmCacheCode = viper.GetBool(FlagEvmCodeCache)

type Cache struct {
	enable bool
	cache  *lru.Cache
}

func NewCache() *Cache {
	c, err := lru.New(CodeCacheSize)
	if err != nil {
		return nil
	}

	return &Cache{
		cache:  c,
		enable: isEvmCacheCode,
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
