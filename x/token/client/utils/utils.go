package utils

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	//"github.com/okex/exchain/x/token/internal/types"
	"io/ioutil"
)

// ModifyDefaultBondDenomProposalJSON defines a ManageTreasureProposal with a deposit used to parse
// manage treasures proposals from a JSON file.
type ModifyDefaultBondDenomProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	DenomName   string       `json:"denom_name" yaml:"denom_name"`
	Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
}

// ParseModifyDefaultBondDenomProposalJSON parses json from proposal file to ModifyDefaultBondDenomProposalJSON struct
func ParseModifyDefaultBondDenomProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ModifyDefaultBondDenomProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
