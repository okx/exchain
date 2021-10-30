package websocket

import (
	"sync"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	pushservice "github.com/okex/exchain/x/stream/pushservice/types"
	"github.com/okex/exchain/x/stream/types"
)

type cache struct {
	depthBooksMap map[string]pushservice.BookRes
	lock          sync.RWMutex
}

var (
	singletonCache *cache
	once           sync.Once
)

func InitialCache(ctx sdk.Context, orderKeeper types.OrderKeeper, dexKeeper types.DexKeeper, logger log.Logger) {
	once.Do(func() {
		size := 200
		tokenPairs := dexKeeper.GetTokenPairs(ctx)
		logger.Debug("initial websocket cache", "tokenPairs", tokenPairs)
		depthBooksMap := make(map[string]pushservice.BookRes, len(tokenPairs))
		for _, tokenPair := range tokenPairs {
			depthBook := orderKeeper.GetDepthBookCopy(tokenPair.Name())
			bookRes := pushservice.ConvertBookRes(tokenPair.Name(), orderKeeper, depthBook, size)
			depthBooksMap[tokenPair.Name()] = bookRes
		}
		logger.Debug("initial websocket cache", "depthbook", depthBooksMap)
		singletonCache = &cache{
			depthBooksMap: depthBooksMap,
		}
	})
}

func GetDepthBookFromCache(product string) (depthBook pushservice.BookRes, ok bool) {
	singletonCache.lock.RLock()
	defer singletonCache.lock.RUnlock()
	depthBook, ok = singletonCache.depthBooksMap[product]
	return
}

func UpdateDepthBookCache(product string, bookRes pushservice.BookRes) {
	singletonCache.lock.Lock()
	defer singletonCache.lock.Unlock()
	singletonCache.depthBooksMap[product] = bookRes
}
