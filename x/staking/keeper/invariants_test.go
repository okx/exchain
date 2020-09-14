package keeper

import (
	"fmt"
	"testing"

	"github.com/okex/okexchain/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDelegatorAddSharesInvariant(t *testing.T) {
	initPower := int64(1000000)
	ctx, _, mKeeper := CreateTestInput(t, false, initPower)
	k := mKeeper.Keeper
	dAddr := Addrs[0]
	vAddr1, vAddr2 := sdk.ValAddress(Addrs[1]), sdk.ValAddress(Addrs[2])
	vPk1, vPk2 := PKs[1], PKs[2]

	// create validator
	msgCreateValidator1 := NewTestMsgCreateValidator(vAddr1, vPk1, types.DefaultMinSelfDelegation)
	validator1 := types.NewValidator(msgCreateValidator1.ValidatorAddress, msgCreateValidator1.PubKey,
		msgCreateValidator1.Description, msgCreateValidator1.MinSelfDelegation.Amount)
	k.SetValidator(ctx, validator1)
	k.SetValidatorByConsAddr(ctx, validator1)
	k.SetNewValidatorByPowerIndex(ctx, validator1)
	// add shares of equal value of msd for validator itself
	defaultMinSelfDelegationToken1 := sdk.NewDecCoinFromDec(k.BondDenom(ctx), validator1.MinSelfDelegation)
	err := k.AddSharesAsMinSelfDelegation(ctx, msgCreateValidator1.DelegatorAddress, &validator1, defaultMinSelfDelegationToken1)
	require.Nil(t, err)

	msgCreateValidator2 := NewTestMsgCreateValidator(vAddr2, vPk2, types.DefaultMinSelfDelegation)
	validator2 := types.NewValidator(msgCreateValidator2.ValidatorAddress, msgCreateValidator2.PubKey,
		msgCreateValidator2.Description, msgCreateValidator2.MinSelfDelegation.Amount)
	k.SetValidator(ctx, validator2)
	k.SetValidatorByConsAddr(ctx, validator2)
	k.SetNewValidatorByPowerIndex(ctx, validator2)
	// add shares of equal value of msd for validator itself
	defaultMinSelfDelegationToken2 := sdk.NewDecCoinFromDec(k.BondDenom(ctx), validator2.MinSelfDelegation)
	err = k.AddSharesAsMinSelfDelegation(ctx, msgCreateValidator2.DelegatorAddress, &validator2, defaultMinSelfDelegationToken2)
	require.Nil(t, err)

	// deposit
	delegateAmount, sdkErr := sdk.ParseDecCoin(fmt.Sprintf("100%s", k.BondDenom(ctx)))
	require.Nil(t, sdkErr)
	err = k.Delegate(ctx, dAddr, delegateAmount)
	require.Nil(t, err)

	// add shares
	vals, sdkErr := k.GetValidatorsToAddShares(ctx, []sdk.ValAddress{vAddr1, vAddr2})
	require.Nil(t, sdkErr)
	delegator, found := k.GetDelegator(ctx, dAddr)
	require.True(t, found)
	totalTokens := delegator.Tokens.Add(delegator.TotalDelegatedTokens)
	shares, sdkErr := k.AddSharesToValidators(ctx, dAddr, vals, totalTokens)
	require.Nil(t, sdkErr)
	lenVals := len(vals)
	valAddrs := make([]sdk.ValAddress, lenVals)
	for i := 0; i < lenVals; i++ {
		valAddrs[i] = vals[i].OperatorAddress
	}
	delegator.ValidatorAddresses = valAddrs
	delegator.Shares = shares
	k.SetDelegator(ctx, delegator)

	// sanity check pass
	invariantFunc := DelegatorAddSharesInvariant(k)
	require.NotNil(t, invariantFunc)
	_, broken := invariantFunc(ctx)
	require.False(t, broken)

	// sanity check pass failed
	validator, found := k.GetValidator(ctx, vAddr1)
	require.True(t, found)
	validator.DelegatorShares = validator.DelegatorShares.Add(sdk.OneDec())
	k.SetValidator(ctx, validator)
	_, broken = invariantFunc(ctx)
	require.True(t, broken)

}
