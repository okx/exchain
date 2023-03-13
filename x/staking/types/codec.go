package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/system"
)

// RegisterCodec registers concrete types for codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, system.Chain+"/staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, system.Chain+"/staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgEditValidatorCommissionRate{}, system.Chain+"/staking/MsgEditValidatorCommissionRate", nil)
	cdc.RegisterConcrete(MsgDestroyValidator{}, system.Chain+"/staking/MsgDestroyValidator", nil)
	cdc.RegisterConcrete(MsgDeposit{}, system.Chain+"/staking/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, system.Chain+"/staking/MsgWithdraw", nil)
	cdc.RegisterConcrete(MsgAddShares{}, system.Chain+"/staking/MsgAddShares", nil)
	cdc.RegisterConcrete(ProposeValidatorProposal{}, ProposeValidatorProposalName, nil)
	cdc.RegisterConcrete(MsgDepositMinSelfDelegation{}, system.Chain+"/staking/MsgDepositMinSelfDelegation", nil)
}

// ModuleCdc is generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
