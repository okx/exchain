package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/okex/exchain/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName

	DefaultIbcTimeout            = uint64(86400000000000) // 1 day
	DefaultAutoDeploymentEnabled = false
)

var (
	KeyEnableAutoDeployment = []byte("EnableAutoDeployment")
	KeyIbcTimeout           = []byte("IbcTimeout")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the module parameters
type Params struct {
	EnableAutoDeployment bool   `json:"enable_auto_deployment" yaml:"enable_auto_deployment"`
	IbcTimeout           uint64 `json:"ibc_timeout" yaml:"ibc_timeout"`
}

// NewParams creates a new Params instance
func NewParams(enableAutoDeployment bool, ibcTimeout uint64) Params {
	return Params{
		EnableAutoDeployment: enableAutoDeployment,
		IbcTimeout:           ibcTimeout,
	}
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		EnableAutoDeployment: DefaultAutoDeploymentEnabled,
		IbcTimeout:           DefaultIbcTimeout,
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
		params.NewParamSetPair(KeyEnableAutoDeployment, &p.EnableAutoDeployment, validateBool),
		params.NewParamSetPair(KeyIbcTimeout, &p.IbcTimeout, validateUint64),
	}
}

// Validate performs basic validation on erc20 parameters.
func (p Params) Validate() error {
	if err := validateUint64(p.IbcTimeout); err != nil {
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
