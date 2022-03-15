package backend

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

var ErrLruNotInitialized = errors.New("lru has not been Initialized")
var ErrLruDataNotFound = errors.New("lru : not found")
var ErrLruDataWrongType = errors.New("lru : wrong type")

const (
	FlagApiBackendBlockLruCache = "rpc-block-cache"
	FlagApiBackendTxLruCache    = "rpc-tx-cache"
)

type LruCache struct {
	lruTx        *lru.Cache
	lruBlock     *lru.Cache
	lruBlockInfo *lru.Cache
}

func NewLruCache() *LruCache {
	blockLruSize := viper.GetInt(FlagApiBackendBlockLruCache)
	txLruSize := viper.GetInt(FlagApiBackendTxLruCache)
	//init lru cache for tx
	lruTx, err := lru.New(txLruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Tx, err :" + err.Error()))
	}
	//init lru cache for block
	lruBlock, err := lru.New(blockLruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Block, err :" + err.Error()))
	}
	//init lru cache for blockinfo
	lruBlockInfo, err := lru.New(blockLruSize)
	if err != nil {
		panic(errors.New("Failed to init LRU for Block, err :" + err.Error()))
	}
	return &LruCache{
		lruTx:        lruTx,
		lruBlock:     lruBlock,
		lruBlockInfo: lruBlockInfo,
	}
}

func (lc *LruCache) GetBlockByNumber(number uint64) (*watcher.Block, error) {
	hash, err := lc.GetBlockHash(number)
	if err != nil {
		return nil, err
	}
	return lc.GetBlockByHash(hash)
}
func (lc *LruCache) GetBlockByHash(hash common.Hash) (*watcher.Block, error) {
	data, ok := lc.lruBlock.Get(hash)
	if !ok {
		return nil, ErrLruDataNotFound
	}
	res, ok := data.(*watcher.Block)
	if !ok {
		return nil, ErrLruDataWrongType
	}
	return res, nil
}
func (lc *LruCache) AddOrUpdateBlock(hash common.Hash, block *watcher.Block) {
	lc.lruBlock.PeekOrAdd(hash, block)
	lc.AddOrUpdateBlockHash(uint64(block.Number), hash)
	if block.Transactions != nil {
		txs, ok := block.Transactions.([]*watcher.Transaction)
		if ok {
			for _, tx := range txs {
				lc.AddOrUpdateTransaction(tx.Hash, tx)
			}
		}
	}
}
func (lc *LruCache) GetTransaction(hash common.Hash) (*watcher.Transaction, error) {
	data, ok := lc.lruTx.Get(hash)
	if !ok {
		return nil, ErrLruDataNotFound
	}
	tx, ok := data.(*watcher.Transaction)
	if !ok {
		return nil, ErrLruDataWrongType
	}
	return tx, nil
}
func (lc *LruCache) AddOrUpdateTransaction(hash common.Hash, tx *watcher.Transaction) {
	lc.lruTx.PeekOrAdd(hash, tx)
}
func (lc *LruCache) GetBlockHash(number uint64) (common.Hash, error) {
	data, ok := lc.lruBlockInfo.Get(number)
	if !ok {
		return common.Hash{}, ErrLruDataNotFound
	}
	dataHash, ok := data.(common.Hash)
	if !ok {
		return common.Hash{}, ErrLruDataWrongType
	}
	return dataHash, nil
}
func (lc *LruCache) AddOrUpdateBlockHash(number uint64, hash common.Hash) {
	lc.lruBlockInfo.PeekOrAdd(number, hash)
}
