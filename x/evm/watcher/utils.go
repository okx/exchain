package watcher

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// NewTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewTransaction(tx *evmtypes.MsgEthereumTx, txHash, blockHash common.Hash, blockNumber, index uint64) (*Transaction, error) {
	// Verify signature and retrieve sender address
	err := tx.VerifySig(tx.ChainID(), int64(blockNumber))
	if err != nil {
		return nil, err
	}

	rpcTx := &Transaction{
		From:     common.HexToAddress(tx.GetFrom()),
		Gas:      hexutil.Uint64(tx.Data.GasLimit),
		GasPrice: (*hexutil.Big)(tx.Data.Price),
		Hash:     txHash,
		Input:    hexutil.Bytes(tx.Data.Payload),
		Nonce:    hexutil.Uint64(tx.Data.AccountNonce),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Data.Amount),
		V:        (*hexutil.Big)(tx.Data.V),
		R:        (*hexutil.Big)(tx.Data.R),
		S:        (*hexutil.Big)(tx.Data.S),
	}

	if blockHash != (common.Hash{}) {
		rpcTx.BlockHash = &blockHash
		rpcTx.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		rpcTx.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	return rpcTx, nil
}
