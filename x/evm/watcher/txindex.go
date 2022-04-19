package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"sort"
)

type TxInfo struct {
	TxHash  ethcmn.Hash
	Index   uint64
	GasUsed uint64
}

func (w *Watcher) addTxsToBlock() {
	sort.Slice(w.txInfoCollector, func(i, j int) bool {
		return w.txInfoCollector[i].Index < w.txInfoCollector[j].Index
	})
	for _, txInfo := range w.txInfoCollector {
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
