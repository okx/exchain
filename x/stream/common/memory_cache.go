package common

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/auth"
	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/backend"
	"github.com/okex/exchain/x/dex/types"
)

type Cache struct {
	// Flush at EndBlock
	transactions      []*backend.Transaction
	newTokenPairMap   []*types.TokenPair
	tokenPairChanged  bool
	updatedAccAddress map[string]struct{}
	swapInfos         []*backend.SwapInfo
	newSwapTokenPairs []*ammswap.SwapTokenPair
	claimInfos        []*backend.ClaimInfo
}

func NewCache() *Cache {
	return &Cache{
		transactions:      make([]*backend.Transaction, 0, 2000),
		newTokenPairMap:   []*types.TokenPair{},
		tokenPairChanged:  false,
		updatedAccAddress: make(map[string]struct{}),
		swapInfos:         make([]*backend.SwapInfo, 0, 2000),
		newSwapTokenPairs: make([]*ammswap.SwapTokenPair, 0, 2000),
		claimInfos:        make([]*backend.ClaimInfo, 0, 2000),
	}
}

// Reset temporary cache, called at BeginBlock
func (c *Cache) Reset() {
	c.transactions = make([]*backend.Transaction, 0, 2000)
	c.newTokenPairMap = []*types.TokenPair{}
	c.tokenPairChanged = false
	c.updatedAccAddress = make(map[string]struct{})
	c.swapInfos = make([]*backend.SwapInfo, 0, 2000)
	c.newSwapTokenPairs = make([]*ammswap.SwapTokenPair, 0, 2000)
	c.claimInfos = make([]*backend.ClaimInfo, 0, 2000)
}

func (c *Cache) AddTransaction(transaction *backend.Transaction) {
	c.transactions = append(c.transactions, transaction)
}

func (c *Cache) GetTransactions() []*backend.Transaction {
	return c.transactions
}

// AddNewTokenPair adds a new token pair into cache
func (c *Cache) AddNewTokenPair(tokenPair *types.TokenPair) {
	c.newTokenPairMap = append(c.newTokenPairMap, tokenPair)
}

// GetNewTokenPairs returns new token pairs from cache
func (c *Cache) GetNewTokenPairs() []*types.TokenPair {
	return c.newTokenPairMap
}

// SetTokenPairChanged sets tokenPairChanged
func (c *Cache) SetTokenPairChanged(changed bool) {
	c.tokenPairChanged = changed
}

// GetTokenPairChanged gets tokenPairChanged
func (c *Cache) GetTokenPairChanged() bool {
	return c.tokenPairChanged
}

func (c *Cache) AddUpdatedAccount(acc auth.Account) {
	c.updatedAccAddress[acc.GetAddress().String()] = struct{}{}
}

func (c *Cache) GetUpdatedAccAddress() (accs []sdk.AccAddress) {
	for acc := range c.updatedAccAddress {
		addr, err := sdk.AccAddressFromBech32(acc)
		if err == nil {
			accs = append(accs, addr)
		}
	}
	return accs
}

// AddSwapInfo appends swapInfo to cache SwapInfos
func (c *Cache) AddSwapInfo(swapInfo *backend.SwapInfo) {
	c.swapInfos = append(c.swapInfos, swapInfo)
}

// nolint
func (c *Cache) GetSwapInfos() []*backend.SwapInfo {
	return c.swapInfos
}

// AddNewSwapTokenPairs appends swapTokenPair to cache newSwapTokenPairs
func (c *Cache) AddNewSwapTokenPair(swapTokenPair *ammswap.SwapTokenPair) {
	c.newSwapTokenPairs = append(c.newSwapTokenPairs, swapTokenPair)
}

// nolint
func (c *Cache) GetNewSwapTokenPairs() []*ammswap.SwapTokenPair {
	return c.newSwapTokenPairs
}

// AddClaimInfo appends claimInfo to cache ClaimInfos
func (c *Cache) AddClaimInfo(claimInfo *backend.ClaimInfo) {
	c.claimInfos = append(c.claimInfos, claimInfo)
}

// nolint
func (c *Cache) GetClaimInfos() []*backend.ClaimInfo {
	return c.claimInfos
}
