package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"

	//"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/upgrade/types"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestKeeper_AddNewVersionInfo(t *testing.T) {
	ctx, keeper := testPrepare(t)
	upgradeConfig := proto.NewAppUpgradeConfig(1, proto.NewProtocolDefinition(1, "software1", 1024, sdk.NewDecWithPrec(75, 2)))

	// success case
	versionInfo := types.NewVersionInfo(upgradeConfig, true)
	require.Panics(t, func() {
		keeper.getVersionInfoSuccessResult(ctx, versionInfo.UpgradeInfo.ProtocolDef.Version)
	})
	require.NotPanics(t, func() {
		keeper.AddNewVersionInfo(ctx, versionInfo)
	})
	require.Equal(t, uint64(1), keeper.getVersionInfoSuccessResult(ctx, versionInfo.UpgradeInfo.ProtocolDef.Version))

	// fail case
	versionInfo = types.NewVersionInfo(upgradeConfig, false)
	require.Panics(t, func() {
		keeper.getVersionInfoFailResult(ctx, versionInfo.UpgradeInfo.ProtocolDef.Version, versionInfo.UpgradeInfo.ProposalID)
	})
	require.NotPanics(t, func() {
		keeper.AddNewVersionInfo(ctx, versionInfo)
	})
	require.Equal(t, uint64(1), keeper.getVersionInfoFailResult(ctx, versionInfo.UpgradeInfo.ProtocolDef.Version, versionInfo.UpgradeInfo.ProposalID))
}

func TestSignal(t *testing.T) {
	ctx, keeper := testPrepare(t)
	description := staking.NewDescription("moniker1", "identity1", "website1", "details1")
	validator := staking.NewValidator(sdk.ValAddress(accAddrs[0]), pubKeys[0], description)
	validator.Status = sdk.Bonded
	validator.Tokens = reserve

	keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
	require.True(t, keeper.GetSignal(ctx, 1, validator.GetConsAddr().String()))
	require.False(t, keeper.GetSignal(ctx, 2, validator.GetConsAddr().String()))

	validator = staking.NewValidator(sdk.ValAddress(accAddrs[1]), pubKeys[1], description)
	validator.Status = sdk.Bonded
	validator.Tokens = reserve

	keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
	require.True(t, keeper.DeleteSignal(ctx, 1, validator.GetConsAddr().String()))
	require.False(t, keeper.DeleteSignal(ctx, 1, validator.GetConsAddr().String()))

	keeper.ClearSignals(ctx, 1)

	require.NotEqual(t, len(GetSignalKey(1, validator.GetConsAddr().String())), 0)
	require.NotEqual(t, len(GetFailedVersionKey(1, 1)), 0)
	require.NotEqual(t, len(GetProposalIDKey(1)), 0)
	require.NotEqual(t, len(GetSignalPrefixKey(1)), 0)
	require.NotEqual(t, len(GetSuccessVersionKey(1)), 0)

}

func TestKeeper_GetAppUpgradeConfig(t *testing.T) {
	ctx, keeper := testPrepare(t)

	require.NoError(t, keeper.SetAppUpgradeConfig(ctx, 1, 1, 100, "software1"))
	require.Error(t, keeper.SetAppUpgradeConfig(ctx, 2, 2, 200, "software2"))

	_, found := keeper.GetAppUpgradeConfig(ctx)
	require.True(t, found)
}

func TestUintToHexString(t *testing.T) {
	require.Equal(t, UintToHexString(15), "000000000000000f")
}
