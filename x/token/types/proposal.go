package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	govtypes "github.com/okex/exchain/x/gov/types"
	"strings"
)

const (
	// proposalTypeManageTreasures defines the type for a ManageContractTreasures
	proposalTypeModifyDefaultBondDenom = "ModifyDefaultBondDenom"

	// RouterKey uses module name for routing
	//RouterKey = ModuleName
)

func init() {
	govtypes.RegisterProposalType(proposalTypeModifyDefaultBondDenom)
	govtypes.RegisterProposalTypeCodec(ModifyDefaultBondDenomProposal{}, "okexchain/token/ModifyDefaultBondDenomProposal")
}

var (
	_ govtypes.Content = (*ModifyDefaultBondDenomProposal)(nil)
)

// ModifyDefaultBondDenomProposal - structure for the proposal to add or delete treasures
type ModifyDefaultBondDenomProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	DenomName   string `json:"denom_name" yaml:"denom_name"`
}

// NewModifyDefaultBondDenomProposal creates a new instance of ModifyDefaultBondDenomProposal
func NewModifyDefaultBondDenomProposal(title, description string, denomName string,
) ModifyDefaultBondDenomProposal {
	return ModifyDefaultBondDenomProposal{
		Title:       title,
		Description: description,
		DenomName:   denomName,
	}
}

// GetTitle returns title of a manage treasures proposal object
func (mp ModifyDefaultBondDenomProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of a manage treasures proposal object
func (mp ModifyDefaultBondDenomProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of a manage treasures proposal object
func (mp ModifyDefaultBondDenomProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of a manage treasures proposal object
func (mp ModifyDefaultBondDenomProposal) ProposalType() string {
	return proposalTypeModifyDefaultBondDenom
}

// ValidateBasic validates a manage treasures proposal
func (mp ModifyDefaultBondDenomProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(mp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(mp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is longer than the maximum title length")
	}

	if len(mp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(mp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is longer than the maximum description length")
	}

	if mp.ProposalType() != proposalTypeModifyDefaultBondDenom {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	return nil
}

// String returns a human readable string representation of a ModifyDefaultBondDenomProposal
func (mp ModifyDefaultBondDenomProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ModifyDefaultBondDenomProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 DenomName:				%s
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.DenomName),
	)

	return strings.TrimSpace(builder.String())
}
