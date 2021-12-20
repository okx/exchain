package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	exported2 "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/vesting/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "cosmos-sdk/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "cosmos-sdk/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "cosmos-sdk/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(&PeriodicVestingAccount{}, "cosmos-sdk/PeriodicVestingAccount", nil)

	exported2.RegisterConcreteAccountInfo(uint(exported2.BaseVestingAcc), func() exported2.MptAccount{
		return &BaseVestingAccount{}
	})
	exported2.RegisterConcreteAccountInfo(uint(exported2.ContinuousVestingAcc), func() exported2.MptAccount{
		return &ContinuousVestingAccount{}
	})
	exported2.RegisterConcreteAccountInfo(uint(exported2.DelayedVestingAcc), func() exported2.MptAccount{
		return &DelayedVestingAccount{}
	})
	exported2.RegisterConcreteAccountInfo(uint(exported2.PeriodicVestingAcc), func() exported2.MptAccount{
		return &PeriodicVestingAccount{}
	})
}

// VestingCdc module wide codec
var VestingCdc *codec.Codec

func init() {
	VestingCdc = codec.New()
	RegisterCodec(VestingCdc)
	VestingCdc.Seal()
}
