package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
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

// LastValidatorPower is needed for validator set update logic
type LastValidatorPower struct {
	Address sdk.ValAddress
	Power   int64
}

// NewLastValidatorPower creates a new instance of LastValidatorPower
func NewLastValidatorPower(valAddr sdk.ValAddress, power int64) LastValidatorPower {
	return LastValidatorPower{
		valAddr,
		power,
	}
}

// NewGenesisState creates a new object of GenesisState
func NewGenesisState(params Params, validators Validators, delegators []Delegator) GenesisState {
	return GenesisState{
		Params:     params,
		Validators: validators.Export(),
		Delegators: delegators,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidatorExported is designed for Validator export
type ValidatorExported struct {
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

// Import converts validator exported format to inner one by filling the zero-value of Tokens and Commission
func (ve ValidatorExported) Import() Validator {
	consPk, err := sdk.GetConsPubKeyBech32(ve.ConsPubKey)
	if err != nil {
		panic(fmt.Sprintf("failed. consensus pubkey is parsed error: %s", err.Error()))
	}

	return Validator{
		ve.OperatorAddress,
		consPk,
		ve.Jailed,
		ve.Status,
		sdk.NewInt(0),
		ve.DelegatorShares,
		ve.Description,
		ve.UnbondingHeight,
		ve.UnbondingCompletionTime,
		NewCommission(sdk.NewDec(1), sdk.NewDec(1), sdk.NewDec(0)),
		ve.MinSelfDelegation,
	}
}

// ConsAddress returns the TM validator address of exported validator
func (ve ValidatorExported) ConsAddress() sdk.ConsAddress {
	consPk, err := sdk.GetConsPubKeyBech32(ve.ConsPubKey)
	if err != nil {
		panic(fmt.Sprintf("failed. consensus pubkey is parsed error: %s", err.Error()))
	}

	return sdk.ConsAddress(consPk.Address())
}

// IsBonded checks if the exported validator status equals Bonded
func (ve ValidatorExported) IsBonded() bool {
	return ve.Status.Equal(sdk.Bonded)
}

// ProxyDelegatorKeyExported is designed for ProxyDelegatorKey export
type ProxyDelegatorKeyExported struct {
	DelAddr   sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ProxyAddr sdk.AccAddress `json:"proxy_address" yaml:"proxy_address"`
}

// NewProxyDelegatorKeyExported creates a new object of ProxyDelegatorKeyExported
func NewProxyDelegatorKeyExported(delAddr, proxyAddr sdk.AccAddress) ProxyDelegatorKeyExported {
	return ProxyDelegatorKeyExported{
		delAddr,
		proxyAddr,
	}
}

// VotesExported is designed for types.Votes export
type VotesExported struct {
	VoterAddress     sdk.AccAddress `json:"voter_address" yaml:"voter_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Votes            Votes          `json:"votes" yaml:"votes"`
}

// NewVoteExported creates a new object of VotesExported
func NewVoteExported(voterAddr sdk.AccAddress, valAddr sdk.ValAddress, votes Votes) VotesExported {
	return VotesExported{
		voterAddr,
		valAddr,
		votes,
	}
}
