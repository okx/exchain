package keeper

import (
	"testing"

	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okx/exchain/libs/tendermint/types"
	"github.com/okx/exchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func TestHandleChangeDistributionTypeProposal(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	tmtypes.UnittestOnlySetMilestoneVenus2Height(0)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)
	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	require.Equal(t, types.DistributionTypeOffChain, dk.GetDistributionType(ctx))

	//distribution type proposal ok
	proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
	tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
	err := HandleChangeDistributionTypeProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, types.DistributionTypeOnChain, dk.GetDistributionType(ctx))

	//same
	err = HandleChangeDistributionTypeProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, types.DistributionTypeOnChain, dk.GetDistributionType(ctx))
}

func TestHandleWithdrawRewardEnabledProposal(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	tmtypes.UnittestOnlySetMilestoneVenus2Height(0)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)
	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	require.Equal(t, true, dk.GetWithdrawRewardEnabled(ctx))

	//set withdraw reward proposal false
	proposal := types.NewWithdrawRewardEnabledProposal("title", "description", false)
	tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
	err := HandleWithdrawRewardEnabledProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, false, dk.GetWithdrawRewardEnabled(ctx))

	//set withdraw reward proposal true
	proposal.Enabled = true
	err = HandleWithdrawRewardEnabledProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, true, dk.GetWithdrawRewardEnabled(ctx))

	//set withdraw reward proposal true, same
	err = HandleWithdrawRewardEnabledProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, true, dk.GetWithdrawRewardEnabled(ctx))
}
