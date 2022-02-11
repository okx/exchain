package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

// ModuleCdc defines the evm module's codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the
// feemarket module
func RegisterCodec(cdc *codec.Codec) {

}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
