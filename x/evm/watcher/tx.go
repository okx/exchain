package watcher

import (
	"fmt"
	syslog "log"
	"runtime/debug"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	"github.com/okex/exchain/x/evm/types"
)

type WatchTx interface {
	GetTxWatchMessage() WatchMessage
	GetTransaction() *Transaction
	GetTxHash() common.Hash
	GetFailedReceipts(cumulativeGas, gasUsed uint64) *TransactionReceipt
	GetIndex() uint64
}

func (w *Watcher) RecordTxAndFailedReceipt(tx tm.TxEssentials, resp *tm.ResponseDeliverTx, txDecoder sdk.TxDecoder) {
	if !w.Enabled() {
		return
	}

	realTx, err := w.getRealTx(tx, txDecoder)
	if err != nil {
		return
	}
	watchTx := w.createWatchTx(realTx)
	switch realTx.GetType() {
	case sdk.EvmTxType:
		if watchTx == nil {
			return
		}
		w.saveTx(watchTx)
		if resp != nil && !resp.IsOK() {
			w.saveFailedReceipts(watchTx, uint64(resp.GasUsed))
		}
	case sdk.StdTxType:
		w.blockStdTxs = append(w.blockStdTxs, common.BytesToHash(realTx.TxHash()))
		txResult := &ctypes.ResultTx{
			Hash:     tx.TxHash(),
			Height:   int64(w.height),
			TxResult: *resp,
			Tx:       tx.GetRaw(),
		}
		w.saveStdTxResponse(txResult)
	}
}

func (w *Watcher) getRealTx(tx tm.TxEssentials, txDecoder sdk.TxDecoder) (sdk.Tx, error) {
	var err error
	realTx, _ := tx.(sdk.Tx)
	if realTx == nil {
		realTx, err = txDecoder(tx.GetRaw())
		if err != nil {
			return nil, err
		}
	}

	return realTx, nil
}

func (w *Watcher) createWatchTx(realTx sdk.Tx) WatchTx {
	var txMsg WatchTx
	switch realTx.GetType() {
	case sdk.EvmTxType:
		evmTx, err := w.extractEvmTx(realTx)
		if err != nil {
			return nil
		}
		txMsg = NewEvmTx(evmTx, common.BytesToHash(evmTx.TxHash()), w.blockHash, w.height, w.evmTxIndex)
		w.evmTxIndex++
	}

	return txMsg
}

func (w *Watcher) extractEvmTx(sdkTx sdk.Tx) (*types.MsgEthereumTx, error) {
	var ok bool
	var evmTx *types.MsgEthereumTx
	// stdTx should only have one tx
	msg := sdkTx.GetMsgs()
	if len(msg) <= 0 {
		return nil, fmt.Errorf("can not extract evm tx, len(msg) <= 0")
	}
	if evmTx, ok = msg[0].(*types.MsgEthereumTx); !ok {
		return nil, fmt.Errorf("sdktx is not evm tx %v", sdkTx)
	}

	return evmTx, nil
}

func (w *Watcher) saveTx(tx WatchTx) {
	if w == nil || tx == nil {
		return
	}
	if w.InfuraKeeper != nil {
		ethTx := tx.GetTransaction()
		if ethTx != nil {
			w.InfuraKeeper.OnSaveTransaction(*ethTx)
		}
	}
	if txWatchMessage := tx.GetTxWatchMessage(); txWatchMessage != nil {
		w.batch = append(w.batch, txWatchMessage)
	}
	w.blockTxs = append(w.blockTxs, tx.GetTxHash())
	syslog.Printf("lcm, append tx(%s) to block(%d), count=%d\n", tx.GetTxHash().Hex(), w.height, len(w.blockTxs))
	if tx.GetTxHash().Hex() == "0x98b98b00d52d1ebc840c82f3da20cf84b291ac45fa34f65eb819c0de3c15e473" || tx.GetTxHash().Hex() == "0xf418fcc3a488ffa8903faf0d768d99db0f424615de2910a2f684d22e2981ef3c" || tx.GetTxHash().Hex() == "0xc225143c648432e4598fc7958f6399fb9ae711f29bf65c77eaacc760bb76ea72" {
		syslog.Printf("lcm block(%d) tx = %v\n", w.height, tx)
		syslog.Printf("lcm block(%d) debug(%s)", w.height, string(debug.Stack()))
	}
}

func (w *Watcher) saveFailedReceipts(watchTx WatchTx, gasUsed uint64) {
	if w == nil || watchTx == nil {
		return
	}
	w.UpdateCumulativeGas(watchTx.GetIndex(), gasUsed)
	receipt := watchTx.GetFailedReceipts(w.cumulativeGas[watchTx.GetIndex()], gasUsed)
	if w.InfuraKeeper != nil {
		w.InfuraKeeper.OnSaveTransactionReceipt(*receipt)
	}
	wMsg := NewMsgTransactionReceipt(*receipt, watchTx.GetTxHash())
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

// SaveParallelTx saves parallel transactions and transactionReceipts to watcher
func (w *Watcher) SaveParallelTx(realTx sdk.Tx, resultData *types.ResultData, resp tm.ResponseDeliverTx) {

	if !w.Enabled() {
		return
	}

	switch realTx.GetType() {
	case sdk.EvmTxType:
		msgs := realTx.GetMsgs()
		evmTx, ok := msgs[0].(*types.MsgEthereumTx)
		if !ok {
			return
		}
		watchTx := NewEvmTx(evmTx, common.BytesToHash(evmTx.TxHash()), w.blockHash, w.height, w.evmTxIndex)
		w.evmTxIndex++
		w.saveTx(watchTx)

		// save transactionReceipts
		if resp.IsOK() && resultData != nil {
			w.SaveTransactionReceipt(TransactionSuccess, evmTx, watchTx.GetTxHash(), watchTx.GetIndex(), resultData, uint64(resp.GasUsed))
		} else {
			w.saveFailedReceipts(watchTx, uint64(resp.GasUsed))
		}
	case sdk.StdTxType:
		w.blockStdTxs = append(w.blockStdTxs, common.BytesToHash(realTx.TxHash()))
		txResult := &ctypes.ResultTx{
			Hash:     realTx.TxHash(),
			Height:   int64(w.height),
			TxResult: resp,
			Tx:       realTx.GetRaw(),
		}
		w.saveStdTxResponse(txResult)
	}
}
