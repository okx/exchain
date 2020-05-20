package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order/types"
)

func TestOrderIDsMapInsertAndRemove(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	orderIDsMap := keeper.diskCache.OrderIDsMap

	// Test insert to same key
	order1 := mockOrder("ID1-1", types.TestTokenPair, types.BuyOrder, "0.5", "1.1")
	order2 := mockOrder("ID1-2", types.TestTokenPair, types.BuyOrder, "0.5", "1.2")
	keeper.GetDiskCache().insertOrder(order1)
	keeper.GetDiskCache().insertOrder(order2)
	key1 := types.FormatOrderIDsKey(order1.Product, order1.Price, order1.Side)
	require.EqualValues(t, 1, len(orderIDsMap.Data))
	require.EqualValues(t, "ID1-1", orderIDsMap.Data[key1][0])
	require.EqualValues(t, "ID1-2", orderIDsMap.Data[key1][1])

	// Test insert to different key
	order3 := mockOrder("ID1-3", types.TestTokenPair, types.SellOrder, "0.5", "1.3")
	keeper.GetDiskCache().insertOrder(order3)
	key2 := types.FormatOrderIDsKey(order3.Product, order3.Price, order3.Side)
	require.EqualValues(t, 2, len(orderIDsMap.Data))
	require.EqualValues(t, "ID1-3", orderIDsMap.Data[key2][0])

	// check update keys
	updatedItemKeys := keeper.GetDiskCache().GetUpdatedOrderIDKeys()
	require.Equal(t, 2, len(updatedItemKeys))
	require.Equal(t, key1, updatedItemKeys[0])
	require.Equal(t, key2, updatedItemKeys[1])
	keeper.Cache2Disk(ctx)

	// Test Remove
	keeper.GetDiskCache().removeOrder(order1)
	require.EqualValues(t, "ID1-2", orderIDsMap.Data[key1][0])
	// check update keys
	updatedItemKeys = keeper.GetDiskCache().GetUpdatedOrderIDKeys()
	require.Equal(t, 2, len(updatedItemKeys))
	require.Equal(t, key1, updatedItemKeys[0])

	// remove all
	keeper.GetDiskCache().removeOrder(order2)
	keeper.GetDiskCache().removeOrder(order3)
	require.EqualValues(t, 0, len(orderIDsMap.Data))
	require.EqualValues(t, 0, len(orderIDsMap.Data[key1]))
	require.EqualValues(t, 0, len(orderIDsMap.Data[key2]))
}

func TestRemoveOrderFromDepthBook(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// mock orders, DepthBook, and orderIDsMap
	orders := []*types.Order{
		mockOrder("ID1-1", types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		mockOrder("ID1-2", types.TestTokenPair, types.BuyOrder, "10.1", "1.0"),
		mockOrder("ID1-3", types.TestTokenPair, types.SellOrder, "10.1", "1.1"),
	}
	depthBook := &types.DepthBook{}

	for i := 0; i < 3; i++ {
		depthBook.InsertOrder(orders[i])
		keeper.GetDiskCache().insertOrder(orders[i])
	}
	keeper.SetDepthBook(types.TestTokenPair, depthBook)
	keeper.Cache2Disk(ctx)
	require.Equal(t, 2, len(depthBook.Items))
	require.Equal(t, sdk.MustNewDecFromStr("10.1"), depthBook.Items[0].Price)
	require.Equal(t, sdk.MustNewDecFromStr("1.0"), depthBook.Items[0].BuyQuantity)
	require.Equal(t, sdk.MustNewDecFromStr("1.1"), depthBook.Items[0].SellQuantity)

	// remove orders[2]
	keeper.RemoveOrderFromDepthBook(orders[2], types.FeeTypeOrderDeal)

	// check depth book
	newDepthBook := keeper.GetDepthBookCopy(types.TestTokenPair)
	require.Equal(t, 2, len(newDepthBook.Items))
	require.Equal(t, sdk.MustNewDecFromStr("10.1"), depthBook.Items[0].Price)
	require.Equal(t, sdk.MustNewDecFromStr("1.0"), depthBook.Items[0].BuyQuantity)
	require.Equal(t, sdk.ZeroDec().String(), depthBook.Items[0].SellQuantity.String())

	// check orderIDsMap
	keys := [3]string{}
	for i := 0; i < 3; i++ {
		keys[i] = types.FormatOrderIDsKey(orders[i].Product, orders[i].Price, orders[i].Side)
	}
	newOrderIDsMap := keeper.diskCache.GetOrderIDsMapCopy()
	require.Equal(t, 2, len(newOrderIDsMap.Data))
	require.Equal(t, 0, len(newOrderIDsMap.Data[keys[2]]))
	require.Equal(t, "ID1-1", newOrderIDsMap.Data[keys[0]][0])
	require.Equal(t, "ID1-2", newOrderIDsMap.Data[keys[1]][0])

	// check update keys
	updatedBookKeys := keeper.GetDiskCache().GetUpdatedDepthbookKeys()
	updatedItemKeys := keeper.GetDiskCache().GetUpdatedOrderIDKeys()
	require.Equal(t, 1, len(updatedBookKeys))
	require.Equal(t, types.TestTokenPair, updatedBookKeys[0])
	require.Equal(t, 3, len(updatedItemKeys))
	require.Equal(t, keys[2], updatedItemKeys[2])
}
