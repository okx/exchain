package cli

import (
	"io/ioutil"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type (
	// ChangeDistributionTypeProposalJSON defines a ChangeDistributionTypeProposal with a deposit
	ChangeDistributionTypeProposalJSON struct {
		Title       string       `json:"title" yaml:"title"`
		Description string       `json:"description" yaml:"description"`
		Type        uint32       `json:"type" yaml:"type"`
		Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
	}

	// WithdrawRewardEnabledProposalJSON defines a WithdrawRewardEnabledProposal with a deposit
	WithdrawRewardEnabledProposalJSON struct {
		Title       string       `json:"title" yaml:"title"`
		Description string       `json:"description" yaml:"description"`
		Enabled     bool         `json:"enabled" yaml:"enabled"`
		Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`
	}
)

// ParseChangeDistributionTypeProposalJSON reads and parses a ChangeDistributionTypeProposalJSON from a file.
func ParseChangeDistributionTypeProposalJSON(cdc *codec.Codec, proposalFile string) (ChangeDistributionTypeProposalJSON, error) {
	proposal := ChangeDistributionTypeProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

// ParseWithdrawRewardEnabledProposalJSON reads and parses a WithdrawRewardEnabledProposalJSON from a file.
func ParseWithdrawRewardEnabledProposalJSON(cdc *codec.Codec, proposalFile string) (WithdrawRewardEnabledProposalJSON, error) {
	proposal := WithdrawRewardEnabledProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
