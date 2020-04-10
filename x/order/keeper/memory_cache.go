package keeper

import (
	"github.com/okex/okchain/x/order/types"
)

type Cache struct {
	// Reset at BeginBlock
	updatedOrderIDs  []string
	blockMatchResult *types.BlockMatchResult

	params *types.Params

	// for statistic
	cancelNum      int64 // canceled orders num in this block
	expireNum      int64 // expired orders num in this block
	partialFillNum int64 // partially filled orders num in this block
	fullFillNum    int64 // fully filled orders num in this block
}

func NewCache() *Cache {
	return &Cache{
		updatedOrderIDs:  []string{},
		blockMatchResult: nil,
		params:           nil,
	}
}

// Reset temporary cache, called at BeginBlock
func (c *Cache) reset() {
	c.updatedOrderIDs = []string{}
	c.blockMatchResult = &types.BlockMatchResult{}
	c.params = nil

	c.cancelNum = 0
	c.expireNum = 0
	c.fullFillNum = 0
	c.partialFillNum = 0
}

// updatedOrderIDs
func (c *Cache) addUpdatedOrderID(orderID string) {
	c.updatedOrderIDs = append(c.updatedOrderIDs, orderID)
}

// blockMatchResult
func (c *Cache) setBlockMatchResult(result *types.BlockMatchResult) {
	c.blockMatchResult = result
}

func (c *Cache) IncreaseExpireNum() int64 {
	c.expireNum++
	return c.expireNum
}

// --------

// --------
func (c *Cache) DecreaseCancelNum() int64 {
	c.cancelNum--
	return c.cancelNum
}

func (c *Cache) IncreaseCancelNum() int64 {
	c.cancelNum++
	return c.cancelNum
}

func (c *Cache) DecreaseFullFillNum() int64 {
	c.fullFillNum--
	return c.fullFillNum
}

func (c *Cache) IncreaseFullFillNum() int64 {
	c.fullFillNum++
	return c.fullFillNum
}

func (c *Cache) DecreasePartialFillNum() int64 {
	c.partialFillNum--
	return c.partialFillNum
}

func (c *Cache) IncreasePartialFillNum() int64 {
	c.partialFillNum++
	return c.partialFillNum
}

// ======================
func (c *Cache) getBlockMatchResult() *types.BlockMatchResult {
	return c.blockMatchResult
}

// params
func (c *Cache) SetParams(params *types.Params) {
	c.params = params
}

func (c *Cache) GetParams() *types.Params {
	return c.params
}

func (c *Cache) getUpdatedOrderIDs() []string {
	return c.updatedOrderIDs
}

func (c *Cache) GetFullFillNum() int64 {
	return c.fullFillNum
}

func (c *Cache) GetCancelNum() int64 {
	return c.cancelNum
}

func (c *Cache) GetExpireNum() int64 {
	return c.expireNum
}

func (c *Cache) GetPartialFillNum() int64 {
	return c.partialFillNum
}
