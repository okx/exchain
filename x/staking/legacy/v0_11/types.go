package v0_11

import (
	"time"

	"github.com/okex/okchain/x/staking/legacy/v0_10"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "staking"

var DefaultMinSelfDelegation = sdk.NewDec(10000)

type (
	// GenesisState - all staking state that must be provided at genesis
	GenesisState struct {
		Params               Params                            `json:"params" yaml:"params"`
		LastTotalPower       sdk.Int                           `json:"last_total_power" yaml:"last_total_power"`
		LastValidatorPowers  []v0_10.LastValidatorPower        `json:"last_validator_powers" yaml:"last_validator_powers"`
		Validators           []v0_10.ValidatorExported         `json:"validators" yaml:"validators"`
		Delegators           []v0_10.Delegator                 `json:"delegators" yaml:"delegators"`
		UnbondingDelegations []v0_10.UndelegationInfo          `json:"unbonding_delegations" yaml:"unbonding_delegations"`
		AllShares            []SharesExported                  `json:"all_shares" yaml:"all_shares"`
		ProxyDelegatorKeys   []v0_10.ProxyDelegatorKeyExported `json:"proxy_delegator_keys" yaml:"proxy_delegator_keys"`
		Exported             bool                              `json:"exported" yaml:"exported"`
	}

	// Params defines the high level settings for staking
	Params struct {
		// time duration of unbonding
		UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
		// note: we need to be a bit careful about potential overflow here, since this is user-determined
		// maximum number of validators (max uint16 = 65535)
		MaxValidators uint16 `json:"max_bonded_validators" yaml:"max_bonded_validators"`
		// epoch for validator update
		Epoch              uint16 `json:"epoch" yaml:"epoch"`
		MaxValsToAddShares uint16 `json:"max_validators_to_add_shares" yaml:"max_validators_to_add_shares"`
		// bondable coin denomination
		BondDenom string `json:"bond_denom" yaml:"bond_denom"`
		// limited amount of delegate
		MinDelegation sdk.Dec `json:"min_delegation" yaml:"min_delegation"`
		// validator's self declared minimum self delegation
		MinSelfDelegation sdk.Dec `json:"min_self_delegation" yaml:"min_self_delegation"`
	}

	Shares = sdk.Dec
	// SharesExported is designed for types.Shares export
	SharesExported struct {
		DelAddress       sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
		Shares           Shares         `json:"shares" yaml:"shares"`
	}
)
