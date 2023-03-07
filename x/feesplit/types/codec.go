package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/system"
)

// ModuleCdc defines the feesplit module's codec
var ModuleCdc = codec.New()

const (
	// Amino names
	registerFeeSplitName = system.Chain+"/MsgRegisterFeeSplit"
	updateFeeSplitName   = system.Chain+"/MsgUpdateFeeSplit"
	cancelFeeSplitName   = system.Chain+"/MsgCancelFeeSplit"
	sharesProposalName   = system.Chain+"/feesplit/SharesProposal"
)

// NOTE: This is required for the GetSignBytes function
func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all the necessary types and interfaces for the
// feesplit module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterFeeSplit{}, registerFeeSplitName, nil)
	cdc.RegisterConcrete(MsgUpdateFeeSplit{}, updateFeeSplitName, nil)
	cdc.RegisterConcrete(MsgCancelFeeSplit{}, cancelFeeSplitName, nil)
	cdc.RegisterConcrete(FeeSplitSharesProposal{}, sharesProposalName, nil)
}
