package backend

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/okbchain/x/evm/watcher"
)

type Cache interface {
	GetBlockByNumber(number uint64, fullTx bool) (*watcher.Block, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (*watcher.Block, error)
	AddOrUpdateBlock(hash common.Hash, block *watcher.Block, fullTx bool)
	GetTransaction(hash common.Hash) (*watcher.Transaction, error)
	AddOrUpdateTransaction(hash common.Hash, tx *watcher.Transaction)
	GetBlockHash(number uint64) (common.Hash, error)
	AddOrUpdateBlockHash(number uint64, hash common.Hash)
}
