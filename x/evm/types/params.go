package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/okex/exchain/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace       = ModuleName
	DefaultMaxGasLimitPerTx = 30000000
)

// Parameter keys
var (
	ParamStoreKeyEnableCreate                = []byte("EnableCreate")
	ParamStoreKeyEnableCall                  = []byte("EnableCall")
	ParamStoreKeyExtraEIPs                   = []byte("EnableExtraEIPs")
	ParamStoreKeyContractDeploymentWhitelist = []byte("EnableContractDeploymentWhitelist")
	ParamStoreKeyContractBlockedList         = []byte("EnableContractBlockedList")
	ParamStoreKeyMaxGasLimitPerTx            = []byte("MaxGasLimitPerTx")
	ParamStoreKeyIbcDenom                    = []byte("IbcDenom")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the EVM module parameters
type Params struct {
	// EnableCreate toggles state transitions that use the vm.Create function
	EnableCreate bool `json:"enable_create" yaml:"enable_create"`
	// EnableCall toggles state transitions that use the vm.Call function
	EnableCall bool `json:"enable_call" yaml:"enable_call"`
	// ExtraEIPs defines the additional EIPs for the vm.Config
	ExtraEIPs []int `json:"extra_eips" yaml:"extra_eips"`
	// EnableContractDeploymentWhitelist controls the authorization of contract deployer
	EnableContractDeploymentWhitelist bool `json:"enable_contract_deployment_whitelist" yaml:"enable_contract_deployment_whitelist"`
	// EnableContractBlockedList controls the availability of contracts
	EnableContractBlockedList bool `json:"enable_contract_blocked_list" yaml:"enable_contract_blocked_list"`
	// MaxGasLimit defines the max gas limit in transaction
	MaxGasLimitPerTx uint64 `json:"max_gas_limit_per_tx" yaml:"max_gas_limit_per_tx"`
}

// NewParams creates a new Params instance
func NewParams(enableCreate, enableCall, enableContractDeploymentWhitelist, enableContractBlockedList bool, maxGasLimitPerTx uint64,
	extraEIPs ...int) Params {
	return Params{
		EnableCreate:                      enableCreate,
		EnableCall:                        enableCall,
		ExtraEIPs:                         extraEIPs,
		EnableContractDeploymentWhitelist: enableContractDeploymentWhitelist,
		EnableContractBlockedList:         enableContractBlockedList,
		MaxGasLimitPerTx:                  maxGasLimitPerTx,
	}
}

// DefaultParams returns default evm parameters
func DefaultParams() Params {
	return Params{
		EnableCreate:                      false,
		EnableCall:                        false,
		ExtraEIPs:                         []int(nil), // TODO: define default values
		EnableContractDeploymentWhitelist: false,
		EnableContractBlockedList:         false,
		MaxGasLimitPerTx:                  DefaultMaxGasLimitPerTx,
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
		params.NewParamSetPair(ParamStoreKeyEnableCreate, &p.EnableCreate, validateBool),
		params.NewParamSetPair(ParamStoreKeyEnableCall, &p.EnableCall, validateBool),
		params.NewParamSetPair(ParamStoreKeyExtraEIPs, &p.ExtraEIPs, validateEIPs),
		params.NewParamSetPair(ParamStoreKeyContractDeploymentWhitelist, &p.EnableContractDeploymentWhitelist, validateBool),
		params.NewParamSetPair(ParamStoreKeyContractBlockedList, &p.EnableContractBlockedList, validateBool),
		params.NewParamSetPair(ParamStoreKeyMaxGasLimitPerTx, &p.MaxGasLimitPerTx, validateUint64),
	}
}

// Validate performs basic validation on evm parameters.
func (p Params) Validate() error {
	return validateEIPs(p.ExtraEIPs)
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateEIPs(i interface{}) error {
	eips, ok := i.([]int)
	if !ok {
		return fmt.Errorf("invalid EIP slice type: %T", i)
	}

	for _, eip := range eips {
		if !vm.ValidEip(eip) {
			return fmt.Errorf("EIP %d is not activateable", eip)
		}
	}

	return nil
}

func validateUint64(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
