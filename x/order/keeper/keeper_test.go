package keeper

import (
	"testing"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order/types"
	token "github.com/okex/okchain/x/token/types"
)

func TestKeeper_Cache2Disk(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
	feeParams.OrderExpireBlocks = 1
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock order
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "8", "1")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	require.EqualValues(t, 0, len(keeper.GetDepthBookFromDB(ctx, types.TestTokenPair).Items))

	// flush
	keeper.Cache2Disk(ctx)
	require.EqualValues(t, 1, len(keeper.GetDepthBookFromDB(ctx, types.TestTokenPair).Items))

	keeper.RemoveOrderFromDepthBook(order, types.FeeTypeOrderCancel)
	require.EqualValues(t, 1, keeper.cache.CancelNum)
	keeper.RemoveOrderFromDepthBook(order, types.FeeTypeOrderExpire)
	require.EqualValues(t, 1, keeper.cache.ExpireNum)
}

func TestCache(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
	feeParams.OrderExpireBlocks = 1
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	err = keeper.AddCollectedFees(ctx,
		sdk.DecCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("0.25920000")}},
		testInput.TestAddrs[0], types.FeeTypeOrderExpire, true)
	require.Nil(t, err)
	tokenKeeper := keeper.GetTokenKeeper()
	require.NotNil(t, tokenKeeper)

	// mock order
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)
	require.EqualValues(t, 1, keeper.GetCache().Params.OrderExpireBlocks)
	require.EqualValues(t, 1, keeper.GetOperationMetric().OpenNum)

	// current cache
	require.EqualValues(t, 1, keeper.diskCache.OpenNum)
	require.EqualValues(t, 1, keeper.diskCache.StoreOrderNum)
	copycache := keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(copycache.Items))
	require.EqualValues(t, 0, len(keeper.GetDepthBookCopy("ABC").Items))
	require.EqualValues(t, 1, len(keeper.GetDiskCache().getDepthBook(types.TestTokenPair).Items))
	require.EqualValues(t, 0, len(keeper.GetDepthBookFromDB(ctx, types.TestTokenPair).Items))
	require.EqualValues(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))

	// flush
	keeper.Cache2Disk(ctx)
	require.EqualValues(t, 1, len(keeper.GetDiskCache().getDepthBook(types.TestTokenPair).Items))
	require.EqualValues(t, 1, len(keeper.GetDepthBookFromDB(ctx, types.TestTokenPair).Items))

	// new cache
	keeper.cache = NewCache()
	keeper.diskCache = newDiskCache()
	require.EqualValues(t, 0, keeper.diskCache.OpenNum)
	require.EqualValues(t, 0, keeper.diskCache.StoreOrderNum)
	require.Nil(t, keeper.GetDiskCache().getDepthBook(types.TestTokenPair))
	keeper.ResetCache(ctx)
	require.EqualValues(t, 1, keeper.diskCache.OpenNum)
	require.EqualValues(t, 1, keeper.diskCache.StoreOrderNum)
	require.EqualValues(t, 1, len(keeper.GetDiskCache().getDepthBook(types.TestTokenPair).Items))

	params := keeper.GetParams(ctx)
	require.EqualValues(t, 1, params.OrderExpireBlocks)

	keeper.SetMetric()
	met := keeper.GetMetric()
	require.NotNil(t, met)
}

func TestKeeper_LockCoins(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
	feeParams.OrderExpireBlocks = 1
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock order
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "8", "1")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	//not enough coin to lock
	err = keeper.LockCoins(ctx, testInput.TestAddrs[0], sdk.DecCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("99.7408")}}, token.LockCoinsTypeQuantity)
	require.NotNil(t, err)

	//lock coin
	err = keeper.LockCoins(ctx, testInput.TestAddrs[0], sdk.DecCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("90")}}, token.LockCoinsTypeQuantity)
	require.Nil(t, err)

	//not enough coin to placeorder
	err = keeper.PlaceOrder(ctx, order)
	require.NotNil(t, err)

	keeper.UnlockCoins(ctx, testInput.TestAddrs[0], sdk.DecCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("90")}}, token.LockCoinsTypeQuantity)

	//placeorder success
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)
}

func TestKeeper_BurnLockedCoins(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
	feeParams.OrderExpireBlocks = 1
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock order
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "10", "1")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	keeper.BalanceAccount(ctx, testInput.TestAddrs[0],
		sdk.MustParseCoins(common.NativeToken, "8"),
		sdk.MustParseCoins(common.NativeToken, "9"),
	)

	keeper.CancelOrder(ctx, order, nil)
	require.EqualValues(t, 0, len(keeper.GetLastClosedOrderIDs(ctx)))

	account := keeper.GetCoins(ctx, testInput.TestAddrs[0])
	require.NotNil(t, account)
	require.EqualValues(t, "99.00000000", account[0].Amount.String())
}

func TestLastPrice(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	price := keeper.GetLastPrice(ctx, "xxxb_"+common.NativeToken)
	require.Equal(t, sdk.ZeroDec(), price)

	keeper.SetLastPrice(ctx, "xxxb_"+common.NativeToken, sdk.MustNewDecFromStr("9.9"))
	price = keeper.GetLastPrice(ctx, "xxxb_"+common.NativeToken)
	require.Equal(t, sdk.MustNewDecFromStr("9.9"), price)

	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), keeper.GetLastPrice(ctx, types.TestTokenPair))
	keeper.SetLastPrice(ctx, types.TestTokenPair, sdk.MustNewDecFromStr("1.234"))
	require.EqualValues(t, sdk.MustNewDecFromStr("1.234"), keeper.GetLastPrice(ctx, types.TestTokenPair))

}

func TestLastExpiredBlockHeight(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	require.EqualValues(t, 0, keeper.GetLastExpiredBlockHeight(ctx))
	keeper.SetLastExpiredBlockHeight(ctx, 100)
	require.EqualValues(t, 100, keeper.GetLastExpiredBlockHeight(ctx))
}

func TestBlockOrderNum(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	require.EqualValues(t, 0, keeper.GetBlockOrderNum(ctx, 10))
	keeper.SetBlockOrderNum(ctx, 10, 1)
	require.EqualValues(t, 1, keeper.GetBlockOrderNum(ctx, 10))
	keeper.DropBlockOrderNum(ctx, 10)
	require.EqualValues(t, 0, keeper.GetBlockOrderNum(ctx, 10))
}

func TestExpireBlockHeight(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	var blockHeight int64 = 10
	expireBlockHeight := []int64{1, 2}
	require.EqualValues(t, []int64{}, keeper.GetExpireBlockHeight(ctx, blockHeight))
	keeper.SetExpireBlockHeight(ctx, blockHeight, expireBlockHeight)
	require.EqualValues(t, expireBlockHeight, keeper.GetExpireBlockHeight(ctx, blockHeight))
	keeper.DropExpireBlockHeight(ctx, blockHeight)
	require.EqualValues(t, []int64{}, keeper.GetExpireBlockHeight(ctx, blockHeight))
}

func TestBlockMatchResult(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper

	require.Nil(t, keeper.GetBlockMatchResult())
	blockMatchResult := &types.BlockMatchResult{
		BlockHeight: 1,
		ResultMap:   make(map[string]types.MatchResult),
		TimeStamp:   1,
	}
	keeper.SetBlockMatchResult(blockMatchResult)
	require.EqualValues(t, blockMatchResult, keeper.GetBlockMatchResult())

}

func TestDropOrder(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx.WithBlockHeight(10)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	order2 := keeper.GetOrder(ctx, order.OrderID)
	require.EqualValues(t, order2, order)

	keeper.DropOrder(ctx, order.OrderID)
	require.Nil(t, keeper.GetOrder(ctx, order.OrderID))
}

func TestKeeper_UpdateOrder(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx.WithBlockHeight(10)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	//xxb_okt:10.00000000:BUY
	require.EqualValues(t, "ID0000000010-1", keeper.diskCache.OrderIDsMap.Data["xxb_"+common.NativeToken+":10.00000000:BUY"][0])

	order.Price = sdk.MustNewDecFromStr("11.0")

	//update order
	keeper.UpdateOrder(order, ctx)
	require.EqualValues(t, 1, len(keeper.GetUpdatedOrderIDs()))
	require.EqualValues(t, 1, len(keeper.GetUpdatedDepthbookKeys()))
	require.EqualValues(t, 0, len(keeper.GetProductPriceOrderIDs("abc")))
	require.EqualValues(t, 0, len(keeper.GetProductPriceOrderIDs(types.TestTokenPair)))
	require.EqualValues(t, 0, len(keeper.GetProductPriceOrderIDsFromDB(ctx, types.TestTokenPair)))

	require.EqualValues(t, "ID0000000010-1", keeper.cache.UpdatedOrderIDs[0])

	order.Status = types.OrderStatusFilled
	keeper.UpdateOrder(order, ctx)
	require.EqualValues(t, 1, keeper.cache.FullFillNum)

}

func TestKeeper_SendFeesToProductOwner(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
	feeParams.OrderExpireBlocks = 1
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// place order err
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "108", "1")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.NotNil(t, err)

	// mock order
	order = mockOrder("", types.TestTokenPair, types.BuyOrder, "8", "1")
	order.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	fee := GetOrderNewFee(order)

	dealFee := sdk.DecCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("0.2592")}}
	require.EqualValues(t, fee, dealFee)

	err = keeper.SendFeesToProductOwner(ctx, dealFee, order.Sender, types.FeeTypeOrderDeal, order.Product)
	require.Nil(t, err)
}

func TestKeeper_GetBestBidAndAsk(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	ask, _ := keeper.GetBestBidAndAsk(ctx, types.TestTokenPair)
	require.EqualValues(t, sdk.MustNewDecFromStr("0"), ask)
}

func TestKeeper_InsertOrderIntoDepthBook(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper

	// mock order
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "8", "1")
	order.Sender = testInput.TestAddrs[0]
	keeper.InsertOrderIntoDepthBook(order)
	require.EqualValues(t, 1, len(keeper.diskCache.DepthBookMap.Data))
}

func TestFilterDelistedProducts(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx.WithBlockHeight(10)

	tokenPair := dex.GetBuiltInTokenPair()
	err := keeper.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	productsList := []string{
		"xxb_yyb",
		types.TestTokenPair,
		"ttb_qqb",
	}

	expectedProductsList := []string{
		types.TestTokenPair,
	}

	cleanProducts := keeper.FilterDelistedProducts(ctx, productsList)
	require.EqualValues(t, expectedProductsList, cleanProducts)
}
