package types

import (
	"fmt"
	"github.com/okx/okbchain/libs/system"
	"strings"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

const (
	maxAddressListLength = 100
	// proposalTypeManageContractDeploymentWhitelist defines the type for a ManageContractDeploymentWhitelist
	proposalTypeManageContractDeploymentWhitelist = "ManageContractDeploymentWhitelist"
	// proposalTypeManageContractBlockedList defines the type for a ManageContractBlockedListProposal
	proposalTypeManageContractBlockedList = "ManageContractBlockedList"
	// proposalTypeManageContractMethodBlockedList defines the type for a ManageContractMethodBlockedList
	proposalTypeManageContractMethodBlockedList = "ManageContractMethodBlockedList"
	// proposalTypeManageSysContractAddress defines the type for a ManageSysContractAddress
	proposalTypeManageSysContractAddress = "ManageSysContractAddress"
	proposalTypeManageContractByteCode   = "ManageContractByteCode"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeManageContractDeploymentWhitelist)
	govtypes.RegisterProposalType(proposalTypeManageContractBlockedList)
	govtypes.RegisterProposalType(proposalTypeManageContractMethodBlockedList)
	govtypes.RegisterProposalType(proposalTypeManageSysContractAddress)
	govtypes.RegisterProposalType(proposalTypeManageContractByteCode)
	govtypes.RegisterProposalTypeCodec(ManageContractDeploymentWhitelistProposal{}, system.Chain+"/evm/ManageContractDeploymentWhitelistProposal")
	govtypes.RegisterProposalTypeCodec(ManageContractBlockedListProposal{}, system.Chain+"/evm/ManageContractBlockedListProposal")
	govtypes.RegisterProposalTypeCodec(ManageContractMethodBlockedListProposal{}, system.Chain+"/evm/ManageContractMethodBlockedListProposal")
	govtypes.RegisterProposalTypeCodec(ManageSysContractAddressProposal{}, system.Chain+"/evm/ManageSysContractAddressProposal")
	govtypes.RegisterProposalTypeCodec(ManageContractByteCodeProposal{}, system.Chain+"/evm/ManageContractBytecode")
}

var (
	_ govtypes.Content = (*ManageContractDeploymentWhitelistProposal)(nil)
	_ govtypes.Content = (*ManageContractBlockedListProposal)(nil)
	_ govtypes.Content = (*ManageContractMethodBlockedListProposal)(nil)
	_ govtypes.Content = (*ManageSysContractAddressProposal)(nil)
	_ govtypes.Content = (*ManageContractByteCodeProposal)(nil)
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

// ManageContractMethodBlockedListProposal - structure for the proposal to add or delete a contract method from blocked list
type ManageContractMethodBlockedListProposal struct {
	Title        string              `json:"title" yaml:"title"`
	Description  string              `json:"description" yaml:"description"`
	ContractList BlockedContractList `json:"contract_addresses" yaml:"contract_addresses"`
	IsAdded      bool                `json:"is_added" yaml:"is_added"`
}

// NewManageContractMethodBlockedListProposal creates a new instance of ManageContractMethodBlockedListProposal
func NewManageContractMethodBlockedListProposal(title, description string, contractList BlockedContractList, isAdded bool,
) ManageContractMethodBlockedListProposal {
	return ManageContractMethodBlockedListProposal{
		Title:        title,
		Description:  description,
		ContractList: contractList,
		IsAdded:      isAdded,
	}
}

// GetTitle returns title of a manage contract blocked list proposal object
func (mp ManageContractMethodBlockedListProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of a manage contract blocked list proposal object
func (mp ManageContractMethodBlockedListProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of a manage contract blocked list proposal object
func (mp ManageContractMethodBlockedListProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of a manage contract blocked list proposal object
func (mp ManageContractMethodBlockedListProposal) ProposalType() string {
	return proposalTypeManageContractMethodBlockedList
}

// ValidateBasic validates a manage contract blocked list proposal
func (mp ManageContractMethodBlockedListProposal) ValidateBasic() sdk.Error {
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

	if mp.ProposalType() != proposalTypeManageContractMethodBlockedList {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	contractAddrLen := len(mp.ContractList)
	if contractAddrLen == 0 {
		return ErrEmptyAddressList
	}

	if contractAddrLen > maxAddressListLength {
		return ErrOversizeAddrList(contractAddrLen)
	}

	if err := mp.ContractList.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// String returns a human readable string representation of a ManageContractMethodBlockedListProposal
func (mp ManageContractMethodBlockedListProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ManageContractMethodBlockedListProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 IsAdded:				%t
 ContractList:
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.IsAdded),
	)

	for i := 0; i < len(mp.ContractList); i++ {
		builder.WriteString("\t\t\t\t\t\t")
		builder.WriteString(mp.ContractList[i].String())
		builder.Write([]byte{'\n'})
	}

	return strings.TrimSpace(builder.String())
}

// FixShortAddr is to fix the short address problem in the OKBC test-net.
// The normal len(BlockedContract.Address) should be 20,
// but there are some BlockedContract.Address in OKBC test-net that have a length of 4.
// The fix is to pad the leading bits of the short address with zeros until the length is 20.
func (mp *ManageContractMethodBlockedListProposal) FixShortAddr() {
	for i := 0; i < len(mp.ContractList); i++ {
		if len(mp.ContractList[i].Address) < 20 {
			validAddress := make([]byte, 20-len(mp.ContractList[i].Address), 20)
			validAddress = append(validAddress, mp.ContractList[i].Address...)
			mp.ContractList[i].Address = validAddress
		}
	}
}

type ManageSysContractAddressProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	// Contract Address
	ContractAddr sdk.AccAddress `json:"contract_address" yaml:"contract_address"`
	IsAdded      bool           `json:"is_added" yaml:"is_added"`
}

// NewManageSysContractAddressProposal creates a new instance of NewManageSysContractAddressProposal
func NewManageSysContractAddressProposal(title, description string, addr sdk.AccAddress, isAdded bool,
) ManageSysContractAddressProposal {
	return ManageSysContractAddressProposal{
		Title:        title,
		Description:  description,
		ContractAddr: addr,
		IsAdded:      isAdded,
	}
}

// GetTitle returns title of a manage system contract address proposal object
func (mp ManageSysContractAddressProposal) GetTitle() string {
	return mp.Title
}

// GetDescription returns description of a manage system contract address proposal object
func (mp ManageSysContractAddressProposal) GetDescription() string {
	return mp.Description
}

// ProposalRoute returns route key of a manage system contract address proposal object
func (mp ManageSysContractAddressProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns type of a manage system contract address proposal object
func (mp ManageSysContractAddressProposal) ProposalType() string {
	return proposalTypeManageSysContractAddress
}

// ValidateBasic validates a manage system contract address proposal
func (mp ManageSysContractAddressProposal) ValidateBasic() sdk.Error {
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

	if mp.ProposalType() != proposalTypeManageSysContractAddress {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	if mp.IsAdded && mp.ContractAddr.Empty() {
		return govtypes.ErrInvalidProposalContent("is_added true, contract address required")
	}

	return nil
}

// String returns a human readable string representation of a ManageSysContractAddressProposal
func (mp ManageSysContractAddressProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ManageSysContractAddressProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 ContractAddr:          %s
 IsAdded:				%t
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.ContractAddr.String(), mp.IsAdded),
	)
	return strings.TrimSpace(builder.String())
}

type ManageContractByteCodeProposal struct {
	Title              string         `json:"title" yaml:"title"`
	Description        string         `json:"description" yaml:"description"`
	Contract           sdk.AccAddress `json:"contract" yaml:"contract"`
	SubstituteContract sdk.AccAddress `json:"substitute_contract" yaml:"substitute_contract"`
}

func NewManageContractByteCodeProposal(title, description string, Contract sdk.AccAddress, SubstituteContract sdk.AccAddress) ManageContractByteCodeProposal {
	return ManageContractByteCodeProposal{
		Title:              title,
		Description:        description,
		Contract:           Contract,
		SubstituteContract: SubstituteContract,
	}
}

func (mp ManageContractByteCodeProposal) GetTitle() string {
	return mp.Title
}

func (mp ManageContractByteCodeProposal) GetDescription() string {
	return mp.Description
}

func (mp ManageContractByteCodeProposal) ProposalRoute() string {
	return RouterKey
}

func (mp ManageContractByteCodeProposal) ProposalType() string {
	return proposalTypeManageContractByteCode
}

func (mp ManageContractByteCodeProposal) ValidateBasic() sdk.Error {

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

	if mp.ProposalType() != proposalTypeManageContractByteCode {
		return govtypes.ErrInvalidProposalType(mp.ProposalType())
	}

	return nil
}
func (mp ManageContractByteCodeProposal) String() string {
	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(`ManageContractByteCodeProposal:
 Title:					%s
 Description:        	%s
 Type:                	%s
 Contract:          %s
 SubstituteContract:				%s
`,
			mp.Title, mp.Description, mp.ProposalType(), mp.Contract.String(), mp.SubstituteContract.String()),
	)
	return strings.TrimSpace(builder.String())
}
