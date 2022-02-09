package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	paramtypes "github.com/okex/exchain/x/params"
	"gopkg.in/yaml.v2"
)

// Parameter keys
var (
	ParamStoreKeyNoBaseFee                = []byte("NoBaseFee")
	ParamStoreKeyBaseFeeChangeDenominator = []byte("BaseFeeChangeDenominator")
	ParamStoreKeyElasticityMultiplier     = []byte("ElasticityMultiplier")
	ParamStoreKeyInitialBaseFee           = []byte("InitialBaseFee")
	ParamStoreKeyEnableHeight             = []byte("EnableHeight")
)

// Params defines the EVM module parameters
type Params struct {
	// no base fee forces the EIP-1559 base fee to 0 (needed for 0 price calls)
	NoBaseFee bool `json:"no_base_fee,omitempty"`
	// base fee change denominator bounds the amount the base fee can change
	// between blocks.
	BaseFeeChangeDenominator uint32 `json:"base_fee_change_denominator,omitempty"`
	// elasticity multiplier bounds the maximum gas limit an EIP-1559 block may
	// have.
	ElasticityMultiplier uint32 `json:"elasticity_multiplier,omitempty"`
	// initial base fee for EIP-1559 blocks.
	InitialBaseFee int64 `json:"initial_base_fee,omitempty"`
	// height at which the base fee calculation is enabled.
	EnableHeight int64 `json:"enable_height,omitempty"`
}

// NewParams creates a new Params instance
func NewParams(noBaseFee bool, baseFeeChangeDenom, elasticityMultiplier uint32, initialBaseFee, enableHeight int64) Params {
	return Params{
		NoBaseFee:                noBaseFee,
		BaseFeeChangeDenominator: baseFeeChangeDenom,
		ElasticityMultiplier:     elasticityMultiplier,
		InitialBaseFee:           initialBaseFee,
		EnableHeight:             enableHeight,
	}
}

// DefaultParams returns default evm parameters
func DefaultParams() Params {
	return Params{
		NoBaseFee:                false,
		BaseFeeChangeDenominator: params.BaseFeeChangeDenominator,
		ElasticityMultiplier:     params.ElasticityMultiplier,
		InitialBaseFee:           params.InitialBaseFee,
		EnableHeight:             0,
	}
}

// String implements the fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyNoBaseFee, &p.NoBaseFee, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFeeChangeDenominator, &p.BaseFeeChangeDenominator, validateBaseFeeChangeDenominator),
		paramtypes.NewParamSetPair(ParamStoreKeyElasticityMultiplier, &p.ElasticityMultiplier, validateElasticityMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyInitialBaseFee, &p.InitialBaseFee, validateInitialBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableHeight, &p.EnableHeight, validateEnableHeight),
	}
}

// Validate performs basic validation on fee market parameters.
func (p Params) Validate() error {
	if p.BaseFeeChangeDenominator == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	if p.InitialBaseFee < 0 {
		return fmt.Errorf("initial base fee cannot be negative: %d", p.InitialBaseFee)
	}

	if p.EnableHeight < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", p.EnableHeight)
	}

	return nil
}
func (p *Params) IsBaseFeeEnabled(height int64) bool {
	return !p.NoBaseFee && height >= p.EnableHeight
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBaseFeeChangeDenominator(i interface{}) error {
	value, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	return nil
}

func validateElasticityMultiplier(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateInitialBaseFee(i interface{}) error {
	value, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value < 0 {
		return fmt.Errorf("initial base fee cannot be negative: %d", value)
	}

	return nil
}

func validateEnableHeight(i interface{}) error {
	value, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", value)
	}

	return nil
}

// GenesisState defines the feemarket module's genesis state.
type GenesisState struct {
	// params defines all the paramaters of the module.
	Params Params `json:"params"`
	// base fee is the exported value from previous software version.
	// Zero by default.
	BaseFee sdk.Int `json:"base_fee"`
	// block gas is the amount of gas used on the last block before the upgrade.
	// Zero by default.
	BlockGas uint64 `json:"block_gas,omitempty"`
}
