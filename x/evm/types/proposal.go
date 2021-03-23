package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
)

const (
	maxAddressListLength = 100
	// proposalTypeManageContractDeploymentWhitelist defines the type for a ManageContractDeploymentWhitelist
	proposalTypeManageContractDeploymentWhitelist = "ManageContractDeploymentWhitelist"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeManageContractDeploymentWhitelist)
	govtypes.RegisterProposalTypeCodec(ManageContractDeploymentWhitelistProposal{}, "okexchain/evm/ManageContractDeploymentWhitelistProposal")
}

var _ govtypes.Content = (*ManageContractDeploymentWhitelistProposal)(nil)

// ManageContractDeploymentWhitelistProposal - structure for the proposal to add or delete deployer addresses from whitelist
type ManageContractDeploymentWhitelistProposal struct {
	Title            string           `json:"title" yaml:"title"`
	Description      string           `json:"description" yaml:"description"`
	DistributorAddrs []sdk.AccAddress `json:"distributor_addresses" yaml:"distributor_addresses"`
	IsAdded          bool             `json:"is_added" yaml:"is_added"`
}

// NewManageContractDeploymentWhitelistProposal creates a new instance of ManageContractDeploymentWhitelistProposal
func NewManageContractDeploymentWhitelistProposal(title, description string, distributorAddrs []sdk.AccAddress, isAdded bool,
) ManageContractDeploymentWhitelistProposal {
	return ManageContractDeploymentWhitelistProposal{
		Title:            title,
		Description:      description,
		DistributorAddrs: distributorAddrs,
		IsAdded:          isAdded,
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

	distributorAddrLen := len(mp.DistributorAddrs)
	if distributorAddrLen == 0 {
		return ErrEmptyAddressList
	}

	if distributorAddrLen > maxAddressListLength {
		return ErrOversizeAddrList(distributorAddrLen)
	}

	if isAddrDuplicated(mp.DistributorAddrs) {
		return ErrDuplicatedAddr
	}

	return nil
}

// String returns a human readable string representation of a ManageContractDeploymentWhitelistProposal
func (mp ManageContractDeploymentWhitelistProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ManageContractDeploymentWhitelistProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 IsAdded:				%t
 DistributorAddrs:
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.IsAdded),
	)

	for i := 0; i < len(mp.DistributorAddrs); i++ {
		builder.WriteString("\t\t\t\t\t\t")
		builder.WriteString(mp.DistributorAddrs[i].String())
		builder.Write([]byte{'\n'})
	}

	return strings.TrimSpace(builder.String())
}

func isAddrDuplicated(addrs []sdk.AccAddress) bool {
	lenAddrs := len(addrs)
	filter := make(map[string]struct{}, lenAddrs)
	for i := 0; i < lenAddrs; i++ {
		key := addrs[i].String()
		if _, ok := filter[key]; ok {
			return true
		}
		filter[key] = struct{}{}
	}

	return false
}
