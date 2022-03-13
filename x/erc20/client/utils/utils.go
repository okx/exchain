package utils

import (
	"io/ioutil"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

type (
	// TokenMappingProposalJSON defines a TokenMappingProposalJSON proposals from a JSON file.
	TokenMappingProposalJSON struct {
		Title       string `json:"title" yaml:"title"`
		Description string `json:"description" yaml:"description"`
		Denom       string `json:"denom" yaml:"denom"`
		Contract    string `json:"contract" yaml:"contract"`
	}
)

// ParseTokenMappingProposalJSON parses json from proposal file to ParseTokenMappingProposalJSON
// struct
func ParseTokenMappingProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal TokenMappingProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
