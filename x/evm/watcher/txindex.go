package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"sort"
)

const (
	DefaultTxResultChanBuffer = 20
)

type TxIndex struct {
	TxHash ethcmn.Hash
	Index  uint64
}

func (w *Watcher) addTxsToBlock() {
	sort.Slice(w.txsInBlock, func(i int, j int) bool {
		return w.txsInBlock[i].Index < w.txsInBlock[j].Index
	})
	for _, txIndex := range w.txsInBlock {
		w.blockTxs = append(w.blockTxs, txIndex.TxHash)
	}
}

type TxResult struct {
	TxMsg     WatchMessage
	TxReceipt WatchMessage
	Index     uint64
	GasUsed   uint64
	TxHash    ethcmn.Hash
}

func (w *Watcher) txResultRoutine() {
	w.txResultChan = make(chan TxResult, DefaultTxResultChanBuffer)

	for result := range w.txResultChan {
		if result.TxMsg != nil {
			w.txsAndReceipts = append(w.txsAndReceipts, result.TxMsg)
		}
		if result.TxReceipt != nil {
			w.txsAndReceipts = append(w.txsAndReceipts, result.TxReceipt)
		}
		//w.UpdateCumulativeGas(result.Index, result.GasUsed)
		w.txsInBlock = append(w.txsInBlock, TxIndex{TxHash: result.TxHash, Index: result.Index})
	}
}
