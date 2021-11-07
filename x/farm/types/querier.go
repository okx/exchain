package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

const (
	QueryPool             = "pool"
	QueryPools            = "pools"
	QueryEarnings         = "earnings"
	QueryLockInfo         = "lock-info"
	QueryParameters       = "parameters"
	QueryWhitelist        = "whitelist"
	QueryAccount          = "account"
	QueryAccountsLockedTo = "accounts-locked-to"
	QueryPoolNum          = "pool-num"
)

// QueryPoolParams defines the params for the following queries:
// - 'custom/farm/pool'
type QueryPoolParams struct {
	PoolName string
}

// NewQueryPoolParams creates a new instance of QueryPoolParams
func NewQueryPoolParams(poolName string) QueryPoolParams {
	return QueryPoolParams{
		PoolName: poolName,
	}
}

// QueryPoolsParams defines the params for the following queries:
// - 'custom/farm/pools'
type QueryPoolsParams struct {
	Page, Limit int
}

// NewQueryPoolsParams creates a new instance of QueryPoolsParams
func NewQueryPoolsParams(page, limit int) QueryPoolsParams {
	return QueryPoolsParams{
		Page:  page,
		Limit: limit,
	}
}

// QueryPoolAccountParams defines the params for the following queries:
// - 'custom/farm/earnings'
// - 'custom/farm/lock-info'
type QueryPoolAccountParams struct {
	PoolName   string
	AccAddress sdk.AccAddress
}

// NewQueryPoolAccountParams creates a new instance of QueryPoolAccountParams
func NewQueryPoolAccountParams(poolName string, accAddr sdk.AccAddress) QueryPoolAccountParams {
	return QueryPoolAccountParams{
		PoolName:   poolName,
		AccAddress: accAddr,
	}
}

// QueryAccountParams defines the params for the following queries:
// - 'custom/farm/account'
type QueryAccountParams struct {
	AccAddress sdk.AccAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams
func NewQueryAccountParams(accAddr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{
		AccAddress: accAddr,
	}
}
