package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

// ModuleCdc defines the erc20 module's codec
var ModuleCdc = codec.New()

const (
	TokenMappingProposalName = "okexchain/erc20/TokenMappingProposal"
)

// RegisterCodec registers all the necessary types and interfaces for the
// erc20 module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(TokenMappingProposal{}, TokenMappingProposalName, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
