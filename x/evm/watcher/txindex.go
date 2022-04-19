package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"sort"
)

const (
	DefaultTxResultChanBuffer = 20
)

type TxInfo struct {
	TxHash  ethcmn.Hash
	Index   uint64
	GasUsed uint64
}

func (w *Watcher) addTxsToBlock() {
	sort.Slice(w.txsInBlock, func(i, j int) bool {
		return w.txsInBlock[i].Index < w.txsInBlock[j].Index
	})
	sort.Slice(w.txReceipts, func(i, j int) bool {
		return w.txReceipts[i].TransactionIndex < w.txReceipts[j].TransactionIndex
	})
	for _, txInfo := range w.txsInBlock {
		w.blockTxs = append(w.blockTxs, txInfo.TxHash)
		w.updateCumulativeGas(txInfo.Index, txInfo.GasUsed)
	}
	for _, receipt := range w.txReceipts {
		receipt.CumulativeGasUsed = hexutil.Uint64(w.cumulativeGas[uint64(receipt.TransactionIndex)])
		w.batch = append(w.batch, &MsgTransactionReceipt{txHash: receipt.TxHash.Bytes(), baseLazyMarshal: newBaseLazyMarshal(receipt)})
	}
}

type TxResult struct {
	TxMsg     WatchMessage
	TxReceipt *TransactionReceipt
	Index     uint64
	GasUsed   uint64
	TxHash    ethcmn.Hash
}

func (w *Watcher) txResultRoutine() {
	w.txResultChan = make(chan TxResult, DefaultTxResultChanBuffer)

	for result := range w.txResultChan {
		if result.TxMsg != nil {
			w.txs = append(w.txs, result.TxMsg)
		}
		if result.TxReceipt != nil {
			w.txReceipts = append(w.txReceipts, result.TxReceipt)
		}
		w.txsInBlock = append(w.txsInBlock, TxInfo{TxHash: result.TxHash, Index: result.Index, GasUsed: result.GasUsed})
	}
}
