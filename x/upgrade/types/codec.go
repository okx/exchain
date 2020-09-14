package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is a generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on code
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(AppUpgradeProposal{}, "okexchain/upgrade/AppUpgradeProposal", nil)
}
