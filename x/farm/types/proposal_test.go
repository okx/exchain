package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	govTypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewManageWhiteListProposal(t *testing.T) {
	tests := []struct {
		title       string
		description string
		poolName    string
		isAdded     bool
		errCode     uint32
	}{
		{
			"title",
			"description",
			"pool",
			true,
			sdk.CodeOK,
		},
		{
			"",
			"description",
			"pool",
			true,
			govTypes.CodeInvalidContent,
		},
		{
			common.GetFixedLengthRandomString(govTypes.MaxTitleLength + 1),
			"description",
			"pool",
			true,
			govTypes.CodeInvalidContent,
		},
		{
			"title",
			"",
			"pool",
			true,
			govTypes.CodeInvalidContent,
		},
		{
			"title",
			common.GetFixedLengthRandomString(govTypes.MaxDescriptionLength + 1),
			"pool",
			true,
			govTypes.CodeInvalidContent,
		},
		{
			"title",
			"description",
			"",
			true,
			govTypes.CodeInvalidContent,
		},
	}

	for _, test := range tests {
		proposal := NewManageWhiteListProposal(test.title, test.description, test.poolName, true)

		require.Equal(t, test.title, proposal.GetTitle())
		require.Equal(t, test.description, proposal.GetDescription())
		require.Equal(t, RouterKey, proposal.ProposalRoute())
		require.Equal(t, proposalTypeManageWhiteList, proposal.ProposalType())

		err := proposal.ValidateBasic()
		if test.errCode != sdk.CodeOK {
			require.Error(t, err)
			testCode(t, err, test.errCode)
		}

		require.NotPanics(t, func() {
			_ = proposal.String()
		})
	}
}
