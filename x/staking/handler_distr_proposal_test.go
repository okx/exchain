package staking

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	keep "github.com/okex/exchain/x/staking/keeper"
)

func TestEditValidatorCommission(t *testing.T) {
	ctx, _, mKeeper := CreateTestInput(t, false, SufficientInitPower)
	tmtypes.UnittestOnlySetMilestoneVenus3Height(-1)
	keeper := mKeeper.Keeper
	_ = setInstantUnbondPeriod(keeper, ctx)
	handler := NewHandler(keeper)

	newRate, _ := sdk.NewDecFromStr("0.5")
	msgEditValidator := NewMsgEditValidatorCommissionRate(sdk.ValAddress(keep.Addrs[0]), newRate)
	require.Nil(t, msgEditValidator.ValidateBasic())

	// validator not exist
	got, err := handler(ctx, msgEditValidator)
	require.NotNil(t, err, "%v", got)

	//create validator
	validatorAddr := sdk.ValAddress(keep.Addrs[0])
	msgCreateValidator := NewTestMsgCreateValidator(validatorAddr, keep.PKs[0], DefaultMSD)
	got, err = handler(ctx, msgCreateValidator)
	require.Nil(t, err, "expected create-validator to be ok, got %v", got)

	// must end-block
	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))
	SimpleCheckValidator(t, ctx, keeper, validatorAddr, DefaultMSD, sdk.Bonded,
		SharesFromDefaultMSD, false)

	// invalid rate
	newRate, _ = sdk.NewDecFromStr("-0.5")
	msgEditValidator = NewMsgEditValidatorCommissionRate(validatorAddr, newRate)
	require.NotNil(t, msgEditValidator.ValidateBasic())
	got, err = handler(ctx, msgEditValidator)
	require.NotNil(t, err)

	// normal rate
	newRate, _ = sdk.NewDecFromStr("0.5")
	msgEditValidator = NewMsgEditValidatorCommissionRate(validatorAddr, newRate)
	require.Nil(t, msgEditValidator.ValidateBasic())
	got, err = handler(ctx, msgEditValidator)
	require.Nil(t, err)

	// normal rate,
	newRate, _ = sdk.NewDecFromStr("0.7")
	msgEditValidator = NewMsgEditValidatorCommissionRate(validatorAddr, newRate)
	require.Nil(t, msgEditValidator.ValidateBasic())
	got, err = handler(ctx, msgEditValidator)
	require.NotNil(t, err)

	// normal rate,
	ctx.SetBlockTime(time.Now())
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	msgEditValidator = NewMsgEditValidatorCommissionRate(validatorAddr, newRate)
	require.Nil(t, msgEditValidator.ValidateBasic())
	got, err = handler(ctx, msgEditValidator)
	require.Nil(t, err)

}
