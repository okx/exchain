package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryPool       = "pool"
	QueryPools      = "pools"
	QueryEarnings   = "earnings"
	QueryParameters = "parameters"
	QueryWhitelist  = "whitelist"
	QueryAccount    = "account"
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

// QueryEarningsParams defines the params for the following queries:
// - 'custom/farm/earnings'
type QueryEarningsParams struct {
	PoolName   string
	AccAddress sdk.AccAddress
}

// NewQueryEarningsParams creates a new instance of QueryEarningsParams
func NewQueryEarningsParams(poolName string, accAddr sdk.AccAddress) QueryEarningsParams {
	return QueryEarningsParams{
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
