package rest

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	lru "github.com/hashicorp/golang-lru"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/wasm/watcher"
	"strconv"
	"sync"
)

const (
	restCacheSize = 40960
)

var (
	once   sync.Once
	gCache *cache
)

type config struct {
	Size int
}

type cache struct {
	*lru.Cache
}

func newCache(conf config) *cache {
	if !watcher.Enable() {
		return nil
	}
	lruCache, err := lru.New(conf.Size)
	if err != nil {
		return nil
	}

	return &cache{
		lruCache,
	}
}

func cacheInst() *cache {
	once.Do(func() {
		gCache = newCache(config{Size: restCacheSize})
	})

	return gCache
}

func (c *cache) set(key common.Hash, value proto.Message) {
	if c == nil {
		return
	}
	c.Add(key, value)
}

func (c *cache) get(key common.Hash) (proto.Message, error) {
	if c == nil {
		return nil, fmt.Errorf("cache is unavaliable")
	}
	value, ok := c.Cache.Get(key)
	if ok {
		ret, ok := value.(proto.Message)
		if ok {
			return ret, nil
		}

		return ret, fmt.Errorf("wrong format")
	}

	return nil, fmt.Errorf("not found")
}

// grpcFn grpc call back
type grpcFn func() (proto.Message, error)

// queryWithCache f query from grpc then save it to lru.
func queryWithCache(request proto.Message, f grpcFn) (proto.Message, error) {
	var res proto.Message
	cacheKey := buildKey(request)
	cacheRes, err := cacheInst().get(cacheKey)
	if err == nil {
		res = cacheRes
	} else {
		res, err = f()

		if err != nil {
			return nil, err
		}
		cacheInst().set(cacheKey, res)
	}
	return res, nil
}

// buildKey consider proto.MessageName / request / height as hash key
func buildKey(request proto.Message) common.Hash {
	if !watcher.Enable() {
		return common.Hash{}
	}
	return sha256.Sum256([]byte(proto.MessageName(request) + request.String() + strconv.FormatInt(global.GetGlobalHeight(), 10)))
}
