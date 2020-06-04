package utils

import (
	"io/ioutil"

	"github.com/okex/okchain/x/token/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CertifiedTokenProposalJSON defines a CertifiedTokenProposal with a deposit used
// to parse CertifiedToken proposals from a JSON file.
type CertifiedTokenProposalJSON struct {
	Title       string               `json:"title" yaml:"title"`
	Description string               `json:"description" yaml:"description"`
	Proposer    sdk.AccAddress       `json:"proposer" yaml:"proposer"`
	Token       types.CertifiedToken `json:"token" yaml:"token"`
	Deposit     sdk.DecCoins         `json:"deposit" yaml:"deposit"`
}

// ParseCertifiedTokenProposalJSON parse json from proposal file to CertifiedTokenProposalJSON struct
func ParseCertifiedTokenProposalJSON(cdc *codec.Codec, proposalFilePath string) (proposal CertifiedTokenProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
