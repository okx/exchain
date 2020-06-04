package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewCertifiedTokenProposal(t *testing.T) {
	title := "Withdraw coins"
	description := "Want to get some coins as reward"
	addr, err := sdk.AccAddressFromBech32("okchain1v853tq96n9ghvyxlvqyxyj97589clccr33yr7a")
	require.Nil(t, err)
	token := CertifiedToken{
		Description: "Bitcoin in testnet，1:1 anchoring with Bitcoin",
		Symbol:      "btc",
		WholeName:   "Bitcoin",
		TotalSupply: "21000000",
		Owner:       addr,
		Mintable:    false,
	}
	proposal := NewCertifiedTokenProposal(title, description, token)

	require.Equal(t, title, proposal.GetTitle())
	require.Equal(t, description, proposal.GetDescription())
	require.Equal(t, RouterKey, proposal.ProposalRoute())
	require.Equal(t, ProposalTypeCertifiedToken, proposal.ProposalType())
	require.Nil(t, proposal.ValidateBasic())
	require.NotPanics(t, func() {
		_ = proposal.String()
	})

	proposal.Title = ""
	require.Error(t, proposal.ValidateBasic())
	proposal.Title = title
	proposal.Description = ""
	require.Error(t, proposal.ValidateBasic())
	proposal.Description = description
	proposal.Token.Description = `DescriptionDescriptionDescriptionDescriptionDescriptionDescription
DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription
DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription
DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription
DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription`
	require.Error(t, proposal.ValidateBasic())
	proposal.Token.Description = "Bitcoin in testnet，1:1 anchoring with Bitcoin"
	proposal.Token.Symbol = ""
	require.Error(t, proposal.ValidateBasic())
	proposal.Token.Symbol = "2398sojncdijwoeidf2dd"
	require.Error(t, proposal.ValidateBasic())
	proposal.Token.Symbol = "btc"
	proposal.Token.Owner = sdk.AccAddress{}
	require.Error(t, proposal.ValidateBasic())
	proposal.Token.Owner = addr
	proposal.Token.WholeName = "9834f09 sadkjnviaer  asefuh9w38 awfnoeadf ekfnre"
	require.Error(t, proposal.ValidateBasic())
	proposal.Token.WholeName = "Bitcoin"
	proposal.Token.TotalSupply = "999999999999999"
	require.Error(t, proposal.ValidateBasic())
}
