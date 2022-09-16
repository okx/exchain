package distribution

import (
	"testing"

	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

func makeChangeDistributionTypeProposal(distrType uint32) govtypes.Proposal {
	return govtypes.Proposal{Content: types.NewChangeDistributionTypeProposal(
		"Test",
		"description",
		distrType,
	)}
}

func TestChangeDistributionTypeProposalHandlerPassed(t *testing.T) {
	ctx, _, k, _, _ := keeper.CreateTestInputDefault(t, false, 10)
	tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
	//init status, distribution off chain
	queryDistrType := k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOffChain)

	//set same type
	proposal := makeChangeDistributionTypeProposal(types.DistributionTypeOffChain)
	hdlr := NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOffChain)

	//set diff type, first
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOnChain)

	//set diff type, second
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOffChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOffChain)

	//set diff type, third
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOnChain)

	//set same type
	proposal = makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr = NewDistributionProposalHandler(k)
	require.NoError(t, hdlr(ctx, &proposal))
	queryDistrType = k.GetDistributionType(ctx)
	require.Equal(t, queryDistrType, types.DistributionTypeOnChain)
}
