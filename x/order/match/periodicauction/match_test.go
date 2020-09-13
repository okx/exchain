package periodicauction

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/okex/okexchain/x/dex"
	orderkeeper "github.com/okex/okexchain/x/order/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/order/types"
)

type BookItemTestData struct {
	Price        string
	BuyQuantity  string
	SellQuantity string
}

type MatchTestData struct {
	items          []BookItemTestData
	pricePrecision int64
	refPrice       string
	output         string
}

func runPeriodicAuctionMatchPriceTest(t *testing.T, testData *MatchTestData, check bool) {
	book := types.DepthBook{}
	for _, input := range testData.items {
		price, err := sdk.NewDecFromStr(input.Price)
		require.Nil(t, err)
		buy, err := sdk.NewDecFromStr(input.BuyQuantity)
		require.Nil(t, err)
		sell, err := sdk.NewDecFromStr(input.SellQuantity)
		require.Nil(t, err)
		book.Items = append(book.Items, types.DepthBookItem{
			Price:        price,
			BuyQuantity:  buy,
			SellQuantity: sell,
		})
	}
	refPrice, err := sdk.NewDecFromStr(testData.refPrice)
	require.Nil(t, err)
	needres, err := sdk.NewDecFromStr(testData.output)
	require.Nil(t, err)
	bestPrice, _ := periodicAuctionMatchPrice(&book, testData.pricePrecision, refPrice)
	if check {
		if !needres.Equal(bestPrice) {
			t.Fatalf("need:%s calc:%s\n", needres.String(), bestPrice.String())
		}
	}

}

func TestPeriodicAuctionMatchPriceRandomData(t *testing.T) {
	for i := 0; i < 10000; i++ {
		n := rand.Int()%200 + 2
		data := MatchTestData{
			pricePrecision: 1,
			refPrice:       "1",
			output:         "98",
		}
		for j := 0; j < n; j++ {
			data.items = append(data.items, BookItemTestData{
				Price:        strconv.Itoa(100 - j),
				BuyQuantity:  strconv.Itoa(rand.Intn(2) * rand.Intn(10)),
				SellQuantity: strconv.Itoa(rand.Intn(2) * rand.Intn(10)),
			})
		}
		runPeriodicAuctionMatchPriceTest(t, &data, false)
	}
}

func TestPeriodicAuctionMatchPrice(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"100", "150", "0",
		}, {
			"99", "10", "0",
		}, {
			"98", "0", "250",
		}, {
			"97", "0", "50",
		},
		},
		pricePrecision: 1,
		refPrice:       "1",
		output:         "98",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceRule0(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"100", "0", "40",
		}, {
			"99", "0", "30",
		}, {
			"98", "80", "0",
		}, {
			"97", "70", "0",
		},
		},
		pricePrecision: 1,
		refPrice:       "1",
		output:         "1",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceRule1(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"100", "150", "0",
		}, {
			"98", "150", "250",
		}, {
			"97", "0", "50",
		},
		},
		pricePrecision: 1,
		refPrice:       "1",
		output:         "98",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceRule2(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"102", "30", "0",
		}, {
			"101", "10", "0",
		}, {
			"99", "50", "0",
		}, {
			"98", "0", "10",
		}, {
			"97", "0", "50",
		}, {
			"96", "15", "0",
		}, {
			"95", "0", "50",
		},
		},
		pricePrecision: 1,
		refPrice:       "1",
		output:         "97",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceRule3A(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"102", "60", "0",
		}, {
			"100", "0", "20",
		}, {
			"95", "0", "30",
		},
		},
		pricePrecision: 1,
		refPrice:       "100",
		output:         "102",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceRule3B(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"102", "30", "0",
		}, {
			"97", "20", "0",
		}, {
			"95", "0", "60",
		},
		},
		pricePrecision: 1,
		refPrice:       "96",
		output:         "95",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceRule3C(t *testing.T) {
	data := MatchTestData{
		items: []BookItemTestData{{
			"100", "25", "0",
		}, {
			"98", "0", "35",
		}, {
			"97", "35", "0",
		}, {
			"95", "0", "25",
		},
		},
		pricePrecision: 1,
		refPrice:       "99",
		output:         "99",
	}
	runPeriodicAuctionMatchPriceTest(t, &data, true)
}

func TestPeriodicAuctionMatchPriceByEmptyDepthBook(t *testing.T) {
	depthBook := &types.DepthBook{}
	bestPrice, maxExecution := periodicAuctionMatchPrice(depthBook, 10, sdk.MustNewDecFromStr("10.0"))

	require.EqualValues(t, sdk.ZeroDec(), bestPrice)
	require.EqualValues(t, sdk.ZeroDec(), maxExecution)
}

func TestPreMatchProcessing(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)

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
		depthBook.InsertOrder(orders[i])
	}

	buyAmountSum, sellAmountSum := preMatchProcessing(depthBook)
	require.Equal(t, sdk.NewDec(3), buyAmountSum[1])
	require.Equal(t, sdk.NewDec(4), sellAmountSum[0])
}

func TestExecRule0(t *testing.T) {
	buyAmountSum := []sdk.Dec{sdk.NewDec(0.0), sdk.NewDec(3.0), sdk.NewDec(3.0)}
	sellAmountSum := []sdk.Dec{sdk.NewDec(4.0), sdk.NewDec(3.0), sdk.NewDec(3.0)}

	maxExecution, execution := execRule0(buyAmountSum, sellAmountSum)
	require.EqualValues(t, sdk.NewDec(3), maxExecution)
	require.EqualValues(t, sdk.NewDec(3), execution[1])
}

func TestExecRule1(t *testing.T) {
	maxExecution := sdk.NewDec(3.0)
	execution := []sdk.Dec{sdk.NewDec(0.0), sdk.NewDec(3.0), sdk.NewDec(3.0)}

	indexesRule1 := execRule1(maxExecution, execution)
	require.EqualValues(t, []int{1, 2}, indexesRule1)
}

func TestExecRule2(t *testing.T) {
	buyAmountSum := []sdk.Dec{sdk.NewDec(0.0), sdk.NewDec(3.0), sdk.NewDec(3.0)}
	sellAmountSum := []sdk.Dec{sdk.NewDec(4.0), sdk.NewDec(3.0), sdk.NewDec(3.0)}
	indexesRule1 := []int{1, 2}

	indexesRule2, _ := execRule2(buyAmountSum, sellAmountSum, indexesRule1)
	require.EqualValues(t, []int{1, 2}, indexesRule2)
}

func TestExecRule3(t *testing.T) {
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

	indexesRule1 := []int{1, 2}
	refPrice := sdk.NewDec(10.0)
	pricePrecision := tokenPair.MaxPriceDigit
	indexesRule2 := []int{1, 2}
	imbalance := []sdk.Dec{sdk.NewDec(0.0), sdk.NewDec(0.0)}

	bestPrice := execRule3(depthBook, indexesRule1[0], refPrice, pricePrecision, indexesRule2, imbalance)
	require.EqualValues(t, sdk.NewDec(10), bestPrice)
}

func TestBestPriceFromRefPrice(t *testing.T) {
	minPrice, err := sdk.NewDecFromStr("10.1")
	require.EqualValues(t, nil, err)
	maxPrice, err := sdk.NewDecFromStr("9.9")
	require.EqualValues(t, nil, err)
	refPrice, err := sdk.NewDecFromStr("10.0")
	require.EqualValues(t, nil, err)

	bestPrice := bestPriceFromRefPrice(minPrice, maxPrice, refPrice)
	require.EqualValues(t, sdk.NewDec(10), bestPrice)
}

func TestExpireOrdersInExpiredBlock(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

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

	keeper.DropExpiredOrdersByBlockHeight(ctx, ctx.BlockHeight())

	order := keeper.GetOrder(ctx, "ID0000000000-1")
	require.NotEqual(t, nil, order)
	require.EqualValues(t, int64(types.OrderStatusExpired), order.Status)
}

func TestMarkCurBlockToFeatureExpireBlockList(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()

	markCurBlockToFutureExpireBlockList(ctx, keeper)
	expiredBlocks := keeper.GetExpireBlockHeight(ctx, ctx.BlockHeight()+feeParams.OrderExpireBlocks)
	require.EqualValues(t, 0, expiredBlocks[0])
}

func TestCleanLastBlockClosedOrders(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

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

	keeper.SetLastClosedOrderIDs(ctx, []string{orders[0].OrderID})

	cleanLastBlockClosedOrders(ctx, keeper)

	order := keeper.GetOrder(ctx, orders[0].OrderID)
	require.EqualValues(t, (*types.Order)(nil), order)
}

func TestCacheExpiredBlockToCurrentHeight(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

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

	keeper.DropExpiredOrdersByBlockHeight(ctx, ctx.BlockHeight())
	keeper.SetExpireBlockHeight(ctx, ctx.BlockHeight(), []int64{ctx.BlockHeight()})

	cacheExpiredBlockToCurrentHeight(ctx, keeper)

	num := keeper.GetCache().GetExpireNum()
	require.EqualValues(t, len(orders), int(num))
}

func TestCleanupExpiredOrders(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()

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

	keeper.SetLastClosedOrderIDs(ctx, []string{orders[0].OrderID})
	keeper.ExpireOrder(ctx, orders[1], ctx.Logger())

	cleanupExpiredOrders(ctx, keeper)

	expiredBlocks := keeper.GetExpireBlockHeight(ctx, ctx.BlockHeight()+
		feeParams.OrderExpireBlocks)
	require.EqualValues(t, true, expiredBlocks[0] == ctx.BlockHeight())

	lastClosedOrderIDs := keeper.GetLastClosedOrderIDs(ctx)
	require.EqualValues(t, true, lastClosedOrderIDs[0] == orders[0].OrderID)

	num := keeper.GetCache().GetExpireNum()
	require.EqualValues(t, 1, int(num))
}

func TestMatchOrders(t *testing.T) {
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	matchOrders(ctx, keeper)

	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("10.2"), depthBook.Items[0].Price)
	require.EqualValues(t, sdk.ZeroDec(), depthBook.Items[0].BuyQuantity)
	require.EqualValues(t, sdk.MustNewDecFromStr("1"), depthBook.Items[0].SellQuantity)

	orders = []*types.Order{
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
		depthBook.InsertOrder(orders[i])
	}
	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})
	lockProduct(ctx, keeper, ctx.Logger(), types.TestTokenPair, updatedProductsBasePrice[types.TestTokenPair],
		sdk.ZeroDec(), sdk.ZeroDec())

	matchOrders(ctx, keeper)

	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("10.2"), depthBook.Items[0].Price)
	require.EqualValues(t, sdk.ZeroDec(), depthBook.Items[0].BuyQuantity)
	require.EqualValues(t, sdk.MustNewDecFromStr("2"), depthBook.Items[0].SellQuantity)
}

func TestMatchOrdersByEmptyBlock(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	matchOrders(ctx, keeper)
	require.EqualValues(t, int64(0), keeper.GetBlockOrderNum(ctx, ctx.BlockHeight()))
	require.EqualValues(t, false, keeper.AnyProductLocked(ctx))
}

func TestCalcMatchPriceAndExecution(t *testing.T) {
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	products := keeper.GetDiskCache().GetUpdatedDepthbookKeys()
	keeper.GetDexKeeper().SortProducts(ctx, products) // sort products

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, products)
	matchResult, ok := updatedProductsBasePrice[types.TestTokenPair]
	require.EqualValues(t, ok, true)
	require.EqualValues(t, matchResult.BlockHeight, ctx.BlockHeight())
	require.EqualValues(t, sdk.MustNewDecFromStr("10.0"), matchResult.Price)
	require.EqualValues(t, sdk.MustNewDecFromStr("3.0"), matchResult.Quantity)
}

func TestLockProduct(t *testing.T) {
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})

	lockProduct(ctx, keeper, ctx.Logger(), types.TestTokenPair, updatedProductsBasePrice[types.TestTokenPair],
		sdk.ZeroDec(), sdk.ZeroDec())

	require.EqualValues(t, true, keeper.IsProductLocked(ctx, types.TestTokenPair))
}

func TestExecuteMatchedUpdatedProduct(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})

	blockRemainDeals := executeMatchedUpdatedProduct(ctx, keeper, updatedProductsBasePrice, &feeParams,
		1000, types.TestTokenPair, ctx.Logger())

	require.EqualValues(t, 997, int(blockRemainDeals))
}

func TestExecuteMatchedUpdatedProductByLimitedBlockDeals(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})

	blockRemainDeals := executeMatchedUpdatedProduct(ctx, keeper, updatedProductsBasePrice, &feeParams,
		0, types.TestTokenPair, ctx.Logger())

	require.EqualValues(t, 0, int(blockRemainDeals))
	require.EqualValues(t, true, keeper.IsProductLocked(ctx, types.TestTokenPair))
}

func TestExecuteLockedProduct(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})
	lockProduct(ctx, keeper, ctx.Logger(), types.TestTokenPair, updatedProductsBasePrice[types.TestTokenPair],
		sdk.ZeroDec(), sdk.ZeroDec())

	lockMap := keeper.GetDexKeeper().GetLockedProductsCopy(ctx)

	blockRemainDeals := executeLockedProduct(ctx, keeper, updatedProductsBasePrice, lockMap, &feeParams,
		1000, types.TestTokenPair, ctx.Logger())

	require.EqualValues(t, 997, int(blockRemainDeals))
}

func TestExecuteLockedProductByLimitedBlockDeals(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})
	lockProduct(ctx, keeper, ctx.Logger(), types.TestTokenPair, updatedProductsBasePrice[types.TestTokenPair],
		sdk.ZeroDec(), sdk.ZeroDec())

	lockMap := keeper.GetDexKeeper().GetLockedProductsCopy(ctx)

	blockRemainDeals := executeLockedProduct(ctx, keeper, updatedProductsBasePrice, lockMap, &feeParams,
		0, types.TestTokenPair, ctx.Logger())

	require.EqualValues(t, 0, int(blockRemainDeals))
}

func TestExecuteLockedProductByLargeLockQuantity(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	feeParams := types.DefaultTestParams()
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, []string{types.TestTokenPair})
	lockProduct(ctx, keeper, ctx.Logger(), types.TestTokenPair, updatedProductsBasePrice[types.TestTokenPair],
		sdk.ZeroDec(), sdk.ZeroDec())

	lockMap := keeper.GetDexKeeper().GetLockedProductsCopy(ctx)
	lock := lockMap.Data[types.TestTokenPair]
	lock.Quantity = lock.Quantity.Add(sdk.NewDec(100))

	blockRemainDeals := executeLockedProduct(ctx, keeper, updatedProductsBasePrice, lockMap, &feeParams,
		1000, types.TestTokenPair, ctx.Logger())

	require.EqualValues(t, 997, int(blockRemainDeals))
}

func TestExecuteMatch(t *testing.T) {
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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	products := keeper.GetDiskCache().GetUpdatedDepthbookKeys()
	keeper.GetDexKeeper().SortProducts(ctx, products) // sort products
	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, products)
	lockMap := keeper.GetDexKeeper().GetLockedProductsCopy(ctx)
	for product := range lockMap.Data {
		products = append(products, product)
	}
	keeper.GetDexKeeper().SortProducts(ctx, products) // sort products

	executeMatch(ctx, keeper, products, updatedProductsBasePrice, lockMap)

	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 1, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("10.2"), depthBook.Items[0].Price)
	require.EqualValues(t, sdk.ZeroDec(), depthBook.Items[0].BuyQuantity)
	require.EqualValues(t, sdk.MustNewDecFromStr("1"), depthBook.Items[0].SellQuantity)
}

func TestCleanupOrdersWhoseTokenPairHaveBeenDelisted(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

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
		require.EqualValues(t, nil, err)
		depthBook.InsertOrder(orders[i])
	}

	cleanupOrdersWhoseTokenPairHaveBeenDelisted(ctx, keeper)

	depthBook = keeper.GetDepthBookCopy(types.TestTokenPair)
	require.EqualValues(t, 0, len(depthBook.Items))

}
