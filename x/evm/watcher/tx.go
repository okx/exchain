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
	atomic.AddInt64(&w.recordingTxsCount, 1)
	index := w.txIndexInBlock
	w.dispatchTxJob(func() {
		w.recordTxsAndReceipts(deliverTx, index, txDecoder)
	})
	w.txIndexInBlock++
}

func (w *Watcher) recordTxsAndReceipts(deliverTx *DeliverTx, index uint64, txDecoder sdk.TxDecoder) {
	defer atomic.AddInt64(&w.recordingTxsCount, -1)
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
	w.txMutex.Lock()
	defer w.txMutex.Unlock()

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
	w.saveEvmTx(evmTx, txHash, index)
	w.saveTransactionReceipt(TransactionFailed, evmTx, txHash, index, &types.ResultData{}, gasUsed)
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
	w.saveEvmTx(evmTx, evmResultData.TxHash, index)
	w.saveTransactionReceipt(TransactionSuccess, evmTx, evmResultData.TxHash, index, &evmResultData, gasUsed)
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

func (w *Watcher) saveEvmTx(msg *types.MsgEthereumTx, txHash ethcmn.Hash, index uint64) {
	wMsg := NewMsgEthTx(msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.txsAndReceipts = append(w.txsAndReceipts, wMsg)
	}
	w.txsInBlock = append(w.txsInBlock, TxIndex{TxHash: txHash, Index: index})
}

func (w *Watcher) saveTransactionReceipt(status uint32, msg *types.MsgEthereumTx, txHash ethcmn.Hash, txIndex uint64, data *types.ResultData, gasUsed uint64) {
	w.UpdateCumulativeGas(txIndex, gasUsed)
	wMsg := NewMsgTransactionReceipt(status, msg, txHash, w.blockHash, txIndex, w.height, data, w.cumulativeGas[txIndex], gasUsed)
	if wMsg != nil {
		w.txsAndReceipts = append(w.txsAndReceipts, wMsg)
	}
}
