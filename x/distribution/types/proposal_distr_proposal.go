package types

import (
	"fmt"
	"strings"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

const (
	// ProposalTypeChangeDistributionType defines the type for a ChangeDistributionTypeProposal
	ProposalTypeChangeDistributionType = "ChangeDistributionType"

	// ProposalTypeWithdrawRewardEnabled defines the type for a WithdrawRewardEnabledProposal
	ProposalTypeWithdrawRewardEnabled = "WithdrawRewardEnabled"

	// ProposalTypeRewardTruncatePrecision defines the type for a RewardTruncatePrecision
	ProposalTypeRewardTruncatePrecision = "RewardTruncatePrecision"
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

// Assert RewardTruncatePrecisionProposal implements govtypes.Content at compile-time
var _ govtypes.Content = RewardTruncatePrecisionProposal{}

// RewardTruncatePrecisionProposal
type RewardTruncatePrecisionProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Precision   int64  `json:"precision" yaml:"precision"`
}

// NewRewardTruncatePrecisionProposal creates a reward truncate precision proposal.
func NewRewardTruncatePrecisionProposal(title, description string, precision int64) RewardTruncatePrecisionProposal {
	return RewardTruncatePrecisionProposal{title, description, precision}
}

// GetTitle returns the title of a set withdraw reward enabled proposal.
func (proposal RewardTruncatePrecisionProposal) GetTitle() string { return proposal.Title }

// GetDescription returns the description of a set withdraw reward enabled proposal.
func (proposal RewardTruncatePrecisionProposal) GetDescription() string { return proposal.Description }

// GetDescription returns the routing key of a set withdraw reward enabled proposal.
func (proposal RewardTruncatePrecisionProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a set withdraw reward enabled proposal.
func (proposal RewardTruncatePrecisionProposal) ProposalType() string {
	return ProposalTypeRewardTruncatePrecision
}

// ValidateBasic runs basic stateless validity checks
func (proposal RewardTruncatePrecisionProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, proposal)
	if err != nil {
		return err
	}

	if proposal.Precision < 0 || proposal.Precision > sdk.Precision {
		return ErrCodeRewardTruncatePrecision()
	}

	return nil
}

// String implements the Stringer interface.
func (proposal RewardTruncatePrecisionProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Reward Truncate Precision Proposal:
  Title:       %s
  Description: %s
  Precision:   %d
`, proposal.Title, proposal.Description, proposal.Precision))
	return b.String()
}
