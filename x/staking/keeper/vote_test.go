package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestVote(t *testing.T) {
	ctx, _, keeper := CreateTestInput(t, false, 0)

	votes := []sdk.Dec{
		sdk.OneDec(),
		sdk.OneDec(),
		sdk.OneDec().MulInt64(100),
	}

	keeper.Keeper.SetVote(ctx, addrDels[0], addrVals[0], votes[0])
	keeper.Keeper.SetVote(ctx, addrDels[1], addrVals[0], votes[1])
	keeper.Keeper.SetVote(ctx, addrDels[2], addrVals[0], votes[2])

	var votesExportedSlice []types.VotesExported
	keeper.IterateVotes(ctx,
		func(_ int64, voterAddr sdk.AccAddress, valAddr sdk.ValAddress, votes types.Votes) (stop bool) {
			votesExportedSlice = append(votesExportedSlice, types.NewVoteExported(voterAddr, valAddr, votes))
			return false
		})

	require.True(t, len(votesExportedSlice) == 3)

	vote1, found := keeper.GetVote(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, vote1.Equal(sdk.OneDec()))

	vote2, found := keeper.GetVote(ctx, addrDels[0], addrVals[2])
	require.False(t, found)
	require.True(t, vote2.IsNil())

	keeper.DeleteVote(ctx, addrVals[0], addrDels[0])
	keeper.DeleteVote(ctx, addrVals[1], addrDels[0])
	_, found = keeper.GetVote(ctx, addrDels[0], addrVals[0])
	require.False(t, found)
}
