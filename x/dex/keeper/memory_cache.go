package keeper

import (
	"github.com/okex/okchain/x/dex/types"
	ordertypes "github.com/okex/okchain/x/order/types"
)

type Cache struct {
	tokenPairMap     map[string]*types.TokenPair
	newTokenPairMap  []*types.TokenPair
	tokenPairChanged bool
	lockMap          *ordertypes.ProductLockMap
}

func NewCache() *Cache {
	return &Cache{
		tokenPairMap:     make(map[string]*types.TokenPair),
		tokenPairChanged: false,
		lockMap:          ordertypes.NewProductLockMap(),
	}
}

func (c *Cache) Reset() {
	c.newTokenPairMap = []*types.TokenPair{}
	c.tokenPairChanged = false
}

//add a new token pair into cache
func (c *Cache) AddTokenPair(tokenPair *types.TokenPair) {
	c.tokenPairMap[c.genTokenPairKey(tokenPair)] = tokenPair
	c.tokenPairChanged = true
}

//generate token pair key
func (c *Cache) genTokenPairKey(t *types.TokenPair) string {
	return t.BaseAssetSymbol + "_" + t.QuoteAssetSymbol
}

func (c *Cache) GetTokenPair(product string) (*types.TokenPair, bool) {
	tokenPair, ok := c.tokenPairMap[product]
	return tokenPair, ok
}

//add full token pairs
func (c *Cache) PrepareTokenPairs(tokenPairs []*types.TokenPair) {
	for _, v := range tokenPairs {
		c.AddTokenPair(v)
	}
}

//get all token pairs from cache
func (c *Cache) GetAllTokenPairs() []*types.TokenPair {
	tokenPairs := make([]*types.TokenPair, 0, len(c.tokenPairMap))
	for _, v := range c.tokenPairMap {
		tokenPairs = append(tokenPairs, v)
	}
	return tokenPairs
}

//delete token pair cache
func (c *Cache) DeleteTokenPair(targetPair *types.TokenPair) {
	delete(c.tokenPairMap, c.genTokenPairKey(targetPair))
	c.tokenPairChanged = true
}

// delete token pair by token pair's name from cache
func (c *Cache) DeleteTokenPairByName(tokenPairName string) {
	delete(c.tokenPairMap, tokenPairName)
	c.tokenPairChanged = true
}

func (c *Cache) TokenPairCount() int {
	return len(c.tokenPairMap)
}

//add a new token pair into cache
func (c *Cache) AddNewTokenPair(tokenPair *types.TokenPair) {
	c.newTokenPairMap = append(c.newTokenPairMap, tokenPair)
	c.tokenPairChanged = true
}

//get new token pairs from cache
func (c *Cache) GetNewTokenPair() []*types.TokenPair {
	return c.newTokenPairMap
}
