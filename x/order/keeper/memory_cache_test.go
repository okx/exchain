package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"time"

	"github.com/okex/okchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestCache_GetCancelNum(t *testing.T) {
	cache := NewCache()
	cache.addUpdatedOrderID("ID0000000010-1")
	require.EqualValues(t, 1, len(cache.UpdatedOrderIDs))

	cache.addUpdatedOrderID("ID0000000010-2")
	require.EqualValues(t, 2, len(cache.getUpdatedOrderIDs()))

	cache.IncreaseCancelNum()
	cache.IncreaseExpireNum()
	cache.IncreaseFullFillNum()
	cache.IncreasePartialFillNum()

	cache.DecreaseCancelNum()
	cache.DecreaseFullFillNum()
	cache.DecreasePartialFillNum()

	require.EqualValues(t, 0, cache.GetCancelNum())
	require.EqualValues(t, 1, cache.GetExpireNum())
	require.EqualValues(t, 0, cache.GetFullFillNum())
	require.EqualValues(t, 0, cache.GetPartialFillNum())

	feeParams := types.DefaultTestParams()
	cache.SetParams(&feeParams)

	require.EqualValues(t, types.DefaultOrderExpireBlocks, cache.Params.OrderExpireBlocks)

	res := types.BlockMatchResult{
		BlockHeight: 0,
		ResultMap:   nil,
		TimeStamp:   time.Now().Unix(),
	}
	cache.setBlockMatchResult(&res)

	require.NotEqual(t, 0, cache.getBlockMatchResult().TimeStamp)

	cache.reset()
	require.Nil(t, cache.GetParams())

}

func TestCache_Clone(t *testing.T) {
	cache := NewCache()
	cloneCache := cache.Clone()
	require.EqualValues(t, len(cache.UpdatedOrderIDs), len(cloneCache.UpdatedOrderIDs))

	cache.addUpdatedOrderID("ID0000000010-1")
	cache.addUpdatedOrderID("ID0000000010-2")
	require.EqualValues(t, len(cache.UpdatedOrderIDs), len(cloneCache.UpdatedOrderIDs) + 2)

	cache.Params = &types.Params{
		OrderExpireBlocks: 0,
		MaxDealsPerBlock:  100,
		FeePerBlock:       sdk.DecCoin{},
		TradeFeeRate:      sdk.Dec{},
	}
	require.Nil(t, cloneCache.Params)
	cloneCache.Params = &types.Params{
		OrderExpireBlocks: 0,
		MaxDealsPerBlock:  10000,
		FeePerBlock:       sdk.DecCoin{},
		TradeFeeRate:      sdk.Dec{},
	}
	require.NotEqual(t, cache.Params.MaxDealsPerBlock, cloneCache.Params.MaxDealsPerBlock)

	cache.BlockMatchResult = &types.BlockMatchResult{
		BlockHeight: 200,
		ResultMap:   nil,
		TimeStamp:   0,
	}
	require.Nil(t, cloneCache.BlockMatchResult)
	cloneCache.BlockMatchResult = &types.BlockMatchResult{
		BlockHeight: 20,
		ResultMap:   nil,
		TimeStamp:   0,
	}
	require.NotEqual(t, cache.BlockMatchResult.BlockHeight, cloneCache.BlockMatchResult.BlockHeight)
}