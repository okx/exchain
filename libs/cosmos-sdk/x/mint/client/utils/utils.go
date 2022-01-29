package utils

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"io/ioutil"
)

// ManageContractMethodBlockedListProposalJSON defines a ManageContractMethodBlockedListProposal with a deposit used to parse
// manage contract method blocked list proposals from a JSON file.
type ManageTreasureProposalJSON struct {
	Title       string           `json:"title" yaml:"title"`
	Description string           `json:"description" yaml:"description"`
	Treasures   []types.Treasure `json:"treasures" yaml:"treasures"`
	IsAdded     bool             `json:"is_added" yaml:"is_added"`
	Deposit     sdk.SysCoins     `json:"deposit" yaml:"deposit"`
}

// ParseManageContractMethodBlockedListProposalJSON parses json from proposal file to ManageContractBlockedListProposalJSON struct
func ParseManageTreasureProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageTreasureProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
