package keeper

import (
	"testing"

	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_AnyProductLocked(t *testing.T) {
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

	require.EqualValues(t, false, keeper.IsProductLocked(ctx, types.TestTokenPair))

	//lock product
	keeper.SetProductLock(ctx, types.TestTokenPair, &types.ProductLock{})
	require.EqualValues(t, true, keeper.IsProductLocked(ctx, types.TestTokenPair))
	require.EqualValues(t, true, keeper.AnyProductLocked(ctx))

	//unlock product
	keeper.UnlockProduct(ctx, types.TestTokenPair)
	require.EqualValues(t, false, keeper.IsProductLocked(ctx, types.TestTokenPair))
}
