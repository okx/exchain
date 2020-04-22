// DONTCOVER
// nolint
package v0_36

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v034staking "github.com/okex/okchain/x/staking/legacy/v0_34"
	"github.com/okex/okchain/x/staking/types"

	"github.com/tendermint/tendermint/crypto"
)

const (
	ModuleName = "staking"
)

type (
	Validator struct {
		OperatorAddress         sdk.ValAddress          `json:"operator_address" yaml:"operator_address"`
		ConsPubKey              crypto.PubKey           `json:"consensus_pubkey" yaml:"consensus_pubkey"`
		Jailed                  bool                    `json:"jailed" yaml:"jailed"`
		Status                  sdk.BondStatus          `json:"status" yaml:"status"`
		DelegatorShares         sdk.Dec                 `json:"delegator_shares" yaml:"delegator_shares"`
		Description             v034staking.Description `json:"description" yaml:"description"`
		UnbondingHeight         int64                   `json:"unbonding_height" yaml:"unbonding_height"`
		UnbondingCompletionTime time.Time               `json:"unbonding_time" yaml:"unbonding_time"`
		MinSelfDelegation       sdk.Int                 `json:"min_self_delegation" yaml:"min_self_delegation"`
	}

	bechValidator struct {
		OperatorAddress         sdk.ValAddress          `json:"operator_address" yaml:"operator_address"`
		ConsPubKey              string                  `json:"consensus_pubkey" yaml:"consensus_pubkey"`
		Jailed                  bool                    `json:"jailed" yaml:"jailed"`
		Status                  sdk.BondStatus          `json:"status" yaml:"status"`
		DelegatorShares         sdk.Dec                 `json:"delegator_shares" yaml:"delegator_shares"`
		Description             v034staking.Description `json:"description" yaml:"description"`
		UnbondingHeight         int64                   `json:"unbonding_height" yaml:"unbonding_height"`
		UnbondingCompletionTime time.Time               `json:"unbonding_time" yaml:"unbonding_time"`
		MinSelfDelegation       sdk.Int                 `json:"min_self_delegation" yaml:"min_self_delegation"`
	}

	Validators []Validator

	GenesisState struct {
		Params               Params                           `json:"params" yaml:"params"`
		LastTotalPower       sdk.Int                          `json:"last_total_power" yaml:"last_total_power"`
		LastValidatorPowers  []v034staking.LastValidatorPower `json:"last_validator_powers" yaml:"last_validator_powers"`
		Validators           Validators                       `json:"validators" yaml:"validators"`
		Delegators           []Delegator                      `json:"delegators" yaml:"delegators"`
		UnbondingDelegations []UndelegationInfo               `json:"unbonding_delegations" yaml:"unbonding_delegations"`
		Votes                []VotesExported                  `json:"votes" yaml:"votes"`
		ProxyDelegatorKeys   []ProxyDelegatorKeyExported      `json:"proxy_delegator_keys" yaml:"proxy_delegator_keys"`
		Exported             bool                             `json:"exported" yaml:"exported"`
	}

	ProxyDelegatorKeyExported struct {
		DelAddr   sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
		ProxyAddr sdk.AccAddress `json:"proxy_address" yaml:"proxy_address"`
	}

	VotesExported struct {
		VoterAddress     sdk.AccAddress `json:"voter_address" yaml:"voter_address"`
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
		Votes            Votes          `json:"votes" yaml:"votes"`
	}
	Votes = sdk.Dec

	UndelegationInfo struct {
		DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
		Quantity         sdk.Int        `json:"quantity" yaml:"quantity"`
		CompletionTime   time.Time      `json:"completion_time"`
	}

	Delegator struct {
		DelegatorAddress   sdk.AccAddress   `json:"delegator_address" yaml:"delegator_address"`
		ValidatorAddresses []sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
		Shares             sdk.Dec          `json:"shares" yaml:"shares"`
		Tokens             sdk.Int          `json:"tokens" yaml:"tokens"` // delegated tokens (not include. self-delegation)
	}

	Params struct {
		// time duration of unbonding
		UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
		// maximum number of validators (max uint16 = 65535)
		MaxValidators          uint16  `json:"max_validators" yaml:"max_validators"`
		// note: we need to be a bit careful about potential overflow here, since this is user-determined
		BondDenom     string `json:"bond_denom" yaml:"bond_denom"` // bondable coin denomination
		Epoch         uint16 `json:"epoch" yaml:"epoch"`           //epoch for validator update
		MaxValsToVote uint16 `json:"max_vals_to_vote" yaml:"max_vals_to_vote"`
		//limited amount of delegate
		MinDelegation sdk.Dec `json:"min_delegation" yaml:"min_delegation"`
	}
)

func NewGenesisState(
	params v034staking.Params, lastTotalPower sdk.Int, lastValPowers []v034staking.LastValidatorPower,
	validators Validators, delegations v034staking.Delegations,
	ubds []v034staking.UnbondingDelegation, reds []v034staking.Redelegation, exported bool,
) GenesisState {

	newParams := Params{
		UnbondingTime:          params.UnbondingTime,
		MaxValidators:          params.MaxValidators,
		BondDenom:              params.BondDenom,
		Epoch:                  types.DefaultEpoch,
		MaxValsToVote:          types.DefaultMaxValsToVote,
		MinDelegation:          types.DefaultMinDelegation,
	}

	var vals Validators
	for _, val := range validators {
		validator := val
		shares := sdk.ZeroDec()
		for _, delegation := range delegations {
			if delegation.ValidatorAddress.Equals(validator.OperatorAddress) {
				shares = shares.Add(delegation.Shares)
			}
		}
		validator.DelegatorShares = shares
		validator.MinSelfDelegation = sdk.NewIntFromBigInt(shares.Int)

		vals = append(vals, validator)
	}

	return GenesisState{
		Params:               newParams,
		LastTotalPower:       lastTotalPower,
		LastValidatorPowers:  lastValPowers,
		Validators:           vals,
		Delegators:           nil,
		UnbondingDelegations: nil,
		Votes:                nil,
		ProxyDelegatorKeys:   nil,
		Exported:             exported,
	}
}

func (v Validator) MarshalJSON() ([]byte, error) {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(v.ConsPubKey)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSON(bechValidator{
		OperatorAddress:         v.OperatorAddress,
		ConsPubKey:              bechConsPubKey,
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		DelegatorShares:         v.DelegatorShares,
		Description:             v.Description,
		UnbondingHeight:         v.UnbondingHeight,
		UnbondingCompletionTime: v.UnbondingCompletionTime,
		MinSelfDelegation:       v.MinSelfDelegation,
	})
}

func (v *Validator) UnmarshalJSON(data []byte) error {
	bv := &bechValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetConsPubKeyBech32(bv.ConsPubKey)
	if err != nil {
		return err
	}
	*v = Validator{
		OperatorAddress:         bv.OperatorAddress,
		ConsPubKey:              consPubKey,
		Jailed:                  bv.Jailed,
		Status:                  bv.Status,
		DelegatorShares:         bv.DelegatorShares,
		Description:             bv.Description,
		UnbondingHeight:         bv.UnbondingHeight,
		UnbondingCompletionTime: bv.UnbondingCompletionTime,
		MinSelfDelegation:       bv.MinSelfDelegation,
	}
	return nil
}
