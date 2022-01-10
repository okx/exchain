package common

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/dex/types"
)

type Cache struct {
	// Flush at EndBlock
	newTokenPairMap   []*types.TokenPair
	tokenPairChanged  bool
	updatedAccAddress map[string]struct{}
	newSwapTokenPairs []*ammswap.SwapTokenPair
}

func NewCache() *Cache {
	return &Cache{
		newTokenPairMap:   []*types.TokenPair{},
		tokenPairChanged:  false,
		updatedAccAddress: make(map[string]struct{}),
		newSwapTokenPairs: make([]*ammswap.SwapTokenPair, 0, 2000),
	}
}

// Reset temporary cache, called at BeginBlock
func (c *Cache) Reset() {
	c.newTokenPairMap = []*types.TokenPair{}
	c.tokenPairChanged = false
	c.updatedAccAddress = make(map[string]struct{})
	c.newSwapTokenPairs = make([]*ammswap.SwapTokenPair, 0, 2000)
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

// AddNewSwapTokenPairs appends swapTokenPair to cache newSwapTokenPairs
func (c *Cache) AddNewSwapTokenPair(swapTokenPair *ammswap.SwapTokenPair) {
	c.newSwapTokenPairs = append(c.newSwapTokenPairs, swapTokenPair)
}

// nolint
func (c *Cache) GetNewSwapTokenPairs() []*ammswap.SwapTokenPair {
	return c.newSwapTokenPairs
}
