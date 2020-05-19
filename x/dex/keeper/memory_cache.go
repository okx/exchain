package keeper

import (
	"github.com/okex/okchain/x/dex/types"
	ordertypes "github.com/okex/okchain/x/order/types"
)

// Cache caches data
type Cache struct {
	tokenPairMap     map[string]*types.TokenPair
	newTokenPairMap  []*types.TokenPair
	tokenPairChanged bool
	lockMap          *ordertypes.ProductLockMap
}

// NewCache returns instance of Cache
func NewCache() *Cache {
	return &Cache{
		tokenPairMap:     make(map[string]*types.TokenPair),
		tokenPairChanged: false,
		lockMap:          ordertypes.NewProductLockMap(),
	}
}

// Reset clears cache
func (c *Cache) Reset() {
	c.newTokenPairMap = []*types.TokenPair{}
	c.tokenPairChanged = false
}

// AddTokenPair adds a new token pair into cache
func (c *Cache) AddTokenPair(tokenPair *types.TokenPair) {
	c.tokenPairMap[c.genTokenPairKey(tokenPair)] = tokenPair
	c.tokenPairChanged = true
}

//generate token pair key
func (c *Cache) genTokenPairKey(t *types.TokenPair) string {
	return t.BaseAssetSymbol + "_" + t.QuoteAssetSymbol
}

// GetTokenPair returns token pair from cache
func (c *Cache) GetTokenPair(product string) (*types.TokenPair, bool) {
	tokenPair, ok := c.tokenPairMap[product]
	return tokenPair, ok
}

// PrepareTokenPairs adds all token pairs
func (c *Cache) PrepareTokenPairs(tokenPairs []*types.TokenPair) {
	for _, v := range tokenPairs {
		c.AddTokenPair(v)
	}
}

// GetAllTokenPairs returns all token pairs from cache
func (c *Cache) GetAllTokenPairs() []*types.TokenPair {
	tokenPairs := make([]*types.TokenPair, 0, len(c.tokenPairMap))
	for _, v := range c.tokenPairMap {
		tokenPairs = append(tokenPairs, v)
	}
	return tokenPairs
}

// DeleteTokenPair deletes token pair cache
func (c *Cache) DeleteTokenPair(targetPair *types.TokenPair) {
	delete(c.tokenPairMap, c.genTokenPairKey(targetPair))
	c.tokenPairChanged = true
}

// DeleteTokenPairByName deletes token pair by token pair's name from cache
func (c *Cache) DeleteTokenPairByName(tokenPairName string) {
	delete(c.tokenPairMap, tokenPairName)
	c.tokenPairChanged = true
}

// TokenPairCount returns count of token pair
func (c *Cache) TokenPairCount() int {
	return len(c.tokenPairMap)
}

// AddNewTokenPair adds a new token pair into cache
func (c *Cache) AddNewTokenPair(tokenPair *types.TokenPair) {
	c.newTokenPairMap = append(c.newTokenPairMap, tokenPair)
	c.tokenPairChanged = true
}

// GetNewTokenPair returns new token pairs from cache
func (c *Cache) GetNewTokenPair() []*types.TokenPair {
	return c.newTokenPairMap
}

// nolint
func (c *Cache) DepthCopy() *Cache {
	cache := Cache{
		tokenPairMap:     nil,
		newTokenPairMap:  nil,
		tokenPairChanged: c.tokenPairChanged,
		lockMap:          nil,
	}

	if c.tokenPairMap != nil {
		cpTokenPairMap := make(map[string]*types.TokenPair)
		for k, v := range c.tokenPairMap{
			cpTokenPairMap[k] = &types.TokenPair{
				BaseAssetSymbol:  v.BaseAssetSymbol,
				QuoteAssetSymbol: v.QuoteAssetSymbol,
				InitPrice:        v.InitPrice,
				MaxPriceDigit:    v.MaxPriceDigit,
				MaxQuantityDigit: v.MaxQuantityDigit,
				MinQuantity:      v.MinQuantity,
				ID:               v.ID,
				Delisting:        v.Delisting,
				Owner:            v.Owner,
				Deposits:         v.Deposits,
				BlockHeight:      v.BlockHeight,
			}
		}

		cache.tokenPairMap = cpTokenPairMap
	}

	if c.newTokenPairMap != nil {
		cpNewTokenPairMap := make([]*types.TokenPair, len(c.newTokenPairMap))
		for _, v := range c.newTokenPairMap{
			cpNewTokenPairMap = append(cpNewTokenPairMap, &types.TokenPair{
				BaseAssetSymbol:  v.BaseAssetSymbol,
				QuoteAssetSymbol: v.QuoteAssetSymbol,
				InitPrice:        v.InitPrice,
				MaxPriceDigit:    v.MaxPriceDigit,
				MaxQuantityDigit: v.MaxQuantityDigit,
				MinQuantity:      v.MinQuantity,
				ID:               v.ID,
				Delisting:        v.Delisting,
				Owner:            v.Owner,
				Deposits:         v.Deposits,
				BlockHeight:      v.BlockHeight,
			})
		}

		cache.newTokenPairMap = cpNewTokenPairMap
	}

	if c.lockMap != nil {
		cpData := make(map[string]*ordertypes.ProductLock)
		for k, v := range c.lockMap.Data{
			cpData[k] = &ordertypes.ProductLock{
				BlockHeight:  v.BlockHeight,
				Price:        v.Price,
				Quantity:     v.Quantity,
				BuyExecuted:  v.BuyExecuted,
				SellExecuted: v.SellExecuted,
			}
		}

		cache.lockMap = &ordertypes.ProductLockMap{
			Data: cpData,
		}
	}

	return &cache
}
