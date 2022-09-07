package mempool

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"

	"github.com/okex/exchain/libs/tendermint/libs/clist"
)

type GasTxQueue struct {
	txPriceBump  int64
	sortedTxs    *clist.CList
	sortedTxsMap sync.Map
	bcTxs        *clist.CList
	bcTxsMap     sync.Map

	*AddressRecord
}

func NewGasTxQueue(txPriceBump int64) *GasTxQueue {
	q := &GasTxQueue{
		txPriceBump: txPriceBump,
		sortedTxs:   clist.New(),
		bcTxs:       clist.New(),
	}
	q.AddressRecord = newAddressRecord(q)
	return q
}

func (q *GasTxQueue) Len() int {
	return q.sortedTxs.Len()
}

func (q *GasTxQueue) Insert(memTx *mempoolTx) error {
	/*
		1. insert tx list
		2. insert address record
		3. insert tx map
	*/
	ele := q.AddressRecord.checkRepeatedAndAddItem(memTx, q.txPriceBump, q.sortedTxs.InsertElement)
	if ele == nil {
		return fmt.Errorf("failed to replace tx for acccount %s with nonce %d, "+
			"the provided gas price %d is not bigger enough", memTx.from, memTx.realTx.GetNonce(), memTx.realTx.GetGasPrice())
	}
	txHash := txOrTxHashToKey(memTx.tx, memTx.realTx.TxHash(), memTx.height)

	q.sortedTxsMap.Store(txHash, ele)

	ele2 := q.bcTxs.PushBack(memTx)
	ele2.Address = memTx.from
	q.bcTxsMap.Store(txHash, ele2)
	return nil
}

func (q *GasTxQueue) Remove(element *clist.CElement) {
	q.removeElement(element)
	q.AddressRecord.DeleteItem(element)
}

func (q *GasTxQueue) RemoveByKey(key [32]byte) (ele *clist.CElement) {
	ele = q.removeElementByKey(key)
	if ele != nil {
		q.AddressRecord.DeleteItem(ele)
	}
	return
}

func (q *GasTxQueue) Front() *clist.CElement {
	return q.sortedTxs.Front()
}

func (q *GasTxQueue) Back() *clist.CElement {
	return q.sortedTxs.Back()
}

func (q *GasTxQueue) BroadcastFront() *clist.CElement {
	return q.bcTxs.Front()
}

func (q *GasTxQueue) BroadcastLen() int {
	return q.bcTxs.Len()
}

func (q *GasTxQueue) TxsWaitChan() <-chan struct{} {
	return q.sortedTxs.WaitChan()
}

func (q *GasTxQueue) Load(hash [sha256.Size]byte) (*clist.CElement, bool) {
	v, ok := q.sortedTxsMap.Load(hash)
	if !ok {
		return nil, false
	}
	return v.(*clist.CElement), true
}

func (q *GasTxQueue) CleanItems(address string, nonce uint64) {
	q.AddressRecord.CleanItems(address, nonce, q.removeElement)
}

func (q *GasTxQueue) reorganizeElements(items []*clist.CElement) {
	if len(items) == 0 {
		return
	}
	// When inserting, strictly order by nonce, otherwise tx will not appear according to nonce,
	// resulting in execution failure
	sort.Slice(items, func(i, j int) bool { return items[i].Nonce < items[j].Nonce })

	for _, item := range items[1:] {
		q.sortedTxs.DetachElement(item)
		item.NewDetachPrev()
		item.NewDetachNext()
	}

	for _, item := range items {
		q.sortedTxs.InsertElement(item)
	}
}

func (q *GasTxQueue) removeElement(element *clist.CElement) {
	q.sortedTxs.Remove(element)
	element.DetachPrev()

	tx := element.Value.(*mempoolTx).tx
	txHash := txKey(tx)
	q.sortedTxsMap.Delete(txHash)

	if v, ok := q.bcTxsMap.LoadAndDelete(txHash); ok {
		ele := v.(*clist.CElement)
		q.bcTxs.Remove(ele)
		ele.DetachPrev()
	}
}

func (q *GasTxQueue) removeElementByKey(key [32]byte) (ret *clist.CElement) {
	if v, ok := q.sortedTxsMap.LoadAndDelete(key); ok {
		ret = v.(*clist.CElement)
		q.sortedTxs.Remove(ret)
		ret.DetachPrev()

		if v, ok := q.bcTxsMap.LoadAndDelete(key); ok {
			ele := v.(*clist.CElement)
			q.bcTxs.Remove(ele)
			ele.DetachPrev()
		}
	}

	return
}
