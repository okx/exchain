package types

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
)

var (
	verifySigCache *Cache
)

const FlagSigCacheSize = "signature-cache-size"

func InitSignatureCache() {
	lruCache, err := lru.New(viper.GetInt(FlagSigCacheSize))
	if err != nil {
		panic(err)
	}
	verifySigCache = &Cache{
		data: lruCache,
	}
}

type Cache struct {
	data *lru.Cache
}

func (c *Cache) Get(key string) (*ethSigCache, bool) {
	// validate key
	if !validateKey(key) {
		return nil, false
	}
	// get cache
	value, ok := c.data.Get(key)
	if ok {
		sigCache, ok := value.(*ethSigCache)
		if ok {
			return sigCache, true
		}
	}
	return nil, false
}

func (c *Cache) Add(key string, value *ethSigCache) {
	// validate key
	if !validateKey(key) {
		return
	}
	// add cache
	c.data.Add(key, value)
}

func validateKey(key string) bool {
	if key == "" {
		return false
	}
	return true
}
