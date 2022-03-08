package backend_cache

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

const (
	MsgLruNotInitialized = "lru has not been Initialized"
	MsgLruDataNotFound   = "lru : not found"
	MsgLruDataWrongType  = "lru : wrong type"
)

var (
	LruKeyTx        = "LruKeyTx"
	LruKeyBlock     = "LruKeyBlock"
	LruKeyBlockInfo = "LruKeyBlockInfo"
)

type BackendLruCache struct {
	lruMap map[string]*lru.Cache
}

func NewBackendLruCache() *BackendLruCache {

	lruMap := map[string]*lru.Cache{}
	watcherLruSize := viper.GetInt(watcher.FlagFastQueryLru)
	//init lru cache for tx
	lruCache, err := lru.New(watcherLruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Tx, err :" + err.Error()))
	}
	lruMap[LruKeyTx] = lruCache
	//init lru cache for block
	lruCache, err = lru.New(watcherLruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Block, err :" + err.Error()))
	}
	lruMap[LruKeyBlock] = lruCache
	//init lru cache for blockinfo
	lruCache, err = lru.New(watcherLruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Block, err :" + err.Error()))
	}
	lruMap[LruKeyBlockInfo] = lruCache
	return &BackendLruCache{lruMap: lruMap}
}

func (alc *BackendLruCache) getDataFromLru(lruKey string, cacheKey interface{}) (interface{}, error) {

	lru, ok := alc.lruMap[lruKey]
	if !ok {
		return nil, errors.New(MsgLruNotInitialized)
	}

	cacheCode, ok := lru.Get(cacheKey)
	if ok {
		return cacheCode, nil
	}
	return nil, errors.New(MsgLruDataNotFound)
}
func (alc *BackendLruCache) addDataToLru(lruKey string, cacheKey interface{}, cacheData interface{}) {
	lru, ok := alc.lruMap[lruKey]
	if ok {
		lru.PeekOrAdd(cacheKey, cacheData)
	}
}
func (alc *BackendLruCache) GetBlockByNumber(number uint64) (*rpctypes.Block, error) {
	hash, err := alc.GetBlockHash(number)
	if err != nil {
		return nil, err
	}
	return alc.GetBlockByHash(hash)
}
func (alc *BackendLruCache) GetBlockByHash(hash common.Hash) (*rpctypes.Block, error) {
	data, err := alc.getDataFromLru(LruKeyBlock, hash)
	if err != nil {
		return nil, err
	}
	res := data.(*rpctypes.Block)
	if res != nil {
		return res, nil
	}
	return nil, errors.New(MsgLruDataWrongType)
}
func (alc *BackendLruCache) AddOrUpdateBlock(hash common.Hash, block *rpctypes.Block) {
	alc.addDataToLru(LruKeyBlock, hash, block)
	alc.AddOrUpdateBlockHash(uint64(block.Number), hash)
	if block.Transactions != nil {
		txs, ok := block.Transactions.([]*rpctypes.Transaction)
		if ok {
			for _, tx := range txs {
				alc.AddOrUpdateTransaction(tx.Hash, tx)
			}
		}
	}
}
func (alc *BackendLruCache) GetTransaction(hash common.Hash) (*rpctypes.Transaction, error) {
	data, err := alc.getDataFromLru(LruKeyTx, hash)
	if err != nil {
		return nil, err
	}
	tx := data.(*rpctypes.Transaction)
	if alc == nil {
		return nil, errors.New(MsgLruDataWrongType)
	}
	return tx, nil
}
func (alc *BackendLruCache) AddOrUpdateTransaction(hash common.Hash, tx *rpctypes.Transaction) {
	alc.addDataToLru(LruKeyTx, hash, tx)
}
func (alc *BackendLruCache) GetBlockHash(number uint64) (common.Hash, error) {
	data, err := alc.getDataFromLru(LruKeyBlockInfo, number)
	if err != nil {
		return common.Hash{}, err
	}
	dataHash := data.(common.Hash)
	if alc == nil {
		return common.Hash{}, errors.New(MsgLruDataWrongType)
	}
	return dataHash, nil
}
func (alc *BackendLruCache) AddOrUpdateBlockHash(number uint64, hash common.Hash) {
	alc.addDataToLru(LruKeyBlockInfo, number, hash)
}
