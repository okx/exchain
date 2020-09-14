package keeper

import (
	"github.com/okex/okexchain/x/order/types"
	"github.com/willf/bitset"
)

// Cache stores some caches that will not be written to disk
type Cache struct {
	// Reset at BeginBlock
	updatedOrderIDs  []string
	blockMatchResult *types.BlockMatchResult
	handlerTxMsgResult []bitset.BitSet

	params *types.Params

	// for statistic
	cancelNum      int64 // canceled orders num in this block
	expireNum      int64 // expired orders num in this block
	partialFillNum int64 // partially filled orders num in this block
	fullFillNum    int64 // fully filled orders num in this block
}

// nolint
func NewCache() *Cache {
	return &Cache{
		updatedOrderIDs:  []string{},
		blockMatchResult: nil,
		params:           nil,
	}
}

// reset resets temporary cache, called at BeginBlock
func (c *Cache) reset() {
	c.updatedOrderIDs = []string{}
	c.blockMatchResult = &types.BlockMatchResult{}
	c.handlerTxMsgResult = []bitset.BitSet{}
	c.params = nil

	c.cancelNum = 0
	c.expireNum = 0
	c.fullFillNum = 0
	c.partialFillNum = 0
}

func (c *Cache) addUpdatedOrderID(orderID string) {
	c.updatedOrderIDs = append(c.updatedOrderIDs, orderID)
}

func (c *Cache) setBlockMatchResult(result *types.BlockMatchResult) {
	c.blockMatchResult = result
}

func (c *Cache) addTxHandlerMsgResult(resultSet bitset.BitSet) {
	c.handlerTxMsgResult = append(c.handlerTxMsgResult, resultSet)
}

// nolint
func (c *Cache) IncreaseExpireNum() int64 {
	c.expireNum++
	return c.expireNum
}

// --------

// nolint
func (c *Cache) DecreaseCancelNum() int64 {
	c.cancelNum--
	return c.cancelNum
}

// nolint
func (c *Cache) IncreaseCancelNum() int64 {
	c.cancelNum++
	return c.cancelNum
}

// nolint
func (c *Cache) DecreaseFullFillNum() int64 {
	c.fullFillNum--
	return c.fullFillNum
}

// nolint
func (c *Cache) IncreaseFullFillNum() int64 {
	c.fullFillNum++
	return c.fullFillNum
}

// nolint
func (c *Cache) DecreasePartialFillNum() int64 {
	c.partialFillNum--
	return c.partialFillNum
}

// nolint
func (c *Cache) IncreasePartialFillNum() int64 {
	c.partialFillNum++
	return c.partialFillNum
}

func (c *Cache) getBlockMatchResult() *types.BlockMatchResult {
	return c.blockMatchResult
}

// nolint
func (c *Cache) SetParams(params *types.Params) {
	c.params = params
}

// nolint
func (c *Cache) GetParams() *types.Params {
	return c.params
}

func (c *Cache) getUpdatedOrderIDs() []string {
	return c.updatedOrderIDs
}

// toggleCopyTxHandlerMsgResult: copy and reset the handlerTxMsgResult
func (c *Cache) toggleCopyTxHandlerMsgResult() []bitset.BitSet {
	txResultCopy := make([]bitset.BitSet, 0, len(c.handlerTxMsgResult))
	txResultCopy = append(txResultCopy, c.handlerTxMsgResult...)
	c.handlerTxMsgResult = []bitset.BitSet{}

	return txResultCopy
}

// nolint
func (c *Cache) GetFullFillNum() int64 {
	return c.fullFillNum
}

// nolint
func (c *Cache) GetCancelNum() int64 {
	return c.cancelNum
}

// nolint
func (c *Cache) GetExpireNum() int64 {
	return c.expireNum
}

// nolint
func (c *Cache) GetPartialFillNum() int64 {
	return c.partialFillNum
}
