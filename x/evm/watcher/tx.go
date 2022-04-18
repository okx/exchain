package watcher

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
	"sync/atomic"
)

type DeliverTx struct {
	Req  tm.TxEssentials
	Resp *tm.ResponseDeliverTx
}

func (w *Watcher) RecordABCIMessage(deliverTx *DeliverTx, txDecoder sdk.TxDecoder) {
	if !w.Enabled() {
		return
	}

	defer atomic.AddInt64(&w.recordedTxsCount, 1)
	index := w.txIndexInBlock
	w.dispatchTxJob(func() {
		w.recordTxsAndReceipts(deliverTx, index, txDecoder)
	})
	w.txIndexInBlock++
}

func (w *Watcher) recordTxsAndReceipts(deliverTx *DeliverTx, index uint64, txDecoder sdk.TxDecoder) {
	defer atomic.AddInt64(&w.recordedTxsCount, -1)
	if deliverTx == nil || deliverTx.Req == nil || deliverTx.Resp == nil {
		w.log.Error("watch parse abci message error", "input", deliverTx)
		return
	}

	realTx, err := txDecoder(deliverTx.Req.GetRaw())
	if err != nil {
		w.log.Error("watch decode deliver tx", "error", err)
		return
	}

	if realTx.GetType() != sdk.EvmTxType {
		return
	}

	if deliverTx.Resp.Code != errors.SuccessABCICode {
		w.saveEvmTxAndFailedReceipt(realTx, index, uint64(deliverTx.Resp.GasUsed))
		return
	}

	w.saveEvmTxAndSuccessReceipt(realTx, deliverTx.Resp.Data, index, uint64(deliverTx.Resp.GasUsed))
}

func (w *Watcher) saveEvmTxAndFailedReceipt(sdkTx sdk.Tx, index, gasUsed uint64) {
	evmTx, err := w.extractEvmTx(sdkTx)
	if err != nil {
		return
	}
	txHash := ethcmn.BytesToHash(evmTx.TxHash())
	ethcmn.BytesToHash(evmTx.TxHash())

	w.saveTxAndReceipt(evmTx, txHash, index, TransactionFailed, &types.ResultData{}, gasUsed)
}

func (w *Watcher) saveEvmTxAndSuccessReceipt(sdkTx sdk.Tx, resultData []byte, index, gasUsed uint64) {
	evmTx, err := w.extractEvmTx(sdkTx)
	if err != nil && evmTx != nil {
		return
	}
	evmResultData, err := types.DecodeResultData(resultData)
	if err != nil {
		w.log.Error("save evm tx and success receipt error", "height", w.height, "index", index, "error", err)
		return
	}

	w.saveTxAndReceipt(evmTx, evmResultData.TxHash, index, TransactionSuccess, &evmResultData, gasUsed)
}

func (w *Watcher) extractEvmTx(sdkTx sdk.Tx) (*types.MsgEthereumTx, error) {
	var ok bool
	var evmTx *types.MsgEthereumTx
	for _, msg := range sdkTx.GetMsgs() {
		if evmTx, ok = msg.(*types.MsgEthereumTx); !ok {
			return nil, fmt.Errorf("sdktx is not evm tx %v", sdkTx)
		}
	}

	return evmTx, nil
}

func (w *Watcher) saveTxAndReceiptEx(msg *types.MsgEthereumTx, txHash ethcmn.Hash, index uint64,
	receiptStatus uint32, data *types.ResultData, gasUsed uint64) {
	w.txMutex.Lock()
	defer w.txMutex.Unlock()

	wMsg := NewMsgEthTx(msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.txsAndReceipts = append(w.txsAndReceipts, wMsg)
	}
	txReceipt := NewMsgTransactionReceipt(receiptStatus, msg, txHash, w.blockHash, index, w.height, data, w.cumulativeGas[index], gasUsed)
	if txReceipt != nil {
		w.txsAndReceipts = append(w.txsAndReceipts, txReceipt)
	}
	w.UpdateCumulativeGas(index, gasUsed)
	w.txsInBlock = append(w.txsInBlock, TxIndex{TxHash: txHash, Index: index})
}

func (w *Watcher) saveTxAndReceipt(msg *types.MsgEthereumTx, txHash ethcmn.Hash, index uint64,
	receiptStatus uint32, data *types.ResultData, gasUsed uint64) {

	wMsg := NewMsgEthTx(msg, txHash, w.blockHash, w.height, index)
	txReceipt := NewMsgTransactionReceipt(receiptStatus, msg, txHash, w.blockHash, index, w.height, data, w.cumulativeGas[index], gasUsed)
	select {
	case w.txResultChan <- TxResult{
		TxMsg:     wMsg,
		TxReceipt: txReceipt,
		Index:     index,
		TxHash:    txHash,
		GasUsed:   gasUsed,
	}:
	default:
		w.log.Error("save to tx too busy.")
		go func() {
			w.txResultChan <- TxResult{
				TxMsg:     wMsg,
				TxReceipt: txReceipt,
				Index:     index,
				TxHash:    txHash,
				GasUsed:   gasUsed,
			}
		}()
	}
}
