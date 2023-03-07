package mempool

import (
	"crypto/sha256"
	"sync"

	"github.com/okx/okbchain/libs/tendermint/libs/clist"
	"github.com/okx/okbchain/libs/tendermint/types"
)

type ITransactionQueue interface {
	Len() int
	Insert(tx *mempoolTx) error
	Remove(element *clist.CElement)
	RemoveByKey(key [sha256.Size]byte) *clist.CElement
	Front() *clist.CElement
	Back() *clist.CElement
	BroadcastFront() *clist.CElement
	BroadcastLen() int
	Load(hash [sha256.Size]byte) (*clist.CElement, bool)
	TxsWaitChan() <-chan struct{}

	AddressRecorder
}

type AddressRecorder interface {
	GetAddressList() []string
	GetAddressNonce(address string) (uint64, bool)
	GetAddressTxsCnt(address string) int
	GetAddressTxs(address string, max int) types.Txs
	CleanItems(address string, nonce uint64)
}

type BaseTxQueue struct {
	txs    *clist.CList // FIFO list
	txsMap sync.Map     //txKey -> CElement

	*AddressRecord
}

func NewBaseTxQueue() *BaseTxQueue {
	return &BaseTxQueue{
		txs:           clist.New(),
		AddressRecord: newAddressRecord(nil),
	}
}

func (q *BaseTxQueue) Len() int {
	return q.txs.Len()
}

func (q *BaseTxQueue) Insert(tx *mempoolTx) error {
	/*
		1. insert tx list
		2. insert address record
		3. insert tx map
	*/
	ele := q.txs.PushBack(tx)

	q.AddressRecord.AddItem(ele.Address, ele)
	q.txsMap.Store(txKey(ele.Value.(*mempoolTx).tx), ele)
	return nil
}

func (q *BaseTxQueue) Remove(element *clist.CElement) {
	q.removeElement(element)
	q.AddressRecord.DeleteItem(element)
}

func (q *BaseTxQueue) RemoveByKey(key [32]byte) (ele *clist.CElement) {
	ele = q.removeElementByKey(key)
	if ele != nil {
		q.AddressRecord.DeleteItem(ele)
	}
	return
}

func (q *BaseTxQueue) Front() *clist.CElement {
	return q.txs.Front()
}

func (q *BaseTxQueue) Back() *clist.CElement {
	return q.txs.Back()
}

func (q *BaseTxQueue) BroadcastFront() *clist.CElement {
	return q.txs.Front()
}

func (q *BaseTxQueue) BroadcastLen() int {
	return q.txs.Len()
}

func (q *BaseTxQueue) TxsWaitChan() <-chan struct{} {
	return q.txs.WaitChan()
}

func (q *BaseTxQueue) Load(hash [sha256.Size]byte) (*clist.CElement, bool) {
	v, ok := q.txsMap.Load(hash)
	if !ok {
		return nil, false
	}
	return v.(*clist.CElement), true
}

func (q *BaseTxQueue) removeElement(element *clist.CElement) {
	q.txs.Remove(element)
	element.DetachPrev()

	tx := element.Value.(*mempoolTx).tx
	txHash := txKey(tx)
	q.txsMap.Delete(txHash)
}

func (q *BaseTxQueue) removeElementByKey(key [32]byte) *clist.CElement {
	if v, ok := q.txsMap.LoadAndDelete(key); ok {
		element := v.(*clist.CElement)
		q.txs.Remove(element)
		element.DetachPrev()
		return element
	}
	return nil
}

func (q *BaseTxQueue) CleanItems(address string, nonce uint64) {
	q.AddressRecord.CleanItems(address, nonce, q.removeElement)
}
