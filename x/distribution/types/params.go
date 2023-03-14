package types

import (
	"fmt"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/params"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = ModuleName
)

// Parameter keys
var (
	ParamStoreKeyCommunityTax            = []byte("communitytax")
	ParamStoreKeyWithdrawAddrEnabled     = []byte("withdrawaddrenabled")
	ParamStoreKeyDistributionType        = []byte("distributiontype")
	ParamStoreKeyWithdrawRewardEnabled   = []byte("withdrawrewardenabled")
	ParamStoreKeyRewardTruncatePrecision = []byte("rewardtruncateprecision")
)

// Params defines the set of distribution parameters.
type Params struct {
	CommunityTax            sdk.Dec `json:"community_tax" yaml:"community_tax"`
	WithdrawAddrEnabled     bool    `json:"withdraw_addr_enabled" yaml:"withdraw_addr_enabled"`
	DistributionType        uint32  `json:"distribution_type" yaml:"distribution_type"`
	WithdrawRewardEnabled   bool    `json:"withdraw_reward_enabled" yaml:"withdraw_reward_enabled"`
	RewardTruncatePrecision int64   `json:"reward_truncate_precision" yaml:"reward_truncate_precision"`
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default distribution parameters
func DefaultParams() Params {
	return Params{
		CommunityTax:            sdk.NewDecWithPrec(2, 2), // 2%
		WithdrawAddrEnabled:     true,
		DistributionType:        0,
		WithdrawRewardEnabled:   true,
		RewardTruncatePrecision: 0,
	}
}

// String returns a human readable string representation of Params
func (p Params) String() string {
	return fmt.Sprintf(`Distribution Params:
  Community Tax:          %s
  Withdraw Addr Enabled:  %t
  Distribution Type: %d
  Withdraw Reward Enabled: %t
  Reward Truncate Precision: %d`,
		p.CommunityTax, p.WithdrawAddrEnabled, p.DistributionType, p.WithdrawRewardEnabled, p.RewardTruncatePrecision)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyCommunityTax, &p.CommunityTax, validateCommunityTax),
		params.NewParamSetPair(ParamStoreKeyWithdrawAddrEnabled, &p.WithdrawAddrEnabled, validateWithdrawAddrEnabled),
		//new params for distribution proposal
		params.NewParamSetPair(ParamStoreKeyDistributionType, &p.DistributionType, validateDistributionType),
		params.NewParamSetPair(ParamStoreKeyWithdrawRewardEnabled, &p.WithdrawRewardEnabled, validateWithdrawRewardEnabled),
		params.NewParamSetPair(ParamStoreKeyRewardTruncatePrecision, &p.RewardTruncatePrecision, validateRewardTruncatePrecision),
	}
}

// ValidateBasic performs basic validation on distribution parameters.
func (p Params) ValidateBasic() error {
	if p.CommunityTax.IsNegative() || p.CommunityTax.GT(sdk.OneDec()) {
		return fmt.Errorf(
			"community tax should non-negative and less than one: %s", p.CommunityTax,
		)
	}

	return nil
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("community tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("community tax must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("community tax too large: %s", v)
	}

	return nil
}

func validateWithdrawAddrEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDistributionType(i interface{}) error {
	distributionType, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if distributionType != DistributionTypeOnChain && distributionType != DistributionTypeOffChain {
		return fmt.Errorf("invalid distribution type: %d", distributionType)
	}

	return nil
}

func validateWithdrawRewardEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateRewardTruncatePrecision(i interface{}) error {
	precision, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if precision < 0 || precision > sdk.Precision {
		return fmt.Errorf("invalid parameter precision: %d", precision)
	}

	return nil
}

// NewParams creates a new instance of Params
func NewParams(communityTax sdk.Dec, withdrawAddrEnabled bool, distributionType uint32, withdrawRewardEnabled bool, rewardTruncatePrecision int64) Params {
	return Params{
		CommunityTax:            communityTax,
		WithdrawAddrEnabled:     withdrawAddrEnabled,
		DistributionType:        distributionType,
		WithdrawRewardEnabled:   withdrawRewardEnabled,
		RewardTruncatePrecision: rewardTruncatePrecision,
	}
}

// MarshalYAML implements the text format for yaml marshaling
func (p Params) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}
