package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io/ioutil"
)

// ManageWhiteListProposalJSON defines a ManageWhiteListProposalJSON with a deposit used to parse manage white list
// proposals from a JSON file.
type ManageWhiteListProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	PoolName    string       `json:"pool_name" yaml:"pool_name"`
	IsAdded     bool         `json:"is_added" yaml:"is_added"`
	Deposit     sdk.DecCoins `json:"deposit" yaml:"deposit"`
}

// ParseManageWhiteListProposalJSON parse json from proposal file to ManageWhiteListProposalJSON struct
func ParseManageWhiteListProposalJSON(cdc *codec.Codec, proposalFilePath string) (proposal ManageWhiteListProposalJSON,
	err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
