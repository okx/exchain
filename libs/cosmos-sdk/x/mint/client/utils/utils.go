package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/pkg/errors"
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

	defer parseRecover(contents, &err)

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}

func parseRecover(contents []byte, err *error) {
	if r := recover(); r != nil {
		*err = errors.New(fmt.Sprintf("Please check the file:\n%s\nFailed to parse the proposal json:%s",
			string(contents), r))
	}
}

// ModifyNextBlockUpdateProposalJSON defines a ModifyNextBlockUpdateProposal with a deposit used to parse
// manage treasures proposals from a JSON file.
type ModifyNextBlockUpdateProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
	BlockNum    uint64       `json:"block_num" yaml:"block_num"`
}

// ParseModifyNextBlockUpdateProposalJSON parses json from proposal file to ModifyNextBlockUpdateProposalJSON struct
func ParseModifyNextBlockUpdateProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ModifyNextBlockUpdateProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
