package keeper

import (
	"testing"

	"github.com/okex/okexchain/x/dex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/order/types"
	token "github.com/okex/okexchain/x/token/types"
)

func TestTryPlaceOrder(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 10)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	// mock order
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "1.0", "10.1")
	order.Sender = testInput.TestAddrs[0]

	// not enough balance
	_, err = keeper.TryPlaceOrder(ctx, order)
	require.Error(t, err)

	// test new order fee
	//feeParams.NewOrder = sdk.MustNewDecFromStr("0.01")
	order.Quantity = sdk.MustNewDecFromStr("9.0")
	fee, err := keeper.TryPlaceOrder(ctx, order)
	require.Nil(t, err)

	order.RecordOrderNewFee(fee)
	require.Equal(t, "0.259200000000000000"+common.NativeToken, order.GetExtraInfoWithKey(types.OrderExtraInfoKeyNewFee))

	order = mockOrder("", types.TestTokenPair, types.BuyOrder, "1.0", "9")
	order.Sender = testInput.TestAddrs[0]
	_, err = keeper.TryPlaceOrder(ctx, order)
	keeper.UnlockCoins(ctx, testInput.TestAddrs[0], sdk.SysCoins{{Denom: common.NativeToken, Amount: sdk.MustNewDecFromStr("9")}}, token.LockCoinsTypeQuantity)
	require.Error(t, err)
}

func TestPlaceOrderAndCancelOrder(t *testing.T) {
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

	// check result & order
	require.EqualValues(t, types.FormatOrderID(10, 1), order.OrderID)
	require.EqualValues(t, 1, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	expectCoins := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("89.7408")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	// check depth book
	depthBook := keeper.GetDepthBookCopy(order.Product)
	require.Equal(t, 1, len(depthBook.Items))
	require.Equal(t, sdk.MustNewDecFromStr("10.0"), depthBook.Items[0].Price)
	require.Equal(t, sdk.MustNewDecFromStr("1.0"), depthBook.Items[0].BuyQuantity)
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids map
	orderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()
	require.Equal(t, 1, len(orderIDsMap.Data))
	require.Equal(t, types.FormatOrderID(10, 1),
		orderIDsMap.Data[types.FormatOrderIDsKey(order.Product, order.Price, order.Side)][0])
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))
	// other check
	require.EqualValues(t, 1, keeper.diskCache.openNum)
	require.EqualValues(t, 1, keeper.diskCache.storeOrderNum)

	// Test cancel order
	ctx = ctx.WithBlockHeight(11)
	fee := keeper.CancelOrder(ctx, order, ctx.Logger())
	// check result
	require.Equal(t, "0.000001000000000000"+common.NativeToken, fee.String())
	// check order status
	require.EqualValues(t, types.OrderStatusCancelled, order.Status)
	require.Equal(t, "", order.GetExtraInfoWithKey(types.OrderExtraInfoKeyCancelFee))
	// check account balance
	acc = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	expectCoins = sdk.SysCoins{
		// 100 - 0.002
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.999999")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	// check fee pool
	feeCollector := testInput.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	require.EqualValues(t, "0.000001000000000000"+common.NativeToken, collectedFees.String())
	// check depth book
	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 0, len(depthBook.Items))
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("9.8"), types.BuyOrder)
	orderIDs := keeper.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))
	// check updated order ids
	updatedOrderIDs := keeper.GetUpdatedOrderIDs()
	require.EqualValues(t, order.OrderID, updatedOrderIDs[0])
	// check closed order id
	closedOrderIDs := keeper.GetDiskCache().GetClosedOrderIDs()
	require.Equal(t, 1, len(closedOrderIDs))
	require.Equal(t, order.OrderID, closedOrderIDs[0])
	// other check
	require.EqualValues(t, 0, keeper.diskCache.openNum)
	require.EqualValues(t, 1, keeper.cache.cancelNum)
}

func TestPlaceOrderAndExpireOrder(t *testing.T) {
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

	// check result & order
	require.EqualValues(t, types.FormatOrderID(10, 1), order.OrderID)
	require.EqualValues(t, 1, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	expectCoins := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("89.7408")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	// check depth book
	depthBook := keeper.GetDepthBookCopy(order.Product)
	require.Equal(t, 1, len(depthBook.Items))
	require.Equal(t, sdk.MustNewDecFromStr("10.0"), depthBook.Items[0].Price)
	require.Equal(t, sdk.MustNewDecFromStr("1.0"), depthBook.Items[0].BuyQuantity)
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids map
	orderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()
	require.Equal(t, 1, len(orderIDsMap.Data))
	require.Equal(t, types.FormatOrderID(10, 1),
		orderIDsMap.Data[types.FormatOrderIDsKey(order.Product, order.Price, order.Side)][0])
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))
	// other check
	require.EqualValues(t, 1, keeper.diskCache.openNum)
	require.EqualValues(t, 1, keeper.diskCache.storeOrderNum)

	// Test expire order
	ctx = ctx.WithBlockHeight(11)
	keeper.ExpireOrder(ctx, order, ctx.Logger())
	// check order status
	require.EqualValues(t, types.OrderStatusExpired, order.Status)
	require.Equal(t, "", order.GetExtraInfoWithKey(types.OrderExtraInfoKeyExpireFee))
	// check account balance
	acc = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	expectCoins = sdk.SysCoins{
		// 100 - 0.002
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.999999")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	// check fee pool
	feeCollector := testInput.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	require.EqualValues(t, "0.000001000000000000"+common.NativeToken, collectedFees.String())
	// check depth book
	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 0, len(depthBook.Items))
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("9.8"), types.BuyOrder)
	orderIDs := keeper.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))
	// check updated order ids
	updatedOrderIDs := keeper.GetUpdatedOrderIDs()
	require.EqualValues(t, order.OrderID, updatedOrderIDs[0])
	// check closed order id
	keeper.Cache2Disk(ctx)
	closedOrderIDs := keeper.GetLastClosedOrderIDs(ctx)
	require.Equal(t, 1, len(closedOrderIDs))
	require.Equal(t, order.OrderID, closedOrderIDs[0])
	// other check
	require.EqualValues(t, 0, keeper.diskCache.openNum)
	require.EqualValues(t, 1, keeper.cache.expireNum)
}
