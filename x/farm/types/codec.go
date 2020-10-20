package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePool{}, "okexchain/farm/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgDestroyPool{}, "okexchain/farm/MsgDestroyPool", nil)
	cdc.RegisterConcrete(MsgSetWhite{}, "okexchain/farm/MsgSetWhite", nil)
	cdc.RegisterConcrete(MsgLock{}, "okexchain/farm/MsgLock", nil)
	cdc.RegisterConcrete(MsgUnlock{}, "okexchain/farm/MsgUnlock", nil)
	cdc.RegisterConcrete(MsgClaim{}, "okexchain/farm/MsgClaim", nil)
	cdc.RegisterConcrete(MsgProvide{}, "okexchain/farm/MsgProvide", nil)
	cdc.RegisterConcrete(ManageWhiteListProposal{}, "okexchain/farm/ManageWhiteListProposal", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
