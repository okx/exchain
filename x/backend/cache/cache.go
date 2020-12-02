package cache

import "github.com/okex/okexchain/x/backend/types"

// Cache defines struct to store data in memory
type Cache struct {
	// Flush at EndBlock
	Transactions []*types.Transaction

	// persist in memory
	LatestTicker map[string]*types.Ticker

	// swap infos, flush at EndBlocker
	swapInfos []*types.SwapInfo
	numTxs    int
}

// NewCache return  cache pointer address, called at NewKeeper
func NewCache() *Cache {
	return &Cache{
		Transactions: make([]*types.Transaction, 0, 2000),
		LatestTicker: make(map[string]*types.Ticker),
		swapInfos:    make([]*types.SwapInfo, 0, 2000),
		numTxs:       0,
	}
}

// Flush temporary cache, called at EndBlock
func (c *Cache) Flush() {
	c.Transactions = make([]*types.Transaction, 0, 2000)
	c.swapInfos = make([]*types.SwapInfo, 0, 2000)
	c.numTxs = 0
}

// AddTransaction append transaction to cache Transactions
func (c *Cache) AddTransaction(transaction []*types.Transaction) {
	c.Transactions = append(c.Transactions, transaction...)
}

// nolint
func (c *Cache) GetTransactions() []*types.Transaction {
	return c.Transactions
}

// AddSwapInfo appends swapInfo to cache SwapInfos
func (c *Cache) AddSwapInfo(swapInfo *types.SwapInfo) {
	c.swapInfos = append(c.swapInfos, swapInfo)
}

// nolint
func (c *Cache) GetSwapInfos() []*types.SwapInfo {
	return c.swapInfos
}

func (c *Cache) AddNumTxs(num int) {
	c.numTxs += num
}

func (c *Cache) GetNumTxs() int {
	return c.numTxs
}
