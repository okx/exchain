package cache

import "github.com/okex/okexchain/x/backend/types"

// Cache defines struct to store data in memory
type Cache struct {
	// Flush at EndBlock
	Transactions []*types.Transaction

	// persist in memory
	LatestTicker map[string]*types.Ticker

	// swap infos, flush at EndBlocker
	swapInfos  []*types.SwapInfo
	claimInfos []*types.ClaimInfo
}

// NewCache return  cache pointer address, called at NewKeeper
func NewCache() *Cache {
	return &Cache{
		Transactions: make([]*types.Transaction, 0, 2000),
		LatestTicker: make(map[string]*types.Ticker),
		swapInfos:    make([]*types.SwapInfo, 0, 2000),
		claimInfos:   make([]*types.ClaimInfo, 0, 2000),
	}
}

// Flush temporary cache, called at EndBlock
func (c *Cache) Flush() {
	c.Transactions = make([]*types.Transaction, 0, 2000)
	c.swapInfos = make([]*types.SwapInfo, 0, 2000)
	c.claimInfos = make([]*types.ClaimInfo, 0, 2000)
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

// AddClaimInfo appends claimInfo to cache ClaimInfos
func (c *Cache) AddClaimInfo(claimInfo *types.ClaimInfo) {
	c.claimInfos = append(c.claimInfos, claimInfo)
}

// nolint
func (c *Cache) GetClaimInfos() []*types.ClaimInfo {
	return c.claimInfos
}
