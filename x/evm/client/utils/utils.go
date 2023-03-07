package utils

import (
	"github.com/okx/okbchain/x/evm/types"
	"io/ioutil"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
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
	// ManageContractMethodBlockedListProposalJSON defines a ManageContractMethodBlockedListProposal with a deposit used to parse
	// manage contract method blocked list proposals from a JSON file.
	ManageContractMethodBlockedListProposalJSON struct {
		Title        string                    `json:"title" yaml:"title"`
		Description  string                    `json:"description" yaml:"description"`
		ContractList types.BlockedContractList `json:"contract_addresses" yaml:"contract_addresses"`
		IsAdded      bool                      `json:"is_added" yaml:"is_added"`
		Deposit      sdk.SysCoins              `json:"deposit" yaml:"deposit"`
	}

	// ManageSysContractAddressProposalJSON defines a ManageSysContractAddressProposal with a deposit used to parse
	// manage system contract address proposals from a JSON file.
	ManageSysContractAddressProposalJSON struct {
		Title       string `json:"title" yaml:"title"`
		Description string `json:"description" yaml:"description"`
		// Contract Address
		ContractAddr sdk.AccAddress `json:"contract_address" yaml:"contract_address"`
		IsAdded      bool           `json:"is_added" yaml:"is_added"`
		Deposit      sdk.SysCoins   `json:"deposit" yaml:"deposit"`
	}

	ManageContractByteCodeProposalJSON struct {
		Title              string         `json:"title" yaml:"title"`
		Description        string         `json:"description" yaml:"description"`
		Contract           sdk.AccAddress `json:"contract" yaml:"contract"`
		SubstituteContract sdk.AccAddress `json:"substitute_contract" yaml:"substitute_contract"`
		Deposit            sdk.SysCoins   `json:"deposit" yaml:"deposit"`
	}

	ResponseBlockContract struct {
		Address      string                `json:"address" yaml:"address"`
		BlockMethods types.ContractMethods `json:"block_methods" yaml:"block_methods"`
	}

	ResponseSysContractAddress struct {
		Address string `json:"address" yaml:"address"`
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

// ParseManageContractMethodBlockedListProposalJSON parses json from proposal file to ManageContractBlockedListProposalJSON struct
func ParseManageContractMethodBlockedListProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageContractMethodBlockedListProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}

// ManageSysContractAddressProposalJSON parses json from proposal file to ManageSysContractAddressProposal struct
func ParseManageSysContractAddressProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageSysContractAddressProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}

// ParseManageContractBytecodeProposalJSON parses json from proposal file to ManageContractByteCodeProposalJSON struct
func ParseManageContractBytecodeProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageContractByteCodeProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
