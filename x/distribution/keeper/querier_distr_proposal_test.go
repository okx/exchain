package keeper

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func getQueriedValidatorOutstandingRewards(t *testing.T, ctx sdk.Context, querier sdk.Querier,
	validatorAddr sdk.ValAddress) (outstandingRewards sdk.DecCoins) {
	bz, err := amino.MarshalJSON(types.NewQueryValidatorCommissionParams(validatorAddr))
	require.NoError(t, err)
	result, err := querier(ctx, []string{types.QueryValidatorOutstandingRewards}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)
	err = amino.UnmarshalJSON(result, &outstandingRewards)
	require.NoError(t, err)

	return outstandingRewards
}

func getQueriedValidatorCommission(t *testing.T, ctx sdk.Context, querier sdk.Querier,
	validatorAddr sdk.ValAddress) (validatorCommission sdk.DecCoins) {
	bz, err := amino.MarshalJSON(types.NewQueryValidatorCommissionParams(validatorAddr))
	require.NoError(t, err)

	result, err := querier(ctx, []string{types.QueryValidatorCommission}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	err = amino.UnmarshalJSON(result, &validatorCommission)
	require.NoError(t, err)

	return validatorCommission
}

func getQueriedDelegatorTotalRewards(t *testing.T, ctx sdk.Context, querier sdk.Querier,
	delegatorAddr sdk.AccAddress) (response types.QueryDelegatorTotalRewardsResponse) {
	bz, err := amino.MarshalJSON(types.NewQueryDelegatorParams(delegatorAddr))
	require.NoError(t, err)

	result, err := querier(ctx, []string{types.QueryDelegatorTotalRewards}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	err = amino.UnmarshalJSON(result, &response)
	require.NoError(t, err)

	return response
}

func getQueriedDelegationRewards(t *testing.T, ctx sdk.Context, querier sdk.Querier,
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) (rewards sdk.DecCoins) {
	bz, err := amino.MarshalJSON(types.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr))
	require.NoError(t, err)

	result, err := querier(ctx, []string{types.QueryDelegationRewards}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	err = amino.UnmarshalJSON(result, &rewards)
	require.NoError(t, err)

	return rewards
}

func getQueriedCommunityPool(t *testing.T, ctx sdk.Context, querier sdk.Querier) (ptr []byte) {
	result, err := querier(ctx, []string{types.QueryCommunityPool}, abci.RequestQuery{Data: nil})
	require.NoError(t, err)

	err = amino.UnmarshalJSON(result, &ptr)
	require.NoError(t, err)

	return
}

func TestRewards(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInputDefault(t, false, 1000)
	querier := NewQuerier(keeper)

	keeper.SetDistributionType(ctx, types.DistributionTypeOnChain)

	//try twice, do nothing
	commissionBefore := keeper.GetValidatorAccumulatedCommission(ctx, valOpAddr1)
	require.True(t, keeper.HasValidatorOutstandingRewards(ctx, valOpAddr1))
	commissionAfter := keeper.GetValidatorAccumulatedCommission(ctx, valOpAddr1)
	require.Equal(t, commissionBefore, commissionAfter)

	//test outstanding rewards query
	outstandingRewards := sdk.DecCoins{{Denom: "mytoken", Amount: sdk.NewDec(3)}, {Denom: "myothertoken", Amount: sdk.NewDecWithPrec(3, 7)}}
	keeper.SetValidatorOutstandingRewards(ctx, valOpAddr1, outstandingRewards)
	require.Equal(t, outstandingRewards, getQueriedValidatorOutstandingRewards(t, ctx, querier, valOpAddr1))

	// test validator commission query
	commission := sdk.DecCoins{{Denom: "token1", Amount: sdk.NewDec(4)}, {Denom: "token2", Amount: sdk.NewDec(2)}}
	keeper.SetValidatorAccumulatedCommission(ctx, valOpAddr1, commission)
	retCommission := getQueriedValidatorCommission(t, ctx, querier, valOpAddr1)
	require.Equal(t, commission, retCommission)

	// test delegator's total rewards query
	delegateAmount, sdkErr := sdk.ParseDecCoin(fmt.Sprintf("100%s", sk.BondDenom(ctx)))
	require.Nil(t, sdkErr)
	dAddr1 := TestDelAddrs[0]
	err := sk.Delegate(ctx, dAddr1, delegateAmount)
	require.Nil(t, err)

	ctx.SetBlockTime(time.Now())
	// add shares
	vals, sdkErr := sk.GetValidatorsToAddShares(ctx, TestValAddrs)
	require.Nil(t, sdkErr)
	delegator, found := sk.GetDelegator(ctx, dAddr1)
	require.True(t, found)
	totalTokens := delegator.Tokens.Add(delegator.TotalDelegatedTokens)
	shares, sdkErr := sk.AddSharesToValidators(ctx, dAddr1, vals, totalTokens)
	require.Nil(t, sdkErr)
	lenVals := len(vals)
	valAddrs := make([]sdk.ValAddress, lenVals)
	for i := 0; i < lenVals; i++ {
		valAddrs[i] = vals[i].OperatorAddress
	}
	delegator.ValidatorAddresses = valAddrs
	delegator.Shares = shares
	sk.SetDelegator(ctx, delegator)

	sk.AfterDelegationModified(ctx, dAddr1, delegator.ValidatorAddresses)

	//types.NewDelegationDelegatorReward(TestValAddrs[0], nil)
	expect := types.NewQueryDelegatorTotalRewardsResponse(
		[]types.DelegationDelegatorReward{
			types.NewDelegationDelegatorReward(TestValAddrs[0], nil),
			types.NewDelegationDelegatorReward(TestValAddrs[1], nil),
			types.NewDelegationDelegatorReward(TestValAddrs[2], nil),
			types.NewDelegationDelegatorReward(TestValAddrs[3], nil)},
		nil)

	delRewards := getQueriedDelegatorTotalRewards(t, ctx, querier, dAddr1)
	require.Equal(t, expect, delRewards)

	// test delegation rewards query
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now())
	DoEditValidator(t, ctx, sk, TestValAddrs[0], newRate)
	require.NoError(t, err)

	staking.EndBlocker(ctx, sk)

	val := sk.Validator(ctx, valOpAddr1)
	rewards := getQueriedDelegationRewards(t, ctx, querier, dAddr1, TestValAddrs[0])
	require.True(t, rewards.IsZero())
	initial := int64(1000000)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}
	sdk.NewDec(initial)

	keeper.AllocateTokensToValidator(ctx, val, tokens)
	rewards = getQueriedDelegationRewards(t, ctx, querier, dAddr1, TestValAddrs[0])
	require.True(t, rewards.AmountOf(sdk.DefaultBondDenom).LT(sdk.NewDec(initial/2)))
	require.True(t, rewards.AmountOf(sdk.DefaultBondDenom).GT(sdk.NewDec(initial/2-1)))

	// test delegator's total rewards query
	delRewards = getQueriedDelegatorTotalRewards(t, ctx, querier, dAddr1)
	wantDelRewards := types.NewQueryDelegatorTotalRewardsResponse(
		[]types.DelegationDelegatorReward{
			types.NewDelegationDelegatorReward(TestValAddrs[0], rewards),
			types.NewDelegationDelegatorReward(TestValAddrs[1], nil),
			types.NewDelegationDelegatorReward(TestValAddrs[2], nil),
			types.NewDelegationDelegatorReward(TestValAddrs[3], nil)},
		rewards)

	require.Equal(t, wantDelRewards, delRewards)

	// currently community pool hold nothing so we should return null
	communityPool := getQueriedCommunityPool(t, ctx, querier)
	require.Nil(t, communityPool)
}
