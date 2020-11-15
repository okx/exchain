package order

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/order/types"
	token "github.com/okex/okexchain/x/token/types"
)

func TestEndBlockerPeriodicMatch(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	feeParams := types.DefaultTestParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	mapp.dexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders
	orders := []*types.Order{
		types.MockOrder(types.FormatOrderID(startHeight, 1), types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.MockOrder(types.FormatOrderID(startHeight, 2), types.TestTokenPair, types.SellOrder, "10.0", "0.5"),
		types.MockOrder(types.FormatOrderID(startHeight, 3), types.TestTokenPair, types.SellOrder, "10.0", "2.5"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	orders[1].Sender = addrKeysSlice[1].Address
	orders[2].Sender = addrKeysSlice[1].Address
	for i := 0; i < 3; i++ {
		err := k.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	// subtract all okb of addr0
	// 100 - 10 - 0.2592
	err = k.LockCoins(ctx, addrKeysSlice[0].Address, sdk.SysCoins{{Denom: common.NativeToken,
		Amount: sdk.MustNewDecFromStr("89.7408")}}, token.LockCoinsTypeQuantity)
	require.NoError(t, err)

	// call EndBlocker to execute periodic match
	EndBlocker(ctx, k)

	// check order status
	order0 := k.GetOrder(ctx, orders[0].OrderID)
	order1 := k.GetOrder(ctx, orders[1].OrderID)
	order2 := k.GetOrder(ctx, orders[2].OrderID)
	require.EqualValues(t, types.OrderStatusFilled, order0.Status)
	require.EqualValues(t, types.OrderStatusFilled, order1.Status)
	require.EqualValues(t, types.OrderStatusOpen, order2.Status)
	require.EqualValues(t, sdk.MustNewDecFromStr("2"), order2.RemainQuantity)

	// check depth book
	depthBook := k.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), depthBook.Items[0].Price)
	require.True(sdk.DecEq(t, sdk.ZeroDec(), depthBook.Items[0].BuyQuantity))
	require.EqualValues(t, sdk.MustNewDecFromStr("2"), depthBook.Items[0].SellQuantity)

	depthBookDB := k.GetDepthBookFromDB(ctx, types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBookDB.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), depthBookDB.Items[0].Price)
	require.True(sdk.DecEq(t, sdk.ZeroDec(), depthBookDB.Items[0].BuyQuantity))
	require.EqualValues(t, sdk.MustNewDecFromStr("2"), depthBookDB.Items[0].SellQuantity)

	// check product price - order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.SellOrder)
	orderIDs := k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 1, len(orderIDs))
	require.EqualValues(t, order2.OrderID, orderIDs[0])
	key = types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.BuyOrder)
	orderIDs = k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))

	// check block match result
	result := k.GetBlockMatchResult()
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), result.ResultMap[types.TestTokenPair].Price)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), result.ResultMap[types.TestTokenPair].Quantity)
	require.EqualValues(t, 3, len(result.ResultMap[types.TestTokenPair].Deals))
	require.EqualValues(t, order0.OrderID, result.ResultMap[types.TestTokenPair].Deals[0].OrderID)
	require.EqualValues(t, order1.OrderID, result.ResultMap[types.TestTokenPair].Deals[1].OrderID)
	require.EqualValues(t, order2.OrderID, result.ResultMap[types.TestTokenPair].Deals[2].OrderID)
	// check closed order id
	closedOrderIDs := k.GetLastClosedOrderIDs(ctx)
	require.Equal(t, 2, len(closedOrderIDs))
	require.Equal(t, orders[0].OrderID, closedOrderIDs[0])
	require.Equal(t, orders[1].OrderID, closedOrderIDs[1])

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("0.2592")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100.999")), // 100 + 1 * (1 - 0.001)
	}
	expectCoins1 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("109.7308")), // 100 + 10 * (1-0.001) - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97")),         // 100 - 0.5 - 2.5
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// check fee pool
	feeCollector := mapp.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	require.EqualValues(t, "", collectedFees.String())
}

func TestEndBlockerPeriodicMatchBusyProduct(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	feeParams := types.DefaultTestParams()
	feeParams.MaxDealsPerBlock = 2
	k.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	mapp.dexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders
	orders := []*types.Order{
		types.MockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.MockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "0.5"),
		types.MockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "2.5"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	orders[1].Sender = addrKeysSlice[1].Address
	orders[2].Sender = addrKeysSlice[1].Address
	for i := 0; i < 3; i++ {
		err := k.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	// ------- call EndBlocker at height 10 -------//
	EndBlocker(ctx, k)

	// check product lock
	lock := k.GetDexKeeper().GetLockedProductsCopy(ctx).Data[types.TestTokenPair]
	require.NotNil(t, lock)
	require.EqualValues(t, 10, lock.BlockHeight)
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), lock.Price)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), lock.Quantity)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), lock.BuyExecuted)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.5"), lock.SellExecuted)

	// check order status
	order0 := k.GetOrder(ctx, orders[0].OrderID)
	order1 := k.GetOrder(ctx, orders[1].OrderID)
	order2 := k.GetOrder(ctx, orders[2].OrderID)
	require.EqualValues(t, types.OrderStatusFilled, order0.Status)
	require.EqualValues(t, types.OrderStatusFilled, order1.Status)
	require.EqualValues(t, types.OrderStatusOpen, order2.Status)
	require.EqualValues(t, sdk.MustNewDecFromStr("2.5"), order2.RemainQuantity)

	// check depth book
	depthBook := k.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), depthBook.Items[0].Price)
	require.True(sdk.DecEq(t, sdk.ZeroDec(), depthBook.Items[0].BuyQuantity))
	require.EqualValues(t, sdk.MustNewDecFromStr("2.5"), depthBook.Items[0].SellQuantity)

	// check product price - order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.SellOrder)
	orderIDs := k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 1, len(orderIDs))
	require.EqualValues(t, order2.OrderID, orderIDs[0])
	key = types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.BuyOrder)
	orderIDs = k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))

	// check block match result
	result := k.GetBlockMatchResult()
	require.EqualValues(t, 10, result.ResultMap[types.TestTokenPair].BlockHeight)
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), result.ResultMap[types.TestTokenPair].Price)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), result.ResultMap[types.TestTokenPair].Quantity)
	require.EqualValues(t, 2, len(result.ResultMap[types.TestTokenPair].Deals))
	require.EqualValues(t, order0.OrderID, result.ResultMap[types.TestTokenPair].Deals[0].OrderID)
	require.EqualValues(t, order1.OrderID, result.ResultMap[types.TestTokenPair].Deals[1].OrderID)

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("90")),    // 100 - 10
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100.999")), // 100 + 1 * (1 - 0.001)
	}
	expectCoins1 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("104.7358")), // 100 + 5 * (1 - 0.001) - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97")),         // 100 - 0.5 - 2.5
	}

	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// ------- call EndBlock at height 11, continue filling ------- //
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(11)
	BeginBlocker(ctx, k)
	EndBlocker(ctx, k)

	// check product lock
	lock = k.GetDexKeeper().GetLockedProductsCopy(ctx).Data[types.TestTokenPair]
	require.Nil(t, lock)

	// check order status
	order2 = k.GetOrder(ctx, orders[2].OrderID)
	require.EqualValues(t, types.OrderStatusOpen, order2.Status)
	require.EqualValues(t, sdk.MustNewDecFromStr("2.0"), order2.RemainQuantity)

	// check depth book
	depthBook = k.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, sdk.MustNewDecFromStr("2.0"), depthBook.Items[0].SellQuantity)

	// check block match result
	result = k.GetBlockMatchResult()
	require.EqualValues(t, 10, result.ResultMap[types.TestTokenPair].BlockHeight)
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), result.ResultMap[types.TestTokenPair].Price)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), result.ResultMap[types.TestTokenPair].Quantity)
	require.EqualValues(t, 1, len(result.ResultMap[types.TestTokenPair].Deals))
	require.EqualValues(t, order2.OrderID, result.ResultMap[types.TestTokenPair].Deals[0].OrderID)

	// check account balance
	acc0 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("90")),    // 100 - 10
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100.999")), // 100 + 1 * (1 - 0.001)
	}
	expectCoins1 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("109.7308")), // 100 + 10 * (1 - 0.001) - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97")),         // 100 - 0.5 - 2.5
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())
}

func TestEndBlockerDropExpireData(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	feeParams := types.DefaultTestParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock orders
	orders := []*types.Order{
		types.MockOrder("", types.TestTokenPair, types.BuyOrder, "9.8", "1.0"),
		types.MockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "1.5"),
		types.MockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	orders[1].Sender = addrKeysSlice[1].Address
	orders[2].Sender = addrKeysSlice[0].Address
	for i := 0; i < 3; i++ {
		err := k.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	EndBlocker(ctx, k) // update blockMatchResult, updatedOrderIds

	// check before expire: order, blockOrderNum, blockMatchResult, updatedOrderIDs
	require.NotNil(t, k.GetOrder(ctx, orders[1].OrderID))
	require.EqualValues(t, 3, k.GetBlockOrderNum(ctx, 10))
	blockMatchResult := k.GetBlockMatchResult()
	require.NotNil(t, blockMatchResult)
	updatedOrderIDs := k.GetUpdatedOrderIDs()
	require.EqualValues(t, []string{orders[2].OrderID, orders[1].OrderID}, updatedOrderIDs)

	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(11)
	EndBlocker(ctx, k)
	// call EndBlocker to expire orders
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10 + feeParams.OrderExpireBlocks)
	param := types.DefaultTestParams()
	mapp.orderKeeper.SetParams(ctx, &param)
	EndBlocker(ctx, k)

	order0 := k.GetOrder(ctx, orders[0].OrderID)
	order1 := k.GetOrder(ctx, orders[1].OrderID)
	order2 := k.GetOrder(ctx, orders[2].OrderID)

	require.EqualValues(t, types.OrderStatusExpired, order0.Status)
	require.EqualValues(t, types.OrderStatusPartialFilledExpired, order1.Status)
	require.Nil(t, order2)

	// call EndBlocker to drop expire orders
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(11 + feeParams.OrderExpireBlocks)
	EndBlocker(ctx, k)

	// check after expire: order, blockOrderNum, blockMatchResult, updatedOrderIDs
	require.Nil(t, k.GetOrder(ctx, orders[0].OrderID))
	require.Nil(t, k.GetOrder(ctx, orders[1].OrderID))
	require.EqualValues(t, 0, k.GetBlockOrderNum(ctx, 10))
}

// test order expire when product is busy
func TestEndBlockerExpireOrdersBusyProduct(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	feeParams := types.DefaultTestParams()

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	// mock orders
	orders := []*types.Order{
		types.MockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "2.0"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	err = k.PlaceOrder(ctx, orders[0])
	require.NoError(t, err)
	EndBlocker(ctx, k)
	// call EndBlocker at 86400 + 9
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).
		WithBlockHeight(9 + feeParams.OrderExpireBlocks)
	EndBlocker(ctx, k)

	// call EndBlocker at 86400 + 10, lock product
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).
		WithBlockHeight(10 + feeParams.OrderExpireBlocks)
	lock := &types.ProductLock{
		Price:        sdk.MustNewDecFromStr("10.0"),
		Quantity:     sdk.MustNewDecFromStr("1.0"),
		BuyExecuted:  sdk.MustNewDecFromStr("1.0"),
		SellExecuted: sdk.MustNewDecFromStr("1.0"),
	}
	k.SetProductLock(ctx, types.TestTokenPair, lock)
	EndBlocker(ctx, k)

	// check order
	order := k.GetOrder(ctx, orders[0].OrderID)
	require.EqualValues(t, types.OrderStatusOpen, order.Status)
	require.EqualValues(t, 9+feeParams.OrderExpireBlocks, k.GetLastExpiredBlockHeight(ctx))

	// call EndBlocker at 86400 + 11, unlock product
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).
		WithBlockHeight(11 + feeParams.OrderExpireBlocks)
	k.UnlockProduct(ctx, types.TestTokenPair)
	EndBlocker(ctx, k)

	// check order
	order = k.GetOrder(ctx, orders[0].OrderID)
	require.EqualValues(t, types.OrderStatusExpired, order.Status)
	require.EqualValues(t, 11+feeParams.OrderExpireBlocks, k.GetLastExpiredBlockHeight(ctx))
}

func TestEndBlockerExpireOrders(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 3)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	feeParams := types.DefaultTestParams()

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	tokenPairDex := dex.GetBuiltInTokenPair()
	err = mapp.dexKeeper.SaveTokenPair(ctx, tokenPairDex)
	require.Nil(t, err)
	mapp.dexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	mapp.orderKeeper.SetParams(ctx, &feeParams)
	EndBlocker(ctx, k)

	// mock orders
	orders := []*types.Order{
		types.MockOrder(types.FormatOrderID(startHeight, 1), types.TestTokenPair, types.BuyOrder, "9.8", "1.0"),
		types.MockOrder(types.FormatOrderID(startHeight, 2), types.TestTokenPair, types.SellOrder, "10.0", "1.0"),
		types.MockOrder(types.FormatOrderID(startHeight, 3), types.TestTokenPair, types.BuyOrder, "10.0", "0.5"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	orders[1].Sender = addrKeysSlice[1].Address
	orders[2].Sender = addrKeysSlice[2].Address
	for i := 0; i < 3; i++ {
		err := k.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	EndBlocker(ctx, k)

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 := sdk.SysCoins{
		// 100 - 9.8 - 0.2592 = 89.9408
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("89.9408")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 := sdk.SysCoins{
		// 100 + 10 * 0.5 * (1 - 0.001) - 0.2592 = 104.7408
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("104.7358")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("99")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// check depth book
	depthBook := k.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 2, len(depthBook.Items))

	// call EndBlocker to expire orders
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).
		WithBlockHeight(startHeight + feeParams.OrderExpireBlocks)

	EndBlocker(ctx, k)

	// check order status
	order0 := k.GetOrder(ctx, orders[0].OrderID)
	order1 := k.GetOrder(ctx, orders[1].OrderID)
	require.EqualValues(t, types.OrderStatusExpired, order0.Status)
	require.EqualValues(t, types.OrderStatusPartialFilledExpired, order1.Status)

	// check depth book
	depthBook = k.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 0, len(depthBook.Items))
	// check order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("9.8"), types.BuyOrder)
	orderIDs := k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))
	// check updated order ids
	updatedOrderIDs := k.GetUpdatedOrderIDs()
	require.EqualValues(t, 2, len(updatedOrderIDs))
	require.EqualValues(t, orders[0].OrderID, updatedOrderIDs[0])
	// check closed order id
	closedOrderIDs := k.GetDiskCache().GetClosedOrderIDs()
	require.Equal(t, 2, len(closedOrderIDs))
	require.Equal(t, orders[0].OrderID, closedOrderIDs[0])

	// check account balance
	acc0 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.7408")), // 100 - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 = sdk.SysCoins{
		// 100 + 10 * 0.5 * (1 - 0.001) - 0.2592
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("104.7358")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("99.5")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// check fee pool
	feeCollector := mapp.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	// 0.2592 + 0.2592
	require.EqualValues(t, "0.51840000"+common.NativeToken, collectedFees.String())
}

func TestEndBlockerCleanupOrdersWhoseTokenPairHaveBeenDelisted(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	feeParams := types.DefaultTestParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	// mock orders
	orders := []*types.Order{
		types.MockOrder(types.FormatOrderID(startHeight, 1), types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.MockOrder(types.FormatOrderID(startHeight, 2), types.TestTokenPair, types.SellOrder, "10.0", "0.5"),
		types.MockOrder(types.FormatOrderID(startHeight, 3), types.TestTokenPair, types.SellOrder, "10.0", "2.5"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	orders[1].Sender = addrKeysSlice[1].Address
	orders[2].Sender = addrKeysSlice[1].Address
	for i := 0; i < 3; i++ {
		err := k.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	// call EndBlocker to execute periodic match
	EndBlocker(ctx, k)

	// check depth book
	depthBook := k.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 0, len(depthBook.Items))

	depthBookDB := k.GetDepthBookFromDB(ctx, types.TestTokenPair)
	require.EqualValues(t, 0, len(depthBookDB.Items))

	// check product price - order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.SellOrder)
	orderIDs := k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))

	key = types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.BuyOrder)
	orderIDs = k.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))

	// check closed order id
	closedOrderIDs := k.GetLastClosedOrderIDs(ctx)
	require.Equal(t, 3, len(closedOrderIDs))
	require.Equal(t, orders[0].OrderID, closedOrderIDs[0])
	require.Equal(t, orders[1].OrderID, closedOrderIDs[1])
	require.Equal(t, orders[2].OrderID, closedOrderIDs[2])

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("100")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("100")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// check fee pool
	feeCollector := mapp.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	require.EqualValues(t, "", collectedFees.String())
}

func TestFillPrecision(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	BeginBlocker(ctx, k)

	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock orders
	orderIdx := 0
	roundN := 1 // Need more balance to make a large round
	orderNums := 20
	var orders []*types.Order

	for j := 0; j < roundN; j++ {
		rand.Seed(time.Now().Unix())
		price := float64(25000+rand.Intn(5000)) / 10000

		for i := 0; i < orderNums; i++ {
			var buyPrice string
			var sellPrice string
			var quantity string

			rand.Seed(time.Now().Unix() + int64(orderIdx))

			// Test Same precision of price and quantity
			quantity = strconv.FormatFloat(float64(rand.Intn(99999))/100000, 'f', 4, 64)
			buyPrice = strconv.FormatFloat(price+0.0001, 'f', 4, 64)
			sellPrice = strconv.FormatFloat(price, 'f', 4, 64)

			tmp, err := strconv.ParseFloat(quantity, 64)
			if tmp == 0.0 {
				continue
			}

			orderIdx += 1
			buyOrder := types.MockOrder(types.FormatOrderID(startHeight, int64(orderIdx)), types.TestTokenPair, types.BuyOrder, buyPrice, quantity)
			orderIdx += 1
			sellOrder := types.MockOrder(types.FormatOrderID(startHeight, int64(orderIdx)), types.TestTokenPair, types.SellOrder, sellPrice, quantity)

			buyOrder.Sender = addrKeysSlice[0].Address
			sellOrder.Sender = addrKeysSlice[1].Address

			orders = append(orders, buyOrder, sellOrder)
			err = k.PlaceOrder(ctx, buyOrder)
			require.NoError(t, err)

			orders = append(orders, sellOrder)
			err = k.PlaceOrder(ctx, sellOrder)
		}
	}
	// call EndBlocker to execute periodic match
	EndBlocker(ctx, k)

	N := len(orders) / 1000
	for i := 0; i < N; i++ {
		ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight + int64(1+i))
		BeginBlocker(ctx, k)
		EndBlocker(ctx, k)
	}

	invaFunc := keeper.ModuleAccountInvariant(mapp.orderKeeper)
	_, isInval := invaFunc(ctx)
	require.EqualValues(t, false, isInval)
}

func buildRandomOrderMsg(addr sdk.AccAddress) MsgNewOrders {
	price := strconv.Itoa(rand.Intn(10) + 100)
	orderItems := []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, price, "1.0"),
	}
	msg := types.NewMsgNewOrders(addr, orderItems)
	return msg

}

func TestEndBlocker(t *testing.T) {
	mapp, addrKeysSlice := getMockAppWithBalance(t, 2, 100000000)
	k := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	feeParams := types.DefaultTestParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	handler := NewOrderHandler(k)

	blockHeight := startHeight
	for i := 0; i < 100000; i++ {
		msg := buildRandomOrderMsg(addrKeysSlice[0].Address)
		result, err := handler(ctx, msg)
		if (i+1)%1000 == 0 {
			blockHeight = blockHeight + 1
			ctx = ctx.WithBlockHeight(blockHeight)
		}
		require.Nil(t, err)
		require.EqualValues(t, "", result.Log)
	}
	// call EndBlocker to execute periodic match
	EndBlocker(ctx, k)

	quantityList := [3]string{"200", "500", "1000"}
	for _, quantity := range quantityList {
		startTime := time.Now()
		blockHeight = blockHeight + 1
		ctx = ctx.WithBlockHeight(blockHeight)
		orderItems := []types.OrderItem{
			types.NewOrderItem(types.TestTokenPair, types.SellOrder, "100", quantity),
		}
		msg := types.NewMsgNewOrders(addrKeysSlice[1].Address, orderItems)
		handler(ctx, msg)
		EndBlocker(ctx, k)
		fmt.Println(time.Since(startTime))
		fmt.Println(k.GetOrder(ctx, types.FormatOrderID(blockHeight, 1)))
	}
}
