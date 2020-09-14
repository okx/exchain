package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgList{}, "okexchain/dex/MsgList", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "okexchain/dex/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "okexchain/dex/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "okexchain/dex/MsgTransferTradingPairOwnership", nil)
	cdc.RegisterConcrete(DelistProposal{}, "okexchain/dex/DelistProposal", nil)
	cdc.RegisterConcrete(MsgCreateOperator{}, "okexchain/dex/CreateOperator", nil)
	cdc.RegisterConcrete(MsgUpdateOperator{}, "okexchain/dex/UpdateOperator", nil)
}

// ModuleCdc represents generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
