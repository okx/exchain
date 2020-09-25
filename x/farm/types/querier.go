package types

const (
	QueryPool      = "pool"
	QueryPools     = "pools"
	QueryEarnings  = "earnings"
	QueryParams    = "params"
	QueryWhitelist = "whitelist"
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
