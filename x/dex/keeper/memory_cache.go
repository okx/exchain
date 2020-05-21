package keeper

//// Cache caches data
//type Cache struct {
//	//tokenPairMap     map[string]*types.TokenPair
//	//newTokenPairMap  []*types.TokenPair
//	//tokenPairChanged bool
//
//	//lockMap          *ordertypes.ProductLockMap
//}

//// NewCache returns instance of Cache
//func NewCache() *Cache {
//	return &Cache{
//		//tokenPairMap:     make(map[string]*types.TokenPair),
//		//tokenPairChanged: false,
//		//lockMap:          ordertypes.NewProductLockMap(),
//	}
//}

//// Reset clears cache
//func (c *Cache) Reset() {
//	//c.newTokenPairMap = []*types.TokenPair{}
//	//c.tokenPairChanged = false
//}

//// AddTokenPair adds a new token pair into cache
//func (c *Cache) AddTokenPair(tokenPair *types.TokenPair) {
//	c.tokenPairMap[c.genTokenPairKey(tokenPair)] = tokenPair
//	//c.tokenPairChanged = true
//}

////generate token pair key
//func (c *Cache) genTokenPairKey(t *types.TokenPair) string {
//	return t.BaseAssetSymbol + "_" + t.QuoteAssetSymbol
//}

//// GetTokenPair returns token pair from cache
//func (c *Cache) GetTokenPair(product string) (*types.TokenPair, bool) {
//	tokenPair, ok := c.tokenPairMap[product]
//	return tokenPair, ok
//}

//// PrepareTokenPairs adds all token pairs
//func (c *Cache) PrepareTokenPairs(tokenPairs []*types.TokenPair) {
//	for _, v := range tokenPairs {
//		c.AddTokenPair(v)
//	}
//}

//// GetAllTokenPairs returns all token pairs from cache
//func (c *Cache) GetAllTokenPairs() []*types.TokenPair {
//	tokenPairs := make([]*types.TokenPair, 0, len(c.tokenPairMap))
//	for _, v := range c.tokenPairMap {
//		tokenPairs = append(tokenPairs, v)
//	}
//	return tokenPairs
//}

//// DeleteTokenPair deletes token pair cache
//func (c *Cache) DeleteTokenPair(targetPair *types.TokenPair) {
//	delete(c.tokenPairMap, c.genTokenPairKey(targetPair))
//	//c.tokenPairChanged = true
//}

//// DeleteTokenPairByName deletes token pair by token pair's name from cache
//func (c *Cache) DeleteTokenPairByName(tokenPairName string) {
//	delete(c.tokenPairMap, tokenPairName)
//	//c.tokenPairChanged = true
//}

//// TokenPairCount returns count of token pair
//func (c *Cache) TokenPairCount() int {
//	return len(c.tokenPairMap)
//}

//// AddNewTokenPair adds a new token pair into cache
//func (c *Cache) AddNewTokenPair(tokenPair *tyGetNewTokenPairpes.TokenPair) {
//	//c.newTokenPairMap = append(c.newTokenPairMap, tokenPair)
//	//c.tokenPairChanged = true
//}

//// GetNewTokenPair returns new token pairs from cache
//func (c *Cache) GetNewTokenPair() []*types.TokenPair {
//	return c.newTokenPairMap
//}
