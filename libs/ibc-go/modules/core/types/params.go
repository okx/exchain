package types

import (
	"fmt"

	paramtypes "github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

// DefaultIbcEnabled enabled
const DefaultIbcEnabled = false

// KeyIbcEnabled is store's key for IbcEnabled Params
var KeyIbcEnabled = []byte("IbcEnabled")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the ibc module
func NewParams(enableIbc bool) Params {
	return Params{
		EnableIbc: enableIbc,
	}
}

// DefaultParams is the default parameter configuration for the ibc module
func DefaultParams() Params {
	return Params{DefaultIbcEnabled}
}

// Validate all ibc module parameters
func (p Params) Validate() error {
	return validateEnabled(p.EnableIbc)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIbcEnabled, &p.EnableIbc, validateEnabled),
	}
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
