package proto

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAppUpgradeConfig_String(t *testing.T) {
	appUpgradeConfig := NewAppUpgradeConfig(1, NewProtocolDefinition(1, "http://web.abc", 100, sdk.NewDecWithPrec(75, 2)))
	require.NotEqual(t, len(appUpgradeConfig.String()), 0)
}

func TestDefaultUpgradeConfig(t *testing.T) {
	ctx, keeper := createTestInput(t)
	appUpgradeConfig := DefaultUpgradeConfig("http://web.abc")
	require.NotPanics(t, func() { keeper.SetUpgradeConfig(ctx, appUpgradeConfig) })

	_, found := keeper.GetUpgradeConfig(ctx)
	require.True(t, found)

	require.NotPanics(t, func() { keeper.ClearUpgradeConfig(ctx) })
	_, found = keeper.GetUpgradeConfig(ctx)
	require.False(t, found)
}

func TestCurrentVersion(t *testing.T) {
	ctx, keeper := createTestInput(t)

	require.Equal(t, keeper.GetCurrentVersion(ctx), uint64(0))

	require.NotPanics(t, func() { keeper.SetCurrentVersion(ctx, 1) })
	require.Equal(t, keeper.GetCurrentVersion(ctx), uint64(1))
}

func TestLastFailedVersion(t *testing.T) {
	ctx, keeper := createTestInput(t)

	require.Equal(t, keeper.GetLastFailedVersion(ctx), uint64(0))

	require.NotPanics(t, func() { keeper.SetLastFailedVersion(ctx, 1) })
	require.Equal(t, keeper.GetLastFailedVersion(ctx), uint64(1))
}

func TestValidVersion(t *testing.T) {
	ctx, keeper := createTestInput(t)

	require.True(t, keeper.IsValidVersion(ctx, 1))

	require.NotPanics(t, func() { keeper.SetLastFailedVersion(ctx, 1) })
	require.True(t, keeper.IsValidVersion(ctx, 2))
}

func TestVersionByStore(t *testing.T) {
	ctx, keeper := createTestInput(t)

	store := ctx.KVStore(keeper.storeKey)

	require.Equal(t, keeper.GetCurrentVersionByStore(store), uint64(0))

	require.NotPanics(t, func() { keeper.SetCurrentVersion(ctx, 1) })
	require.Equal(t, keeper.GetCurrentVersionByStore(store), uint64(1))
}

func TestUpgradeConfigByStore(t *testing.T) {
	ctx, keeper := createTestInput(t)

	store := ctx.KVStore(keeper.storeKey)

	_, found := keeper.GetUpgradeConfigByStore(store)
	require.False(t, found)

	appUpgradeConfig := DefaultUpgradeConfig("http://web.abc")
	require.NotPanics(t, func() { keeper.SetUpgradeConfig(ctx, appUpgradeConfig) })
	_, found = keeper.GetUpgradeConfigByStore(store)
	require.True(t, found)
}