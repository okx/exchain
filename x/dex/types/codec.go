package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgList{}, "okchain/dex/MsgList", nil)
	//cdc.RegisterConcrete(MsgDelist{}, "okchain/dex/MsgDelist", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okchain/dex/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okchain/dex/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "okchain/dex/MsgTransferTradingPairOwnership", nil)
	cdc.RegisterConcrete(DelistProposal{}, "okchain/dex/DelistProposal", nil)
	cdc.RegisterConcrete(MsgCreateOperator{}, "okchan/dex/CreateOperator", nil)
	cdc.RegisterConcrete(MsgUpdateOperator{}, "okchan/dex/UpdateOperator", nil)
}

// ModuleCdc represents generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
