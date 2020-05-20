package keeper

import (
	"encoding/json"
	"github.com/okex/okchain/x/order/types"
)

// Cache stores some caches that will not be written to disk
type Cache struct {
	// Reset at BeginBlock
	UpdatedOrderIDs  []string                `json:"update_order_ids"`
	BlockMatchResult *types.BlockMatchResult `json:"block_match_result, omitempty"`

	Params *types.Params `json:"params, omitempty"`

	// for statistic
	CancelNum      int64 `json:"cancel_number"`       // canceled orders num in this block
	ExpireNum      int64 `json:"expire_number"`       // expired orders num in this block
	PartialFillNum int64 `json:"partial_fill_number"` // partially filled orders num in this block
	FullFillNum    int64 `json:"full_fill_number"`    // fully filled orders num in this block
}

// nolint
func NewCache() *Cache {
	return &Cache{
		UpdatedOrderIDs:  []string{},
		BlockMatchResult: nil,
		Params:           nil,
	}
}

// reset resets temporary cache, called at BeginBlock
func (c *Cache) reset() {
	c.UpdatedOrderIDs = []string{}
	c.BlockMatchResult = &types.BlockMatchResult{}
	c.Params = nil

	c.CancelNum = 0
	c.ExpireNum = 0
	c.FullFillNum = 0
	c.PartialFillNum = 0
}

func (c *Cache) addUpdatedOrderID(orderID string) {
	c.UpdatedOrderIDs = append(c.UpdatedOrderIDs, orderID)
}

func (c *Cache) setBlockMatchResult(result *types.BlockMatchResult) {
	c.BlockMatchResult = result
}

// nolint
func (c *Cache) IncreaseExpireNum() int64 {
	c.ExpireNum++
	return c.ExpireNum
}

// --------

// nolint
func (c *Cache) DecreaseCancelNum() int64 {
	c.CancelNum--
	return c.CancelNum
}

// nolint
func (c *Cache) IncreaseCancelNum() int64 {
	c.CancelNum++
	return c.CancelNum
}

// nolint
func (c *Cache) DecreaseFullFillNum() int64 {
	c.FullFillNum--
	return c.FullFillNum
}

// nolint
func (c *Cache) IncreaseFullFillNum() int64 {
	c.FullFillNum++
	return c.FullFillNum
}

// nolint
func (c *Cache) DecreasePartialFillNum() int64 {
	c.PartialFillNum--
	return c.PartialFillNum
}

// nolint
func (c *Cache) IncreasePartialFillNum() int64 {
	c.PartialFillNum++
	return c.PartialFillNum
}

func (c *Cache) getBlockMatchResult() *types.BlockMatchResult {
	return c.BlockMatchResult
}

// nolint
func (c *Cache) SetParams(params *types.Params) {
	c.Params = params
}

// nolint
func (c *Cache) GetParams() *types.Params {
	return c.Params
}

func (c *Cache) getUpdatedOrderIDs() []string {
	return c.UpdatedOrderIDs
}

// nolint
func (c *Cache) GetFullFillNum() int64 {
	return c.FullFillNum
}

// nolint
func (c *Cache) GetCancelNum() int64 {
	return c.CancelNum
}

// nolint
func (c *Cache) GetExpireNum() int64 {
	return c.ExpireNum
}

// nolint
func (c *Cache) GetPartialFillNum() int64 {
	return c.PartialFillNum
}

// nolint
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
		UpdatedOrderIDs:  nil,
		BlockMatchResult: nil,
		Params:           nil,
		CancelNum:        c.CancelNum,
		ExpireNum:        c.ExpireNum,
		PartialFillNum:   c.PartialFillNum,
		FullFillNum:      c.FullFillNum,
	}

	if c.UpdatedOrderIDs != nil {
		cpUpdatedOrderIDs := make([]string, 0, len(c.UpdatedOrderIDs))
		cpUpdatedOrderIDs = append(cpUpdatedOrderIDs, c.UpdatedOrderIDs...)

		cache.UpdatedOrderIDs = cpUpdatedOrderIDs
	}

	if c.BlockMatchResult != nil {
		cache.BlockMatchResult = &types.BlockMatchResult{
			BlockHeight: c.BlockMatchResult.BlockHeight,
			ResultMap:   nil,
			TimeStamp:   c.BlockMatchResult.TimeStamp,
		}

		if c.BlockMatchResult.ResultMap != nil {
			cpResultMap := make(map[string]types.MatchResult)
			for k, v := range c.BlockMatchResult.ResultMap {
				cpDeals := make([]types.Deal, 0, len(v.Deals))
				cpDeals = append(cpDeals, v.Deals...)

				cpResultMap[k] = types.MatchResult{
					BlockHeight: v.BlockHeight,
					Price:       v.Price,
					Quantity:    v.Quantity,
					Deals:       cpDeals,
				}
			}

			cache.BlockMatchResult.ResultMap = cpResultMap
		}
	}

	if c.Params != nil {
		cache.Params = &types.Params{
			OrderExpireBlocks: c.Params.OrderExpireBlocks,
			MaxDealsPerBlock:  c.Params.MaxDealsPerBlock,
			FeePerBlock:       c.Params.FeePerBlock,
			TradeFeeRate:      c.Params.TradeFeeRate,
		}
	}

	return &cache
}
