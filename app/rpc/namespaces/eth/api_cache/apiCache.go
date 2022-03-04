package api_cache

import (
	"github.com/ethereum/go-ethereum/common"
	rpctypes "github.com/okex/exchain/app/rpc/types"
)

type ApiCache interface {
	GetBlockByNumber(number uint64, fullTx bool) (interface{}, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error)
	UpdateBlock(hash common.Hash, block interface{})
	GetTransaction(hash common.Hash) (*rpctypes.Transaction, error)
	UpdateTransaction(hash common.Hash, tx *rpctypes.Transaction)
	GetBlockHash(number uint64) (common.Hash, error)
	UpdateBlockInfo(number uint64, hash common.Hash)
}
