package types

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
	"strings"

	govtypes "github.com/okex/exchain/x/gov/types"
)

const (
	// ProposalTypeChangeDistributionType defines the type for a ChangeDistributionTypeProposal
	ProposalTypeChangeDistributionType = "ChangeDistributionType"

	// ProposalTypeWithdrawRewardEnabled defines the type for a WithdrawRewardEnabledProposal
	ProposalTypeWithdrawRewardEnabled = "WithdrawRewardEnabled"
)

const (
	DistributionTypeOffChain uint32 = 0
	DistributionTypeOnChain  uint32 = 1
)

// Assert ChangeDistributionTypeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = ChangeDistributionTypeProposal{}

// ChangeDistributionTypeProposal
type ChangeDistributionTypeProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Type        uint32 `json:"type" yaml:"type"`
}

// NewChangeDistributionTypeProposal creates a new change distribution type proposal.
func NewChangeDistributionTypeProposal(title, description string, distrType uint32) ChangeDistributionTypeProposal {
	return ChangeDistributionTypeProposal{title, description, distrType}
}

// GetTitle returns the title of a change distribution type proposal.
func (cdtp ChangeDistributionTypeProposal) GetTitle() string { return cdtp.Title }

// GetDescription returns the description of a change distribution type proposal.
func (cdtp ChangeDistributionTypeProposal) GetDescription() string { return cdtp.Description }

// GetDescription returns the routing key of a change distribution type proposal.
func (cdtp ChangeDistributionTypeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a change distribution type proposal.
func (cdtp ChangeDistributionTypeProposal) ProposalType() string {
	return ProposalTypeChangeDistributionType
}

// ValidateBasic runs basic stateless validity checks
func (cdtp ChangeDistributionTypeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, cdtp)
	if err != nil {
		return err
	}
	if cdtp.Type != DistributionTypeOffChain && cdtp.Type != DistributionTypeOnChain {
		return ErrInvalidDistributionType()
	}

	//will delete it after upgrade venus3
	if !types.HigherThanVenus2(global.GetGlobalHeight()) {
		return ErrCodeNotSupportDistributionProposal()
	}

	return nil
}

// String implements the Stringer interface.
func (cdtp ChangeDistributionTypeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Change Distribution Type Proposal:
  Title:       %s
  Description: %s
  Type:   %d
`, cdtp.Title, cdtp.Description, cdtp.Type))
	return b.String()
}

// Assert WithdrawRewardEnabledProposal implements govtypes.Content at compile-time
var _ govtypes.Content = WithdrawRewardEnabledProposal{}

// WithdrawRewardEnabledProposal
type WithdrawRewardEnabledProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
}

// NewWithdrawRewardEnabledProposal creates a set withdraw reward enabled proposal.
func NewWithdrawRewardEnabledProposal(title, description string, enable bool) WithdrawRewardEnabledProposal {
	return WithdrawRewardEnabledProposal{title, description, enable}
}

// GetTitle returns the title of a set withdraw reward enabled proposal.
func (proposal WithdrawRewardEnabledProposal) GetTitle() string { return proposal.Title }

// GetDescription returns the description of a set withdraw reward enabled proposal.
func (proposal WithdrawRewardEnabledProposal) GetDescription() string { return proposal.Description }

// GetDescription returns the routing key of a set withdraw reward enabled proposal.
func (proposal WithdrawRewardEnabledProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a set withdraw reward enabled proposal.
func (proposal WithdrawRewardEnabledProposal) ProposalType() string {
	return ProposalTypeWithdrawRewardEnabled
}

// ValidateBasic runs basic stateless validity checks
func (proposal WithdrawRewardEnabledProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, proposal)
	if err != nil {
		return err
	}

	//will delete it after upgrade venus3
	if !types.HigherThanVenus2(global.GetGlobalHeight()) {
		return ErrCodeNotSupportWithdrawRewardEnabledProposal()
	}

	return nil
}

// String implements the Stringer interface.
func (proposal WithdrawRewardEnabledProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Withdraw Reward Enabled Proposal:
  Title:       %s
  Description: %s
  Enabled:   %t
`, proposal.Title, proposal.Description, proposal.Enabled))
	return b.String()
}
