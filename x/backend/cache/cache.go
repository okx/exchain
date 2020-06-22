package cache

import "github.com/okex/okchain/x/backend/types"

// Cache defines struct to store data in memory
type Cache struct {
	// Flush at EndBlock
	Transactions []*types.Transaction

	// persist in memory
	LatestTicker map[string]*types.Ticker
	ProductsBuf  []string
}

// NewCache return  cache pointer address, called at NewKeeper
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

// AddTransaction append transaction to cache Transactions
func (c *Cache) AddTransaction(transaction *types.Transaction) {
	if transaction == nil {
		panic("failed. a nil pointer appears")
	}
	c.Transactions = append(c.Transactions, transaction)
}

// nolint
func (c *Cache) GetTransactions() []*types.Transaction {
	return c.Transactions
}
