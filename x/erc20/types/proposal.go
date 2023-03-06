package types

import (
	"fmt"
	"github.com/okx/okbchain/libs/system"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

const (
	// proposalTypeTokenMapping defines the type for a TokenMappingProposal

	proposalTypeTokenMapping          = "TokenMapping"
	proposalTypeProxyContractRedirect = "ProxyContractRedirect"

	proposalTypeContractTemplate = "ContractTemplate"

	ProposalTypeContextTemplateProxy = "proxy"
	ProposalTypeContextTemplateImpl  = "implement"
)

func init() {
	govtypes.RegisterProposalType(proposalTypeTokenMapping)
	govtypes.RegisterProposalType(proposalTypeProxyContractRedirect)
	govtypes.RegisterProposalType(proposalTypeContractTemplate)
	govtypes.RegisterProposalTypeCodec(TokenMappingProposal{}, system.Chain+"/erc20/TokenMappingProposal")
	govtypes.RegisterProposalTypeCodec(ProxyContractRedirectProposal{}, system.Chain+"/erc20/ProxyContractRedirectProposal")
	govtypes.RegisterProposalTypeCodec(ContractTemplateProposal{}, system.Chain+"/erc20/ContractTemplateProposal")
}

var (
	_ govtypes.Content = (*TokenMappingProposal)(nil)
	_ govtypes.Content = (*ContractTemplateProposal)(nil)
	_ govtypes.Content = (*ProxyContractRedirectProposal)(nil)
)

type TokenMappingProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Denom       string `json:"denom" yaml:"denom"`
	Contract    string `json:"contract" yaml:"contract"`
}

func NewTokenMappingProposal(title, description, denom string, contractAddr *common.Address) TokenMappingProposal {
	contract := ""
	if contractAddr != nil {
		contract = contractAddr.Hex()
	}
	return TokenMappingProposal{title, description, denom, contract}
}

func (tp TokenMappingProposal) GetTitle() string       { return tp.Title }
func (tp TokenMappingProposal) GetDescription() string { return tp.Description }
func (tp TokenMappingProposal) ProposalRoute() string  { return RouterKey }
func (tp TokenMappingProposal) ProposalType() string   { return proposalTypeTokenMapping }
func (tp TokenMappingProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(tp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(tp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is longer than the max")
	}

	if len(tp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(tp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is longer than the max")
	}

	if tp.ProposalType() != proposalTypeTokenMapping {
		return govtypes.ErrInvalidProposalType(tp.ProposalType())
	}

	if len(strings.TrimSpace(tp.Denom)) == 0 {
		return govtypes.ErrInvalidProposalContent("denom is required")
	}

	if err := sdk.ValidateDenom(tp.Denom); err != nil {
		return govtypes.ErrInvalidProposalContent("invalid denom")
	}

	if len(strings.TrimSpace(tp.Contract)) > 0 && !common.IsHexAddress(tp.Contract) {
		return govtypes.ErrInvalidProposalContent("invalid contract")
	}

	return nil
}

func (tp TokenMappingProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Token Mapping Proposal:
  Title:       %s
  Description: %s
  Denom:       %s
  Contract:    %s
`, tp.Title, tp.Description, tp.Denom, tp.Contract))

	return b.String()
}

type RedirectType int

const (
	RedirectImplementation = iota
	RedirectOwner
)

var RedirectMap = map[RedirectType]string{
	RedirectImplementation: "ImplementationAddr",
	RedirectOwner:          "OwnerAddr",
}

type ProxyContractRedirectProposal struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Denom       string       `json:"denom" yaml:"denom"`
	Tp          RedirectType `json:"type" yaml:"type"`
	Addr        string       `json:"addr" yaml:"addr"`
}

func NewProxyContractRedirectProposal(title, description, denom string, tp RedirectType, addr *common.Address) ProxyContractRedirectProposal {
	address := ""
	if addr != nil {
		address = addr.Hex()
	}
	return ProxyContractRedirectProposal{title, description, denom, tp, address}
}

func (tp ProxyContractRedirectProposal) GetTitle() string {
	return tp.Title
}

func (tp ProxyContractRedirectProposal) GetDescription() string {
	return tp.Description
}

func (tp ProxyContractRedirectProposal) ProposalRoute() string {
	return RouterKey
}

func (tp ProxyContractRedirectProposal) ProposalType() string {
	return proposalTypeProxyContractRedirect
}

func (tp ProxyContractRedirectProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(tp.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(tp.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is longer than the max")
	}

	if len(tp.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(tp.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is longer than the max")
	}

	if tp.ProposalType() != proposalTypeProxyContractRedirect {
		return govtypes.ErrInvalidProposalType(tp.ProposalType())
	}

	if len(strings.TrimSpace(tp.Denom)) == 0 {
		return govtypes.ErrInvalidProposalContent("denom is required")
	}

	if err := sdk.ValidateDenom(tp.Denom); err != nil {
		return govtypes.ErrInvalidProposalContent("invalid denom")
	}
	switch tp.Tp {
	case RedirectImplementation, RedirectOwner:
	default:
		return govtypes.ErrInvalidProposer()
	}
	if len(strings.TrimSpace(tp.Addr)) > 0 && !common.IsHexAddress(tp.Addr) {
		return govtypes.ErrInvalidProposalContent("invalid contract")
	}
	return nil
}

func (tp ProxyContractRedirectProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Proxy Contract Redirect Proposal:
  		Title:       %s
  		Description: %s
		Denom:       %s
		Tp:          %s
  		Addr:        %s
	`, tp.Title, tp.Description, tp.Denom, RedirectMap[tp.Tp], tp.Addr))

	return b.String()
}

type ContractTemplateProposal struct {
	Title        string `json:"title" yaml:"title"`
	Description  string `json:"description" yaml:"description"`
	ContractType string `json:"contract_type"`
	Contract     string `json:"contract"`
}

func NewContractTemplateProposal(title string, description string, contractType string, contract string) ContractTemplateProposal {
	return ContractTemplateProposal{Title: title, Description: description, ContractType: contractType, Contract: contract}
}

func (b ContractTemplateProposal) GetTitle() string { return b.Title }

func (b ContractTemplateProposal) GetDescription() string { return b.Description }

func (b ContractTemplateProposal) ProposalRoute() string { return RouterKey }

func (b ContractTemplateProposal) ProposalType() string { return proposalTypeContractTemplate }

func (b ContractTemplateProposal) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(b.Title)) == 0 {
		return govtypes.ErrInvalidProposalContent("title is required")
	}
	if len(b.Title) > govtypes.MaxTitleLength {
		return govtypes.ErrInvalidProposalContent("title length is longer than the max")
	}

	if len(b.Description) == 0 {
		return govtypes.ErrInvalidProposalContent("description is required")
	}

	if len(b.Description) > govtypes.MaxDescriptionLength {
		return govtypes.ErrInvalidProposalContent("description length is longer than the max")
	}
	if len(b.ContractType) == 0 || b.ContractType != ProposalTypeContextTemplateProxy && b.ContractType != ProposalTypeContextTemplateImpl {
		return govtypes.ErrInvalidProposalContent("invalid type , should be proxy or implement")
	}

	if b.ProposalType() != proposalTypeContractTemplate {
		return govtypes.ErrInvalidProposalType(b.ProposalType())
	}
	con, err := UnmarshalCompileContract([]byte(b.Contract))
	if nil != err {
		return err
	}
	if err := con.ValidBasic(); nil != err {
		return err
	}
	return nil
}

func (b ContractTemplateProposal) String() string {
	return ""
}
