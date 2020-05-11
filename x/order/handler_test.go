package order

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/supply"
	commonType "github.com/okex/okchain/x/common/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order/types"
	"github.com/okex/okchain/x/token"
	tokentypes "github.com/okex/okchain/x/token/types"
)

func TestEventNewOrders(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	handler := NewOrderHandler(keeper)
	//test multi order fee is 80%
	orderItems := []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.SellOrder, "10.0", "1.0"),
		types.NewOrderItem(types.TestTokenPair+"A", types.BuyOrder, "10.0", "1.0"),
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
	}

	mapp.orderKeeper.SetParams(ctx, &feeParams)
	msg := types.NewMsgNewOrders(addrKeysSlice[0].Address, orderItems)
	result := handler(ctx, msg)

	require.EqualValues(t, 2, len(result.Events[4].Attributes))

}

func TestFeesNewOrders(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	handler := NewOrderHandler(keeper)
	//test multi order fee is 80%
	orderItems := []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair+"a", types.BuyOrder, "10.0", "1.0"),
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
	}
	acc := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("100")), // 100
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())

	mapp.orderKeeper.SetParams(ctx, &feeParams)
	msg := types.NewMsgNewOrders(addrKeysSlice[0].Address, orderItems)
	result := handler(ctx, msg)

	// check account balance
	// multi fee 7958528000
	acc = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("89.79264")), // 100 - 10  - 0.20736
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	require.EqualValues(t, true, result.Code.IsOK())

}

func TestHandleMsgNewOrderInvalid(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	handler := NewOrderHandler(mapp.orderKeeper)

	// not-exist product
	msg := types.NewMsgNewOrder(addrKeysSlice[0].Address, "nobb_"+common.NativeToken, types.BuyOrder, "10.0", "1.0")
	result := handler(ctx, msg)
	orderRes := parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, sdk.CodeUnknownRequest, orderRes[0].Code)

	// invalid price precision
	//msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.01", "1.0")
	//result = handler(ctx, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid quantity precision
	//msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "1.001")
	//result = handler(ctx, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid quantity amount
	//msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "0.09")
	//result = handler(ctx, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// insufficient coins
	msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "10.1")
	result = handler(ctx, msg)
	orderRes = parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, sdk.CodeInsufficientCoins, orderRes[0].Code)

	// check depth book
	depthBook := mapp.orderKeeper.GetDepthBookCopy(types.TestTokenPair)
	require.Equal(t, 0, len(depthBook.Items))
}

func TestValidateMsgNewOrder(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	keeper := mapp.orderKeeper
	feeParams := types.DefaultParams()
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// normal
	msg := types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	result := ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeOK, result.Code)

	// not-exist product
	msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, "nobb_"+common.NativeToken, types.BuyOrder, "10.0", "1.0")
	result = ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid price precision
	//msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.01", "1.0")
	//result = ValidateMsgNewOrder(ctx, keeper, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid quantity precision
	//msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "1.001")
	//result = ValidateMsgNewOrder(ctx, keeper, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid quantity amount
	//msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "0.09")
	//result = ValidateMsgNewOrder(ctx, keeper, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// insufficient coins
	msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "10.1")
	result = ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeInsufficientCoins, result.Code)

	// busy product
	keeper.SetProductLock(ctx, types.TestTokenPair, &types.ProductLock{})
	msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	result = ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeInternal, result.Code)

	// price * quantity over accuracy
	keeper.SetProductLock(ctx, types.TestTokenPair, &types.ProductLock{})
	msg = types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.000001", "1.0001")
	result = ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)
}

// test order cancel without enough okb as fee
func TestHandleMsgCancelOrder2(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	feeParams := types.DefaultParams()
	//feeParams.CancelNative = sdk.MustNewDecFromStr("0.1")
	mapp.orderKeeper.SetParams(ctx, &feeParams)
	tokenPair := dex.GetBuiltInTokenPair()
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// subtract all okb of addr0
	err = keeper.LockCoins(ctx, addrKeysSlice[0].Address, sdk.DecCoins{{Denom: common.NativeToken,
		Amount: sdk.MustNewDecFromStr("99.7408")}}, tokentypes.LockCoinsTypeQuantity)
	require.NoError(t, err)

	// mock orders
	orders := []*types.Order{
		types.MockOrder(types.FormatOrderID(startHeight, 1), types.TestTokenPair, types.SellOrder, "10.0", "2.0"),
	}
	orders[0].Sender = addrKeysSlice[0].Address
	err = keeper.PlaceOrder(ctx, orders[0])
	require.NoError(t, err)

	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight + 1)

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins0 := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("98")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())

	// Start Testing...
	handler := NewOrderHandler(keeper)
	keeper.ResetCache(ctx)

	// Test fully cancel
	msg := types.NewMsgCancelOrder(addrKeysSlice[0].Address, orders[0].OrderID)
	result := handler(ctx, msg)
	// check result
	require.EqualValues(t, sdk.CodeOK, result.Code)
	orderRes := parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, "0.00000100okt", orderRes[0].Message)
	// check account balance
	acc0 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins0 = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("0.25919900")), // no change
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),          // 100 - 0.000001
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	// check fee pool
	feeCollector := mapp.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	require.EqualValues(t, "0.00000100"+common.NativeToken, collectedFees.String())
}

func TestHandleMsgCancelOrderInvalid(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)
	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	// mock orders
	order := types.MockOrder(types.FormatOrderID(startHeight, 1), types.TestTokenPair, types.SellOrder, "10.0", "1.0")
	order.Sender = addrKeysSlice[0].Address
	err = keeper.PlaceOrder(ctx, order)
	require.Nil(t, err)

	EndBlocker(ctx, keeper) // update depthBook, orderIdsMap

	handler := NewOrderHandler(keeper)

	// invalid owner
	msg := types.NewMsgCancelOrder(addrKeysSlice[1].Address, order.OrderID)
	result := handler(ctx, msg)
	orderRes := parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, sdk.CodeUnauthorized, orderRes[0].Code)
	// invalid orderID
	msg = types.NewMsgCancelOrder(addrKeysSlice[1].Address, "InvalidID-0001")
	result = handler(ctx, msg)
	orderRes = parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, sdk.CodeUnknownRequest, orderRes[0].Code)
	// busy product
	keeper.SetProductLock(ctx, order.Product, &types.ProductLock{})
	msg = types.NewMsgCancelOrder(addrKeysSlice[0].Address, order.OrderID)
	result = handler(ctx, msg)
	orderRes = parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, sdk.CodeInternal, orderRes[0].Code)
	keeper.UnlockProduct(ctx, order.Product)

	// normal
	msg = types.NewMsgCancelOrder(addrKeysSlice[0].Address, order.OrderID)
	result = handler(ctx, msg)

	// check result
	require.EqualValues(t, sdk.CodeOK, result.Code)
	orderRes = parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, "0.00000000"+common.NativeToken, orderRes[0].Message)
	// check order status
	order = keeper.GetOrder(ctx, order.OrderID)
	require.EqualValues(t, types.OrderStatusCancelled, order.Status)
	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins0 := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("100")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())

	// invalid order status
	msg = types.NewMsgCancelOrder(addrKeysSlice[0].Address, order.OrderID)
	result = handler(ctx, msg)
	orderRes = parseOrderResult(result)
	require.NotNil(t, orderRes)
	require.EqualValues(t, sdk.CodeInternal, orderRes[0].Code)
}

func TestHandleInvalidMsg(t *testing.T) {
	mapp, _ := getMockApp(t, 0)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)

	handler := NewOrderHandler(keeper)
	var msg token.MsgSend
	result := handler(ctx, msg)
	require.EqualValues(t, commonType.CodeUnknownMsgType, result.Code)
}

const orderKey = "orders"

func getOrderID(result sdk.Result) string {
	var res = ""
	var evs []types.OrderResult
	for i := 0; i < len(result.Events); i++ {
		event := result.Events[i]
		for j := 0; j < len(event.Attributes); j++ {
			attribute := event.Attributes[j]
			if string(attribute.Key) == orderKey {
				res = string(attribute.Value)
				if err := json.Unmarshal([]byte(res), &evs); err == nil {
					for k := 0; k < len(evs); k++ {
						res = evs[k].OrderID
					}
				}

			}

		}
	}
	return res
}

func getOrderIDList(result sdk.Result) []string {
	var res []string
	for i := 0; i < len(result.Events); i++ {
		event := result.Events[i]
		var evs []types.OrderResult
		for j := 0; j < len(event.Attributes); j++ {
			attribute := event.Attributes[j]
			if string(attribute.Key) == orderKey {
				value := string(attribute.Value)
				if err := json.Unmarshal([]byte(value), &evs); err == nil {
					for k := 0; k < len(evs); k++ {
						res = append(res, evs[k].OrderID)
					}
				}
			}

		}
	}
	return res
}

func parseOrderResult(result sdk.Result) []types.OrderResult {
	var evs []types.OrderResult
	for i := 0; i < len(result.Events); i++ {
		event := result.Events[i]

		for j := 0; j < len(event.Attributes); j++ {
			attribute := event.Attributes[j]
			if string(attribute.Key) == orderKey {
				value := string(attribute.Value)

				if err := json.Unmarshal([]byte(value), &evs); err != nil {
					return nil
					//ignore
				}
			}
		}
	}
	return evs
}

func TestHandleMsgMultiNewOrder(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	handler := NewOrderHandler(keeper)

	// Test buy order
	orderItems := []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
	}
	msg := types.NewMsgNewOrders(addrKeysSlice[0].Address, orderItems)
	result := handler(ctx, msg)
	require.Equal(t, "", result.Log)
	// Test order when locked
	keeper.SetProductLock(ctx, types.TestTokenPair, &types.ProductLock{})
	result1 := handler(ctx, msg)
	res1 := parseOrderResult(result1)
	require.EqualValues(t, sdk.CodeInternal, res1[0].Code)
	keeper.UnlockProduct(ctx, types.TestTokenPair)

	//check result & order
	orderID := getOrderID(result)
	require.EqualValues(t, types.FormatOrderID(10, 2), orderID)
	order := keeper.GetOrder(ctx, orderID)
	require.NotNil(t, order)
	require.EqualValues(t, 2, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("79.58528")), // 100 - 10 - 10 - 0.2592 * 2 * 0.8
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	// check depth book
	depthBook := keeper.GetDepthBookCopy(order.Product)
	require.Equal(t, 1, len(depthBook.Items))
	require.Equal(t, sdk.MustNewDecFromStr("10.0"), depthBook.Items[0].Price)
	require.Equal(t, sdk.MustNewDecFromStr("2.0"), depthBook.Items[0].BuyQuantity)
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids map
	orderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()
	require.Equal(t, 1, len(orderIDsMap.Data))
	require.Equal(t, types.FormatOrderID(10, 1),
		orderIDsMap.Data[types.FormatOrderIDsKey(order.Product, order.Price, order.Side)][0])
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))

	// Test sell order
	orderItems = []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.SellOrder, "10.0", "1.0"),
	}
	msg = types.NewMsgNewOrders(addrKeysSlice[0].Address, orderItems)
	result = handler(ctx, msg)

	// check result & order
	orderID = getOrderID(result)
	require.EqualValues(t, types.FormatOrderID(10, 3), orderID)
	order = keeper.GetOrder(ctx, orderID)
	require.NotNil(t, order)
	require.EqualValues(t, 3, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("79.32608")), // 100 - 10 - 10 - 0.2592 * 2 * 0.8 - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("99")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())

	// test new order with fee
	feeParams.FeePerBlock = sdk.NewDecCoinFromDec(types.DefaultFeeDenomPerBlock, sdk.MustNewDecFromStr("0.000002"))
	mapp.orderKeeper.SetParams(ctx, &feeParams)
	orderItems = []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.SellOrder, "10.0", "1.0"),
	}
	msg = types.NewMsgNewOrders(addrKeysSlice[0].Address, orderItems)
	result = handler(ctx, msg)

	orderID = getOrderID(result)
	require.EqualValues(t, types.FormatOrderID(10, 4), orderID)
	// check account balance
	acc = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("78.80768")), // 79.32608 - 0.2592 * 2
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("98")),         // 99 - 1
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())

	feeParams = types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	require.EqualValues(t, 4, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("78.80768")), // 78.80768
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("98")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
}

func TestHandleMsgMultiCancelOrder(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	handler := NewOrderHandler(keeper)

	// Test buy order
	orderItems := []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
	}
	msg := types.NewMsgNewOrders(addrKeysSlice[0].Address, orderItems)
	result := handler(ctx, msg)
	require.Equal(t, "", result.Log)
	// Test order when locked
	keeper.SetProductLock(ctx, types.TestTokenPair, &types.ProductLock{})

	result1 := handler(ctx, msg)

	require.Equal(t, "", result1.Log)
	keeper.UnlockProduct(ctx, types.TestTokenPair)

	// check result & order

	orderID := getOrderID(result)
	require.EqualValues(t, types.FormatOrderID(10, 3), orderID)
	order := keeper.GetOrder(ctx, orderID)
	require.NotNil(t, order)
	require.EqualValues(t, 3, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("69.37792")), // 100 - 10*6 - 0.2592 * 6 * 0.8
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())
	// check depth book
	depthBook := keeper.GetDepthBookCopy(order.Product)
	require.Equal(t, 1, len(depthBook.Items))
	require.Equal(t, sdk.MustNewDecFromStr("10.0"), depthBook.Items[0].Price)
	require.Equal(t, sdk.MustNewDecFromStr("3.0"), depthBook.Items[0].BuyQuantity)
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids map
	orderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()
	require.Equal(t, 1, len(orderIDsMap.Data))
	require.Equal(t, types.FormatOrderID(10, 1),
		orderIDsMap.Data[types.FormatOrderIDsKey(order.Product, order.Price, order.Side)][0])
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))

	// Test cancel order
	orderIDItems := getOrderIDList(result)
	multiCancelMsg := types.NewMsgCancelOrders(addrKeysSlice[0].Address, orderIDItems[:len(orderItems)-1])
	result = handler(ctx, multiCancelMsg)

	require.Equal(t, true, result.Code.IsOK())
	// check result & order

	require.EqualValues(t, 3, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("89.79264")), // 100 - 10 - 10 - 0.2592 * 2 * 0.8 - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())

	// Test cancel order
	orderIDItems = orderIDItems[2:]
	orderIDItems = append(orderIDItems, "")

	multiCancelMsg = types.NewMsgCancelOrders(addrKeysSlice[0].Address, orderIDItems)
	result = handler(ctx, multiCancelMsg)

	require.Equal(t, true, result.Code.IsOK())
	require.Equal(t, "", result.Log)
	// check result & order

	require.EqualValues(t, 3, keeper.GetBlockOrderNum(ctx, 10))
	// check account balance
	acc = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("100")), // 100 - 10 - 10 - 0.2592 * 2 * 0.8 - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())

}

func TestValidateMsgMultiNewOrder(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	keeper := mapp.orderKeeper
	feeParams := types.DefaultParams()
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	orderItems := []types.OrderItem{
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "0.1", "1.0"),
		types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "0.1", "1.0"),
	}

	// normal
	orderItem := types.NewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	msg := types.NewMsgNewOrders(addrKeysSlice[0].Address, append(orderItems, orderItem))
	result := ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeOK, result.Code)

	// not-exist product
	orderItem = types.NewOrderItem("nobb_"+common.NativeToken, types.BuyOrder, "10.0", "1.0")
	msg = types.NewMsgNewOrders(addrKeysSlice[0].Address, append(orderItems, orderItem))
	result = ValidateMsgNewOrders(ctx, keeper, msg)
	require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid price precision
	//orderItem = types.NewMultiNewOrderItem(types.TestTokenPair, types.BuyOrder, "10.01", "1.0")
	//msg = types.NewMsgMultiNewOrder(addrKeysSlice[0].Address, append(orderItems, orderItem))
	//result = ValidateMsgMultiNewOrder(ctx, keeper, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid quantity precision
	//orderItem = types.NewMultiNewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "1.001")
	//msg = types.NewMsgMultiNewOrder(addrKeysSlice[0].Address, append(orderItems, orderItem))
	//result = ValidateMsgMultiNewOrder(ctx, keeper, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	// invalid quantity amount
	//orderItem = types.NewMultiNewOrderItem(types.TestTokenPair, types.BuyOrder, "10.0", "0.09")
	//msg = types.NewMsgMultiNewOrder(addrKeysSlice[0].Address, append(orderItems, orderItem))
	//result = ValidateMsgMultiNewOrder(ctx, keeper, msg)
	//require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)
}

func TestValidateMsgMultiCancelOrder(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	keeper := mapp.orderKeeper
	feeParams := types.DefaultParams()
	keeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	orderIDItems := []string{""}
	multiCancelMsg := types.NewMsgCancelOrders(addrKeysSlice[0].Address, orderIDItems)
	result := ValidateMsgCancelOrders(ctx, keeper, multiCancelMsg)
	require.EqualValues(t, sdk.CodeUnknownRequest, result.Code)

	err = mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	handler := NewOrderHandler(keeper)

	// new order
	msg := types.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	result = handler(ctx, msg)

	// validate true
	orderID := getOrderID(result)
	orderIDItems = []string{orderID}
	multiCancelMsg = types.NewMsgCancelOrders(addrKeysSlice[0].Address, orderIDItems)
	result = ValidateMsgCancelOrders(ctx, keeper, multiCancelMsg)
	require.EqualValues(t, sdk.CodeOK, result.Code)

	// validate empty orderIDItems
	orderIDItems = []string{}
	multiCancelMsg = types.NewMsgCancelOrders(addrKeysSlice[0].Address, orderIDItems)
	result = ValidateMsgCancelOrders(ctx, keeper, multiCancelMsg)
	require.EqualValues(t, sdk.CodeOK, result.Code)

}

func TestHandleMsgCancelOrder(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 3)
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)
	tokenPair := dex.GetBuiltInTokenPair()
	err := mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	tokenPairDex := dex.GetBuiltInTokenPair()
	err = mapp.dexKeeper.SaveTokenPair(ctx, tokenPairDex)
	require.Nil(t, err)

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
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	EndBlocker(ctx, keeper) // update blockMatchResult, updatedOrderIds, depthBook, orderIdsMap

	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight + 1)

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins0 := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("89.9408")), // 100 - 9.8 - 0.2592
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 := sdk.DecCoins{
		// 100 + 10 * 0.5 * (1 - 0.001) - 0.2592
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("104.7358")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("99")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// check depth book
	depthBook := keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 2, len(depthBook.Items))
	// check order ids
	key := types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("9.8"), types.BuyOrder)
	orderIDs := keeper.GetProductPriceOrderIDs(key)
	require.EqualValues(t, orders[0].OrderID, orderIDs[0])
	key = types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10.0"), types.SellOrder)
	orderIDs = keeper.GetProductPriceOrderIDs(key)
	require.EqualValues(t, orders[1].OrderID, orderIDs[0])

	// Start Testing...
	keeper.ResetCache(ctx)
	handler := NewOrderHandler(keeper)

	// Test fully cancel
	msg := types.NewMsgCancelOrder(addrKeysSlice[0].Address, orders[0].OrderID)
	result := handler(ctx, msg)
	for i := 0; i < len(result.Events); i++ {
		fmt.Println(i)
		for j := 0; j < len(result.Events[i].Attributes); j++ {
			arr := result.Events[i].Attributes[j]
			fmt.Println(string(arr.Key), string(arr.Value))
		}
	}

	orderRes := parseOrderResult(result)
	// check result
	require.EqualValues(t, sdk.CodeOK, result.Code)
	require.EqualValues(t, "0.00000100"+common.NativeToken, orderRes[0].Message)
	// check order status
	orders[0] = keeper.GetOrder(ctx, orders[0].OrderID)
	require.EqualValues(t, types.OrderStatusCancelled, orders[0].Status)
	// check account balance
	acc0 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins0 = sdk.DecCoins{
		// 100 - 0.002
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.999999")), // 100 - 9.8 + 9.8 - 0.000001
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	// check fee pool
	feeCollector := mapp.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees := feeCollector.GetCoins()
	require.EqualValues(t, "0.00000100"+common.NativeToken, collectedFees.String()) // 0.002+0.002
	// check depth book
	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBook.Items))
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedDepthbookKeys()))
	// check order ids
	key = types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("9.8"), types.BuyOrder)
	orderIDs = keeper.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))
	require.Equal(t, 1, len(keeper.GetDiskCache().GetUpdatedOrderIDKeys()))
	// check updated order ids
	updatedOrderIDs := keeper.GetUpdatedOrderIDs()
	require.EqualValues(t, orders[0].OrderID, updatedOrderIDs[0])
	// check closed order id
	closedOrderIDs := keeper.GetDiskCache().GetClosedOrderIDs()
	require.Equal(t, 1, len(closedOrderIDs))
	require.Equal(t, orders[0].OrderID, closedOrderIDs[0])

	// Test partially cancel
	msg = types.NewMsgCancelOrder(addrKeysSlice[1].Address, orders[1].OrderID)
	result = handler(ctx, msg)
	// check result
	require.EqualValues(t, sdk.CodeOK, result.Code)
	// check order status
	orders[1] = keeper.GetOrder(ctx, orders[1].OrderID)
	require.EqualValues(t, types.OrderStatusPartialFilledCancelled, orders[1].Status)
	// check account balance
	acc1 = mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	expectCoins1 = sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("104.994999")), // 99.999999 + 5 * (1 - 0.001)
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("99.5")),
	}
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())
	// check fee pool, partially cancel, no fees
	feeCollector = mapp.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	collectedFees = feeCollector.GetCoins()
	require.EqualValues(t, "0.00000200"+common.NativeToken, collectedFees.String())
	// check order ids
	key = types.FormatOrderIDsKey(types.TestTokenPair, sdk.MustNewDecFromStr("10"), types.SellOrder)
	orderIDs = keeper.GetProductPriceOrderIDs(key)
	require.EqualValues(t, 0, len(orderIDs))
}

func TestFeesTable(t *testing.T) {
	//test xxb_okt
	orders0 := []*types.Order{
		types.MockOrder(types.FormatOrderID(10, 1), types.TestTokenPair, types.BuyOrder, "10", "1.0"),
		types.MockOrder(types.FormatOrderID(10, 2), types.TestTokenPair, types.BuyOrder, "10", "2.0"),
		types.MockOrder(types.FormatOrderID(10, 1), types.TestTokenPair, types.SellOrder, "10.0", "1"),
		types.MockOrder(types.FormatOrderID(10, 2), types.TestTokenPair, types.SellOrder, "10.0", "2"),
	}
	expectCoins0 := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("169.98")), // 200 - 10 -20 - 0.2592*10000/259200*2
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("202.997")),  // 200 + (3 - 3*0.001)
	}

	//test btc-b19_okt
	orders1 := []*types.Order{
		types.MockOrder(types.FormatOrderID(10, 1), "btc-b19_okt", types.BuyOrder, "10", "1"),
		types.MockOrder(types.FormatOrderID(10, 2), "btc-b19_okt", types.SellOrder, "10", "1"),
	}
	expectCoins1 := sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc-b19", sdk.MustNewDecFromStr("100.999")),         //100 + (1 - 1*0.0001)
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("189.99")), // 200 - 10 - 0.2592*10000/259200
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}

	//test btc-b19_xxb
	orders2 := []*types.Order{
		types.MockOrder(types.FormatOrderID(10, 1), "btc-b19_xxb", types.BuyOrder, "11", "1"),
		types.MockOrder(types.FormatOrderID(10, 2), "btc-b19_xxb", types.SellOrder, "11", "1"),
	}
	expectCoins2 := sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc-b19", sdk.MustNewDecFromStr("100.999")),        //100 + (1 - 1*0.0001)
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.99")), // 100 - 0.2592*10000/259200
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("189")),     //200 - 11
	}

	//test btc-b19_xxb match order on 800 block
	expectCoins3 := sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc-b19", sdk.MustNewDecFromStr("100.999")),          //100 + (1 - 1*0.0001)
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.9992")), // 100 - 0.2592*800/259200
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("189")),       //200 - 11
	}

	//test btc-a8a_xxb 2 match orders
	orders4 := []*types.Order{
		types.MockOrder(types.FormatOrderID(10, 1), "btc-a8a_xxb", types.BuyOrder, "11", "1"),
		types.MockOrder(types.FormatOrderID(10, 2), "btc-a8a_xxb", types.BuyOrder, "11", "2"),
		types.MockOrder(types.FormatOrderID(10010, 1), "btc-a8a_xxb", types.SellOrder, "11", "1"),
		types.MockOrder(types.FormatOrderID(10010, 2), "btc-a8a_xxb", types.SellOrder, "11", "2"),
	}
	expectCoins4 := sdk.DecCoins{
		sdk.NewDecCoinFromDec("btc-a8a", sdk.MustNewDecFromStr("102.997")),        //100 +(2 - 2 * 0.001) + (1 - 1*0.0001)
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.98")), // 100 - 0.2592*10000/259200*2
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("167")),     //200 - 11 - 11*2
	}

	tests := []struct {
		baseasset   string
		quoteasset  string
		orders      []*types.Order
		balance     sdk.DecCoins
		blockheight int64
	}{
		{common.TestToken, common.NativeToken, orders0, expectCoins0, 10000},
		{"btc-b19", "okt", orders1, expectCoins1, 10000},
		{"btc-b19", "xxb", orders2, expectCoins2, 10000},
		{"btc-b19", "xxb", orders2, expectCoins3, 800},
		{"btc-a8a", "xxb", orders4, expectCoins4, 10000},
	}

	for i, tc := range tests {
		expectCoins := handleOrders(t, tc.baseasset, tc.quoteasset, tc.orders, tc.blockheight)
		require.EqualValues(t, tc.balance.String(), expectCoins.String(), "test: %v", i)
	}
}

func handleOrders(t *testing.T, baseasset string, quoteasset string, orders []*types.Order, blockheight int64) sdk.DecCoins {
	TestTokenPairOwner := "okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya"
	addr, err := sdk.AccAddressFromBech32(TestTokenPairOwner)
	require.Nil(t, err)
	mapp, addrKeysSlice := getMockApp(t, len(orders))
	keeper := mapp.orderKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	var startHeight int64 = 10
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	feeParams := types.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	//init balance account0 & account1
	decCoins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s", 100, baseasset, 100, quoteasset))
	require.Nil(t, err)
	_, err = mapp.bankKeeper.AddCoins(ctx, addrKeysSlice[0].Address, decCoins)
	require.Nil(t, err)
	_, err = mapp.bankKeeper.AddCoins(ctx, addrKeysSlice[1].Address, decCoins)
	require.Nil(t, err)
	//init token pair
	tokenPair := dex.TokenPair{
		BaseAssetSymbol:  baseasset,
		QuoteAssetSymbol: quoteasset,
		InitPrice:        sdk.MustNewDecFromStr("10.0"),
		MaxPriceDigit:    8,
		MaxQuantityDigit: 8,
		MinQuantity:      sdk.MustNewDecFromStr("0"),
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
	}

	err = mapp.dexKeeper.SaveTokenPair(ctx, &tokenPair)
	require.Nil(t, err)
	acc := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	require.NotNil(t, acc)
	//place buy order
	for i := 0; i < len(orders)/2; i++ {
		orders[i].Sender = addrKeysSlice[0].Address
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	EndBlocker(ctx, keeper)
	//update blockheight
	ctx = mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(startHeight + blockheight)
	//place sell order
	for i := len(orders) / 2; i < len(orders); i++ {
		orders[i].Sender = addrKeysSlice[1].Address
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	EndBlocker(ctx, keeper)

	// check account balance
	acc0 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	acc1 := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[1].Address)
	require.NotNil(t, acc1)
	return acc0.GetCoins()
}
