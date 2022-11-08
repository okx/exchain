package dydx

import (
	"container/list"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type OrderQueueIterator interface {
	Next() *WrapOrder
}

type orderQueueIterator struct {
	q    *OrderQueue
	cur  *list.Element
	done bool
}

func (it *orderQueueIterator) Next() *WrapOrder {
	it.q.mtx.RLock()
	defer it.q.mtx.RUnlock()

	if it.done {
		return nil
	}

	if it.cur == nil {
		it.cur = it.q.list.Front()
	} else {
		it.cur = it.cur.Next()
	}
	if it.cur == nil {
		it.done = true
		return nil
	}
	return it.cur.Value.(*WrapOrder)
}

type OrderQueue struct {
	list *list.List
	mtx  sync.RWMutex
	m    map[common.Hash]*list.Element

	book *DepthBook
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		list: list.New(),
		m:    make(map[common.Hash]*list.Element),
		book: NewDepthBook(),
	}
}

func (q *OrderQueue) Book() *DepthBook {
	return q.book
}

func (q *OrderQueue) RLock() {
	q.mtx.RLock()
}

func (q *OrderQueue) RUnlock() {
	q.mtx.RUnlock()
}

func (q *OrderQueue) NewIterator() OrderQueueIterator {
	return &orderQueueIterator{q: q}
}

func (q *OrderQueue) Enqueue(v *WrapOrder) bool {
	if v == nil {
		return false
	}
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if o := q.m[v.Hash()]; o != nil {
		return false
	}
	q.m[v.Hash()] = q.list.PushBack(v)

	_ = q.book.Insert(v)
	return true
}

func (q *OrderQueue) Dequeue() *WrapOrder {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	e := q.list.Front()
	if e != nil {
		o := q.list.Remove(e).(*WrapOrder)
		delete(q.m, o.Hash())
		q.book.DeleteByHash(o.Hash())
		return o
	}
	return nil
}

func (q *OrderQueue) Get(hash common.Hash) *WrapOrder {
	q.mtx.RLock()
	defer q.mtx.RUnlock()

	e, ok := q.m[hash]
	if !ok {
		return nil
	}
	return e.Value.(*WrapOrder)
}

func (q *OrderQueue) GetAllOrderHashes() [][32]byte {
	q.mtx.RLock()
	defer q.mtx.RUnlock()

	var orderHashes [][32]byte
	if len(q.m) > 0 {
		orderHashes = make([][32]byte, 0, len(q.m))
		for k := range q.m {
			orderHashes = append(orderHashes, k)
		}
	}
	return orderHashes
}

func (q *OrderQueue) Delete(hash common.Hash) *WrapOrder {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	e, ok := q.m[hash]
	if !ok {
		return nil
	}
	delete(q.m, hash)
	q.book.DeleteByHash(hash)
	return q.list.Remove(e).(*WrapOrder)
}

func (q *OrderQueue) Len() int {
	q.mtx.RLock()
	defer q.mtx.RUnlock()

	return q.list.Len()
}
