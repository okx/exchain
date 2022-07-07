package cli

import (
	"io/ioutil"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type (
	// ChangeDistributionTypeProposalJSON defines a ChangeDistributionTypeProposal with a deposit
	ChangeDistributionTypeProposalJSON struct {
		Title       string       `json:"title" yaml:"title"`
		Description string       `json:"description" yaml:"description"`
		Type        uint32       `json:"type" yaml:"type"`
		Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
	}
)

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
