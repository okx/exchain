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
	cache := &Cache{
		UpdatedOrderIDs:  []string{"ID00000000054-1","ID00000000056-2"},
		BlockMatchResult: &types.BlockMatchResult{
			BlockHeight: 1000,
			ResultMap:   nil,
			TimeStamp:   10000000240,
		},
		Params:           &types.Params{
			OrderExpireBlocks: 20,
			MaxDealsPerBlock:  30,
			FeePerBlock:       sdk.NewDecCoinFromDec("okt", sdk.MustNewDecFromStr("1000000")),
			TradeFeeRate:      sdk.MustNewDecFromStr("10000"),
		},
		CancelNum:        1,
		ExpireNum:        2,
		PartialFillNum:   3,
		FullFillNum:      4,
	}

	resultMap := make(map[string]types.MatchResult)
	resultMap["okxxxxxx1"] = types.MatchResult{
		BlockHeight: 100001,
		Price:       sdk.MustNewDecFromStr("10000"),
		Quantity:    sdk.MustNewDecFromStr("1000"),
		Deals:       []types.Deal{
			{"ID0000000040-1","BUY",sdk.MustNewDecFromStr("1000"),"1000okt" },
			{"ID0000000040-2","BUY",sdk.MustNewDecFromStr("1000"),"1000okt" },
			{"ID0000000040-3","BUY",sdk.MustNewDecFromStr("1000"),"1000okt" },
			{"ID0000000040-4","BUY",sdk.MustNewDecFromStr("1000"),"1000okt" },
		},
	}
	cache.BlockMatchResult.ResultMap = resultMap

	cloneCopy := cache.Clone()
	require.NotNil(t, cloneCopy)

	depthCopy := cache.DepthCopy()
	require.NotNil(t, depthCopy)

	require.EqualValues(t, len(cache.BlockMatchResult.ResultMap), len(cloneCopy.BlockMatchResult.ResultMap))
	require.EqualValues(t, len(cache.BlockMatchResult.ResultMap), len(depthCopy.BlockMatchResult.ResultMap))
	require.EqualValues(t, len(cache.BlockMatchResult.ResultMap["okxxxxxx1"].Deals), len(cloneCopy.BlockMatchResult.ResultMap["okxxxxxx1"].Deals))
	require.EqualValues(t, len(cache.BlockMatchResult.ResultMap["okxxxxxx1"].Deals), len(depthCopy.BlockMatchResult.ResultMap["okxxxxxx1"].Deals))

	addMatchResult := types.MatchResult{
		BlockHeight: 100002,
		Price:       sdk.MustNewDecFromStr("10000"),
		Quantity:    sdk.MustNewDecFromStr("1000"),
		Deals:       []types.Deal{
			{"ID0000000050-1","SELL",sdk.MustNewDecFromStr("1000"),"1000okt" },
			{"ID0000000050-2","SELL",sdk.MustNewDecFromStr("1000"),"1000okt" },
			{"ID0000000050-3","SELL",sdk.MustNewDecFromStr("1000"),"1000okt" },
			{"ID0000000050-4","SELL",sdk.MustNewDecFromStr("1000"),"1000okt" },
		},
	}
	cache.BlockMatchResult.ResultMap["okxxxxxx2"] = addMatchResult

	require.EqualValues(t, len(cache.BlockMatchResult.ResultMap), len(cloneCopy.BlockMatchResult.ResultMap) + 1)
	require.EqualValues(t, len(cache.BlockMatchResult.ResultMap), len(depthCopy.BlockMatchResult.ResultMap) + 1)
}