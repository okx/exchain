package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/global"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	govtypes "github.com/okex/exchain/x/gov/types"
	"strings"
)

const (
	// proposalTypeManageTreasures defines the type for a ManageContractTreasures
	proposalTypeManageTreasures = "ManageTreasures"

	// proposalTypeModifyNextBlockUpdate defines the type for a ModifyNextBlockUpdate
	proposalTypeModifyNextBlockUpdate = "ModifyNextBlockUpdate"

	// RouterKey uses module name for routing
	RouterKey = ModuleName
)

func init() {
	govtypes.RegisterProposalType(proposalTypeManageTreasures)
	govtypes.RegisterProposalType(proposalTypeModifyNextBlockUpdate)
	govtypes.RegisterProposalTypeCodec(ManageTreasuresProposal{}, "okexchain/mint/ManageTreasuresProposal")
	govtypes.RegisterProposalTypeCodec(ModifyNextBlockUpdateProposal{}, "okexchain/mint/ModifyNextBlockUpdateProposal")
}

var (
	_ govtypes.Content = (*ManageTreasuresProposal)(nil)
	_ govtypes.Content = (*ModifyNextBlockUpdateProposal)(nil)
)

// ManageTreasuresProposal - structure for the proposal to add or delete treasures
type ManageTreasuresProposal struct {
	Title       string     `json:"title" yaml:"title"`
	Description string     `json:"description" yaml:"description"`
	Treasures   []Treasure `json:"treasures" yaml:"treasures"`
	IsAdded     bool       `json:"is_added" yaml:"is_added"`
}

// NewManageTreasuresProposal creates a new instance of ManageTreasuresProposal
func NewManageTreasuresProposal(title, description string, treasures []Treasure, isAdded bool,
) ManageTreasuresProposal {
	return ManageTreasuresProposal{
		Title:       title,
		Description: description,
		Treasures:   treasures,
		IsAdded:     isAdded,
	}
}

// GetTitle returns title of a manage treasures proposal object
func (mp ManageTreasuresProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of a manage treasures proposal object
func (mp ManageTreasuresProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of a manage treasures proposal object
func (mp ManageTreasuresProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of a manage treasures proposal object
func (mp ManageTreasuresProposal) ProposalType() string {
	return proposalTypeManageTreasures
}

// ValidateBasic validates a manage treasures proposal
func (mp ManageTreasuresProposal) ValidateBasic() sdk.Error {
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

	if mp.ProposalType() != proposalTypeManageTreasures {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	if treasuresLen := len(mp.Treasures); treasuresLen == 0 {
		return ErrEmptyTreasures
	}

	if isTreasureDuplicated(mp.Treasures) {
		return ErrDuplicatedTreasure
	}
	if err := ValidateTreasures(mp.Treasures); err != nil {
		return ErrTreasuresInternal(err)
	}
	return nil
}

// String returns a human readable string representation of a ManageTreasuresProposal
func (mp ManageTreasuresProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ManageTreasuresProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 IsAdded:				%t
 Treasures:
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.IsAdded),
	)

	for i := 0; i < len(mp.Treasures); i++ {
		builder.WriteString("\t\t\t\t\t\t")
		builder.WriteString("Address:")
		builder.WriteString(mp.Treasures[i].Address.String())
		builder.WriteString("\t\tProportion:")
		builder.WriteString(mp.Treasures[i].Proportion.String())
		builder.Write([]byte{'\n'})
	}

	return strings.TrimSpace(builder.String())
}

// ModifyNextBlockUpdateProposal - structure for the proposal modify next block update
type ModifyNextBlockUpdateProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	BlockNum    uint64 `json:"block_num" yaml:"block_num"`
}

// NewModifyNextBlockUpdateProposal creates a new instance of ModifyNextBlockUpdateProposal
func NewModifyNextBlockUpdateProposal(title, description string, blockNum uint64) ModifyNextBlockUpdateProposal {
	return ModifyNextBlockUpdateProposal{
		Title:       title,
		Description: description,
		BlockNum:    blockNum,
	}
}

// GetTitle returns title of the proposal object
func (mp ModifyNextBlockUpdateProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of proposal object
func (mp ModifyNextBlockUpdateProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of the proposal object
func (mp ModifyNextBlockUpdateProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of the proposal object
func (mp ModifyNextBlockUpdateProposal) ProposalType() string {
	return proposalTypeModifyNextBlockUpdate
}

// ValidateBasic validates the proposal
func (mp ModifyNextBlockUpdateProposal) ValidateBasic() sdk.Error {
	if global.GetGlobalHeight() > 0 && !tmtypes.HigherThanVenus5(global.GetGlobalHeight()) {
		return ErrNotReachedVenus5Height
	}

	if len(strings.TrimSpace(mp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(mp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is bigger than max title length")
	}

	if len(mp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(mp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is bigger than max description length")
	}

	if mp.ProposalType() != proposalTypeModifyNextBlockUpdate {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	return nil
}

// String returns a human readable string representation of a ModifyNextBlockUpdateProposal
func (mp ModifyNextBlockUpdateProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ModifyNextBlockUpdateProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 BlockNum:				%d
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.BlockNum),
	)

	return strings.TrimSpace(builder.String())
}
