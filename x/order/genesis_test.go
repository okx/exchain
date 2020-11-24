package order

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/order/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types/time"
)

func TestValidateGenesis(t *testing.T) {
	genesisState := DefaultGenesisState()
	err := ValidateGenesis(genesisState)
	require.NoError(t, err)
}

func TestExportGenesis(t *testing.T) {
	testInput := keeper.CreateTestInput(t)
	ctx := testInput.Ctx
	orderKeeper := testInput.OrderKeeper

	params := types.DefaultParams()
	params.OrderExpireBlocks = 1234
	params.MaxDealsPerBlock = 2345

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.NoError(t, err)
	product := fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)

	var orders []*types.Order
	order1 := types.NewOrder("txHash",
		testInput.TestAddrs[0],
		product,
		types.SellOrder,
		sdk.NewDec(123),
		sdk.NewDec(456),
		time.Now().Unix(),
		5,
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1)))
	order1.OrderID = types.FormatOrderID(10, 1)
	order1.FilledAvgPrice = sdk.ZeroDec()

	orders = append(orders, order1)

	initGenesis := GenesisState{
		Params:     params,
		OpenOrders: orders,
	}

	InitGenesis(ctx, orderKeeper, initGenesis)
	require.Equal(t, params, *orderKeeper.GetParams(ctx))
	// 0x11
	require.Equal(t, order1, orderKeeper.GetOrder(ctx, order1.OrderID))
	// 0x12
	depthBook := &types.DepthBook{}
	depthBook.InsertOrder(initGenesis.OpenOrders[0])
	require.Equal(t, depthBook, orderKeeper.GetDepthBookFromDB(ctx, product))
	// 0x13
	key1 := types.FormatOrderIDsKey(product, initGenesis.OpenOrders[0].Price, initGenesis.OpenOrders[0].Side)
	require.Equal(t, []string{initGenesis.OpenOrders[0].OrderID}, orderKeeper.GetProductPriceOrderIDsFromDB(ctx, key1))
	// 0x14
	require.Equal(t, tokenPair.InitPrice, orderKeeper.GetLastPrice(ctx, product))
	// 0x15
	require.Equal(t, []int64{10}, orderKeeper.GetExpireBlockHeight(ctx, initGenesis.Params.OrderExpireBlocks+10))
	// 0x16
	require.Equal(t, int64(1), orderKeeper.GetBlockOrderNum(ctx, 10))
	// 0x17
	require.Equal(t, 0, len(orderKeeper.GetLastClosedOrderIDs(ctx)))
	// 0x18
	require.Equal(t, int64(0), orderKeeper.GetLastExpiredBlockHeight(ctx))
	// 0x19
	require.Equal(t, int64(1), orderKeeper.GetOpenOrderNum(ctx))
	// 0x20
	require.Equal(t, int64(1), orderKeeper.GetStoreOrderNum(ctx))

	exportGenesis := ExportGenesis(ctx, orderKeeper)
	require.Equal(t, params, exportGenesis.Params)
	require.Equal(t, order1, exportGenesis.OpenOrders[0])

	params.MaxDealsPerBlock = 1
	params.FeePerBlock = sdk.NewDecCoinFromDec(common.NativeToken, sdk.OneDec())
	params.TradeFeeRate = sdk.NewDecWithPrec(5, 2)
	params.OrderExpireBlocks = 3333
	orderKeeper.SetParams(ctx, &params)

	order2 := types.NewOrder("txHash",
		testInput.TestAddrs[0],
		product,
		types.BuyOrder,
		sdk.NewDec(2),
		sdk.NewDec(5),
		time.Now().Unix(),
		5,
		sdk.NewDecCoinFromDec(types.DefaultFeeDenomPerBlock, sdk.NewDec(1)))
	order2.FilledAvgPrice = sdk.ZeroDec()
	ctx = ctx.WithBlockHeight(1000)
	err = orderKeeper.PlaceOrder(ctx, order2)
	require.NoError(t, err)
	orderKeeper.Cache2Disk(ctx)

	exportGenesis = ExportGenesis(ctx, orderKeeper)
	require.Equal(t, params, exportGenesis.Params)
	require.Equal(t, order1, orderKeeper.GetOrder(ctx, order1.OrderID))
	require.Equal(t, order2, orderKeeper.GetOrder(ctx, order2.OrderID))
	require.Equal(t, orderKeeper.GetOrder(ctx, exportGenesis.OpenOrders[0].OrderID), exportGenesis.OpenOrders[0])
	require.Equal(t, orderKeeper.GetOrder(ctx, exportGenesis.OpenOrders[1].OrderID), exportGenesis.OpenOrders[1])

	newTestInput := keeper.CreateTestInput(t)
	newCtx := newTestInput.Ctx
	newOrderKeeper := newTestInput.OrderKeeper
	err = newTestInput.DexKeeper.SaveTokenPair(newCtx, tokenPair)
	require.NoError(t, err)
	InitGenesis(newCtx, newOrderKeeper, exportGenesis)
	require.Equal(t, exportGenesis.Params, *newOrderKeeper.GetParams(newCtx))
	// 0x11
	require.Equal(t, order1, newOrderKeeper.GetOrder(newCtx, order1.OrderID))
	require.Equal(t, order2, newOrderKeeper.GetOrder(newCtx, order2.OrderID))
	// 0x12
	depthBook = &types.DepthBook{}
	depthBook.InsertOrder(exportGenesis.OpenOrders[0])
	depthBook.InsertOrder(exportGenesis.OpenOrders[1])
	require.Equal(t, depthBook, newOrderKeeper.GetDepthBookFromDB(newCtx, product))
	// 0x13
	key1 = types.FormatOrderIDsKey(product, exportGenesis.OpenOrders[0].Price, exportGenesis.OpenOrders[0].Side)
	key2 := types.FormatOrderIDsKey(product, exportGenesis.OpenOrders[1].Price, exportGenesis.OpenOrders[1].Side)
	require.Equal(t, []string{exportGenesis.OpenOrders[0].OrderID}, newOrderKeeper.GetProductPriceOrderIDsFromDB(newCtx, key1))
	require.Equal(t, []string{exportGenesis.OpenOrders[1].OrderID}, newOrderKeeper.GetProductPriceOrderIDsFromDB(newCtx, key2))
	// 0x14
	require.Equal(t, tokenPair.InitPrice, newOrderKeeper.GetLastPrice(newCtx, product))
	// 0x15
	require.Equal(t, []int64{10}, newOrderKeeper.GetExpireBlockHeight(newCtx, exportGenesis.Params.OrderExpireBlocks+10))
	require.Equal(t, []int64{1000}, newOrderKeeper.GetExpireBlockHeight(newCtx, exportGenesis.Params.OrderExpireBlocks+1000))
	// 0x16
	require.Equal(t, int64(1), newOrderKeeper.GetBlockOrderNum(newCtx, 10))
	require.Equal(t, int64(1), newOrderKeeper.GetBlockOrderNum(newCtx, 1000))
	// 0x17
	require.Equal(t, 0, len(newOrderKeeper.GetLastClosedOrderIDs(newCtx)))
	// 0x18
	require.Equal(t, int64(0), newOrderKeeper.GetLastExpiredBlockHeight(newCtx))
	// 0x19
	require.Equal(t, int64(2), newOrderKeeper.GetOpenOrderNum(newCtx))
	// 0x20
	require.Equal(t, int64(2), newOrderKeeper.GetStoreOrderNum(newCtx))
}
