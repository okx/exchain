package types

import (
	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
