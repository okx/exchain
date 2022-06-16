package types

import (
	"fmt"
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	govtypes "github.com/okex/exchain/x/gov/types"
)

const (
	// ProposalTypeCommunityPoolSpend defines the type for a CommunityPoolSpendProposal
	ProposalTypeCommunityPoolSpend = "CommunityPoolSpend"
	// ProposalTypeChangeDistributionModel defines the type for a ChangeDistributionTypeProposal
	ProposalTypeChangeDistributionModel = "ChangeDistributionModel"
)

const (
	DistributionTypeOffChain uint32 = 0
	DistributionTypeOnChain  uint32 = 1
)

// Assert CommunityPoolSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = CommunityPoolSpendProposal{}

// Assert ChangeDistributionTypeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = ChangeDistributionTypeProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeCommunityPoolSpend)
	govtypes.RegisterProposalType(ProposalTypeChangeDistributionModel)
	govtypes.RegisterProposalTypeCodec(CommunityPoolSpendProposal{}, "okexchain/distribution/CommunityPoolSpendProposal")
	govtypes.RegisterProposalTypeCodec(ChangeDistributionTypeProposal{}, "okexchain/distribution/ChangeDistributionTypeProposal")
}

// CommunityPoolSpendProposal spends from the community pool
type CommunityPoolSpendProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Recipient   sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Amount      sdk.SysCoins   `json:"amount" yaml:"amount"`
}

// NewCommunityPoolSpendProposal creates a new community pool spned proposal.
func NewCommunityPoolSpendProposal(title, description string, recipient sdk.AccAddress, amount sdk.SysCoins) CommunityPoolSpendProposal {
	return CommunityPoolSpendProposal{title, description, recipient, amount}
}

// GetTitle returns the title of a community pool spend proposal.
func (csp CommunityPoolSpendProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (csp CommunityPoolSpendProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp CommunityPoolSpendProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp CommunityPoolSpendProposal) ProposalType() string { return ProposalTypeCommunityPoolSpend }

// ValidateBasic runs basic stateless validity checks
func (csp CommunityPoolSpendProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, csp)
	if err != nil {
		return err
	}
	if !csp.Amount.IsValid() {
		return ErrInvalidProposalAmount()
	}
	if csp.Recipient.Empty() {
		return ErrEmptyProposalRecipient()
	}
	return nil
}

// String implements the Stringer interface.
func (csp CommunityPoolSpendProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Recipient:   %s
  Amount:      %s
`, csp.Title, csp.Description, csp.Recipient, csp.Amount))
	return b.String()
}

// ChangeDistributionTypeProposal
type ChangeDistributionTypeProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Type        uint32 `json:"type" yaml:"type"`
}

// NewChangeDistributionModelProposal creates a new change distribution model proposal.
func NewChangeDistributionModelProposal(title, description string, model uint32) ChangeDistributionTypeProposal {
	return ChangeDistributionTypeProposal{title, description, model}
}

// GetTitle returns the title of a change distribution model proposal.
func (cdmp ChangeDistributionTypeProposal) GetTitle() string { return cdmp.Title }

// GetDescription returns the description of a change distribution model proposal.
func (cdmp ChangeDistributionTypeProposal) GetDescription() string { return cdmp.Description }

// GetDescription returns the routing key of a change distribution model proposal.
func (cdmp ChangeDistributionTypeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a change distribution model proposal.
func (cdmp ChangeDistributionTypeProposal) ProposalType() string {
	return ProposalTypeChangeDistributionModel
}

// ValidateBasic runs basic stateless validity checks
func (cdmp ChangeDistributionTypeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, cdmp)
	if err != nil {
		return err
	}
	if cdmp.Type != DistributionTypeOffChain && cdmp.Type != DistributionTypeOnChain {
		return ErrInvalidDistributionModelType()
	}

	return nil
}

// String implements the Stringer interface.
func (cdmp ChangeDistributionTypeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Type:   %d
`, cdmp.Title, cdmp.Description, cdmp.Type))
	return b.String()
}
