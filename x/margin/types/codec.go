package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDexDeposit{}, "okchain/margin/MsgDexDeposit", nil)
	cdc.RegisterConcrete(MsgDexWithdraw{}, "okchain/margin/MsgDexWithdraw", nil)
	//cdc.RegisterConcrete(MsgDexSet{}, "okchain/margin/MsgDexSet", nil)
	cdc.RegisterConcrete(MsgDexSave{}, "okchain/margin/MsgDexSave", nil)
	cdc.RegisterConcrete(MsgDexReturn{}, "okchain/margin/MsgDexReturn", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/margin/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgBorrow{}, "okchain/margin/MsgBorrow", nil)
	cdc.RegisterConcrete(MsgRepay{}, "okchain/margin/MsgRepay", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okchain/margin/MsgWithdraw", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
