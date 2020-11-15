package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestNewCommunityPoolSpendProposal(t *testing.T) {
	title := "Withdraw coins"
	description := "Want to get some coins as reward"
	recipient := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	amount := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())
	proposal := NewCommunityPoolSpendProposal(title, description, recipient, sdk.NewCoins(amount))

	require.Equal(t, title, proposal.GetTitle())
	require.Equal(t, description, proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, ProposalTypeCommunityPoolSpend, proposal.ProposalType())
	require.Nil(t, proposal.ValidateBasic())
	require.NotPanics(t, func() {
		_ = proposal.String()
	})

	proposal.Title = ""
	require.Error(t, proposal.ValidateBasic())
	proposal.Title = title
	proposal.Amount = sdk.SysCoins{sdk.SysCoin{Denom: "UNKNOWN", Amount: sdk.OneDec()}}
	require.Error(t, proposal.ValidateBasic())
	proposal.Amount = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt()))
	proposal.Recipient = nil
	require.Error(t, proposal.ValidateBasic())
}
