package utils

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelistProposalJSON defines a DelistProposal with a deposit used
// to parse parameter change proposals from a JSON file.
type DelistProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	BaseAsset   string       `json:"base_asset" yaml:"base_asset"`
	QuoteAsset  string       `json:"quote_asset" yaml:"quote_asset"`
	Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
}

// ParseDelistProposalJSON parse json from proposal file to DelistProposalJSON struct
func ParseDelistProposalJSON(cdc *codec.Codec, proposalFilePath string) (proposal DelistProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
