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

	// ChangeDistributionModelProposalJSON defines a ChangeDistributionModelProposal with a deposit
	ChangeDistributionModelProposalJSON struct {
		Title       string       `json:"title" yaml:"title"`
		Description string       `json:"description" yaml:"description"`
		Model       uint32       `json:"model" yaml:"model"`
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

// ParseChangeDistributionModelProposalJSON reads and parses a ChangeDistributionModelProposalJSON from a file.
func ParseChangeDistributionModelProposalJSON(cdc *codec.Codec, proposalFile string) (ChangeDistributionModelProposalJSON, error) {
	proposal := ChangeDistributionModelProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
