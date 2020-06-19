package staking

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keep "github.com/okex/okchain/x/staking/keeper"
	"github.com/stretchr/testify/require"
)

func TestSanityCheck(t *testing.T) {
	initPower := int64(1000000)

	ctx, _, mKeeper := CreateTestInput(t, false, initPower)
	keeper := mKeeper.Keeper

	dAddr := keep.Addrs[0]
	vAddr1, vAddr2 := sdk.ValAddress(keep.Addrs[1]), sdk.ValAddress(keep.Addrs[2])
	vPk1, vPk2 := keep.PKs[1], keep.PKs[2]

	// create validator
	msgCreateValidator1 := NewTestMsgCreateValidator(vAddr1, vPk1, DefaultMSD)
	got := handleMsgCreateValidator(ctx, msgCreateValidator1, keeper)
	require.True(t, got.IsOK(), "%v", got)

	msgCreateValidator2 := NewTestMsgCreateValidator(vAddr2, vPk2, DefaultMSD)
	got = handleMsgCreateValidator(ctx, msgCreateValidator2, keeper)
	require.True(t, got.IsOK(), "%v", got)

	// delegate
	delegateAmount, err := sdk.ParseDecCoin(fmt.Sprintf("100%s", keeper.BondDenom(ctx)))
	require.Nil(t, err)
	msgDelegate := NewMsgDeposit(dAddr, delegateAmount)
	got = handleMsgDeposit(ctx, msgDelegate, keeper)
	require.True(t, got.IsOK(), "%v", got)

	// vote
	msgVote := NewMsgAddShares(dAddr, []sdk.ValAddress{vAddr1, vAddr2})
	got = handleMsgAddShares(ctx, msgVote, keeper)
	require.True(t, got.IsOK(), "%v", got)

	// sanity check pass
	require.NotPanics(t, func() {
		sanityCheck(ctx, keeper)
	})

	// sanity check pass failed
	validator, found := keeper.GetValidator(ctx, vAddr1)
	require.True(t, found)
	validator.DelegatorShares = validator.DelegatorShares.Add(sdk.OneDec())
	keeper.SetValidator(ctx, validator)
	require.Panics(t, func() {
		sanityCheck(ctx, keeper)
	})

}
