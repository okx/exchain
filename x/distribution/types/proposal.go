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
	// ProposalTypeChangeDistributionModel defines the type for a ChangeDistributionModelProposal
	ProposalTypeChangeDistributionModel = "ChangeDistributionModel"
)

const (
	DistributionModelOffChain uint32 = 0
	DistributionModelOnChain  uint32 = 1
)

// Assert CommunityPoolSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = CommunityPoolSpendProposal{}

// Assert ChangeDistributionModelProposal implements govtypes.Content at compile-time
var _ govtypes.Content = ChangeDistributionModelProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeCommunityPoolSpend)
	govtypes.RegisterProposalType(ProposalTypeChangeDistributionModel)
	govtypes.RegisterProposalTypeCodec(CommunityPoolSpendProposal{}, "okexchain/distribution/CommunityPoolSpendProposal")
	govtypes.RegisterProposalTypeCodec(ChangeDistributionModelProposal{}, "okexchain/distribution/ChangeDistributionModelProposal")
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

// ChangeDistributionModelProposal
type ChangeDistributionModelProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Model       uint32 `json:"model" yaml:"model"`
}

// NewChangeDistributionModelProposal creates a new change distribution model proposal.
func NewChangeDistributionModelProposal(title, description string, model uint32) ChangeDistributionModelProposal {
	return ChangeDistributionModelProposal{title, description, model}
}

// GetTitle returns the title of a change distribution model proposal.
func (cdmp ChangeDistributionModelProposal) GetTitle() string { return cdmp.Title }

// GetDescription returns the description of a change distribution model proposal.
func (cdmp ChangeDistributionModelProposal) GetDescription() string { return cdmp.Description }

// GetDescription returns the routing key of a change distribution model proposal.
func (cdmp ChangeDistributionModelProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a change distribution model proposal.
func (cdmp ChangeDistributionModelProposal) ProposalType() string {
	return ProposalTypeChangeDistributionModel
}

// ValidateBasic runs basic stateless validity checks
func (cdmp ChangeDistributionModelProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, cdmp)
	if err != nil {
		return err
	}
	if cdmp.Model != DistributionModelOffChain && cdmp.Model != DistributionModelOnChain {
		return ErrInvalidDistributionModelType()
	}

	return nil
}

// String implements the Stringer interface.
func (cdmp ChangeDistributionModelProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Model:   %d
`, cdmp.Title, cdmp.Description, cdmp.Model))
	return b.String()
}
