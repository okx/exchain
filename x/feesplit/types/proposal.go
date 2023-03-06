package types

import (
	"fmt"
	"github.com/okx/okbchain/libs/system"
	"strings"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

const (
	// proposalTypeFeeSplitShares defines the type for a FeeSplitProposalShares
	proposalTypeFeeSplitShares = "FeeSplit"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeFeeSplitShares)
	govtypes.RegisterProposalTypeCodec(FeeSplitSharesProposal{}, system.Chain+"/feesplit/SharesProposal")
}

var (
	_ govtypes.Content = (*FeeSplitSharesProposal)(nil)
)

type FeeSplitSharesProposal struct {
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description" yaml:"description"`
	Shares      []Shares `json:"shares" yaml:"shares"`
}

type Shares struct {
	ContractAddr string  `json:"contract_addr" yaml:"contract_addr"`
	Share        sdk.Dec `json:"share" yaml:"share"`
}

func NewFeeSplitSharesProposal(title, description string, shares []Shares) FeeSplitSharesProposal {
	return FeeSplitSharesProposal{title, description, shares}
}

func (fp FeeSplitSharesProposal) GetTitle() string       { return fp.Title }
func (fp FeeSplitSharesProposal) GetDescription() string { return fp.Description }
func (fp FeeSplitSharesProposal) ProposalRoute() string  { return RouterKey }
func (fp FeeSplitSharesProposal) ProposalType() string   { return proposalTypeFeeSplitShares }
func (fp FeeSplitSharesProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(fp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(fp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is longer than the max")
	}

	if len(fp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(fp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is longer than the max")
	}

	if fp.ProposalType() != proposalTypeFeeSplitShares {
		return govtypes.ErrInvalidProposalType(fp.ProposalType())
	}

	if len(fp.Shares) == 0 {
		return govtypes.ErrInvalidProposalContent("fee split shares is required")
	}

	for _, share := range fp.Shares {
		if err := ValidateNonZeroAddress(share.ContractAddr); err != nil {
			return govtypes.ErrInvalidProposalContent("invalid contract address")
		}

		if share.Share.IsNil() || share.Share.IsNegative() || share.Share.GT(sdk.OneDec()) {
			return govtypes.ErrInvalidProposalContent("invalid share: nil, negative, greater than 1")
		}
	}

	return nil
}

func (fp FeeSplitSharesProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Fee Split Shares Proposal:
  Title:       %s
  Description: %s
  Shares:      
`, fp.Title, fp.Description))

	for _, share := range fp.Shares {
		b.WriteString("\t\t\t\t\t\t")
		b.WriteString(fmt.Sprintf("%s: %s", share.ContractAddr, share.Share))
		b.Write([]byte{'\n'})
	}

	return b.String()
}
