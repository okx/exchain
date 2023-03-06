package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/system"
)

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)

	cdc.RegisterConcrete(MsgSubmitProposal{}, system.Chain+"/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, system.Chain+"/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, system.Chain+"/gov/MsgVote", nil)

	cdc.RegisterConcrete(TextProposal{}, system.Chain+"/gov/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, system.Chain+"/gov/SoftwareUpgradeProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)
}
