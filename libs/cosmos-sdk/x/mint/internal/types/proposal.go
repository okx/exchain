package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/system"
	"github.com/okx/okbchain/libs/tendermint/global"
	govtypes "github.com/okx/okbchain/x/gov/types"
	"strings"
)

const (
	// proposalTypeManageTreasures defines the type for a ManageContractTreasures
	proposalTypeManageTreasures = "ManageTreasures"

	// ProposalTypeExtra defines the type for a MintExtraProposal
	ProposalTypeExtra = "MintExtra"

	ActionNextBlockUpdate = "NextBlockUpdate"
	ActionMintedPerBlock  = "MintedPerBlock"

	// RouterKey uses module name for routing
	RouterKey = ModuleName
)

func init() {
	govtypes.RegisterProposalType(proposalTypeManageTreasures)
	govtypes.RegisterProposalType(ProposalTypeExtra)
	govtypes.RegisterProposalTypeCodec(ManageTreasuresProposal{}, system.Chain+"/mint/ManageTreasuresProposal")
	govtypes.RegisterProposalTypeCodec(ExtraProposal{}, system.Chain+"/mint/ExtraProposal")

}

var (
	_ govtypes.Content = (*ManageTreasuresProposal)(nil)
	_ govtypes.Content = (*ExtraProposal)(nil)
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

type ExtraProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Action      string `json:"action" yaml:"action"`
	Extra       string `json:"extra" yaml:"extra"`
}

// NewExtraProposal creates a new extra proposal.
func NewExtraProposal(title, description, action, extra string) ExtraProposal {
	return ExtraProposal{title, description, action, extra}
}

// GetTitle returns the title of a community pool spend proposal.
func (p ExtraProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of a community pool spend proposal.
func (p ExtraProposal) GetDescription() string { return p.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (p ExtraProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (p ExtraProposal) ProposalType() string { return ProposalTypeExtra }

// ValidateBasic runs basic stateless validity checks
func (p ExtraProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(ModuleName, p)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(p.Action)) == 0 {
		return govtypes.ErrInvalidProposalContent("extra proposal's action is required")
	}
	if len(p.Action) > govtypes.MaxExtraActionLength {
		return govtypes.ErrInvalidProposalContent("extra proposal's action length is bigger than max length")
	}
	if len(strings.TrimSpace(p.Extra)) == 0 {
		return govtypes.ErrInvalidProposalContent("extra proposal's extra is required")
	}
	if len(p.Extra) > govtypes.MaxExtraBodyLength {
		return govtypes.ErrInvalidProposalContent("extra proposal's extra body length is bigger than max length")
	}
	switch p.Action {
	case ActionNextBlockUpdate:
		_, err = NewNextBlockUpdate(p.Extra)
		return err
	case ActionMintedPerBlock:
		_, err = NewMintedPerBlockParams(p.Extra)
		return err
	default:
		return ErrUnknownExtraProposalAction
	}

	return nil
}

// String implements the Stringer interface.
func (p ExtraProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Action:   %s
  Extra:      %s
`, p.Title, p.Description, p.Action, p.Extra))
	return b.String()
}

type NextBlockUpdateParams struct {
	BlockNum uint64 `json:"block_num" yaml:"block_num"`
}

func NewNextBlockUpdate(data string) (NextBlockUpdateParams, error) {
	var param NextBlockUpdateParams
	err := json.Unmarshal([]byte(data), &param)
	if err != nil {
		return param, ErrExtraProposalParams("parse json error")
	}

	if global.GetGlobalHeight() > 0 && param.BlockNum <= uint64(global.GetGlobalHeight()) {
		return param, ErrCodeInvalidHeight
	}

	return param, nil
}

type MintedPerBlockParams struct {
	Coin sdk.SysCoin `json:"coin" yaml:"coin"` // minted per block on this proposal.
}

func NewMintedPerBlockParams(jsonData string) (MintedPerBlockParams, error) {
	var param MintedPerBlockParams
	err := json.Unmarshal([]byte(jsonData), &param)
	if err != nil {
		return param, ErrExtraProposalParams("parse json error")
	}

	if param.Coin.Amount.IsNil() || param.Coin.Denom != sdk.DefaultBondDenom {
		return param, ErrExtraProposalParams("coin is nil")
	}

	if param.Coin.IsNegative() {
		return param, ErrExtraProposalParams("coin is negative")
	}

	return param, nil
}
