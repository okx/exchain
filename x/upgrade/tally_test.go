package upgrade

import (
	"testing"

	"github.com/okex/okchain/x/staking/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking"
	//"github.com/okex/okchain/x/staking"
)

func TestTallyPassed(t *testing.T) {
	ctx, keeper, stakingKeeper, _ := testPrepare(t)
	description := staking.NewDescription("moniker1", "identity1", "website1", "details1")
	for i := 0; i < 4; i++ {
		validator := staking.NewValidator(sdk.ValAddress(accAddrs[i]), pubKeys[i], description, types.DefaultMinSelfDelegation)
		validator.Status = sdk.Bonded
		validator.DelegatorShares = sdk.OneDec()
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
	}

	require.True(t, tally(ctx, 1, keeper, sdk.NewDecWithPrec(75, 2)))
}

func TestTallyNotPassed(t *testing.T) {
	ctx, keeper, stakingKeeper, _ := testPrepare(t)
	description := staking.NewDescription("moniker2", "identity2", "website2", "details2")
	for i := 0; i < 4; i++ {
		validator := staking.NewValidator(sdk.ValAddress(accAddrs[i]), pubKeys[i], description, types.DefaultMinSelfDelegation)
		validator.Status = sdk.Bonded
		validator.DelegatorShares = sdk.OneDec()
		stakingKeeper.SetValidator(ctx, validator)
		stakingKeeper.SetValidatorByPowerIndex(ctx, validator)
		if i%2 == 0 {
			keeper.SetSignal(ctx, 1, validator.GetConsAddr().String())
		}
	}
	require.True(t, tally(ctx, 1, keeper, sdk.NewDecWithPrec(5, 2)))
	require.False(t, tally(ctx, 1, keeper, sdk.NewDecWithPrec(75, 2)))
}
