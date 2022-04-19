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
	sort.Slice(w.txsCollector, func(i, j int) bool {
		return w.txsCollector[i].Index < w.txsCollector[j].Index
	})
	for _, txInfo := range w.txsCollector {
		w.blockTxs = append(w.blockTxs, txInfo.TxHash)
		w.updateCumulativeGas(txInfo.Index, txInfo.GasUsed)
	}

	sort.Slice(w.txReceipts, func(i, j int) bool {
		return w.txReceipts[i].TransactionIndex < w.txReceipts[j].TransactionIndex
	})
	for _, receipt := range w.txReceipts {
		receipt.CumulativeGasUsed = hexutil.Uint64(w.cumulativeGas[uint64(receipt.TransactionIndex)])
		w.batch = append(w.batch, &MsgTransactionReceipt{txHash: receipt.TxHash.Bytes(), baseLazyMarshal: newBaseLazyMarshal(receipt)})
	}
}

type TxResult struct {
	TxMsg     WatchMessage
	TxReceipt *TransactionReceipt
	TxHash    ethcmn.Hash
	Index     uint64
	GasUsed   uint64
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
		w.txsCollector = append(w.txsCollector, TxInfo{})
	}
}
