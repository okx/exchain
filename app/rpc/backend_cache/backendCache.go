package backend_cache

import (
	"github.com/ethereum/go-ethereum/common"
	rpctypes "github.com/okex/exchain/app/rpc/types"
)

type BackendCache interface {
	GetBlockByNumber(number uint64) (*rpctypes.Block, error)
	GetBlockByHash(hash common.Hash) (*rpctypes.Block, error)
	AddOrUpdateBlock(hash common.Hash, block *rpctypes.Block)
	GetTransaction(hash common.Hash) (*rpctypes.Transaction, error)
	AddOrUpdateTransaction(hash common.Hash, tx *rpctypes.Transaction)
	GetBlockHash(number uint64) (common.Hash, error)
	AddOrUpdateBlockHash(number uint64, hash common.Hash)
}
