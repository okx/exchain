package types

import (
	"sync/atomic"

	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"

	lru "github.com/hashicorp/golang-lru"
)

var (
	signatureCache *Cache
)

const FlagSigCacheSize = "signature-cache-size"

func init() {
	// used for ut
	defaultCache := &Cache{
		data:      nil,
		readCount: 0,
		hitCount:  0,
	}
	signatureCache = defaultCache
}

func InitSignatureCache() {
	lruCache, err := lru.New(viper.GetInt(FlagSigCacheSize))
	if err != nil {
		panic(err)
	}
	signatureCache = &Cache{
		data: lruCache,
	}
}

func SignatureCache() *Cache {
	return signatureCache
}

type Cache struct {
	data      *lru.Cache
	readCount int64
	hitCount  int64
}

func (c *Cache) Get(key []byte) (string, bool) {
	// validate
	if !c.validate(key) {
		return "", false
	}
	atomic.AddInt64(&c.readCount, 1)
	// get cache
	value, ok := c.data.Get(amino.BytesToStr(key))
	if ok {
		sigCache, ok := value.(string)
		if ok {
			atomic.AddInt64(&c.hitCount, 1)
			return sigCache, true
		}
	}
	return "", false
}

func (c *Cache) Add(key []byte, value string) {
	// validate
	if !c.validate(key) {
		return
	}
	// add cache
	c.data.Add(string(key), value)
}

func (c *Cache) Remove(key []byte) {
	// validate
	if !c.validate(key) {
		return
	}
	c.data.Remove(amino.BytesToStr(key))
}

func (c *Cache) ReadCount() int64 {
	return atomic.LoadInt64(&c.readCount)
}

func (c *Cache) HitCount() int64 {
	return atomic.LoadInt64(&c.hitCount)
}

func (c *Cache) validate(key []byte) bool {
	// validate key
	if len(key) == 0 {
		return false
	}
	// validate lru cache
	if c.data == nil {
		return false
	}
	return true
}
