package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators          = "validators"
	QueryValidator           = "validator"
	QueryUnbondingDelegation = "unbondingDelegation"
	QueryPool                = "pool"
	QueryParameters          = "parameters"
	QueryAddress             = "address"
	QueryForAddress          = "validatorAddress"
	QueryForAccAddress       = "validatorAccAddress"
	QueryProxy               = "proxy"
	QueryValidatorAllShares  = "validatorAllShares"
	QueryDelegator           = "delegator"
)

// QueryDelegatorParams defines the params for the following queries:
// - 'custom/staking/delegatorDelegations'
// - 'custom/staking/delegatorUnbondingDelegations'
// - 'custom/staking/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

// NewQueryDelegatorParams creates a new instance of QueryDelegatorParams
func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// QueryValidatorParams defines the params for the following queries:
// - 'custom/staking/validator'
// - 'custom/staking/validatorDelegations'
// - 'custom/staking/validatorUnbondingDelegations'
// - 'custom/staking/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
}

// NewQueryValidatorParams creates a new instance of QueryValidatorParams
func NewQueryValidatorParams(validatorAddr sdk.ValAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

//// QueryBondsParams defines the params for the following queries:
//// - 'custom/staking/delegation'
//// - 'custom/staking/unbondingDelegation'
//// - 'custom/staking/delegatorValidator'
//type QueryBondsParams struct {
//	DelegatorAddr sdk.AccAddress
//	ValidatorAddr sdk.ValAddress
//}
//
//// NewQueryBondsParams creates a new instance of QueryBondsParams
//func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryBondsParams {
//	return QueryBondsParams{
//		DelegatorAddr: delegatorAddr,
//		ValidatorAddr: validatorAddr,
//	}
//}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/staking/validators'
type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

// NewQueryValidatorsParams creates a new instance of QueryValidatorsParams
func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}
