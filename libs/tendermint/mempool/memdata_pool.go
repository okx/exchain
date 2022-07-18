package mempool

import (
	"container/list"
	"sync"
)

type mempoolTxList struct {
	mtx    sync.Mutex
	stores *list.List
}

func newMempoolTxList() *mempoolTxList {
	return &mempoolTxList{
		stores: list.New(),
	}
}

func (c *mempoolTxList) PushTx(tx *mempoolTx) {
	c.mtx.Lock()
	c.stores.PushBack(tx)
	c.mtx.Unlock()
}

func (c *mempoolTxList) GetTx() *mempoolTx {
	c.mtx.Lock()
	if c.stores.Len() > 0 {
		front := c.stores.Remove(c.stores.Front())
		c.mtx.Unlock()
		return front.(*mempoolTx)
	}
	c.mtx.Unlock()
	return &mempoolTx{}
}

func (c *mempoolTxList) Clear() {
	c.mtx.Lock()
	c.stores.Init()
	c.mtx.Unlock()
}
