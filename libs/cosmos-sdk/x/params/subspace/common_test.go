package subspace_test

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
)

var (
	keyUnbondingTime = []byte("UnbondingTime")
	keyMaxValidators = []byte("MaxValidators")
	keyBondDenom     = []byte("BondDenom")

	key  = sdk.NewKVStoreKey("storekey")
	tkey = sdk.NewTransientStoreKey("transientstorekey")
)

type params struct {
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	MaxValidators uint16        `json:"max_validators" yaml:"max_validators"`
	BondDenom     string        `json:"bond_denom" yaml:"bond_denom"`
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < (24 * time.Hour) {
		return fmt.Errorf("unbonding time must be at least one day")
	}

	return nil
}

func validateMaxValidators(i interface{}) error {
	_, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return errors.New("denom cannot be empty")
	}

	return nil
}

func (p *params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{keyUnbondingTime, &p.UnbondingTime, validateUnbondingTime},
		{keyMaxValidators, &p.MaxValidators, validateMaxValidators},
		{keyBondDenom, &p.BondDenom, validateBondDenom},
	}
}

func paramKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&params{})
}
