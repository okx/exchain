package cache

import "github.com/okex/okchain/x/backend/types"

type Cache struct {
	// Flush at EndBlock
	Transactions []*types.Transaction

	// persist in memory
	LatestTicker map[string]*types.Ticker
	ProductsBuf  []string
}

func NewCache() *Cache {
	return &Cache{
		Transactions: make([]*types.Transaction, 0, 2000),
		LatestTicker: make(map[string]*types.Ticker),
		ProductsBuf:  make([]string, 0, 200),
	}
}

// Flush temporary cache, called at EndBlock
func (c *Cache) Flush() {
	c.Transactions = make([]*types.Transaction, 0, 2000)
}

func (c *Cache) AddTransaction(transaction *types.Transaction) {
	c.Transactions = append(c.Transactions, transaction)
}

func (c *Cache) GetTransactions() []*types.Transaction {
	return c.Transactions
}
