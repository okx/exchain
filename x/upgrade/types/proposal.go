package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"
	govTypes "github.com/okex/okchain/x/gov/types"
)

// ProposalTypeDexList defines the type for a DexList proposal
const (
	ProposalAppUpgrade = "AppUpgrade"
)

func init() {
	govTypes.RegisterProposalType(ProposalAppUpgrade)
	govTypes.RegisterProposalTypeCodec(AppUpgradeProposal{}, "okchain/upgrade/AppUpgradeProposal")
}

type AppUpgradeProposal struct {
	Title              string                   `json:"title" yaml:"titile"`
	Description        string                   `json:"description" yaml:"description"`
	ProtocolDefinition proto.ProtocolDefinition `json:"protocol_definition" yaml:"protocol_definition"`
}

func NewAppUpgradeProposal(title, desc string, definition proto.ProtocolDefinition) *AppUpgradeProposal {
	proposal := &AppUpgradeProposal{
		Title:       title,
		Description: desc,
		ProtocolDefinition: proto.ProtocolDefinition{
			Version:   definition.Version,
			Software:  definition.Software,
			Height:    definition.Height,
			Threshold: definition.Threshold,
		},
	}

	return proposal
}

func (apu AppUpgradeProposal) GetProtocolDefinition() proto.ProtocolDefinition {
	return apu.ProtocolDefinition
}

func (apu *AppUpgradeProposal) SetProtocolDefinition(protocolDefinition proto.ProtocolDefinition) {
	apu.ProtocolDefinition = protocolDefinition
}

// GetTitle returns the title of a dex list proposal.
func (apu AppUpgradeProposal) GetTitle() string { return apu.Title }

// GetDescription returns the description of a dex list proposal.
func (apu AppUpgradeProposal) GetDescription() string { return apu.Description }

// ProposalRoute returns the routing key of a dex list proposal.
func (apu AppUpgradeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a dex list proposal.
func (apu AppUpgradeProposal) ProposalType() string { return ProposalAppUpgrade }

// ValidateBasic validates the parameter change proposal
func (apu AppUpgradeProposal) ValidateBasic() sdk.Error {
	err := govTypes.ValidateAbstract(DefaultCodespace, apu)
	if err != nil {
		return err
	}

	if len(apu.ProtocolDefinition.Software) > 70 || len(apu.ProtocolDefinition.Software) == 0 {
		return ErrInvalidLength(DefaultCodespace, "software description",
			len(apu.ProtocolDefinition.Software), 70)
	}

	if apu.ProtocolDefinition.Height == 0 {
		return ErrZeroSwitchHeight(DefaultCodespace)
	}

	// if threshold not in [0.75,1), then print error
	if apu.ProtocolDefinition.Threshold.LT(sdk.NewDecWithPrec(75, 2)) || apu.ProtocolDefinition.Threshold.GTE(sdk.NewDec(1)) {
		return ErrInvalidUpgradeThreshold(DefaultCodespace, apu.ProtocolDefinition.Threshold)
	}

	return nil
}

func (apu AppUpgradeProposal) String() string {
	return fmt.Sprintf(`Proposal:
  Title:              %s
  Type:               %s
  Version:            %d
  Software:           %s
  Switch Height:      %d
  Threshold:          %s`, apu.Title, apu.ProposalType(), apu.ProtocolDefinition.Version,
		apu.ProtocolDefinition.Software, apu.ProtocolDefinition.Height, apu.ProtocolDefinition.Threshold.String())
}
