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

	// ProposalTypeDistrExtend defines the type for a DistrExtendProposal
	ProposalTypeDistrExtend = "DistrExtend"
)

// Assert CommunityPoolSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = CommunityPoolSpendProposal{}

// Assert DistrExtendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = DistrExtendProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeCommunityPoolSpend)
	govtypes.RegisterProposalType(ProposalTypeDistrExtend)
	govtypes.RegisterProposalType(ProposalTypeChangeDistributionType)
	govtypes.RegisterProposalType(ProposalTypeWithdrawRewardEnabled)
	govtypes.RegisterProposalType(ProposalTypeRewardTruncatePrecision)
	govtypes.RegisterProposalTypeCodec(CommunityPoolSpendProposal{}, "okexchain/distribution/CommunityPoolSpendProposal")
	govtypes.RegisterProposalTypeCodec(DistrExtendProposal{}, "okexchain/distribution/DistrExtendProposal")
	govtypes.RegisterProposalTypeCodec(ChangeDistributionTypeProposal{}, "okexchain/distribution/ChangeDistributionTypeProposal")
	govtypes.RegisterProposalTypeCodec(WithdrawRewardEnabledProposal{}, "okexchain/distribution/WithdrawRewardEnabledProposal")
	govtypes.RegisterProposalTypeCodec(RewardTruncatePrecisionProposal{}, "okexchain/distribution/RewardTruncatePrecisionProposal")
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

type DistrExtendProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Method      string `json:"method" yaml:"method"`
	Params      string `json:"params" yaml:"params"`
}

// NewDistrExtendProposal creates a new distr extend proposal.
func NewDistrExtendProposal(title, description, method, params string) DistrExtendProposal {
	return DistrExtendProposal{title, description, method, params}
}

// GetTitle returns the title of a community pool spend proposal.
func (p DistrExtendProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of a community pool spend proposal.
func (p DistrExtendProposal) GetDescription() string { return p.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (p DistrExtendProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (p DistrExtendProposal) ProposalType() string { return ProposalTypeDistrExtend }

// ValidateBasic runs basic stateless validity checks
func (p DistrExtendProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, p)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(p.Method)) == 0 {
		return govtypes.ErrInvalidProposalContent("method is required")
	}
	if len(p.Method) > govtypes.MaxExtendMethodLength {
		return govtypes.ErrInvalidProposalContent("method length is bigger than max length")
	}

	if len(strings.TrimSpace(p.Params)) == 0 {
		return govtypes.ErrInvalidProposalContent("method is required")
	}
	if len(p.Params) > govtypes.MaxExtendParamsLength {
		return govtypes.ErrInvalidProposalContent("params length is bigger than max length")
	}

	switch p.Method {
	case MethodHello:
		if test, err := NewTestExtend(p.Params); err != nil {
			return err
		} else {
			return test.ValidateBasic()
		}
	}

	return nil
}

// String implements the Stringer interface.
func (p DistrExtendProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Method:   %s
  Params:      %s
`, p.Title, p.Description, p.Method, p.Params))
	return b.String()
}
