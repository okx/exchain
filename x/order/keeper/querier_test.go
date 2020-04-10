package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order/types"
)

func TestQueryOrder(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	orderID1 := "abc"
	order1 := mockOrder(orderID1, types.TestTokenPair, types.BuyOrder, "0.5", "1.1")
	keeper.SetOrder(ctx, orderID1, order1)

	path := []string{types.QueryOrderDetail, orderID1}
	BytesOrder, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	order2 := &types.Order{}
	err2 := keeper.cdc.UnmarshalJSON(BytesOrder, order2)
	require.Nil(t, err2)
	require.EqualValues(t, order1.String(), order2.String())

	// Test query not-existed order
	path = []string{types.QueryOrderDetail, "Non-existedID"}
	_, err = querier(ctx, path, abci.RequestQuery{})
	require.NotNil(t, err)
}

func TestQueryDepthBookV2(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	depthBook := &types.DepthBook{}
	product := types.TestTokenPair
	order1 := mockOrder("", product, types.SellOrder, "0.6", "1.1")
	depthBook.InsertOrder(order1)
	depthBook.InsertOrder(order1)
	order2 := mockOrder("", product, types.SellOrder, "0.5", "1.2")
	depthBook.InsertOrder(order2)
	order3 := mockOrder("", product, types.BuyOrder, "0.4", "1.3")
	depthBook.InsertOrder(order3)
	order4 := mockOrder("", product, types.BuyOrder, "0.3", "1.4")
	depthBook.InsertOrder(order4)
	depthBook.InsertOrder(order4)
	keeper.StoreDepthBook(ctx, product, depthBook)

	// Default query
	path := []string{types.QueryDepthBookV2}
	params := NewQueryDepthBookParams(product, 0)
	data := keeper.cdc.MustMarshalJSON(params)
	req := abci.RequestQuery{
		Data: data,
	}
	bookResBytes, err := querier(ctx, path, req)
	require.Nil(t, err)
	require.NotNil(t, bookResBytes)
	//bookRes := &BookRes{}
	//keeper.cdc.MustUnmarshalJSON(bookResBytes, bookRes)
	//expectBookRes := &BookRes{
	//	Asks: []BookResItem{
	//		{sdk.MustNewDecFromStr("0.5").String(), sdk.MustNewDecFromStr("1.2").String()},
	//		{sdk.MustNewDecFromStr("0.6").String(), sdk.MustNewDecFromStr("2.2").String()},
	//	},
	//	Bids: []BookResItem{
	//		{sdk.MustNewDecFromStr("0.4").String(), sdk.MustNewDecFromStr("1.3").String()},
	//		{sdk.MustNewDecFromStr("0.3").String(), sdk.MustNewDecFromStr("2.8").String()},
	//	},
	//}
	//require.EqualValues(t, expectBookRes, bookRes)
}

func TestQueryDepthBook(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	depthBook := &types.DepthBook{}
	product := types.TestTokenPair
	order1 := mockOrder("", product, types.SellOrder, "0.6", "1.1")
	depthBook.InsertOrder(order1)
	depthBook.InsertOrder(order1)
	order2 := mockOrder("", product, types.SellOrder, "0.5", "1.2")
	depthBook.InsertOrder(order2)
	order3 := mockOrder("", product, types.BuyOrder, "0.4", "1.3")
	depthBook.InsertOrder(order3)
	order4 := mockOrder("", product, types.BuyOrder, "0.3", "1.4")
	depthBook.InsertOrder(order4)
	depthBook.InsertOrder(order4)
	keeper.StoreDepthBook(ctx, product, depthBook)

	// Default query
	path := []string{types.QueryDepthBook}
	params := NewQueryDepthBookParams(product, 0)
	data := keeper.cdc.MustMarshalJSON(params)
	req := abci.RequestQuery{
		Data: data,
	}
	bookResBytes, err := querier(ctx, path, req)
	require.Nil(t, err)
	bookRes := &BookRes{}
	keeper.cdc.MustUnmarshalJSON(bookResBytes, bookRes)
	expectBookRes := &BookRes{
		Asks: []BookResItem{
			{sdk.MustNewDecFromStr("0.5").String(), sdk.MustNewDecFromStr("1.2").String()},
			{sdk.MustNewDecFromStr("0.6").String(), sdk.MustNewDecFromStr("2.2").String()},
		},
		Bids: []BookResItem{
			{sdk.MustNewDecFromStr("0.4").String(), sdk.MustNewDecFromStr("1.3").String()},
			{sdk.MustNewDecFromStr("0.3").String(), sdk.MustNewDecFromStr("2.8").String()},
		},
	}
	require.EqualValues(t, expectBookRes, bookRes)

	// limit size
	params = NewQueryDepthBookParams(product, 1)
	data = keeper.cdc.MustMarshalJSON(params)
	req = abci.RequestQuery{
		Data: data,
	}
	bookResBytes, err = querier(ctx, path, req)
	require.Nil(t, err)
	bookRes = &BookRes{}
	keeper.cdc.MustUnmarshalJSON(bookResBytes, bookRes)
	expectBookRes = &BookRes{
		Asks: []BookResItem{
			{sdk.MustNewDecFromStr("0.5").String(), sdk.MustNewDecFromStr("1.2").String()},
		},
		Bids: []BookResItem{
			{sdk.MustNewDecFromStr("0.4").String(), sdk.MustNewDecFromStr("1.3").String()},
		},
	}
	require.EqualValues(t, expectBookRes, bookRes)

	// invalid request
	req = abci.RequestQuery{
		Data: nil,
	}
	_, err = querier(ctx, path, req)
	require.NotNil(t, err)

	// invalid product
	params = NewQueryDepthBookParams("invalid_product", 0)
	data = keeper.cdc.MustMarshalJSON(params)
	req = abci.RequestQuery{
		Data: data,
	}
	_, err = querier(ctx, path, req)
	require.NotNil(t, err)

	path = []string{types.QueryDepthBookV2}
	bookResBytes, err = querier(ctx, path, req)
	require.Nil(t, err)
	require.NotNil(t, bookResBytes)
}

func TestQueryStore(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	product := types.TestTokenPair
	orders := []*types.Order{
		mockOrder("", product, types.SellOrder, "0.6", "1.1"),
		mockOrder("", product, types.SellOrder, "0.5", "1.2"),
		mockOrder("", product, types.BuyOrder, "0.4", "1.3"),
		mockOrder("", product, types.BuyOrder, "0.3", "1.4"),
	}
	for i := 0; i < 4; i++ {
		orders[i].Sender = testInput.TestAddrs[0]
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	keeper.Cache2Disk(ctx)

	// Default query
	path := []string{types.QueryStore}
	bz, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	ss := &StoreStatistic{}
	keeper.cdc.MustUnmarshalJSON(bz, ss)
	require.EqualValues(t, 4, ss.StoreOrderNum)
	require.EqualValues(t, 4, ss.DepthBookNum[product])

	keeper.SetLastPrice(ctx, product, sdk.MustNewDecFromStr("1.5"))
	keeper.DumpStore(ctx)

	path = []string{types.QueryDepthBookV2}
	_, sdkErr := querier(ctx, path, abci.RequestQuery{})
	require.EqualValues(t, sdk.CodeUnknownRequest, sdkErr.Code())
}

func TestQueryParameters(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	params := &types.Params{
		OrderExpireBlocks: 1000,
		MaxDealsPerBlock:  10000,
		FeePerBlock:       sdk.NewDecCoinFromDec(types.DefaultFeeDenomPerBlock, sdk.NewDec(1)),
		TradeFeeRate:      sdk.MustNewDecFromStr("0.001"),
	}
	keeper.SetParams(ctx, params)
	path := []string{types.QueryParameters}
	byteRes, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	res := &types.Params{}
	err2 := keeper.cdc.UnmarshalJSON(byteRes, res)
	require.Nil(t, err2)
	require.EqualValues(t, params, res)
}

func TestQueryInvalidPath(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	path := []string{"invalid-path"}
	_, err := querier(ctx, path, abci.RequestQuery{})
	require.NotNil(t, err)
	require.EqualValues(t, sdk.CodeUnknownRequest, err.Code())
}
