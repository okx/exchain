package upgrade

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/staking"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestEndBlocker(t *testing.T) {
	ctx, keeper, stakingKeeper, _ := testPrepare(t)

	// set appUpgradeConfig manually
	require.NoError(t, keeper.SetAppUpgradeConfig(ctx, 1, 1, 1024, "software1"))
	// EndBlocker should panic because of no validator
	require.Panics(t, func() {
		EndBlocker(ctx, keeper)
	})

	// create validator
	description := staking.NewDescription("moniker1", "identity1", "website1", "details1")
	validator := staking.NewValidator(sdk.ValAddress(accAddrs[0]), pubKeys[0], description)
	validator.Status = sdk.Bonded
	validator.Tokens = reserve

	// local setting
	ctx = ctx.WithBlockHeader(abci.Header{ProposerAddress: validator.GetConsAddr()})
	stakingKeeper.SetValidator(ctx, validator)
	stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
	stakingKeeper.SetValidatorByConsAddr(ctx, validator)

	require.NotPanics(t, func() {
		EndBlocker(ctx, keeper)
	})

	// check the conditional statement in EndBlocker
	ctx = ctx.WithBlockHeader(abci.Header{Version: abci.Version{Block: 1, App: 1}, ProposerAddress: validator.GetConsAddr()})
	require.NotPanics(t, func() {
		EndBlocker(ctx, keeper)
	})
	// log "Validator has downloaded the latest software"

	ctx = ctx.WithBlockHeader(abci.Header{Version: abci.Version{Block: 1, App: 0}, ProposerAddress: validator.GetConsAddr()})
	require.NotPanics(t, func() {
		EndBlocker(ctx, keeper)
	})
	// log "Validator has restarted the old software"

	// set new app upgrade config
	keeper.ClearUpgradeConfig(ctx)
	require.NoError(t, keeper.SetAppUpgradeConfig(ctx, 2, 2, 2048, "software2"))

	ctx = ctx.WithBlockHeight(2047)
	require.NotPanics(t, func() {
		EndBlocker(ctx, keeper)
	})
	// log "Tally Start" && "Software Upgrade is failure"
}

func TestEndBlockerTallySuccess(t *testing.T) {
	ctx, keeper, stakingKeeper, _ := testPrepare(t)

	require.Equal(t, uint64(0), keeper.GetCurrentVersion(ctx))

	// set appUpgradeConfig manually
	require.NoError(t, keeper.SetAppUpgradeConfig(ctx, 1, 1, 1024, "software1"))
	description := staking.NewDescription("moniker2", "identity2", "website2", "details2")

	// get validator && proposer
	validatorPro := staking.NewValidator(sdk.ValAddress(accAddrs[0]), pubKeys[0], description)
	validatorPro.Status = sdk.Bonded
	validatorPro.Tokens = reserve
	stakingKeeper.SetValidator(ctx, validatorPro)
	stakingKeeper.SetValidatorByPowerIndex(ctx, validatorPro)
	stakingKeeper.SetValidatorByConsAddr(ctx, validatorPro)

	keeper.SetSignal(ctx, 1, validatorPro.GetConsAddr().String())

	// add vote to tally
	for i := 1; i < 4; i++ {
		validator := staking.NewValidator(sdk.ValAddress(accAddrs[i]), pubKeys[i], description)
		validator.Status = sdk.Bonded
		validator.Tokens = reserve
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
	}

	ctx = ctx.WithBlockHeader(abci.Header{Version: abci.Version{Block: 1, App: 1}, ProposerAddress: validatorPro.GetConsAddr()})
	ctx = ctx.WithBlockHeight(1023)
	require.NotPanics(t, func() {
		EndBlocker(ctx, keeper)
	})
	// log "Software Upgrade is successful"
	require.Equal(t, uint64(1), keeper.GetCurrentVersion(ctx))

}
