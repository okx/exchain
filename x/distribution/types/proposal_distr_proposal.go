package types

import (
	"fmt"
	"strings"

	govtypes "github.com/okex/exchain/x/gov/types"
)

const (
	// ProposalTypeChangeDistributionType defines the type for a ChangeDistributionTypeProposal
	ProposalTypeChangeDistributionType = "ChangeDistributionType"
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
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Type:   %d
`, cdtp.Title, cdtp.Description, cdtp.Type))
	return b.String()
}
