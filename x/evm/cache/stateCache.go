package cache

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	lru "github.com/hashicorp/golang-lru"
)

//the default lru cache size is 1kw, that means the max memory size we needs is (32 + 32 + 4) * 10000000, about 700MB
var (
	defaultCacheSize int = 100000000
	gStateCache      *lru.Cache
	once             sync.Once
	stateCache       Monitor
)

func InstanceOfMonitor() *Monitor {
	return &stateCache
}

func instanceOfStateLru() *lru.Cache {
	once.Do(func() {
		var e error = nil
		gStateCache, e = lru.New(defaultCacheSize)
		if e != nil {
			panic(errors.New("Failed to call InstanceOfStateLru cause :" + e.Error()))
		}
	})
	return gStateCache
}

var MatchCounter = 0
var TotalCounter = 0
var SetCounter = 0

func GetStateFromCache(key common.Hash) []byte {
	cache := instanceOfStateLru()
	if cache == nil {
		return nil
	}
	TotalCounter++
	value, ok := cache.Get(key)
	if ok {
		ret, ok := value.([]byte)
		if ok {
			MatchCounter++
			return ret
		}
	}
	return nil
}

func SetStateToCache(key common.Hash, value []byte) {
	cache := instanceOfStateLru()
	if cache == nil {
		return
	}
	SetCounter++
	cache.Add(key, value)
}
