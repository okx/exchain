package v0_9

import (
	"testing"

	upgradeTypes "github.com/okex/okchain/x/upgrade/types"

	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
	paramsTypes "github.com/okex/okchain/x/params/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	oldGenState := v08gov.GenesisState{
		StartingProposalID: 1,
		Proposals: []v08gov.Proposal{
			&v08gov.TextProposal{BasicProposal: v08gov.BasicProposal{ProposalType: v08gov.ProposalTypeText}},
			&v08gov.ParameterProposal{
				BasicProposal: v08gov.BasicProposal{ProposalType: v08gov.ProposalTypeParameterChange},
				Params: v08gov.Params{
					{
						Subspace: "gov",
						Key:      "test",
						Value:    "test",
					},
				},
			},
			&v08gov.AppUpgradeProposal{
				BasicProposal: v08gov.BasicProposal{ProposalType: v08gov.ProposalTypeAppUpgrade},
			},
			&v08gov.DexListProposal{
				BasicProposal: v08gov.BasicProposal{ProposalType: v08gov.ProposalTypeDexList},
			},
		},
		Deposits: v08gov.Deposits{v08gov.Deposit{}},
		Votes:    v08gov.Votes{v08gov.Vote{}},
	}
	newGenState := Migrate(oldGenState)
	require.Equal(t, 3, len(newGenState.Proposals))
	require.Equal(t, 1, len(newGenState.Deposits))
	require.Equal(t, 1, len(newGenState.Votes))
}

func TestRegisterCodec(t *testing.T) {
	cdc := codec.New()
	RegisterCodec(cdc)
	var proposal sdkGovTypes.Content
	proposalBytes := cdc.MustMarshalBinaryBare(&sdkGovTypes.TextProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok := proposal.(sdkGovTypes.TextProposal)
	require.True(t, ok)

	proposal = nil
	proposalBytes = cdc.MustMarshalBinaryBare(&paramsTypes.ParameterChangeProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok = proposal.(paramsTypes.ParameterChangeProposal)
	require.True(t, ok)

	proposal = nil
	proposalBytes = cdc.MustMarshalBinaryBare(&upgradeTypes.AppUpgradeProposal{})
	cdc.MustUnmarshalBinaryBare(proposalBytes, &proposal)
	_, ok = proposal.(upgradeTypes.AppUpgradeProposal)
	require.True(t, ok)
}
