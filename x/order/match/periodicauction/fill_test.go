package periodicauction

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex"
	orderkeeper "github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/order/types"
)

var mockOrder = types.MockOrder

func TestFillDepthBook(t *testing.T) {
	common.InitConfig()
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	keeper.ResetCache(ctx)
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "9.9", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "1.1"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.1", "1.1"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	// check account balance
	acc0 := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("79.5816")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.4816")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97.8")),
	}
	require.EqualValues(t, expectCoins0, acc0.GetCoins())
	require.EqualValues(t, expectCoins1, acc1.GetCoins())

	// call fillDepthBook
	buyExecuted := sdk.ZeroDec()
	sellExecuted := sdk.ZeroDec()
	remainDeals := int64(1000)
	deals, _ := fillDepthBook(ctx, keeper, types.TestTokenPair,
		sdk.MustNewDecFromStr("10.0"), sdk.MustNewDecFromStr("1.0"),
		&buyExecuted, &sellExecuted, remainDeals, &feeParams)
	depthBook := keeper.GetDepthBookCopy(types.TestTokenPair)

	// check depthBook
	expectDepthBook := types.DepthBook{
		Items: []types.DepthBookItem{
			{Price: sdk.MustNewDecFromStr("10.1"), BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.MustNewDecFromStr("1.1")},
			{Price: sdk.MustNewDecFromStr("10.0"), BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.MustNewDecFromStr("0.1")},
			{Price: sdk.MustNewDecFromStr("9.9"), BuyQuantity: sdk.MustNewDecFromStr("1.0"), SellQuantity: sdk.ZeroDec()},
		},
	}
	//require.EqualValues(t, expectDepthBook, *depthBook)
	require.EqualValues(t, 3, len(depthBook.Items))
	for i := 0; i < 3; i++ {
		require.Equal(t, expectDepthBook.Items[i].Price.String(), depthBook.Items[i].Price.String())
		require.Equal(t, expectDepthBook.Items[i].BuyQuantity.String(),
			depthBook.Items[i].BuyQuantity.String())
		require.Equal(t, expectDepthBook.Items[i].SellQuantity.String(),
			depthBook.Items[i].SellQuantity.String())
	}

	// check deals
	require.EqualValues(t, 2, len(deals))
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), deals[0].Quantity)

	// check orderIDsMap
	keys := [4]string{}
	for i, order := range orders {
		keys[i] = types.FormatOrderIDsKey(order.Product, order.Price, order.Side)
	}
	orderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()
	require.EqualValues(t, orders[0].OrderID, orderIDsMap.Data[keys[0]][0])
	require.EqualValues(t, 0, len(orderIDsMap.Data[keys[1]]))
	require.EqualValues(t, orders[2].OrderID, orderIDsMap.Data[keys[2]][0])
	require.EqualValues(t, orders[3].OrderID, orderIDsMap.Data[keys[3]][0])

	// check order status
	order := keeper.GetOrder(ctx, orders[1].OrderID)
	require.EqualValues(t, types.OrderStatusFilled, order.Status)

	// check account balance
	acc0 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("79.8408")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100.999")), // 100 + 1 * (1 - 0.001)
	}
	expectCoins1 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("109.4716")), // 100 + 10 * (1 - 0.001) - 0.004
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97.8")),       // no change
	}
	require.EqualValues(t, expectCoins0, acc0.GetCoins())
	require.EqualValues(t, expectCoins1, acc1.GetCoins())
}

func TestFillDepthBookSecondCase(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "9.8", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "1.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.1"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]
	depthBook := &types.DepthBook{}

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
		depthBook.InsertOrder(orders[i])
	}
	// check account balance
	acc0 := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("79.5816")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.4816")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97.9")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	buyExecuted := sdk.ZeroDec()
	sellExecuted := sdk.ZeroDec()
	remainDeals := int64(1000)
	deals, _ := fillDepthBook(ctx, keeper, types.TestTokenPair,
		sdk.MustNewDecFromStr("10.0"), sdk.MustNewDecFromStr("1.0"),
		&buyExecuted, &sellExecuted, remainDeals, &feeParams)
	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)

	// check depthBook
	expectDepthBook := types.DepthBook{
		Items: []types.DepthBookItem{
			{Price: sdk.MustNewDecFromStr("10.2"), BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.MustNewDecFromStr("1.1")},
			{Price: sdk.MustNewDecFromStr("9.8"), BuyQuantity: sdk.MustNewDecFromStr("1.0"), SellQuantity: sdk.ZeroDec()},
		},
	}
	//require.EqualValues(t, expectDepthBook, *depthBook)
	require.EqualValues(t, 2, len(depthBook.Items))
	for i := 0; i < 2; i++ {
		require.True(t, expectDepthBook.Items[i].Price.Equal(depthBook.Items[i].Price))
		require.True(t, expectDepthBook.Items[i].BuyQuantity.Equal(depthBook.Items[i].BuyQuantity))
		require.True(t,
			expectDepthBook.Items[i].SellQuantity.Equal(depthBook.Items[i].SellQuantity))
	}

	// check deals
	require.EqualValues(t, 2, len(deals))
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), deals[0].Quantity)

	// check orderIDsMap
	keys := [4]string{}
	for i, order := range orders {
		keys[i] = types.FormatOrderIDsKey(order.Product, order.Price, order.Side)
	}
	// call fillDepthBook
	newOrderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()

	// OrderIDsMap will remove empty keys(10.1-BUY, 9.9-SELL) immediately, not by the end of endblock
	require.EqualValues(t, 2, len(newOrderIDsMap.Data))
	require.EqualValues(t, 0, len(newOrderIDsMap.Data[keys[1]]))
	require.EqualValues(t, 0, len(newOrderIDsMap.Data[keys[2]]))

	// check order status
	order := keeper.GetOrder(ctx, orders[1].OrderID)
	require.EqualValues(t, types.OrderStatusFilled, order.Status)

	// check account balance
	acc0 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("79.9408")), // 80.1 + 0.1 - 0.004
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100.999")),   // 100 + 1 * (1 - 0.001)
	}
	expectCoins1 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("109.7308")), // 100 + 10 * (1 - 0.001) - 0.004
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("97.9")),       // no change
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())
}

func TestPartialFillDepthBook(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]
	depthBook := &types.DepthBook{}

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
		depthBook.InsertOrder(orders[i])
	}
	// check account balance
	acc0 := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 := testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("69.1816")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100")),
	}
	expectCoins1 := sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.4816")),
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("96")),
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// call fillDepthBook
	buyExecuted := sdk.ZeroDec()
	sellExecuted := sdk.ZeroDec()
	remainDeals := int64(1)
	// only fill orders[0]
	deals, remainDeals := fillDepthBook(ctx, keeper, types.TestTokenPair,
		sdk.MustNewDecFromStr("10.0"), sdk.MustNewDecFromStr("3.0"),
		&buyExecuted, &sellExecuted, remainDeals, &feeParams)
	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)

	// check result
	require.EqualValues(t, 1, len(deals))
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), deals[0].Quantity)
	require.EqualValues(t, orders[0].OrderID, deals[0].OrderID)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.0"), buyExecuted)
	require.EqualValues(t, sdk.ZeroDec(), sellExecuted)
	require.EqualValues(t, 0, remainDeals)

	// check depthBook
	expectDepthBook := types.DepthBook{
		Items: []types.DepthBookItem{
			{Price: sdk.MustNewDecFromStr("10.2"), BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.MustNewDecFromStr("1.0")},
			{Price: sdk.MustNewDecFromStr("10.1"), BuyQuantity: sdk.MustNewDecFromStr("2.0"), SellQuantity: sdk.ZeroDec()},
			{Price: sdk.MustNewDecFromStr("9.9"), BuyQuantity: sdk.ZeroDec(), SellQuantity: sdk.MustNewDecFromStr("3.0")},
		},
	}
	//require.EqualValues(t, expectDepthBook, *depthBook)
	require.EqualValues(t, 3, len(depthBook.Items))
	for i := 0; i < 3; i++ {
		require.True(t, expectDepthBook.Items[i].Price.Equal(depthBook.Items[i].Price))
		require.True(t, expectDepthBook.Items[i].BuyQuantity.Equal(depthBook.Items[i].BuyQuantity))
		require.True(t,
			expectDepthBook.Items[i].SellQuantity.Equal(depthBook.Items[i].SellQuantity))
	}

	// check order status
	order1 := keeper.GetOrder(ctx, orders[0].OrderID)
	require.EqualValues(t, types.OrderStatusFilled, order1.Status)
	order2 := keeper.GetOrder(ctx, orders[1].OrderID)
	require.EqualValues(t, types.OrderStatusOpen, order2.Status)

	// check account balance
	acc0 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("69.5408")), // 69.7 + 0.1 - 0.004
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("100.999")),   // 100 + 1 * (1 - 0.001)
	}
	expectCoins1 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("99.4816")), // no change
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("96")),        // no change
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())

	// fill orders[1] & orders[2]
	remainDeals = int64(1000)
	deals, remainDeals = fillDepthBook(ctx, keeper, types.TestTokenPair,
		sdk.MustNewDecFromStr("10.0"), sdk.MustNewDecFromStr("3.0"),
		&buyExecuted, &sellExecuted, remainDeals, &feeParams)

	require.EqualValues(t, 2, len(deals))
	require.EqualValues(t, sdk.MustNewDecFromStr("3.0"), buyExecuted)
	require.EqualValues(t, sdk.MustNewDecFromStr("3.0"), sellExecuted)
	require.EqualValues(t, 998, remainDeals)

	// check account balance
	acc0 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[0])
	acc1 = testInput.AccountKeeper.GetAccount(ctx, testInput.TestAddrs[1])
	expectCoins0 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("70")),    // 69.7 + (0.1 - 0.004) * 3
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("102.997")), // 100.999 + 2 * (1 - 0.001)
	}
	expectCoins1 = sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr("129.7108")), // 99.4816 + 0.2592 + 30 * (1 - 0.001)
		sdk.NewDecCoinFromDec(common.TestToken, sdk.MustNewDecFromStr("96")),         // no change
	}
	require.EqualValues(t, expectCoins0.String(), acc0.GetCoins().String())
	require.EqualValues(t, expectCoins1.String(), acc1.GetCoins().String())
}

func TestFillDepthBookByZeroMaxExecution(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock orders, DepthBook, and orderIDsMap
	keeper.ResetCache(ctx)
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "9.9", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "1.1"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.1", "1.1"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	buyExecuted := sdk.ZeroDec()
	sellExecuted := sdk.ZeroDec()
	remainDeals := int64(1000)
	deals, remainDeals := fillDepthBook(ctx, keeper, types.TestTokenPair,
		sdk.NewDec(10), sdk.NewDec(0),
		&buyExecuted, &sellExecuted, remainDeals, &feeParams)

	require.EqualValues(t, remainDeals, int64(1000))
	require.EqualValues(t, 0, len(deals))
}

func TestFillBuyOrders(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]
	depthBook := &types.DepthBook{}

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
		depthBook.InsertOrder(orders[i])
	}

	maxExecution := sdk.NewDec(100.0)
	buyExecution := sdk.ZeroDec()
	bestPrice := sdk.NewDec(10.0)
	blockRemainDeals := int64(1000)
	feeParams := types.DefaultTestParams()

	buyDeals, blockRemainDeals := fillBuyOrders(ctx, keeper, types.TestTokenPair,
		bestPrice, maxExecution, &buyExecution, blockRemainDeals, &feeParams)

	require.EqualValues(t, 2, len(buyDeals))
	require.EqualValues(t, int64(998), blockRemainDeals)
}

func TestFillSellOrders(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]
	depthBook := &types.DepthBook{}

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
		depthBook.InsertOrder(orders[i])
	}

	maxExecution := sdk.NewDec(100.0)
	sellExecution := sdk.ZeroDec()
	bestPrice := sdk.NewDec(10.0)
	blockRemainDeals := int64(1000)
	feeParams := types.DefaultTestParams()

	sellDeals, blockRemainDeals := fillSellOrders(ctx, keeper, types.TestTokenPair,
		bestPrice, maxExecution, &sellExecution, blockRemainDeals, &feeParams)

	require.EqualValues(t, 1, len(sellDeals))
	require.EqualValues(t, int64(999), blockRemainDeals)
}

func TestFillSellOrdersByLimitedMaxDeals(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]
	depthBook := &types.DepthBook{}

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
		depthBook.InsertOrder(orders[i])
	}

	maxExecution := sdk.NewDec(100.0)
	sellExecution := sdk.ZeroDec()
	bestPrice := sdk.NewDec(10.0)
	blockRemainDeals := int64(1)
	feeParams := types.DefaultTestParams()

	_, blockRemainDeals = fillSellOrders(ctx, keeper, types.TestTokenPair,
		bestPrice, maxExecution, &sellExecution, blockRemainDeals, &feeParams)

	require.EqualValues(t, 0, blockRemainDeals)
}

func TestFillOrderByKey(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.EqualValues(t, nil, err)
	}

	fillPrice := sdk.NewDec(10.0)
	needFillAmount := sdk.NewDec(3.0)
	feeParams := types.DefaultTestParams()
	remainDeals := int64(1000)
	key := types.FormatOrderIDsKey(types.TestTokenPair, orders[0].Price, types.BuyOrder)

	deals, filledAmount, filledDealsCnt := fillOrderByKey(ctx, keeper, key, needFillAmount, fillPrice, &feeParams,
		remainDeals)

	require.EqualValues(t, 2, len(deals))
	require.EqualValues(t, sdk.NewDec(3), filledAmount)
	require.EqualValues(t, 2, filledDealsCnt)
}

func TestFillOrderByKeyByNotExistKey(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.EqualValues(t, nil, err)
	}

	fillPrice := sdk.NewDec(10.0)
	needFillAmount := sdk.NewDec(3.0)
	feeParams := types.DefaultTestParams()
	remainDeals := int64(1000)
	key := types.FormatOrderIDsKey(types.TestTokenPair+"_test", orders[0].Price, types.BuyOrder)

	deals, filledAmount, filledDealsCnt := fillOrderByKey(ctx, keeper, key, needFillAmount, fillPrice, &feeParams,
		remainDeals)
	require.EqualValues(t, 0, len(deals))
	require.EqualValues(t, filledAmount, sdk.ZeroDec())
	require.EqualValues(t, filledDealsCnt, int64(0))
}

func TestFillOrder(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.EqualValues(t, nil, err)
	}

	fillPrice := sdk.NewDec(10.0)
	fillQuantity := sdk.NewDec(1.0)
	feeParams := types.DefaultTestParams()

	for _, order := range orders {
		retDeals := fillOrder(order, ctx, keeper, fillPrice, fillQuantity, &feeParams)
		require.NotEmpty(t, retDeals)
	}
}

func TestTransferTokens(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.1", "2.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "9.9", "3.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.2", "1.0"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[2].Sender = testInput.TestAddrs[1]
	orders[3].Sender = testInput.TestAddrs[1]

	for i := 0; i < 4; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.EqualValues(t, nil, err)
	}

	fillPrice := sdk.NewDec(10.0)
	fillQuantity := sdk.NewDec(1.0)
	for _, order := range orders {
		balanceAccount(order, ctx, keeper, fillPrice, fillQuantity)
	}
}

func TestChargeFee(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)
	testInput.DexKeeper.SetOperator(ctx, dex.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
	})

	keeper.ResetCache(ctx)
	orders := []*types.Order{
		mockOrder("", types.TestTokenPair, types.BuyOrder, "9.9", "1.0"),
		mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.0", "1.1"),
		mockOrder("", types.TestTokenPair, types.SellOrder, "10.1", "1.1"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[0]
	orders[1].Status = types.OrderStatusFilled
	orders[2].Sender = testInput.TestAddrs[1]
	orders[2].Status = types.OrderStatusFilled
	orders[3].Sender = testInput.TestAddrs[1]

	fillQuantity := sdk.NewDec(1.0)
	feeParams := types.DefaultTestParams()

	for _, order := range orders {
		retFee, feeReceiver := chargeFee(order, ctx, keeper, fillQuantity, &feeParams)
		require.NotEmpty(t, retFee)
		require.NotEmpty(t, feeReceiver)
	}
}
