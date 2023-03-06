package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/system"
)

// ModuleCdc defines the erc20 module's codec
var ModuleCdc = codec.New()

const (
	TokenMappingProposalName          = system.Chain+"/erc20/TokenMappingProposal"
	ProxyContractRedirectProposalName = system.Chain+"/erc20/ProxyContractRedirectProposal"
	ContractTemplateProposalName      = system.Chain+"/erc20/ContractTemplateProposal"
	CompiledContractProposalName      = system.Chain+"/erc20/Contract"
)

// RegisterCodec registers all the necessary types and interfaces for the
// erc20 module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(TokenMappingProposal{}, TokenMappingProposalName, nil)

	cdc.RegisterConcrete(ProxyContractRedirectProposal{}, ProxyContractRedirectProposalName, nil)
	cdc.RegisterConcrete(ContractTemplateProposal{}, ContractTemplateProposalName, nil)
	cdc.RegisterConcrete(CompiledContract{}, CompiledContractProposalName, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
