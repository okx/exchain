package infura

import evm "github.com/okex/exchain/x/evm/watcher"

const defaultCacheCap = 2000

type Cache struct {
	transactionReceipts []evm.TransactionReceipt
}

func NewCache() *Cache {
	return &Cache{
		transactionReceipts: make([]evm.TransactionReceipt, 0, defaultCacheCap),
	}
}

func (c *Cache) Reset() {
	c.transactionReceipts = make([]evm.TransactionReceipt, 0, defaultCacheCap)
}

func (c *Cache) AddTransactionReceipt(tr evm.TransactionReceipt) {
	c.transactionReceipts = append(c.transactionReceipts, tr)
}

func (c *Cache) GetTransactionReceipts() []evm.TransactionReceipt {
	return c.transactionReceipts
}
