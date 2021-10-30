package types

import (
	"fmt"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/x/params/types"

	govtypes "github.com/okex/exchain/x/gov/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkparams "github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

// Assert ParameterChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = ParameterChangeProposal{}

func init() {
	govtypes.RegisterProposalType(sdkparams.ProposalTypeChange)
	govtypes.RegisterProposalTypeCodec(ParameterChangeProposal{}, "okexchain/params/ParameterChangeProposal")
}

// ParameterChangeProposal is the struct of param change proposal
type ParameterChangeProposal struct {
	sdkparams.ParameterChangeProposal
	Height uint64 `json:"height" yaml:"height"`
}

// NewParameterChangeProposal creates a new instance of ParameterChangeProposal
func NewParameterChangeProposal(title, description string, changes []types.ParamChange, height uint64,
) ParameterChangeProposal {
	return ParameterChangeProposal{
		ParameterChangeProposal: sdkparams.NewParameterChangeProposal(title, description, changes),
		Height:                  height,
	}
}

// ValidateBasic validates the parameter change proposal
func (pcp ParameterChangeProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(pcp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(pcp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is longer than max title length")
	}

	if len(pcp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(pcp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is longer than max DggescriptionLength")
	}

	if pcp.ProposalType() != sdkparams.ProposalTypeChange {
		return govtypes.ErrInvalidProposalType(pcp.ProposalType())
	}

	if len(pcp.Changes) != 1 {
		return ErrInvalidParamsNum(DefaultCodespace, fmt.Sprintf("one proposal can only change one pair of parameter"))
	}

	return sdkparams.ValidateChanges(pcp.Changes)
}
