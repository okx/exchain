package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/okex/exchain/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
)

var (
	ParamStoreKeyIbcDenom = []byte("IbcDenom")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the module parameters
type Params struct {
	IbcDenom string `json:"ibc_denom" yaml:"ibc_denom"`
}

// NewParams creates a new Params instance
func NewParams(ibc_denom string) Params {
	return Params{
		IbcDenom: ibc_denom,
	}
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		IbcDenom: "ibc/DDCD907790B8AA2BF9B2B3B614718FA66BFC7540E832CE3E3696EA717DCEFF49",
	}
}

// String implements the fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyIbcDenom, &p.IbcDenom, validateIbcDenom),
	}
}

// Validate performs basic validation on erc20 parameters.
func (p Params) Validate() error {
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateIbcDenom(i interface{}) error {
	return nil
}
