package backend

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/x/evm/watcher"
)

type Cache interface {
	GetBlockByNumber(number uint64) (*watcher.Block, error)
	GetBlockByHash(hash common.Hash) (*watcher.Block, error)
	AddOrUpdateBlock(hash common.Hash, block *watcher.Block)
	GetTransaction(hash common.Hash) (*watcher.Transaction, error)
	AddOrUpdateTransaction(hash common.Hash, tx *watcher.Transaction)
	GetBlockHash(number uint64) (common.Hash, error)
	AddOrUpdateBlockHash(number uint64, hash common.Hash)
}
