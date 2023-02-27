package types

import (
	"fmt"
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	govtypes "github.com/okex/exchain/x/gov/types"
)

const (
	proposalTypeProposeValidator = "ProposeValidator"
	ProposeValidatorProposalName = "okexchain/staking/ProposeValidatorProposal"
)

var _ govtypes.Content = (*ProposeValidatorProposal)(nil)

func init() {
	govtypes.RegisterProposalType(proposalTypeProposeValidator)
	govtypes.RegisterProposalTypeCodec(ProposeValidatorProposal{}, ProposeValidatorProposalName)
}

// ProposeValidatorProposal - structure for the proposal of proposing validator
type ProposeValidatorProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	IsAdd       bool   `json:"is_add" yaml:"is_add"`
	BlockNum    uint64 `json:"block_num" yaml:"block_num"`
}

// NewProposeValidatorProposal creates a new instance of ProposeValidatorProposal
func NewProposeValidatorProposal(title, description string, isAdd bool, blockNum uint64) ProposeValidatorProposal {
	return ProposeValidatorProposal{
		Title:       title,
		Description: description,
		IsAdd:       isAdd,
		BlockNum:    blockNum,
	}
}

// GetTitle returns title of the proposal object
func (pv ProposeValidatorProposal) GetTitle() string {
	return pv.Title
}

// GetDescription returns description of proposal object
func (pv ProposeValidatorProposal) GetDescription() string {
	return pv.Description
}

// ProposalRoute returns route key of the proposal object
func (pv ProposeValidatorProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of the proposal object
func (pv ProposeValidatorProposal) ProposalType() string {
	return proposalTypeProposeValidator
}

// ValidateBasic validates the proposal
func (pv ProposeValidatorProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(pv.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(pv.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is bigger than max title length")
	}

	if len(pv.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(pv.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is bigger than max description length")
	}

	if pv.ProposalType() != proposalTypeProposeValidator {
		return govtypes.ErrInvalidProposalType(pv.ProposalType())
	}

	return nil
}

// String returns a human readable string representation of a ProposeValidatorProposal
func (pv ProposeValidatorProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ProposeValidatorProposal:
 Title:					%s
 Description:        	%s
 IsADD:                 %t
 Type:                	%s
 BlockNum:				%d
`,
			pv.Title, pv.Description, pv.IsAdd, pv.ProposalType(), pv.BlockNum),
	)

	return strings.TrimSpace(builder.String())
}
