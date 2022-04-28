package watcher

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
)

type WatchTx interface {
	GetTxWatchMessage() WatchMessage
	GetTxHash() common.Hash
	GetFailedReceipts(cumulativeGas, gasUsed uint64) WatchMessage
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
	if watchTx == nil {
		return
	}
	w.saveTx(watchTx)

	if resp != nil && !resp.IsOK() {
		w.saveFailedReceipts(watchTx, uint64(resp.GasUsed))
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

	if txWatchMessage := tx.GetTxWatchMessage(); txWatchMessage != nil {
		w.batch = append(w.batch, txWatchMessage)
	}
	w.blockTxs = append(w.blockTxs, tx.GetTxHash())
}

func (w *Watcher) saveFailedReceipts(watchTx WatchTx, gasUsed uint64) {
	if w == nil || watchTx == nil {
		return
	}
	w.UpdateCumulativeGas(watchTx.GetIndex(), gasUsed)
	if receipts := watchTx.GetFailedReceipts(w.cumulativeGas[watchTx.GetIndex()], gasUsed); receipts != nil {
		w.batch = append(w.batch, receipts)
	}
}
