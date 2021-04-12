package utils

import (
	"github.com/okex/exchain/x/evm/types"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// ManageContractDeploymentWhitelistProposalJSON defines a ManageContractDeploymentWhitelistProposal with a deposit used
	// to parse manage whitelist proposals from a JSON file.
	ManageContractDeploymentWhitelistProposalJSON struct {
		Title            string            `json:"title" yaml:"title"`
		Description      string            `json:"description" yaml:"description"`
		DistributorAddrs types.AddressList `json:"distributor_addresses" yaml:"distributor_addresses"`
		IsAdded          bool              `json:"is_added" yaml:"is_added"`
		Deposit          sdk.SysCoins      `json:"deposit" yaml:"deposit"`
	}
	// ManageContractBlockedListProposalJSON defines a ManageContractBlockedListProposal with a deposit used to parse
	// manage blocked list proposals from a JSON file.
	ManageContractBlockedListProposalJSON struct {
		Title         string            `json:"title" yaml:"title"`
		Description   string            `json:"description" yaml:"description"`
		ContractAddrs types.AddressList `json:"contract_addresses" yaml:"contract_addresses"`
		IsAdded       bool              `json:"is_added" yaml:"is_added"`
		Deposit       sdk.SysCoins      `json:"deposit" yaml:"deposit"`
	}
)

// ParseManageContractDeploymentWhitelistProposalJSON parses json from proposal file to ManageContractDeploymentWhitelistProposalJSON
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

// ParseManageContractBlockedListProposalJSON parses json from proposal file to ManageContractBlockedListProposalJSON struct
func ParseManageContractBlockedListProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageContractBlockedListProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
