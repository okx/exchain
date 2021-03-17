package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
)

const (
	// proposalTypeManageContractDeploymentWhitelist defines the type for a ManageContractDeploymentWhitelist
	proposalTypeManageContractDeploymentWhitelist = "ManageContractDeploymentWhitelist"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeManageContractDeploymentWhitelist)
	govtypes.RegisterProposalTypeCodec(ManageContractDeploymentWhitelistProposal{}, "okexchain/evm/ManageContractDeploymentWhitelistProposal")
}

var _ govtypes.Content = (*ManageContractDeploymentWhitelistProposal)(nil)

// ManageContractDeploymentWhitelistProposal - structure for the proposal to add or delete a deployer address from whitelist
type ManageContractDeploymentWhitelistProposal struct {
	Title           string         `json:"title" yaml:"title"`
	Description     string         `json:"description" yaml:"description"`
	DistributorAddr sdk.AccAddress `json:"distributor_address" yaml:"distributor_address"`
	IsAdded         bool           `json:"is_added" yaml:"is_added"`
}

// NewManageContractDeploymentWhitelistProposal creates a new instance of ManageContractDeploymentWhitelistProposal
func NewManageContractDeploymentWhitelistProposal(title, description string, distributorAddr sdk.AccAddress, isAdded bool,
) ManageContractDeploymentWhitelistProposal {
	return ManageContractDeploymentWhitelistProposal{
		Title:           title,
		Description:     description,
		DistributorAddr: distributorAddr,
		IsAdded:         isAdded,
	}
}

// GetTitle returns title of a manage contract deployment whitelist proposal object
func (mp ManageContractDeploymentWhitelistProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of a manage contract deployment whitelist proposal object
func (mp ManageContractDeploymentWhitelistProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of a manage contract deployment whitelist proposal object
func (mp ManageContractDeploymentWhitelistProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of a manage contract deployment whitelist proposal object
func (mp ManageContractDeploymentWhitelistProposal) ProposalType() string {
	return proposalTypeManageContractDeploymentWhitelist
}

// ValidateBasic validates a manage contract deployment whitelist proposal
func (mp ManageContractDeploymentWhitelistProposal) ValidateBasic() sdk.Error {
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

	if mp.ProposalType() != proposalTypeManageContractDeploymentWhitelist {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	if mp.DistributorAddr.Empty() {
		return ErrEmptyAddress
	}

	return nil
}

// String returns a human readable string representation of a ManageContractDeploymentWhitelistProposal
func (mp ManageContractDeploymentWhitelistProposal) String() string {
	return fmt.Sprintf(`ManageContractDeploymentWhitelistProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 DistributorAddr:		%s
 IsAdded:				%t`,
		mp.Title, mp.Description, mp.ProposalType(), mp.DistributorAddr.String(), mp.IsAdded)
}
