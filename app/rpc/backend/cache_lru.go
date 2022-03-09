package backend

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

const (
	MsgLruNotInitialized = "lru has not been Initialized"
	MsgLruDataNotFound   = "lru : not found"
	MsgLruDataWrongType  = "lru : wrong type"
)
const (
	FlagApiBackendLru = "api-backend-lru"
)

type LruCache struct {
	lruTx        *lru.Cache
	lruBlock     *lru.Cache
	lruBlockInfo *lru.Cache
}

func NewLruCache() *LruCache {
	lruSize := viper.GetInt(FlagApiBackendLru)
	//init lru cache for tx
	lruTx, err := lru.New(lruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Tx, err :" + err.Error()))
	}
	//init lru cache for block
	lruBlock, err := lru.New(lruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Block, err :" + err.Error()))
	}
	//init lru cache for blockinfo
	lruBlockInfo, err := lru.New(lruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Block, err :" + err.Error()))
	}
	return &LruCache{
		lruTx:        lruTx,
		lruBlock:     lruBlock,
		lruBlockInfo: lruBlockInfo,
	}
}

func (alc *LruCache) GetBlockByNumber(number uint64) (*watcher.Block, error) {
	hash, err := alc.GetBlockHash(number)
	if err != nil {
		return nil, err
	}
	return alc.GetBlockByHash(hash)
}
func (alc *LruCache) GetBlockByHash(hash common.Hash) (*watcher.Block, error) {
	data, ok := alc.lruBlock.Get(hash)
	if !ok {
		return nil, errors.New(MsgLruDataNotFound)
	}
	res, ok := data.(*watcher.Block)
	if !ok {
		return nil, errors.New(MsgLruDataWrongType)
	}
	return res, nil
}
func (alc *LruCache) AddOrUpdateBlock(hash common.Hash, block *watcher.Block) {
	alc.lruBlock.PeekOrAdd(hash, block)
	alc.AddOrUpdateBlockHash(uint64(block.Number), hash)
	if block.Transactions != nil {
		txs, ok := block.Transactions.([]*watcher.Transaction)
		if ok {
			for _, tx := range txs {
				alc.AddOrUpdateTransaction(tx.Hash, tx)
			}
		}
	}
}
func (alc *LruCache) GetTransaction(hash common.Hash) (*watcher.Transaction, error) {
	data, ok := alc.lruTx.Get(hash)
	if !ok {
		return nil, errors.New(MsgLruDataNotFound)
	}
	tx, ok := data.(*watcher.Transaction)
	if !ok {
		return nil, errors.New(MsgLruDataWrongType)
	}
	return tx, nil
}
func (alc *LruCache) AddOrUpdateTransaction(hash common.Hash, tx *watcher.Transaction) {
	alc.lruTx.PeekOrAdd(hash, tx)
}
func (alc *LruCache) GetBlockHash(number uint64) (common.Hash, error) {
	data, ok := alc.lruBlockInfo.Get(number)
	if !ok {
		return common.Hash{}, errors.New(MsgLruDataNotFound)
	}
	dataHash, ok := data.(common.Hash)
	if !ok {
		return common.Hash{}, errors.New(MsgLruDataWrongType)
	}
	return dataHash, nil
}
func (alc *LruCache) AddOrUpdateBlockHash(number uint64, hash common.Hash) {
	alc.lruBlockInfo.PeekOrAdd(number, hash)
}
