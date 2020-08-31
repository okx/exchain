package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/dex/types"
)

type Cache struct {
	// Flush at EndBlock
	transactions      []*backend.Transaction
	newTokenPairMap   []*types.TokenPair
	tokenPairChanged  bool
	updatedAccAddress map[string]struct{}
}

func NewCache() *Cache {
	return &Cache{
		transactions:      make([]*backend.Transaction, 0, 2000),
		newTokenPairMap:   []*types.TokenPair{},
		tokenPairChanged:  false,
		updatedAccAddress: make(map[string]struct{}),
	}
}

// Reset temporary cache, called at BeginBlock
func (c *Cache) Reset() {
	c.transactions = make([]*backend.Transaction, 0, 2000)
	c.newTokenPairMap = []*types.TokenPair{}
	c.tokenPairChanged = false
	c.updatedAccAddress = make(map[string]struct{})
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
