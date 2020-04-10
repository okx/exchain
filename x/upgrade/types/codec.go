package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	// just 4 test
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// just 4 test
// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUpgradeConfig{}, "okchain/upgrade/MsgUpgradeConfig", nil)
	cdc.RegisterConcrete(AppUpgradeProposal{}, "okchain/upgrade/AppUpgradeProposal", nil)
}
