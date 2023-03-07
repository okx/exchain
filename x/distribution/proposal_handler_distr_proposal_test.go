package distribution

import (
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/keeper"
	"github.com/okx/okbchain/x/distribution/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
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

func makeRewardTruncatePrecisionProposal(precision int64) govtypes.Proposal {
	return govtypes.Proposal{Content: types.NewRewardTruncatePrecisionProposal(
		"Test",
		"description",
		precision,
	)}
}

func (suite *HandlerSuite) TestRewardTruncatePrecisionProposal() {
	testCases := []struct {
		title          string
		percision      int64
		expectPercison int64
		error          sdk.Error
	}{
		{
			"ok", 0, 0, sdk.Error(nil),
		},
		{
			"ok", 1, 1, sdk.Error(nil),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			ctx, _, dk, _, _ := keeper.CreateTestInputDefault(suite.T(), false, 10)
			require.Equal(suite.T(), int64(0), dk.GetRewardTruncatePrecision(ctx))
			handler := NewDistributionProposalHandler(dk)
			proposal := makeRewardTruncatePrecisionProposal(tc.percision)
			err := handler(ctx, &proposal)
			require.Equal(suite.T(), tc.error, err)
			require.Equal(suite.T(), tc.expectPercison, dk.GetRewardTruncatePrecision(ctx))
		})
	}
}
