package watcher

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher/abci"
)

func (w *Watcher) ReceiveABCIMessage(deliverTx *abci.DeliverTx, txsCount int, txDecoder sdk.TxDecoder) {
	w.txsCount = txsCount
	w.dispatchJob(func() {
		w.parseABCIMessage(deliverTx, txDecoder)
	})
}

func (w *Watcher) parseABCIMessage(deliverTx *abci.DeliverTx, txDecoder sdk.TxDecoder) {
	if deliverTx == nil || deliverTx.Req == nil || deliverTx.Resp == nil {
		w.log.Error("watch parse abci message error", "input", deliverTx)
		return
	}

	realTx, err := txDecoder(deliverTx.Req.Tx)
	if err != nil {
		w.log.Error("watch decode deliver tx", "error", err)
		return
	}

	if deliverTx.Resp.Code != errors.SuccessABCICode {
		w.saveEvmTxAndFailedReceipt(realTx, deliverTx.Index, uint64(deliverTx.Resp.GasUsed))
		return
	}

	w.saveEvmTxAndSuccessReceipt(realTx, deliverTx.Resp.Data, deliverTx.Index, uint64(deliverTx.Resp.GasUsed))
}

func (w *Watcher) extractEvmTx(sdkTx sdk.Tx) (msg *types.MsgEthereumTx, err error) {
	var ok bool
	for _, msg := range sdkTx.GetMsgs() {
		if msg, ok = msg.(*types.MsgEthereumTx); !ok {
			return nil, fmt.Errorf("sdktx is not evm tx %v", sdkTx)
		}
	}

	return
}

func (w *Watcher) saveEvmTx(msg *types.MsgEthereumTx, txHash ethcmn.Hash, index uint64) {
	wMsg := NewMsgEthTx(msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.txs = append(w.batch, wMsg)
	}
	w.UpdateBlockTxs(txHash)
}

func (w *Watcher) saveTransactionReceipt(status uint32, msg *types.MsgEthereumTx, txHash ethcmn.Hash, txIndex uint64, data *types.ResultData, gasUsed uint64) {
	w.UpdateCumulativeGas(txIndex, gasUsed)
	wMsg := NewMsgTransactionReceipt(status, msg, txHash, w.blockHash, txIndex, w.height, data, w.cumulativeGas[txIndex], gasUsed)
	if wMsg != nil {
		w.txReceipts = append(w.batch, wMsg)
	}
}

func (w *Watcher) saveEvmTxAndSuccessReceipt(sdkTx sdk.Tx, resultData []byte, index, gasUsed uint64) {
	evmTx, err := w.extractEvmTx(sdkTx)
	if err != nil {
		return
	}
	evmResultData, err := types.DecodeResultData(resultData)
	w.saveEvmTx(evmTx, evmResultData.TxHash, index)
	w.saveTransactionReceipt(TransactionSuccess, evmTx, evmResultData.TxHash, index, &evmResultData, gasUsed)
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
