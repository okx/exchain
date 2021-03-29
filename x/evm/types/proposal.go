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
	// proposalTypeManageContractBlockedList defines the type for a ManageContractBlockedListProposal
	proposalTypeManageContractBlockedList = "ManageContractBlockedList"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeManageContractDeploymentWhitelist)
	govtypes.RegisterProposalType(proposalTypeManageContractBlockedList)
	govtypes.RegisterProposalTypeCodec(ManageContractDeploymentWhitelistProposal{}, "okexchain/evm/ManageContractDeploymentWhitelistProposal")
	govtypes.RegisterProposalTypeCodec(ManageContractBlockedListProposal{}, "okexchain/evm/ManageContractBlockedListProposal")
}

var (
	_ govtypes.Content = (*ManageContractDeploymentWhitelistProposal)(nil)
	_ govtypes.Content = (*ManageContractBlockedListProposal)(nil)
)

// ManageContractDeploymentWhitelistProposal - structure for the proposal to add or delete deployer addresses from whitelist
type ManageContractDeploymentWhitelistProposal struct {
	Title            string      `json:"title" yaml:"title"`
	Description      string      `json:"description" yaml:"description"`
	DistributorAddrs AddressList `json:"distributor_addresses" yaml:"distributor_addresses"`
	IsAdded          bool        `json:"is_added" yaml:"is_added"`
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

// ManageContractBlockedListProposal - structure for the proposal to add or delete a contract address from blocked list
type ManageContractBlockedListProposal struct {
	Title         string      `json:"title" yaml:"title"`
	Description   string      `json:"description" yaml:"description"`
	ContractAddrs AddressList `json:"contract_addresses" yaml:"contract_addresses"`
	IsAdded       bool        `json:"is_added" yaml:"is_added"`
}

// NewManageContractBlockedListProposal creates a new instance of ManageContractBlockedListProposal
func NewManageContractBlockedListProposal(title, description string, contractAddrs AddressList, isAdded bool,
) ManageContractBlockedListProposal {
	return ManageContractBlockedListProposal{
		Title:         title,
		Description:   description,
		ContractAddrs: contractAddrs,
		IsAdded:       isAdded,
	}
}

// GetTitle returns title of a manage contract blocked list proposal object
func (mp ManageContractBlockedListProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of a manage contract blocked list proposal object
func (mp ManageContractBlockedListProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of a manage contract blocked list proposal object
func (mp ManageContractBlockedListProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of a manage contract blocked list proposal object
func (mp ManageContractBlockedListProposal) ProposalType() string {
	return proposalTypeManageContractBlockedList
}

// ValidateBasic validates a manage contract blocked list proposal
func (mp ManageContractBlockedListProposal) ValidateBasic() sdk.Error {
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

	if mp.ProposalType() != proposalTypeManageContractBlockedList {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	contractAddrLen := len(mp.ContractAddrs)
	if contractAddrLen == 0 {
		return ErrEmptyAddressList
	}

	if contractAddrLen > maxAddressListLength {
		return ErrOversizeAddrList(contractAddrLen)
	}

	if isAddrDuplicated(mp.ContractAddrs) {
		return ErrDuplicatedAddr
	}

	return nil
}

// String returns a human readable string representation of a ManageContractBlockedListProposal
func (mp ManageContractBlockedListProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ManageContractBlockedListProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 IsAdded:				%t
 ContractAddrs:
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.IsAdded),
	)

	for i := 0; i < len(mp.ContractAddrs); i++ {
		builder.WriteString("\t\t\t\t\t\t")
		builder.WriteString(mp.ContractAddrs[i].String())
		builder.Write([]byte{'\n'})
	}

	return strings.TrimSpace(builder.String())
}
