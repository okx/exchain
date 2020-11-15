package types

import "github.com/cosmos/cosmos-sdk/codec"

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
}

// nolint
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
