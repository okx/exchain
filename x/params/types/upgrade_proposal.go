package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	sdkgovtypes "github.com/okex/exchain/libs/cosmos-sdk/x/gov/types"
	govtypes "github.com/okex/exchain/x/gov/types"

	"strings"
)

const (
	ProposalTypeUpgrade = "oKCUpgrade"
	UpgradeRouterKey    = "okcUpgrade"

	QueryUpgrade = "okcUpgrade"

	maxNameLength = 140
)

// Assert ParameterChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = UpgradeProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpgrade)
	govtypes.RegisterProposalTypeCodec(UpgradeProposal{}, "okexchain/params/UpgradeProposal")
}

// UpgradeProposal is the struct of param change proposal
type UpgradeProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"Description" yaml:"description"`

	UpgradeInfo
}

type UpgradeInfo struct {
	Name         string            `json:"name" yaml:"name"`
	ExpectHeight uint64            `json:"expectHeight" yaml:"expectHeight"`
	Config       map[string]string `json:"config,omitempty" yaml:"config,omitempty"`

	// only used in store
	EffectiveHeight uint64 `json:"effectiveHeight,omitempty" yaml:"effectiveHeight,omitempty"`
}

func NewUpgradeProposal(title, description, name string, expectHeight uint64, config map[string]string) UpgradeProposal {
	return UpgradeProposal{
		Title:       title,
		Description: description,
		UpgradeInfo: UpgradeInfo{
			Name:         name,
			ExpectHeight: expectHeight,
			Config:       config,

			EffectiveHeight: 0,
		},
	}
}

func (up UpgradeProposal) GetTitle() string {
	return up.Title
}

func (up UpgradeProposal) GetDescription() string {
	return up.Description
}

func (up UpgradeProposal) ProposalRoute() string {
	return UpgradeRouterKey
}

func (up UpgradeProposal) ProposalType() string {
	return ProposalTypeUpgrade
}

func (up UpgradeProposal) ValidateBasic() sdk.Error {
	if err := sdkgovtypes.ValidateAbstract(up); err != nil {
		return err
	}

	if up.ProposalType() != ProposalTypeUpgrade {
		return govtypes.ErrInvalidProposalType(up.ProposalType())
	}

	if len(strings.TrimSpace(up.Name)) == 0 {
		return govtypes.ErrInvalidProposalContent("name is required")
	}
	if len(up.Name) == maxNameLength {
		return govtypes.ErrInvalidProposalContent("name length is longer than max name length")
	}

	return nil
}

// String implements the Stringer interface.
func (up UpgradeProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Upgrade Proposal:
  Title:       %s
  Description: %s
  Name:        %s
  Height:      %d
  Config:
`, up.Title, up.Description, up.Name, up.ExpectHeight))

	for k, v := range up.Config {
		b.WriteString(fmt.Sprintf("    %s:%s\n", k, v))
	}

	return b.String()
}
