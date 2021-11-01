package keeper

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	//"github.com/okex/exchain/x/common"
	//dexTypes "github.com/okex/exchain/x/dex/types"
	"github.com/okex/exchain/x/gov/types"
)

//func TestKeeper_SubmitProposal(t *testing.T) {
//	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
//
//	// not registered router
//	content := dexTypes.NewDexListProposal("Test", "", Addrs[0], "btc-123",
//		common.NativeToken, sdk.NewDec(1000), 0, 4, 4,
//		sdk.MustNewDecFromStr("0.0001"))
//	proposal, err := keeper.SubmitProposal(ctx, content)
//	require.NotNil(t, err)
//	proposalID := proposal.ProposalID
//	fmt.Println(proposalID)
//}

func TestKeeper_GetProposalID(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ProposalIDKey)
	proposalID, err := keeper.GetProposalID(ctx)
	require.NotNil(t, err)
	require.Equal(t, uint64(0), proposalID)
}

func TestKeeper_GetProposalsFiltered(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// no proposal
	proposals := keeper.GetProposalsFiltered(ctx, nil, nil,
		types.StatusDepositPeriod, 0)
	require.Equal(t, 0, len(proposals))

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	// get proposals by status
	proposals = keeper.GetProposalsFiltered(ctx, nil, nil,
		types.StatusDepositPeriod, 0)
	require.Equal(t, 1, len(proposals))

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.SysCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)}, "")
	require.Nil(t, err)

	proposals = keeper.GetProposalsFiltered(ctx, nil, nil,
		types.StatusDepositPeriod, 0)
	require.Equal(t, 0, len(proposals))

	err, voteFee := keeper.AddVote(ctx, proposalID, Addrs[1], types.OptionYes)
	require.Nil(t, err)
	require.Equal(t, "", voteFee)

	proposals = keeper.GetProposalsFiltered(ctx, Addrs[1], nil, types.StatusNil, 0)
	require.Equal(t, 1, len(proposals))

	proposals = keeper.GetProposalsFiltered(ctx, nil, Addrs[0], types.StatusNil, 0)
	require.Equal(t, 1, len(proposals))

	proposals = keeper.GetProposalsFiltered(ctx, Addrs[1], Addrs[0], types.StatusNil, 0)
	require.Equal(t, 1, len(proposals))
}

func TestKeeper_DeleteProposal(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	require.Panics(t, func() { keeper.DeleteProposal(ctx, 1) })

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID
	keeper.DeleteProposal(ctx, proposalID)
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	require.False(t, ok)
	require.Equal(t, types.Proposal{}, proposal)
}

func TestKeeper_GetProposals(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	_, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	proposals := keeper.GetProposals(ctx)
	require.Equal(t, 2, len(proposals))
}
