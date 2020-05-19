package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace     = ModuleName
	DefaultWithdrawPeriod = time.Hour * 24 * 3
)

// Parameter store keys
var (
	keyWithdrawPeriod = []byte("WithdrawPeriod")
)

// ParamKeyTable for margin module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for margin at genesis
type Params struct {
	WithdrawPeriod time.Duration `json:"withdraw_period"`
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf("Params: \nWithdrawPeriod:%d\n", p.WithdrawPeriod)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: keyWithdrawPeriod, Value: &p.WithdrawPeriod},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() *Params {
	return &Params{
		WithdrawPeriod: DefaultWithdrawPeriod,
	}
}
