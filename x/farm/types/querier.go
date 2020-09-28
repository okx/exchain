package types

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
