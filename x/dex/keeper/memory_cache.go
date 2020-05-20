package keeper

import (
	"encoding/json"
	"github.com/okex/okchain/x/dex/types"
	ordertypes "github.com/okex/okchain/x/order/types"
)

// Cache caches data
type Cache struct {
	TokenPairMap     map[string]*types.TokenPair `json:"token_pair_map"`
	NewTokenPairMap  []*types.TokenPair `json:"new_token_pair_map"`
	TokenPairChanged bool `json:"token_pair_changed"`
	LockMap          *ordertypes.ProductLockMap `json:"lock_map"`
}

// NewCache returns instance of Cache
func NewCache() *Cache {
	return &Cache{
		TokenPairMap:     make(map[string]*types.TokenPair),
		TokenPairChanged: false,
		LockMap:          ordertypes.NewProductLockMap(),
	}
}

// Reset clears cache
func (c *Cache) Reset() {
	c.NewTokenPairMap = []*types.TokenPair{}
	c.TokenPairChanged = false
}

// AddTokenPair adds a new token pair into cache
func (c *Cache) AddTokenPair(tokenPair *types.TokenPair) {
	c.TokenPairMap[c.genTokenPairKey(tokenPair)] = tokenPair
	c.TokenPairChanged = true
}

//generate token pair key
func (c *Cache) genTokenPairKey(t *types.TokenPair) string {
	return t.BaseAssetSymbol + "_" + t.QuoteAssetSymbol
}

// GetTokenPair returns token pair from cache
func (c *Cache) GetTokenPair(product string) (*types.TokenPair, bool) {
	tokenPair, ok := c.TokenPairMap[product]
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
	tokenPairs := make([]*types.TokenPair, 0, len(c.TokenPairMap))
	for _, v := range c.TokenPairMap {
		tokenPairs = append(tokenPairs, v)
	}
	return tokenPairs
}

// DeleteTokenPair deletes token pair cache
func (c *Cache) DeleteTokenPair(targetPair *types.TokenPair) {
	delete(c.TokenPairMap, c.genTokenPairKey(targetPair))
	c.TokenPairChanged = true
}

// DeleteTokenPairByName deletes token pair by token pair's name from cache
func (c *Cache) DeleteTokenPairByName(tokenPairName string) {
	delete(c.TokenPairMap, tokenPairName)
	c.TokenPairChanged = true
}

// TokenPairCount returns count of token pair
func (c *Cache) TokenPairCount() int {
	return len(c.TokenPairMap)
}

// AddNewTokenPair adds a new token pair into cache
func (c *Cache) AddNewTokenPair(tokenPair *types.TokenPair) {
	c.NewTokenPairMap = append(c.NewTokenPairMap, tokenPair)
	c.TokenPairChanged = true
}

// GetNewTokenPair returns new token pairs from cache
func (c *Cache) GetNewTokenPair() []*types.TokenPair {
	return c.NewTokenPairMap
}

func (c *Cache) Clone() *Cache {
	cache := &Cache{}
	bytes, _ := json.Marshal(c)
	err := json.Unmarshal(bytes, cache)
	if err != nil {
		return c.DepthCopy()
	}

	return cache
}

// nolint
func (c *Cache) DepthCopy() *Cache {
	cache := Cache{
		TokenPairMap:     nil,
		NewTokenPairMap:  nil,
		TokenPairChanged: c.TokenPairChanged,
		LockMap:          nil,
	}

	if c.TokenPairMap != nil {
		cpTokenPairMap := make(map[string]*types.TokenPair)
		for k, v := range c.TokenPairMap{
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

		cache.TokenPairMap = cpTokenPairMap
	}

	if c.NewTokenPairMap != nil {
		cpNewTokenPairMap := make([]*types.TokenPair, 0, len(c.NewTokenPairMap))
		for _, v := range c.NewTokenPairMap{
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

		cache.NewTokenPairMap = cpNewTokenPairMap
	}

	if c.LockMap != nil {
		cpData := make(map[string]*ordertypes.ProductLock)
		for k, v := range c.LockMap.Data{
			cpData[k] = &ordertypes.ProductLock{
				BlockHeight:  v.BlockHeight,
				Price:        v.Price,
				Quantity:     v.Quantity,
				BuyExecuted:  v.BuyExecuted,
				SellExecuted: v.SellExecuted,
			}
		}

		cache.LockMap = &ordertypes.ProductLockMap{
			Data: cpData,
		}
	}

	return &cache
}
