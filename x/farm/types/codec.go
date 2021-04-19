package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePool{}, "exchain/farm/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgDestroyPool{}, "exchain/farm/MsgDestroyPool", nil)
	cdc.RegisterConcrete(MsgLock{}, "exchain/farm/MsgLock", nil)
	cdc.RegisterConcrete(MsgUnlock{}, "exchain/farm/MsgUnlock", nil)
	cdc.RegisterConcrete(MsgClaim{}, "exchain/farm/MsgClaim", nil)
	cdc.RegisterConcrete(MsgProvide{}, "exchain/farm/MsgProvide", nil)
	cdc.RegisterConcrete(ManageWhiteListProposal{}, "exchain/farm/ManageWhiteListProposal", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
