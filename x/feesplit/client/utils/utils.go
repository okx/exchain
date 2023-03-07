package utils

import (
	"io/ioutil"

	"github.com/okx/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"
	"github.com/okx/exchain/x/feesplit/types"
)

type FeeSplitSharesProposalJSON struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Shares      []types.Shares `json:"shares" yaml:"shares"`
	Deposit     sdk.SysCoins   `json:"deposit" yaml:"deposit"`
}

// ParseFeeSplitSharesProposalJSON reads and parses a FeeSplitSharesProposalJSON from file
func ParseFeeSplitSharesProposalJSON(cdc *codec.Codec, proposalFile string) (FeeSplitSharesProposalJSON, error) {
	var proposal FeeSplitSharesProposalJSON

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
