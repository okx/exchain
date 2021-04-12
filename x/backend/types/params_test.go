package types

import (
	"github.com/okex/exchain/x/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQueryDealsParams(t *testing.T) {
	want := QueryDealsParams{
		Address: "address",
		Product: "btc_" + common.NativeToken,
		Start:   5,
		End:     9,
		Page:    3,
		PerPage: 17,
		Side:    "side",
	}

	got := NewQueryDealsParams(want.Address, want.Product, want.Start, want.End, want.Page, want.PerPage, want.Side)
	require.Equal(t, want, got)

	want = QueryDealsParams{
		Address: "address",
		Product: "btc_" + common.NativeToken,
		Start:   5,
		End:     9,
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
		Side:    "side",
	}

	got = NewQueryDealsParams(want.Address, want.Product, want.Start, want.End, 0, 0, want.Side)
	require.Equal(t, want, got)
}

func TestNewQueryFeeDetailsParams(t *testing.T) {
	want := QueryFeeDetailsParams{
		Address: "address",
		Page:    10,
		PerPage: 10,
	}
	got := NewQueryFeeDetailsParams(want.Address, want.Page, want.PerPage)
	require.Equal(t, want, got)

	want = QueryFeeDetailsParams{
		Address: "address",
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
	}
	got = NewQueryFeeDetailsParams(want.Address, 0, 0)
	require.Equal(t, want, got)
}

func TestNewQueryKlinesParams(t *testing.T) {
	want := QueryKlinesParams{
		Product:     "okb_btc",
		Granularity: 3,
		Size:        9,
	}
	got := NewQueryKlinesParams(want.Product, want.Granularity, want.Size)
	require.Equal(t, want, got)
}

func TestNewQueryMatchParams(t *testing.T) {
	want := QueryMatchParams{
		Product: "btc_" + common.NativeToken,
		Start:   0,
		End:     10,
		Page:    10,
		PerPage: 30,
	}
	got := NewQueryMatchParams(want.Product, want.Start, want.End, want.Page,
		want.PerPage)
	require.Equal(t, want, got)

	want = QueryMatchParams{
		Product: "btc_" + common.NativeToken,
		Start:   0,
		End:     10,
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
	}
	got = NewQueryMatchParams(want.Product, want.Start, want.End, 0, 0)
	require.Equal(t, want, got)
}

func TestNewQueryOrderListParams(t *testing.T) {
	want := QueryOrderListParams{
		Address:    "address",
		Product:    "okb_bch",
		Page:       1,
		PerPage:    33,
		Start:      2,
		End:        9,
		Side:       "side",
		HideNoFill: true,
	}
	got := NewQueryOrderListParams(want.Address, want.Product, want.Side,
		want.Page, want.PerPage, want.Start, want.End, want.HideNoFill)
	require.Equal(t, want, got)

	want = QueryOrderListParams{
		Address:    "address",
		Product:    "okb_bch",
		Page:       DefaultPage,
		PerPage:    DefaultPerPage,
		Start:      0,
		End:        0,
		Side:       "side",
		HideNoFill: true,
	}
	got = NewQueryOrderListParams(want.Address, want.Product, want.Side,
		0, 0, 0, 0, want.HideNoFill)
	require.Equal(t, DefaultPage, got.Page)
	require.Equal(t, DefaultPerPage, got.PerPage)
	require.Equal(t, int64(0), got.Start)
	require.NotEqual(t, int64(0), got.End)
}

func TestNewQueryTxListParams(t *testing.T) {
	want := QueryTxListParams{
		Address:   "address",
		TxType:    1,
		StartTime: 2384639292,
		EndTime:   2394334444,
		Page:      3,
		PerPage:   44,
	}
	got := NewQueryTxListParams(want.Address, want.TxType, want.StartTime,
		want.EndTime, want.Page, want.PerPage)
	require.Equal(t, want, got)

	want = QueryTxListParams{
		Address:   "address",
		TxType:    1,
		StartTime: 2384639292,
		EndTime:   2394334444,
		Page:      DefaultPage,
		PerPage:   DefaultPerPage,
	}
	got = NewQueryTxListParams(want.Address, want.TxType, want.StartTime,
		want.EndTime, 0, 0)
	require.Equal(t, want, got)
}
