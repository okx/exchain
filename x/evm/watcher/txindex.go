package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"sort"
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
