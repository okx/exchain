package keeper

import (
	"testing"

	"github.com/okex/exchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestGetOrderIDsMapCopy(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper

	keeper.diskCache.setOrderIDs("k1", []string{"v1", "v2"})
	keeper.diskCache.setOrderIDs("k2", []string{"v3", "v4"})
	keeper.diskCache.setOrderIDs("k3", []string{})
	keeper.diskCache.setOrderIDs("k4", nil)
	idMap := keeper.diskCache.GetOrderIDsMapCopy()

	require.EqualValues(t, 2, len(idMap.Data))
	idMap.Data["k1"][0] = "update v1"
	idMap.Data["k1"][1] = "update v2"
	idMap.Data["k2"][0] = "update v3"
	idMap.Data["k2"][1] = "update v4"
	idMap.Data["k3"] = nil

	idMap2 := keeper.diskCache.GetOrderIDsMapCopy()
	require.EqualValues(t, 2, len(idMap2.Data))
	require.EqualValues(t, "v1", idMap2.Data["k1"][0])
	require.EqualValues(t, "v2", idMap2.Data["k1"][1])
	require.EqualValues(t, "v3", idMap2.Data["k2"][0])
	require.EqualValues(t, "v4", idMap2.Data["k2"][1])
}

func TestKeeper_SetBlockOrderNum(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	keeper.SetBlockOrderNum(ctx, 10, 88)
	require.EqualValues(t, 88, keeper.GetBlockOrderNum(ctx, 10))
}

func TestKeeper_DropBlockOrderNum(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	keeper.SetBlockOrderNum(ctx, 10, 88)
	keeper.DropBlockOrderNum(ctx, 10)
	require.NotEqual(t, 88, keeper.GetBlockOrderNum(ctx, 10))
}

func TestKeeper_SetExpireBlockHeight(t *testing.T) {
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

func TestKeeper_StoreOrderIDsMap(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	orderIDs := []string{"ID0000000010-1", "ID0000000010-2", "ID0000000010-3"}
	//key xxb_okt:10.00000000:BUY
	keeper.StoreOrderIDsMap(ctx, "xxb_"+common.NativeToken+":10.00000000:BUY", orderIDs)

	require.EqualValues(t, 0, len(keeper.diskCache.orderIDsMap.Data))
}

func TestFlushCache(t *testing.T) {
	testInput := CreateTestInput(t)
	dcache := newDiskCache()

	//reset
	dcache.reset()

	dcache.setOpenNum(6)
	require.EqualValues(t, 6, dcache.getOpenNum())

	//store order number
	dcache.setStoreOrderNum(88)
	dcache.DecreaseStoreOrderNum(1)
	require.EqualValues(t, 87, dcache.storeOrderNum)

	dcache.setLastPrice(types.TestTokenPair, sdk.MustNewDecFromStr("11.0"))
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "8", "1")
	order.Sender = testInput.TestAddrs[0]

	//insert & remove order
	dcache.insertOrder(order)
	require.EqualValues(t, 7, dcache.getOpenNum())
	require.EqualValues(t, sdk.MustNewDecFromStr("11"), dcache.getLastPrice(types.TestTokenPair))
	require.EqualValues(t, sdk.MustNewDecFromStr("0"), dcache.getLastPrice(types.TestTokenPair+"a"))
	require.EqualValues(t, 1, len(dcache.GetUpdatedOrderIDKeys()))
	require.EqualValues(t, 1, len(dcache.GetUpdatedDepthbookKeys()))
	require.EqualValues(t, 1, len(dcache.getOrderIDs("xxb_"+common.NativeToken+":8.000000000000000000:BUY")))
	dcache.removeOrder(order)
	require.EqualValues(t, 6, dcache.openNum)

	//close order
	dcache.closeOrder("ID0000000010-1")
	require.EqualValues(t, 2, len(dcache.GetClosedOrderIDs()))

	mapdata := dcache.GetOrderIDsMapCopy()
	require.Equal(t, make(map[string][]string), mapdata.Data)

}

func TestGetProductsFromDepthBookMap(t *testing.T) {
	testInput := CreateTestInputWithBalance(t, 1, 100)
	keeper := testInput.OrderKeeper
	product := types.TestTokenPair
	depthBook := types.DepthBook{}
	order := mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	depthBook.InsertOrder(order)
	keeper.SetDepthBook(product, &depthBook)
	productsList := keeper.GetProductsFromDepthBookMap()
	require.EqualValues(t, 1, len(productsList))
	require.EqualValues(t, product, productsList[0])
}
