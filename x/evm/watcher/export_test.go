package watcher

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
)

// For Testing getRealTx in watcher_test package
func (w *Watcher) GetRealTx(tx tm.TxEssentials, txDecoder sdk.TxDecoder) (sdk.Tx, error) {
	return w.getRealTx(tx, txDecoder)
}

func (w *Watcher) CreateWatchTx(realTx sdk.Tx) WatchTx {
	return w.createWatchTx(realTx)
}

func NewTransactionReceipt(status uint32, tx *types.MsgEthereumTx, txHash, blockHash ethcommon.Hash, txIndex, height uint64, data *types.ResultData, cumulativeGas, GasUsed uint64) TransactionReceipt {
	return newTransactionReceipt(status, tx, txHash, blockHash, txIndex, height, data, cumulativeGas, GasUsed)
}

func (w *Watcher) GetBlockHash() ethcommon.Hash {
	return w.blockHash
}

func (w *Watcher) GetHeight() uint64 {
	return w.height
}

func (w *Watcher) GetCumulativeGas() map[uint64]uint64 {
	return w.cumulativeGas
}

func (w *Watcher) GetBatch() []WatchMessage {
	return w.batch
}

func (w *Watcher) GetBlockTxs() []ethcommon.Hash {
	return w.blockTxs
}

func (w *Watcher) GetBlockStdTxs() []ethcommon.Hash {
	return w.blockStdTxs
}

func (w *Watcher) GetHeader() tm.Header {
	return w.header
}

func (m MsgEthTx) GetTransaction() Transaction {
	return *(m.Transaction)
}

func (tx Transaction) GetTx() types.MsgEthereumTx {
	return *(tx.tx)
}

func (tx Transaction) GetHash() ethcommon.Hash {
	return tx.Hash
}

func (m MsgTransactionReceipt) GetTxHash() []byte {
	return m.txHash
}

func (wtx evmTx) SetIndex(index uint64) {
	wtx.index = index
	return
}

// Create the same Watcher as the tested function
// with no evmTxIndex increase
func (w *Watcher) CreateExpectedWatchTx(realTx sdk.Tx) WatchTx {
	var txMsg WatchTx
	switch realTx.GetType() {
	case sdk.EvmTxType:
		evmTx, err := w.extractEvmTx(realTx)
		if err != nil {
			return nil
		}
		txMsg = NewEvmTx(evmTx, ethcommon.BytesToHash(evmTx.TxHash()), w.blockHash, w.height, w.evmTxIndex-1)
	}

	return txMsg
}
