package api_cache

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

const (
	MsgLruNotInitialized = "Lru has not been Initialized"
	MsgLruDataNotFound   = "lru : not found"
	MsgLruDataWrongType  = "lru : wrong type"
)

var (
	LruKeyTx        = "LruKeyTx"
	LruKeyBlock     = "LruKeyBlock"
	LruKeyBlockInfo = "LruKeyBlockInfo"
)

type ApiLruCache struct {
	lruMap map[string]*lru.Cache
}

func NewApiLruCache() *ApiLruCache {

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
	return &ApiLruCache{lruMap: lruMap}
}

func (alc *ApiLruCache) getBytesFromLru(lruKey string, cacheKey interface{}) ([]byte, error) {
	data, err := alc.getDataFromLru(lruKey, cacheKey)
	if err != nil {
		return nil, err
	}
	dataBytes := data.([]byte)
	if alc == nil {
		return nil, errors.New(MsgLruDataWrongType)
	}
	return dataBytes, nil
}
func (alc *ApiLruCache) getDataFromLru(lruKey string, cacheKey interface{}) (interface{}, error) {

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
func (alc *ApiLruCache) addDataToLru(lruKey string, cacheKey interface{}, cacheData interface{}) {
	lru, ok := alc.lruMap[lruKey]
	if ok {
		return
	}
	lru.Add(cacheKey, cacheData)
}
func (alc *ApiLruCache) GetBlockByNumber(number uint64, fullTx bool) (interface{}, error) {
	hash, err := alc.GetBlockHash(number)
	if err != nil {
		return nil, err
	}
	return alc.GetBlockByHash(hash, fullTx)
}
func (alc *ApiLruCache) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	return alc.getDataFromLru(LruKeyBlock, hash)
}
func (alc *ApiLruCache) UpdateBlock(hash common.Hash, block interface{}) {
	alc.addDataToLru(LruKeyBlock, hash, block)
}
func (alc *ApiLruCache) GetTransaction(hash common.Hash) (*rpctypes.Transaction, error) {
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
func (alc *ApiLruCache) UpdateTransaction(hash common.Hash, tx *rpctypes.Transaction) {
	alc.addDataToLru(LruKeyTx, hash, tx)
}

func (alc *ApiLruCache) GetBlockHash(number uint64) (common.Hash, error) {
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
func (alc *ApiLruCache) UpdateBlockInfo(number uint64, hash common.Hash) {
	alc.addDataToLru(LruKeyBlockInfo, number, hash)
}
