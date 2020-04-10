package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/okex/okchain/x/gov/types"
)

const (
	// ProposalTypeDelist defines the type for a Delist proposal
	ProposalTypeDelist = "Delist"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeDelist)
	govtypes.RegisterProposalTypeCodec(DelistProposal{}, "okchain/dex/DelistProposal")

}

// Assert DelistProposal implements govtypes.Content at compile-time
var _ govtypes.Content = (*DelistProposal)(nil)

type DelistProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	BaseAsset   string         `json:"base_asset" yaml:"base_asset"`
	QuoteAsset  string         `json:"quote_asset" yaml:"quote_asset"`
}

func NewDelistProposal(title, description string, proposer sdk.AccAddress, baseAsset, quoteAsset string) DelistProposal {
	return DelistProposal{
		Title:       title,
		Description: description,
		Proposer:    proposer,
		BaseAsset:   baseAsset,
		QuoteAsset:  quoteAsset,
	}
}
func (drp DelistProposal) GetTitle() string {
	return drp.Title
}

func (drp DelistProposal) GetDescription() string {
	return drp.Description
}

func (DelistProposal) ProposalRoute() string {
	return RouterKey
}

func (DelistProposal) ProposalType() string {
	return ProposalTypeDelist
}

func (drp DelistProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(drp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace, "failed to submit delist proposal because title is blank")
	}
	if len(drp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace, fmt.Sprintf("failed to submit delist proposal because title is longer than max length of %d", govtypes.MaxTitleLength))
	}

	if len(drp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace, "failed to submit delist proposal because  description is blank")
	}

	if len(drp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent(DefaultCodespace, fmt.Sprintf("failed to submit delist proposal because  description is longer than max length of %d", govtypes.MaxDescriptionLength))
	}

	if drp.ProposalType() != ProposalTypeDelist {
		return govtypes.ErrInvalidProposalType(DefaultCodespace, drp.ProposalType())
	}

	if drp.Proposer.Empty() {
		return sdk.ErrInvalidAddress(drp.Proposer.String())
	}

	if drp.BaseAsset == drp.QuoteAsset {
		return sdk.ErrInvalidCoins(fmt.Sprintf("failed to submit delist proposal because baseasset is same as quoteasset"))
	}

	return nil
}

func (drp DelistProposal) String() string {
	return fmt.Sprintf(`DelistProposal:
 Title:               %s
 Description:         %s
 Type:                %s
 Proposer:            %s
 ListAsset            %s
 QuoteAsset           %s
`, drp.Title, drp.Description,
		drp.ProposalType(), drp.Proposer,
		drp.BaseAsset, drp.QuoteAsset,
	)
}
