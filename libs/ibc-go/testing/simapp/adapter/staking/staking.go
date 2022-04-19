package staking

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/staking/keeper"
	"github.com/okex/exchain/x/staking/types"
	"time"
)

type StakingKeeper struct {
	keeper.Keeper
}

func (k StakingKeeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	return types2.DefaultUnbondingTime
}

// NewKeeper creates a new staking Keeper instance
func NewStakingKeeper(cdcMarshl *codec.CodecProxy, key sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace) *StakingKeeper {
	// set KeyTable if it has not already been set
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(ParamKeyTable())
	}
	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}
	k := keeper.NewKeeperWithNoParam(cdcMarshl, key, supplyKeeper, paramstore)
	return &StakingKeeper{
		Keeper: k,
	}
}

// ParamKeyTable returns param table for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(newTestParams())
}

type TestParams struct {
	*types.Params
}

func newTestParams() *TestParams {
	p := types.DefaultParams()
	p.UnbondingTime = types2.DefaultUnbondingTime
	ret := &TestParams{
		Params: &p,
	}
	return ret
}

// ParamSetPairs is the implements params.ParamSet
func (p *TestParams) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: types.KeyUnbondingTime, Value: &p.UnbondingTime, ValidatorFn: common.ValidateDurationPositive("unbonding time")},
		{Key: types.KeyMaxValidators, Value: &p.MaxValidators, ValidatorFn: common.ValidateUint16Positive("max validators")},
		{Key: types.KeyEpoch, Value: &p.Epoch, ValidatorFn: common.ValidateUint16Positive("epoch")},
		{Key: types.KeyMaxValsToAddShares, Value: &p.MaxValsToAddShares, ValidatorFn: common.ValidateUint16Positive("max vals to add shares")},
		{Key: types.KeyMinDelegation, Value: &p.MinDelegation, ValidatorFn: common.ValidateDecPositive("min delegation")},
		{Key: types.KeyMinSelfDelegation, Value: &p.MinSelfDelegation, ValidatorFn: common.ValidateDecPositive("min self delegation")},
		{Key: types.KeyHistoricalEntries, Value: &p.HistoricalEntries, ValidatorFn: validateHistoricalEntries},
	}
}

func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
