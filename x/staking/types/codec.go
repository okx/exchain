package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types for codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "exchain/staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "exchain/staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDestroyValidator{}, "exchain/staking/MsgDestroyValidator", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "exchain/staking/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "exchain/staking/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgAddShares{}, "exchain/staking/MsgAddShares", nil)
	cdc.RegisterConcrete(MsgRegProxy{}, "exchain/staking/MsgRegProxy", nil)
	cdc.RegisterConcrete(MsgBindProxy{}, "exchain/staking/MsgBindProxy", nil)
	cdc.RegisterConcrete(MsgUnbindProxy{}, "exchain/staking/MsgUnbindProxy", nil)
}

// ModuleCdc is generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
