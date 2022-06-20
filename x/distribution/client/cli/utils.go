package cli

import (
	"io/ioutil"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type (
	// CommunityPoolSpendProposalJSON defines a CommunityPoolSpendProposal with a deposit
	CommunityPoolSpendProposalJSON struct {
		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Recipient   sdk.AccAddress `json:"recipient" yaml:"recipient"`
		Amount      sdk.SysCoins   `json:"amount" yaml:"amount"`
		Deposit     sdk.SysCoins   `json:"deposit" yaml:"deposit"`
	}

	// ChangeDistributionTypeProposalJSON defines a ChangeDistributionTypeProposal with a deposit
	ChangeDistributionTypeProposalJSON struct {
		Title       string       `json:"title" yaml:"title"`
		Description string       `json:"description" yaml:"description"`
		Type        uint32       `json:"type" yaml:"type"`
		Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
	}
)

// ParseCommunityPoolSpendProposalJSON reads and parses a CommunityPoolSpendProposalJSON from a file.
func ParseCommunityPoolSpendProposalJSON(cdc *codec.Codec, proposalFile string) (CommunityPoolSpendProposalJSON, error) {
	proposal := CommunityPoolSpendProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// ParseChangeDistributionTypeProposalJSON reads and parses a ChangeDistributionTypeProposalJSON from a file.
func ParseChangeDistributionTypeProposalJSON(cdc *codec.Codec, proposalFile string) (ChangeDistributionTypeProposalJSON, error) {
	proposal := ChangeDistributionTypeProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
