package v0_10

import (
	"time"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

const ModuleName = "staking"

type (
	// GenesisState - all staking state that must be provided at genesis
	GenesisState struct {
		Params               Params                      `json:"params" yaml:"params"`
		LastTotalPower       sdk.Int                     `json:"last_total_power" yaml:"last_total_power"`
		LastValidatorPowers  []LastValidatorPower        `json:"last_validator_powers" yaml:"last_validator_powers"`
		Validators           []ValidatorExported         `json:"validators" yaml:"validators"`
		Delegators           []Delegator                 `json:"delegators" yaml:"delegators"`
		UnbondingDelegations []UndelegationInfo          `json:"unbonding_delegations" yaml:"unbonding_delegations"`
		Votes                []VotesExported             `json:"votes" yaml:"votes"`
		ProxyDelegatorKeys   []ProxyDelegatorKeyExported `json:"proxy_delegator_keys" yaml:"proxy_delegator_keys"`
		Exported             bool                        `json:"exported" yaml:"exported"`
	}

	// Params defines the high level settings for staking
	Params struct {
		// time duration of unbonding
		UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
		// note: we need to be a bit careful about potential overflow here, since this is user-determined
		// maximum number of validators (max uint16 = 65535)
		MaxValidators uint16 `json:"max_bonded_validators" yaml:"max_bonded_validators"`
		// epoch for validator update
		Epoch         uint16 `json:"epoch" yaml:"epoch"`
		MaxValsToVote uint16 `json:"max_validators_to_vote" yaml:"max_validators_to_vote"`
		// bondable coin denomination
		BondDenom string `json:"bond_denom" yaml:"bond_denom"`
		// limited amount of delegate
		MinDelegation sdk.Dec `json:"min_delegation" yaml:"min_delegation"`
	}

	// LastValidatorPower is needed for validator set update logic
	LastValidatorPower struct {
		Address sdk.ValAddress
		Power   int64
	}

	// ValidatorExported is designed for Validator export
	ValidatorExported struct {
		OperatorAddress         sdk.ValAddress `json:"operator_address"`
		ConsPubKey              string         `json:"consensus_pubkey"`
		Jailed                  bool           `json:"jailed"`
		Status                  sdk.BondStatus `json:"status" yaml:"status"`
		DelegatorShares         sdk.Dec        `json:"delegator_shares"`
		Description             Description    `json:"description"`
		UnbondingHeight         int64          `json:"unbonding_height"`
		UnbondingCompletionTime time.Time      `json:"unbonding_time"`
		MinSelfDelegation       sdk.Dec        `json:"min_self_delegation"`
	}

	// Description - description fields for a validator
	Description struct {
		Moniker  string `json:"moniker" yaml:"moniker"`   // name
		Identity string `json:"identity" yaml:"identity"` // optional identity signature (ex. UPort or Keybase)
		Website  string `json:"website" yaml:"website"`   // optional website link
		Details  string `json:"details" yaml:"details"`   // optional details
	}

	// Delegator is the struct of delegator info
	Delegator struct {
		DelegatorAddress     sdk.AccAddress   `json:"delegator_address" yaml:"delegator_address"`
		ValidatorAddresses   []sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
		Shares               sdk.Dec          `json:"shares" yaml:"shares"`
		Tokens               sdk.Dec          `json:"tokens" yaml:"tokens"`
		IsProxy              bool             `json:"is_proxy" yaml:"is_proxy"`
		TotalDelegatedTokens sdk.Dec          `json:"total_delegated_tokens" yaml:"total_delegated_tokens"`
		ProxyAddress         sdk.AccAddress   `json:"proxy_address" yaml:"proxy_address"`
	}

	// UndelegationInfo is the struct of the undelegation info
	UndelegationInfo struct {
		DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
		Quantity         sdk.Dec        `json:"quantity" yaml:"quantity"`
		CompletionTime   time.Time      `json:"completion_time"`
	}

	Votes = sdk.Dec
	// VotesExported is designed for types.Votes export
	VotesExported struct {
		VoterAddress     sdk.AccAddress `json:"voter_address" yaml:"voter_address"`
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
		Votes            Votes          `json:"votes" yaml:"votes"`
	}

	// ProxyDelegatorKeyExported is designed for ProxyDelegatorKey export
	ProxyDelegatorKeyExported struct {
		DelAddr   sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
		ProxyAddr sdk.AccAddress `json:"proxy_address" yaml:"proxy_address"`
	}
)
