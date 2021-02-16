package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FarmPoolStatus int

const (
	WhitelistFarmPool = "whitelist"
	NormalFarmPool    = "normal"

	SecondsInDay  = 24 * 60 * 60
	DaysInYear    = 365
	BlockInterval = 3
	BlocksPerDay  = SecondsInADay / BlockInterval

	// farm pool status
	FarmPoolCreated  FarmPoolStatus = 1
	FarmPoolProvided FarmPoolStatus = 2
	FarmPoolYielded  FarmPoolStatus = 3
	FarmPoolFinished FarmPoolStatus = 4

	// query key
	QueryFarmPools      = "farmPools"
	QueryFarmDashboard  = "farmDashboard"
	QueryFarmMaxApy     = "farmMaxApy"
	QueryFarmStakedInfo = "farmStakedInfo"
	QueryFarmFirstPool  = "farmFirstPool"

	// farm sort column
	FarmPoolTotalStaked = "total_staked"
	FarmPoolApy         = "farm_apy"
	FarmPoolStartAt     = "start_at"
	FarmPoolFinishAt    = "finish_at"

	// sort direction
	FarmSortAsc = "asc"
)

// nolint
type QueryFarmPoolsParams struct {
	PoolType      string `json:"pool_type"`
	SortColumn    string `json:"sort_column"`
	SortDirection string `json:"sort_direction"`
	Page          int    `json:"page"`
	PerPage       int    `json:"per_page"`
}

// NewQueryFarmPoolsParams creates a new instance of QueryFarmPoolsParams
func NewQueryFarmPoolsParams(poolType string, sortColumn string, sortDirection string, page int, perPage int) QueryFarmPoolsParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	// default sort by FarmPoolTotalStaked
	if sortColumn == "" {
		sortColumn = FarmPoolTotalStaked
	}
	return QueryFarmPoolsParams{
		PoolType:      poolType,
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
		Page:          page,
		PerPage:       perPage,
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
	UserStaked    sdk.Dec        `json:"user_staked"`
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

func (farmPool FarmPoolResponse) TotalApy() sdk.Dec {
	sum := sdk.ZeroDec()
	for _, coin := range farmPool.FarmApy {
		sum = sum.Add(coin.Amount)
	}
	return sum
}

type FarmInfo struct {
	Symbol    string  `json:"symbol"`
	Claimed   sdk.Dec `json:"claimed"`
	UnClaimed sdk.Dec `json:"unclaimed"`
}

type FarmResponseList []FarmPoolResponse

type FarmResponseListSorter struct {
	FarmPoolList  FarmResponseList
	SortField     string
	SortDirectory string
}

func (s *FarmResponseListSorter) Len() int { return len(s.FarmPoolList) }

func (s *FarmResponseListSorter) Less(i, j int) bool {
	isSortAsc := false
	if s.SortDirectory == FarmSortAsc {
		isSortAsc = true
	}

	switch s.SortField {
	case FarmPoolTotalStaked:
		if isSortAsc {
			return s.FarmPoolList[i].TotalStaked.LT(s.FarmPoolList[j].TotalStaked)
		} else {
			return s.FarmPoolList[i].TotalStaked.GT(s.FarmPoolList[j].TotalStaked)
		}
	case FarmPoolApy:
		if isSortAsc {
			return s.FarmPoolList[i].TotalApy().LT(s.FarmPoolList[j].TotalApy())
		} else {
			return s.FarmPoolList[i].TotalApy().GT(s.FarmPoolList[j].TotalApy())
		}
	case FarmPoolStartAt:
		if isSortAsc {
			return s.FarmPoolList[i].StartAt < s.FarmPoolList[j].StartAt
		} else {
			return s.FarmPoolList[i].StartAt > s.FarmPoolList[j].StartAt
		}
	case FarmPoolFinishAt:
		if isSortAsc {
			return s.FarmPoolList[i].FinishAt < s.FarmPoolList[j].FinishAt
		} else {
			return s.FarmPoolList[i].FinishAt > s.FarmPoolList[j].FinishAt
		}
	}
	return false
}
func (s *FarmResponseListSorter) Swap(i, j int) {
	s.FarmPoolList[i], s.FarmPoolList[j] = s.FarmPoolList[j], s.FarmPoolList[i]
}

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

type ClaimInfo struct {
	Id        uint64 `gorm:"primaryKey`
	PoolName  string `grom:"index;"`
	Address   string `grom:"index;"`
	Claimed   string `gorm:"type:varchar(256)"`
	Timestamp int64  `gorm:"index;"`
}

// nolint
type QueryFarmFirstPoolParams struct {
	PoolName    string `json:"pool_name"`
	Address     string `json:"address"`
	StakeAt     int64  `json:"stake_at"`
	ClaimHeight int64  `json:"claim_height"`
}

// NewQueryFarmFirstPoolParams creates a new instance of QueryFarmFirstPoolParams
func NewQueryFarmFirstPoolParams(poolName string, address string, stakeAt int64, claimHeight int64) QueryFarmFirstPoolParams {
	return QueryFarmFirstPoolParams{
		PoolName:    poolName,
		Address:     address,
		StakeAt:     stakeAt,
		ClaimHeight: claimHeight,
	}
}

type FarmFirstPool struct {
	FarmApy       sdk.Dec `json:"farm_apy"`
	FarmAmount    sdk.Dec `json:"farm_amount"`
	TotalStaked   sdk.Dec `json:"total_staked"`
	ClaimAt       int64   `json:"claim_at"`
	AccountStaked sdk.Dec `json:"account_staked"`
	EstimatedFarm sdk.Dec `json:"estimated_farm"`
	Balance       sdk.Dec `json:"balance"`
}
