package types

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/spf13/viper"
)

const (
	FlagEvmCodeCache = "evm-code-cache"
	CodeCacheSize    = 67108864 // 64 MB
)

var (
	CODE_PREFIX      = []byte{'c'}
	CODE_HASH_PREFIX = []byte{'h'}

	isEvmCacheCode = viper.GetBool(FlagEvmCodeCache)
)

type Cache struct {
	enable bool
	cache  *fastcache.Cache
}

func NewCache() *Cache {
	// c, err := fastcache.New(CodeCacheSize)
	c := fastcache.New(CodeCacheSize)

	return &Cache{
		cache:  c,
		enable: isEvmCacheCode,
	}
}

func (c *Cache) Set(prefix, key, value []byte) {
	if !c.enable {
		return
	}
	c.cache.SetBig(append(prefix, key...), value)
}

func (c *Cache) Get(prefix, key []byte) []byte {
	if !c.enable {
		return nil
	}
	return c.cache.GetBig(nil, append(prefix, key...))
}
