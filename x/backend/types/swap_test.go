package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func TestSwapWatchlistSorter(t *testing.T) {
	watchlist := []SwapWatchlist{
		{
			SwapPair:  "pair0",
			Liquidity: sdk.NewDec(3),
			Volume24h: sdk.NewDec(3),
			FeeApy:    sdk.NewDec(3),
			LastPrice: sdk.NewDec(3),
			Change24h: sdk.NewDec(3),
		},
		{
			SwapPair:  "pair1",
			Liquidity: sdk.NewDec(1),
			Volume24h: sdk.NewDec(1),
			FeeApy:    sdk.NewDec(1),
			LastPrice: sdk.NewDec(1),
			Change24h: sdk.NewDec(1),
		},
		{
			SwapPair:  "pair2",
			Liquidity: sdk.NewDec(2),
			Volume24h: sdk.NewDec(2),
			FeeApy:    sdk.NewDec(2),
			LastPrice: sdk.NewDec(2),
			Change24h: sdk.NewDec(2),
		},
		{
			SwapPair:  "pair3",
			Liquidity: sdk.NewDec(4),
			Volume24h: sdk.NewDec(4),
			FeeApy:    sdk.NewDec(4),
			LastPrice: sdk.NewDec(4),
			Change24h: sdk.NewDec(4),
		},
	}

	tests := []struct {
		sorter *SwapWatchlistSorter
		want   []string
	}{
		{ // sort by liquidity asc
			sorter: &SwapWatchlistSorter{
				Watchlist:     watchlist,
				SortField:     SwapWatchlistLiquidity,
				SortDirectory: SwapWatchlistSortAsc,
			},
			want: []string{"pair1", "pair2", "pair0", "pair3"},
		},
		{ // sort by liquidity desc
			sorter: &SwapWatchlistSorter{
				Watchlist: watchlist,
				SortField: SwapWatchlistLiquidity,
			},
			want: []string{"pair3", "pair0", "pair2", "pair1"},
		},
		{ // sort by volume24h asc
			sorter: &SwapWatchlistSorter{
				Watchlist:     watchlist,
				SortField:     SwapWatchlistVolume24h,
				SortDirectory: SwapWatchlistSortAsc,
			},
			want: []string{"pair1", "pair2", "pair0", "pair3"},
		},
		{ // sort by volume24h desc
			sorter: &SwapWatchlistSorter{
				Watchlist: watchlist,
				SortField: SwapWatchlistVolume24h,
			},
			want: []string{"pair3", "pair0", "pair2", "pair1"},
		},
		{ // sort by fee apy asc
			sorter: &SwapWatchlistSorter{
				Watchlist:     watchlist,
				SortField:     SwapWatchlistFeeApy,
				SortDirectory: SwapWatchlistSortAsc,
			},
			want: []string{"pair1", "pair2", "pair0", "pair3"},
		},
		{ // sort by fee apy desc
			sorter: &SwapWatchlistSorter{
				Watchlist: watchlist,
				SortField: SwapWatchlistFeeApy,
			},
			want: []string{"pair3", "pair0", "pair2", "pair1"},
		},
		{ // sort by last price asc
			sorter: &SwapWatchlistSorter{
				Watchlist:     watchlist,
				SortField:     SwapWatchlistLastPrice,
				SortDirectory: SwapWatchlistSortAsc,
			},
			want: []string{"pair1", "pair2", "pair0", "pair3"},
		},
		{ // sort by last price desc
			sorter: &SwapWatchlistSorter{
				Watchlist: watchlist,
				SortField: SwapWatchlistLastPrice,
			},
			want: []string{"pair3", "pair0", "pair2", "pair1"},
		},
		{ // sort by change24h asc
			sorter: &SwapWatchlistSorter{
				Watchlist:     watchlist,
				SortField:     SwapWatchlistChange24h,
				SortDirectory: SwapWatchlistSortAsc,
			},
			want: []string{"pair1", "pair2", "pair0", "pair3"},
		},
		{ // sort by change24h desc
			sorter: &SwapWatchlistSorter{
				Watchlist: watchlist,
				SortField: SwapWatchlistChange24h,
			},
			want: []string{"pair3", "pair0", "pair2", "pair1"},
		},
	}

	for _, test := range tests {
		sort.Sort(test.sorter)
		var sortedNames []string
		for _, watchlist := range test.sorter.Watchlist {
			sortedNames = append(sortedNames, watchlist.SwapPair)
		}
		require.Equal(t, test.want, sortedNames)
	}
}
