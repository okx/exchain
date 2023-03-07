package types

import (
	"fmt"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/params"

	"gopkg.in/yaml.v2"
)

// Parameter store key
var (
	DefaultEnableFeeSplit  = true
	DefaultDeveloperShares = sdk.NewDecWithPrec(50, 2) // 50%
	// Cost for executing `crypto.CreateAddress` must be at least 36 gas for the
	// contained keccak256(word) operation
	DefaultAddrDerivationCostCreate = uint64(50)

	ParamStoreKeyEnableFeeSplit           = []byte("EnableFeeSplit")
	ParamStoreKeyDeveloperShares          = []byte("DeveloperShares")
	ParamStoreKeyAddrDerivationCostCreate = []byte("AddrDerivationCostCreate")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the feesplit module params
type Params struct {
	// enable_fee_split defines a parameter to enable the feesplit module
	EnableFeeSplit bool `json:"enable_fee_split",yaml:"enable_fee_split"`
	// developer_shares defines the proportion of the transaction fees to be
	// distributed to the registered contract owner
	DeveloperShares sdk.Dec `json:"developer_shares",yaml:"developer_shares"`
	// addr_derivation_cost_create defines the cost of address derivation for
	// verifying the contract deployer at fee registration
	AddrDerivationCostCreate uint64 `json:"addr_derivation_cost_create",yaml:"addr_derivation_cost_create"`
}

// NewParams creates a new Params object
func NewParams(
	enableFeeSplit bool,
	developerShares sdk.Dec,
	addrDerivationCostCreate uint64,
) Params {
	return Params{
		EnableFeeSplit:           enableFeeSplit,
		DeveloperShares:          developerShares,
		AddrDerivationCostCreate: addrDerivationCostCreate,
	}
}

func DefaultParams() Params {
	return Params{
		EnableFeeSplit:           DefaultEnableFeeSplit,
		DeveloperShares:          DefaultDeveloperShares,
		AddrDerivationCostCreate: DefaultAddrDerivationCostCreate,
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
		params.NewParamSetPair(ParamStoreKeyEnableFeeSplit, &p.EnableFeeSplit, validateBool),
		params.NewParamSetPair(ParamStoreKeyDeveloperShares, &p.DeveloperShares, validateShares),
		params.NewParamSetPair(ParamStoreKeyAddrDerivationCostCreate, &p.AddrDerivationCostCreate, validateUint64),
	}
}

func (p Params) Validate() error {
	if err := validateBool(p.EnableFeeSplit); err != nil {
		return err
	}
	if err := validateShares(p.DeveloperShares); err != nil {
		return err
	}
	return validateUint64(p.AddrDerivationCostCreate)
}

func validateUint64(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
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

func validateShares(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %T", i)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %T", i)
	}

	return nil
}
