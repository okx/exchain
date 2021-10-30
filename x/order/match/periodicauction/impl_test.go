package periodicauction

import (
	"testing"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/dex"
	orderkeeper "github.com/okex/exchain/x/order/keeper"
	"github.com/okex/exchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestPaEngine_Run(t *testing.T) {
	testInput := orderkeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	var startHeight int64 = 10

	// mock orders
	orders := []*types.Order{
		types.MockOrder(types.FormatOrderID(startHeight, 1), types.TestTokenPair, types.BuyOrder, "10.0", "1.0"),
		types.MockOrder(types.FormatOrderID(startHeight, 2), types.TestTokenPair, types.SellOrder, "10.0", "0.5"),
		types.MockOrder(types.FormatOrderID(startHeight, 3), types.TestTokenPair, types.SellOrder, "10.0", "2.5"),
	}
	orders[0].Sender = testInput.TestAddrs[0]
	orders[1].Sender = testInput.TestAddrs[1]
	orders[2].Sender = testInput.TestAddrs[1]
	for i := 0; i < 3; i++ {
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}

	engine := &PaEngine{}
	engine.Run(ctx, keeper)

	// check order status
	order0 := keeper.GetOrder(ctx, orders[0].OrderID)
	order1 := keeper.GetOrder(ctx, orders[1].OrderID)
	order2 := keeper.GetOrder(ctx, orders[2].OrderID)
	require.EqualValues(t, types.OrderStatusFilled, order0.Status)
	require.EqualValues(t, types.OrderStatusFilled, order1.Status)
	require.EqualValues(t, types.OrderStatusOpen, order2.Status)
	require.EqualValues(t, sdk.MustNewDecFromStr("2"), order2.RemainQuantity)
}
