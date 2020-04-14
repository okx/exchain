package keeper

import (
	"fmt"
	"testing"

	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order/types"
	token "github.com/okex/okchain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAccountInvariant(t *testing.T) {
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx.WithBlockHeight(10)
	invariant := ModuleAccountInvariant(keeper)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.NoError(t, err)

	order1 := mockOrder("", types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	order1.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order1)
	require.NoError(t, err)

	msg, broken := invariant(ctx)
	require.False(t, broken)
	expectedLockCoins := order1.NeedLockCoins().Add(GetOrderNewFee(order1))
	require.Equal(t, invariantMsg(expectedLockCoins), msg)

	order2 := mockOrder("", types.TestTokenPair, types.SellOrder, "20.0", "3.0")
	order2.Sender = testInput.TestAddrs[0]
	err = keeper.PlaceOrder(ctx, order2)
	require.False(t, broken)
	require.NoError(t, err)

	msg, broken = invariant(ctx)
	require.False(t, broken)
	expectedLockCoins = expectedLockCoins.Add(order2.NeedLockCoins()).Add(GetOrderNewFee(order2))
	require.Equal(t, invariantMsg(expectedLockCoins), msg)

	// cancel order
	ctx = ctx.WithBlockHeight(11)
	keeper.CancelOrder(ctx, order1, ctx.Logger())

	msg, broken = invariant(ctx)
	require.False(t, broken)
	expectedLockCoins = expectedLockCoins.Sub(order1.NeedLockCoins()).Sub(GetOrderNewFee(order1))
	require.Equal(t, invariantMsg(expectedLockCoins), msg)

	// expire order
	ctx = ctx.WithBlockHeight(12)
	keeper.ExpireOrder(ctx, order2, ctx.Logger())

	msg, broken = invariant(ctx)
	require.False(t, broken)
	expectedLockCoins = expectedLockCoins.Sub(order2.NeedLockCoins()).Sub(GetOrderNewFee(order2))
	require.Equal(t, invariantMsg(expectedLockCoins), msg)

	// error case
	err = keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, testInput.TestAddrs[1], token.ModuleName, sdk.MustParseCoins(sdk.DefaultBondDenom, "11.11"))
	require.NoError(t, err)
	_, broken = invariant(ctx)
	require.True(t, broken)
}

func invariantMsg(lockCoins sdk.DecCoins) string {
	return sdk.FormatInvariant(types.ModuleName, "locks",
		fmt.Sprintf("\ttoken ModuleAccount coins: %s\n\tsum of locks amounts:  %s\n",
			lockCoins, lockCoins))
}
