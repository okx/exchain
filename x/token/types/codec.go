package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/system"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTokenIssue{}, system.Chain+"/token/MsgIssue", nil)
	cdc.RegisterConcrete(MsgTokenBurn{}, system.Chain+"/token/MsgBurn", nil)
	cdc.RegisterConcrete(MsgTokenMint{}, system.Chain+"/token/MsgMint", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, system.Chain+"/token/MsgMultiTransfer", nil)
	cdc.RegisterConcrete(MsgSend{}, system.Chain+"/token/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, system.Chain+"/token/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgConfirmOwnership{}, system.Chain+"/token/MsgConfirmOwnership", nil)
	cdc.RegisterConcrete(MsgTokenModify{}, system.Chain+"/token/MsgModify", nil)

	// for test
	//cdc.RegisterConcrete(MsgTokenDestroy{}, system.Chain+"/token/MsgDestroy", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
