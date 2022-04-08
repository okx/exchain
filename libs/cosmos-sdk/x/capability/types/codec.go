package types

import "github.com/okex/exchain/libs/cosmos-sdk/codec"

var ModuleCdc *codec.Codec

func init(){
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
