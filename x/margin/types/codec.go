package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/margin/MsgDexDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okchain/margin/MsgDexWithdraw", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/margin/MsgDexSet", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okchain/margin/MsgDexSave", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/margin/MsgDexReturn", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okchain/margin/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/margin/MsgBorrow", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okchain/margin/MsgRepay", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/margin/MsgWithdraw", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
