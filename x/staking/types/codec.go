package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types for codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "okchain/staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "okchain/staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDestroyValidator{}, "okchain/staking/MsgDestroyValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "okchain/staking/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "okchain/staking/MsgUnDelegate", nil)
	cdc.RegisterConcrete(MsgVote{}, "okchain/staking/MsgVote", nil)
	cdc.RegisterConcrete(MsgRegProxy{}, "okchain/staking/MsgRegProxy", nil)
	cdc.RegisterConcrete(MsgBindProxy{}, "okchain/staking/MsgBindProxy", nil)
	cdc.RegisterConcrete(MsgUnbindProxy{}, "okchain/staking/MsgUnbindProxy", nil)
}

// ModuleCdc is generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
