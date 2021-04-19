package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTokenIssue{}, "exchain/token/MsgIssue", nil)
	cdc.RegisterConcrete(MsgTokenBurn{}, "exchain/token/MsgBurn", nil)
	cdc.RegisterConcrete(MsgTokenMint{}, "exchain/token/MsgMint", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "exchain/token/MsgMultiTransfer", nil)
	cdc.RegisterConcrete(MsgSend{}, "exchain/token/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "exchain/token/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgConfirmOwnership{}, "exchain/token/MsgConfirmOwnership", nil)
	cdc.RegisterConcrete(MsgTokenModify{}, "exchain/token/MsgModify", nil)

	// for test
	//cdc.RegisterConcrete(MsgTokenDestroy{}, "exchain/token/MsgDestroy", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
