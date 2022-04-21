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
	w.dispatchJob(func() {
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

	if deliverTx.Resp.Code != errors.SuccessABCICode {
		w.saveEvmTxAndFailedReceipt(realTx, index, uint64(deliverTx.Resp.GasUsed))
		return
	}

	w.saveEvmTxAndSuccessReceipt(realTx, deliverTx.Resp.Data, index, uint64(deliverTx.Resp.GasUsed))
}

func (w *Watcher) saveEvmTxAndFailedReceipt(sdkTx sdk.Tx, index, gasUsed uint64) {
	evmTx, err := w.extractEvmTx(sdkTx)
	if err != nil {
		w.log.Error("save evm tx and failed receipt error", "height", w.height, "index", index, "error", err)
		return
	}
	txHash := ethcmn.BytesToHash(evmTx.TxHash())

	w.saveTxAndReceipt(evmTx, txHash, index, TransactionFailed, &types.ResultData{}, gasUsed)
}

func (w *Watcher) saveEvmTxAndSuccessReceipt(sdkTx sdk.Tx, resultData []byte, index, gasUsed uint64) {
	evmTx, err := w.extractEvmTx(sdkTx)
	if err != nil {
		w.log.Error("save evm tx and success receipt error", "height", w.height, "index", index, "error", err)
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

func (w *Watcher) saveTxAndReceipt(msg *types.MsgEthereumTx, txHash ethcmn.Hash, index uint64,
	receiptStatus uint32, data *types.ResultData, gasUsed uint64) {
	//	w.txsMutex.Lock()
	//	defer w.txsMutex.Unlock()

	wMsg := NewMsgEthTx(msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.txs = append(w.txs, wMsg)
	}
	txReceipt := newTransactionReceipt(receiptStatus, msg, txHash, w.blockHash, index, w.height, data, gasUsed)
	if txReceipt != nil {
		w.txReceipts = append(w.txReceipts, txReceipt)
	}
	w.txInfoCollector = append(w.txInfoCollector, TxInfo{TxHash: txHash, Index: index, GasUsed: gasUsed})
}
