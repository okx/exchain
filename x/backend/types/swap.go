package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	// watchlist sort column
	SwapWatchlistLiquidity = "liquidity"
	SwapWatchlistVolume24h = "volume24h"
	SwapWatchlistFeeApy    = "fee_apy"
	SwapWatchlistLastPrice = "last_price"
	SwapWatchlistChange24h = "change24h"

	// sort direction
	SwapWatchlistSortAsc = "asc"

	// query key
	QuerySwapWatchlist          = "swapWatchlist"
	QuerySwapTokens             = "swapTokens"
	QuerySwapTokenPairs         = "swapTokenPairs"
	QuerySwapLiquidityHistories = "swapLiquidityHistories"

	// swap business type
	SwapBusinessTypeCreate = "create"
	SwapBusinessTypeAdd    = "add"
	SwapBusinessTypeSwap   = "swap"
)

// nolint
type QuerySwapWatchlistParams struct {
	SortColumn    string `json:"sort_column"`
	SortDirection string `json:"sort_direction"`
	Page          int    `json:"page"`
	PerPage       int    `json:"per_page"`
}

// NewQuerySwapWatchlistParams creates a new instance of QuerySwapWatchlistParams
func NewQuerySwapWatchlistParams(sortColumn string, sortDirection string, page int, perPage int) QuerySwapWatchlistParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QuerySwapWatchlistParams{
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
		Page:          page,
		PerPage:       perPage,
	}
}

type SwapVolumePriceInfo struct {
	Volume    sdk.Dec
	Price24h  sdk.Dec
	Timestamp int64
}

type SwapWatchlist struct {
	SwapPair  string  `json:"swap_pair"`
	Liquidity sdk.Dec `json:"liquidity"`
	Volume24h sdk.Dec `json:"volume24h"`
	FeeApy    sdk.Dec `json:"fee_apy"`
	LastPrice sdk.Dec `json:"last_price"`
	Change24h sdk.Dec `json:"change24h"`
}

type SwapWatchlistSorter struct {
	Watchlist     []SwapWatchlist
	SortField     string
	SortDirectory string
}

func (s *SwapWatchlistSorter) Len() int { return len(s.Watchlist) }

func (s *SwapWatchlistSorter) Less(i, j int) bool {
	isSortAsc := false
	if s.SortDirectory == SwapWatchlistSortAsc {
		isSortAsc = true
	}

	switch s.SortField {
	case SwapWatchlistLiquidity:
		if isSortAsc {
			return s.Watchlist[i].Liquidity.LT(s.Watchlist[j].Liquidity)
		} else {
			return s.Watchlist[i].Liquidity.GT(s.Watchlist[j].Liquidity)
		}
	case SwapWatchlistVolume24h:
		if isSortAsc {
			return s.Watchlist[i].Volume24h.LT(s.Watchlist[j].Volume24h)
		} else {
			return s.Watchlist[i].Volume24h.GT(s.Watchlist[j].Volume24h)
		}
	case SwapWatchlistFeeApy:
		if isSortAsc {
			return s.Watchlist[i].FeeApy.LT(s.Watchlist[j].FeeApy)
		} else {
			return s.Watchlist[i].FeeApy.GT(s.Watchlist[j].FeeApy)
		}
	case SwapWatchlistLastPrice:
		if isSortAsc {
			return s.Watchlist[i].LastPrice.LT(s.Watchlist[j].LastPrice)
		} else {
			return s.Watchlist[i].LastPrice.GT(s.Watchlist[j].LastPrice)
		}
	case SwapWatchlistChange24h:
		if isSortAsc {
			return s.Watchlist[i].Change24h.LT(s.Watchlist[j].Change24h)
		} else {
			return s.Watchlist[i].Change24h.GT(s.Watchlist[j].Change24h)
		}
	}
	return false
}
func (s *SwapWatchlistSorter) Swap(i, j int) {
	s.Watchlist[i], s.Watchlist[j] = s.Watchlist[j], s.Watchlist[i]
}

type SwapInfo struct {
	Address          string `grom:"index;"`
	TokenPairName    string `gorm:"index;"`
	BaseTokenAmount  string `gorm:"type:varchar(40)"`
	QuoteTokenAmount string `gorm:"type:varchar(40)"`
	SellAmount       string `gorm:"type:varchar(40)"`
	BuysAmount       string `gorm:"type:varchar(40)"`
	Price            string `gorm:"type:varchar(40)"`
	Timestamp        int64  `gorm:"index;"`
}

type SwapWhitelist struct {
	Id            uint64 `gorm:"primaryKey`
	TokenPairName string `gorm:"index;type:varchar(128)"`
	Deleted       bool   `gorm:"type:bool"`
	Timestamp     int64  `gorm:""`
}

// nolint
type QuerySwapTokensParams struct {
	BusinessType  string `json:"business_type"`
	Address       string `json:"address"`
	BaseTokenName string `json:"base_token_name"`
}

// NewQuerySwapTokensParams creates a new instance of QueryDexFeesParams
func NewQuerySwapTokensParams(businessType string, address string, baseTokenName string) QuerySwapTokensParams {
	return QuerySwapTokensParams{
		BusinessType:  businessType,
		Address:       address,
		BaseTokenName: baseTokenName,
	}
}

type SwapToken struct {
	Symbol    string  `json:"symbol"`
	Available sdk.Dec `json:"available"`
}

func NewSwapToken(symbol string, available sdk.Dec) SwapToken {
	return SwapToken{
		Symbol:    symbol,
		Available: available,
	}
}

type SwapTokens []SwapToken

type SwapTokensResponse struct {
	NativeToken string     `json:"native_token"`
	Tokens      SwapTokens `json:"tokens"`
}

func (swapTokens SwapTokens) Len() int { return len(swapTokens) }

func (swapTokens SwapTokens) Less(i, j int) bool {
	return swapTokens[i].Available.GT(swapTokens[j].Available)
}

func (swapTokens SwapTokens) Swap(i, j int) {
	swapTokens[i], swapTokens[j] = swapTokens[j], swapTokens[i]
}

// nolint
type QuerySwapLiquidityInfoParams struct {
	Address       string `json:"address"`
	TokenPairName string `json:"token_pair_name"`
}

// NewQuerySwapBuyInfoParams creates a new instance of QuerySwapLiquidityInfoParams
func NewQuerySwapLiquidityInfoParams(address string, tokenPairName string) QuerySwapLiquidityInfoParams {
	return QuerySwapLiquidityInfoParams{
		Address:       address,
		TokenPairName: tokenPairName,
	}
}

type SwapLiquidityInfo struct {
	BasePooledCoin  sdk.SysCoin `json:"base_pooled_coin"`
	QuotePooledCoin sdk.SysCoin `json:"quote_pooled_coin"`
	PoolTokenCoin   sdk.SysCoin `json:"pool_token_coin"`
	PoolTokenRatio  sdk.Dec     `json:"pool_token_ratio"`
}
