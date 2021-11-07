package types

import (
	"fmt"
	"github.com/okex/exchain/x/common"
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDelistProposal_ValidateBasic(t *testing.T) {
	common.InitConfig()
	addr, err := sdk.AccAddressFromBech32(TestTokenPairOwner)
	require.Nil(t, err)

	proposal := NewDelistProposal("proposal", "right delist proposal", addr, "eth", "btc")
	require.Equal(t, "proposal", proposal.GetTitle())
	require.Equal(t, "right delist proposal", proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, proposalTypeDelist, proposal.ProposalType())

	tests := []struct {
		name   string
		drp    DelistProposal
		result bool
	}{
		{"delist-proposal", proposal, true},

		{"no-title", DelistProposal{"", "delist proposal", addr, "eth", "btc"}, false},
		{"no-description", DelistProposal{"proposal", "", addr, "eth", "btc"}, false},
		{"no-proposer", DelistProposal{"proposal", "delist proposal", nil, "eth", "btc"}, false},
		{"no-product", DelistProposal{"proposal", "delist proposal", addr, "btc", "btc"}, false},

		{"long-title", DelistProposal{getLongString(15),
			"right delist proposal", addr, "eth", "btc"}, false},
		{"long-description", DelistProposal{"proposal",
			getLongString(501), addr, "eth", "btc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.result {
				require.Nil(t, tt.drp.ValidateBasic(), "test: %v", tt.name)
			} else {
				require.NotNil(t, tt.drp.ValidateBasic(), "test: %v", tt.name)
			}
		})
	}
}

func getLongString(n int) (s string) {
	str := "0123456789"
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s", s, str)
	}
	fmt.Println(len(s))
	return s
}
