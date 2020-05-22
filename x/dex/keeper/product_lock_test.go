package keeper

import (
	"testing"

	"github.com/okex/okchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestGetLockedProductsCopy(t *testing.T) {

	testInput := createTestInput(t)
	keeper := testInput.DexKeeper
	ctx := testInput.Ctx

	productLock1 := &types.ProductLock{
		BlockHeight: 1,
	}
	productLock2 := &types.ProductLock{
		BlockHeight: 2,
	}
	keeper.LockTokenPair(ctx, "p1", productLock1)
	keeper.LockTokenPair(ctx, "p2", productLock2)

	copy := keeper.GetLockedProductsCopy(ctx)
	require.EqualValues(t, 1, copy.Data["p1"].BlockHeight)
	require.EqualValues(t, 2, copy.Data["p2"].BlockHeight)

	copy.Data["p1"].BlockHeight = 10
	copy.Data["p2"].BlockHeight = 20

	copy2 := keeper.GetLockedProductsCopy(ctx)
	require.EqualValues(t, 1, copy2.Data["p1"].BlockHeight)
	require.EqualValues(t, 2, copy2.Data["p2"].BlockHeight)

	require.NotEqual(t, copy.Data["p1"].BlockHeight, copy2.Data["p1"].BlockHeight)
	require.NotEqual(t, copy.Data["p2"].BlockHeight, copy2.Data["p2"].BlockHeight)
}

func TestProductLockMap(t *testing.T) {
	testInput := createTestInput(t)
	keeper := testInput.DexKeeper
	ctx := testInput.Ctx

	productLock1 := &types.ProductLock{
		BlockHeight: 1,
	}
	productLock2 := &types.ProductLock{
		BlockHeight: 2,
	}
	keeper.LockTokenPair(ctx, "p1", productLock1)
	keeper.LockTokenPair(ctx, "p2", productLock2)

	lockMap := keeper.GetLockedProductsCopy(ctx)
	require.EqualValues(t, 2, len(lockMap.Data))
	require.EqualValues(t, 1, lockMap.Data["p1"].BlockHeight)
	require.True(t, keeper.IsAnyProductLocked(ctx))
	lockMapDb := keeper.LoadProductLocks(ctx)
	require.True(t, lockMapDb.Data["p1"] != nil)
	require.True(t, lockMapDb.Data["p2"] != nil)
	require.EqualValues(t, 2, len(lockMapDb.Data))

	// unlock
	keeper.UnlockTokenPair(ctx, "p1")
	require.False(t, keeper.IsTokenPairLocked(ctx, "p1"))

	lockMap = keeper.GetLockedProductsCopy(ctx)
	require.EqualValues(t, 1, len(lockMap.Data))
	require.EqualValues(t, 2, lockMap.Data["p2"].BlockHeight)
	require.True(t, keeper.IsAnyProductLocked(ctx))

	// unlock all
	keeper.UnlockTokenPair(ctx, "p2")
	require.False(t, keeper.IsAnyProductLocked(ctx))
}
