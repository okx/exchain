package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/okex/exchain/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace      = ModuleName
	IbcDenomDefaultValue   = "ibc/DDCD907790B8AA2BF9B2B3B614718FA66BFC7540E832CE3E3696EA717DCEFF49"
	IbcTimeoutDefaultValue = uint64(86400000000000) // 1 day
)

var (
	KeyIbcDenom             = []byte("IbcDenom")
	KeyEnableAutoDeployment = []byte("EnableAutoDeployment")
	KeyIbcTimeout           = []byte("IbcTimeout")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the module parameters
type Params struct {
	IbcDenom             string `json:"ibc_denom" yaml:"ibc_denom"`
	EnableAutoDeployment bool   `json:"enable_auto_deployment" yaml:"enable_auto_deployment"`
	IbcTimeout           uint64 `json:"ibc_timeout" yaml:"ibc_timeout"`
}

// NewParams creates a new Params instance
func NewParams(ibc_denom string, enableAutoDeployment bool, ibcTimeout uint64) Params {
	return Params{
		IbcDenom:             ibc_denom,
		EnableAutoDeployment: enableAutoDeployment,
		IbcTimeout:           ibcTimeout,
	}
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		IbcDenom:             IbcDenomDefaultValue,
		EnableAutoDeployment: true,
		IbcTimeout:           IbcTimeoutDefaultValue,
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
		params.NewParamSetPair(KeyIbcDenom, &p.IbcDenom, validateIbcDenom),
		params.NewParamSetPair(KeyEnableAutoDeployment, &p.EnableAutoDeployment, validateBool),
		params.NewParamSetPair(KeyIbcTimeout, &p.IbcTimeout, validateUint64),
	}
}

// Validate performs basic validation on erc20 parameters.
func (p Params) Validate() error {
	if err := validateUint64(p.IbcTimeout); err != nil {
		return err
	}
	if err := validateIbcDenom(p.IbcDenom); err != nil {
		return err
	}
	return nil
}

func validateUint64(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
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
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !IsValidIBCDenom(s) {
		return fmt.Errorf("invalid ibc denom: %T", i)
	}
	return nil
}
