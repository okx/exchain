package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FarmPoolStatus int

const (
	WhitelistFarmPool = "whitelist"
	NormalFarmPool    = "normal"

	BlockInterval = 3
	BlocksPerDay  = 24 * 60 * 60 / BlockInterval
	DaysInYear    = 365

	// farm pool status
	FarmPoolCreated  FarmPoolStatus = 1
	FarmPoolProvided FarmPoolStatus = 2
	FarmPoolYielded  FarmPoolStatus = 3
	FarmPoolFinished FarmPoolStatus = 4
)

// nolint
type QueryFarmPoolsParams struct {
	PoolType string `json:"pool_type"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
}

// NewQueryFarmPoolsParams creates a new instance of QueryFarmPoolsParams
func NewQueryFarmPoolsParams(poolType string, page int, perPage int) QueryFarmPoolsParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryFarmPoolsParams{
		PoolType: poolType,
		Page:     page,
		PerPage:  perPage,
	}
}

// nolint
type QueryFarmDashboardParams struct {
	Address string `json:"address"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

// NewQueryFarmDashboardParams creates a new instance of QueryFarmDashboardParams
func NewQueryFarmDashboardParams(address string, page int, perPage int) QueryFarmDashboardParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryFarmDashboardParams{
		Address: address,
		Page:    page,
		PerPage: perPage,
	}
}

type FarmPoolResponse struct {
	PoolName      string         `json:"pool_name"`
	LockSymbol    string         `json:"lock_symbol"`
	YieldSymbol   string         `json:"yield_symbol"`
	TotalStaked   sdk.Dec        `json:"total_staked"`
	StartAt       int64          `json:"start_at"`
	FinishAt      int64          `json:"finish_at"`
	PoolRate      sdk.SysCoins   `json:"pool_rate"`
	FarmApy       sdk.SysCoins   `json:"farm_apy"`
	PoolRatio     sdk.Dec        `json:"pool_ratio"`
	InWhitelist   bool           `json:"in_whitelist"`
	TotalFarmed   sdk.Dec        `json:"total_farmed"`
	FarmedDetails []FarmInfo     `json:"farmed_details"`
	Status        FarmPoolStatus `json:"status"`
}

type FarmInfo struct {
	Symbol    string  `json:"symbol"`
	Claimed   sdk.Dec `json:"claimed"`
	UnClaimed sdk.Dec `json:"unclaimed"`
}

type FarmResponseList []FarmPoolResponse

func (list FarmResponseList) Len() int { return len(list) }
func (list FarmResponseList) Less(i, j int) bool {
	return list[i].TotalStaked.GT(list[j].TotalStaked)
}
func (list FarmResponseList) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

// nolint
type QueryFarmStakedInfoParams struct {
	PoolName string `json:"pool_name"`
	Address  string `json:"address"`
}

// NewQueryFarmStakedInfoParams creates a new instance of QueryFarmStakedInfoParams
func NewQueryFarmStakedInfoParams(poolName string, address string) QueryFarmStakedInfoParams {
	return QueryFarmStakedInfoParams{
		PoolName: poolName,
		Address:  address,
	}
}

type FarmStakedInfo struct {
	PoolName        string  `json:"pool_name"`
	Balance         sdk.Dec `json:"balance"`
	AccountStaked   sdk.Dec `json:"account_staked"`
	PoolTotalStaked sdk.Dec `json:"pool_total_staked"`
	PoolRatio       sdk.Dec `json:"pool_ratio"`
	MinLockAmount   sdk.Dec `json:"min_lock_amount"`
}
