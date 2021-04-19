package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgList{}, "exchain/dex/MsgList", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "exchain/dex/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "exchain/dex/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "exchain/dex/MsgTransferTradingPairOwnership", nil)
	cdc.RegisterConcrete(MsgConfirmOwnership{}, "exchain/dex/MsgConfirmOwnership", nil)
	cdc.RegisterConcrete(DelistProposal{}, "exchain/dex/DelistProposal", nil)
	cdc.RegisterConcrete(MsgCreateOperator{}, "exchain/dex/CreateOperator", nil)
	cdc.RegisterConcrete(MsgUpdateOperator{}, "exchain/dex/UpdateOperator", nil)
}

// ModuleCdc represents generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
