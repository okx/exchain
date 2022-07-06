package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandleChangeDistributionTypeProposal(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(0)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)
	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	require.Equal(t, types.DistributionTypeOffChain, dk.GetDistributionType(ctx))

	//distribution type proposal ok
	proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	err := HandleChangeDistributionTypeProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, types.DistributionTypeOnChain, dk.GetDistributionType(ctx))

	//same
	err = HandleChangeDistributionTypeProposal(ctx, dk, proposal)
	require.Nil(t, err)
	require.Equal(t, types.DistributionTypeOnChain, dk.GetDistributionType(ctx))
}
