package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)

	cdc.RegisterConcrete(MsgSubmitProposal{}, "exchain/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "exchain/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "exchain/gov/MsgVote", nil)

	cdc.RegisterConcrete(TextProposal{}, "exchain/gov/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, "exchain/gov/SoftwareUpgradeProposal", nil)
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
