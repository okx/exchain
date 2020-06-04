package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTokenIssue{}, "okchain/token/MsgIssue", nil)
	cdc.RegisterConcrete(MsgTokenBurn{}, "okchain/token/MsgBurn", nil)
	cdc.RegisterConcrete(MsgTokenMint{}, "okchain/token/MsgMint", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "okchain/token/MsgMultiTransfer", nil)
	cdc.RegisterConcrete(MsgSend{}, "okchain/token/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "okchain/token/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgTokenModify{}, "okchain/token/MsgModify", nil)
	cdc.RegisterConcrete(MsgTokenActive{}, "okchain/token/MsgTokenActive", nil)
	cdc.RegisterConcrete(CertifiedTokenProposal{}, "okchain/token/CertifiedTokenProposal", nil)

	// for test
	//cdc.RegisterConcrete(MsgTokenDestroy{}, "okchain/token/MsgDestroy", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
