package mempool

import "sync"

type checkingCache struct {
	txHashMap map[[32]byte]struct{}
	lock      sync.RWMutex
}

func (c *checkingCache) TryAdd(txHash [32]byte) bool {
	c.lock.RLock()
	_, ok := c.txHashMap[txHash]
	c.lock.RUnlock()
	if ok {
		return false
	}

	c.lock.Lock()
	_, ok = c.txHashMap[txHash]
	if !ok {
		c.txHashMap[txHash] = struct{}{}
		c.lock.Unlock()
		return true
	} else {
		c.lock.Unlock()
		return false
	}
}

func (c *checkingCache) Delete(txHash [32]byte) {
	c.lock.Lock()
	delete(c.txHashMap, txHash)
	c.lock.Unlock()
}

func newCheckingCache() *checkingCache {
	return &checkingCache{
		txHashMap: make(map[[32]byte]struct{}),
	}
}
