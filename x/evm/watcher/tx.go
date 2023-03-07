package watcher

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	tm "github.com/okx/okbchain/libs/tendermint/abci/types"
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	"github.com/okx/okbchain/x/evm/types"
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
	// record ResultTx always
	if resp != nil {
		txResult := &ctypes.ResultTx{
			Hash:     tx.TxHash(),
			Height:   int64(w.height),
			TxResult: *resp,
			Tx:       tx.GetRaw(),
		}
		w.saveStdTxResponse(txResult)
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
			return
		}
		if resp != nil && resp.IsOK() && !w.IsRealEvmTx(resp) { // for evm2cm
			msgs := realTx.GetMsgs()
			if len(msgs) == 0 {
				return
			}
			evmTx, ok := msgs[0].(*types.MsgEthereumTx)
			if !ok {
				return
			}
			w.SaveTransactionReceipt(TransactionSuccess, evmTx, watchTx.GetTxHash(), watchTx.GetIndex(), &types.ResultData{}, uint64(resp.GasUsed))
		}
	case sdk.StdTxType:
		w.blockStdTxs = append(w.blockStdTxs, common.BytesToHash(realTx.TxHash()))
	}
}

func (w *Watcher) IsRealEvmTx(resp *tm.ResponseDeliverTx) bool {
	for _, ev := range resp.Events {
		if ev.Type == sdk.EventTypeMessage {
			for _, attr := range ev.Attributes {
				if string(attr.Key) == sdk.AttributeKeyModule &&
					string(attr.Value) == types.AttributeValueCategory {
					return true
				}
			}
		}
	}
	return false
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

	// record ResultTx always
	txResult := &ctypes.ResultTx{
		Hash:     realTx.TxHash(),
		Height:   int64(w.height),
		TxResult: resp,
		Tx:       realTx.GetRaw(),
	}
	w.saveStdTxResponse(txResult)

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
		if resp.IsOK() {
			if resultData == nil {
				resultData = &types.ResultData{}
			}
			w.SaveTransactionReceipt(TransactionSuccess, evmTx, watchTx.GetTxHash(), watchTx.GetIndex(), resultData, uint64(resp.GasUsed))
		} else {
			w.saveFailedReceipts(watchTx, uint64(resp.GasUsed))
		}
	case sdk.StdTxType:
		w.blockStdTxs = append(w.blockStdTxs, common.BytesToHash(realTx.TxHash()))
	}
}
