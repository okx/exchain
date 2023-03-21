package types

import (
	"fmt"
	"github.com/okx/okbchain/libs/system"
	"strings"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkgovtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/gov/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

type UpgradeStatus uint32

const (
	ProposalTypeUpgrade = "OKBCUpgrade"
	UpgradeRouterKey    = "okbcUpgrade"

	QueryUpgrade = "okbcUpgrade"

	maxNameLength = 140

	UpgradeStatusPreparing        = 0
	UpgradeStatusWaitingEffective = 1
	UpgradeStatusEffective        = 2
)

// Assert ParameterChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = UpgradeProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpgrade)
	govtypes.RegisterProposalTypeCodec(UpgradeProposal{}, system.Chain+"/params/UpgradeProposal")
}

// UpgradeProposal is the struct of param change proposal
type UpgradeProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Name         string `json:"name" yaml:"name"`
	ExpectHeight uint64 `json:"expectHeight" yaml:"expectHeight"`
	Config       string `json:"config,omitempty" yaml:"config,omitempty"`
}

type UpgradeInfo struct {
	Name         string `json:"name" yaml:"name"`
	ExpectHeight uint64 `json:"expectHeight" yaml:"expectHeight"`
	Config       string `json:"config,omitempty" yaml:"config,omitempty"`

	EffectiveHeight uint64        `json:"effectiveHeight,omitempty" yaml:"effectiveHeight,omitempty"`
	Status          UpgradeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

func NewUpgradeProposal(title, description, name string, expectHeight uint64, config string) UpgradeProposal {
	return UpgradeProposal{
		Title:       title,
		Description: description,

		Name:         name,
		ExpectHeight: expectHeight,
		Config:       config,
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
  Config:      %s
`, up.Title, up.Description, up.Name, up.ExpectHeight, up.Config))

	return b.String()
}
