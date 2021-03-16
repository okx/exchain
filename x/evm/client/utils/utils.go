package utils

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ManageContractDeploymentWhitelistProposalJSON defines a ManageContractDeploymentWhitelistProposal with a deposit used
// to parse manage white list proposals from a JSON file.
type ManageContractDeploymentWhitelistProposalJSON struct {
	Title        string       `json:"title" yaml:"title"`
	Description  string       `json:"description" yaml:"description"`
	DistributorAddr string       `json:"deployer_address" yaml:"deployer_address"`
	IsAdded      bool         `json:"is_added" yaml:"is_added"`
	Deposit      sdk.SysCoins `json:"deposit" yaml:"deposit"`
}

// ParseManageContractDeploymentWhitelistProposalJSON parse json from proposal file to ManageContractDeploymentWhitelistProposalJSON
// struct
func ParseManageContractDeploymentWhitelistProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageContractDeploymentWhitelistProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
