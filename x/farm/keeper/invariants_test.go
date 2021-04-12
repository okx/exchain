package keeper

import (
	swaptypes "github.com/okex/exchain/x/ammswap/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvariants(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	keeper.swapKeeper.SetParams(ctx, swaptypes.DefaultParams())
	initPoolsAndLockInfos(t, ctx, keeper)

	_, broken := yieldFarmingAccountInvariant(keeper.Keeper)(ctx)
	require.False(t, broken)
	_, broken = moduleAccountInvariant(keeper.Keeper)(ctx)
	require.False(t, broken)
	_, broken = mintFarmingAccountInvariant(keeper.Keeper)(ctx)
	require.False(t, broken)
}
