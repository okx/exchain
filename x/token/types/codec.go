package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTokenIssue{}, "okexchain/token/MsgIssue", nil)
	cdc.RegisterConcrete(MsgTokenBurn{}, "okexchain/token/MsgBurn", nil)
	cdc.RegisterConcrete(MsgTokenMint{}, "okexchain/token/MsgMint", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "okexchain/token/MsgMultiTransfer", nil)
	cdc.RegisterConcrete(MsgSend{}, "okexchain/token/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "okexchain/token/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgConfirmOwnership{}, "okexchain/token/MsgConfirmOwnership", nil)
	cdc.RegisterConcrete(MsgTokenModify{}, "okexchain/token/MsgModify", nil)

	// for test
	//cdc.RegisterConcrete(MsgTokenDestroy{}, "okexchain/token/MsgDestroy", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
