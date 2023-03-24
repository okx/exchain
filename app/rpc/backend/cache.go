package backend

import (
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/evm/watcher"
)

type Cache interface {
	GetBlockByNumber(number uint64, fullTx bool) (*evmtypes.Block, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (*evmtypes.Block, error)
	AddOrUpdateBlock(hash common.Hash, block *evmtypes.Block, fullTx bool)
	GetTransaction(hash common.Hash) (*watcher.Transaction, error)
	AddOrUpdateTransaction(hash common.Hash, tx *watcher.Transaction)
	GetBlockHash(number uint64) (common.Hash, error)
	AddOrUpdateBlockHash(number uint64, hash common.Hash)
}
