package utils

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint/internal/types"
	"io/ioutil"
)

// ManageTreasuresProposalJSON defines a ManageTreasureProposal with a deposit used to parse
// manage treasures proposals from a JSON file.
type ManageTreasuresProposalJSON struct {
	Title       string           `json:"title" yaml:"title"`
	Description string           `json:"description" yaml:"description"`
	Treasures   []types.Treasure `json:"treasures" yaml:"treasures"`
	IsAdded     bool             `json:"is_added" yaml:"is_added"`
	Deposit     sdk.SysCoins     `json:"deposit" yaml:"deposit"`
}

// ParseManageTreasuresProposalJSON parses json from proposal file to ManageTreasuresProposalJSON struct
func ParseManageTreasuresProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageTreasuresProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}

// ExtraProposalJSON defines a ExtraProposal with a deposit used to parse
// manage treasures proposals from a JSON file.
type ExtraProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
	Action      string       `json:"action" yaml:"action"`
	Extra       string       `json:"extra" yaml:"extra"`
}

// ParseExtraProposalJSON parses json from proposal file to ExtraProposalJSON struct
func ParseExtraProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ExtraProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
