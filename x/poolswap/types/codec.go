package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddLiquidity{}, "okchain/poolswap/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(MsgRemoveLiquidity{}, "okchain/poolswap/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(MsgCreateExchange{}, "okchain/poolswap/MsgCreateExchange", nil)
	cdc.RegisterConcrete(MsgTokenToNativeToken{}, "okchain/poolswap/MsgSwapToken", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
