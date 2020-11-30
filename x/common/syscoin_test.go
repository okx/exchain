package common

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------------------------------------------------
// SysCoin tests

var (
	syscoinTestDenom1 = "eos"
	syscoinTestDenom2 = "okt"
)

func TestCoin(t *testing.T) {
	require.Panics(t, func() { sdk.NewInt64Coin(syscoinTestDenom1, -1) })
	require.Panics(t, func() { sdk.NewCoin(syscoinTestDenom1, sdk.NewInt(-1)) })
	require.Panics(t, func() { sdk.NewInt64Coin(strings.ToUpper(syscoinTestDenom1), 10) })
	require.Panics(t, func() { sdk.NewCoin(strings.ToUpper(syscoinTestDenom1), sdk.NewInt(10)) })
	require.Equal(t, sdk.NewDec(5), sdk.NewInt64Coin(syscoinTestDenom1, 5).Amount)
	require.Equal(t, sdk.NewDec(5), sdk.NewCoin(syscoinTestDenom1, sdk.NewInt(5)).Amount)
}

func TestNewDecCoinFromDec(t *testing.T) {
	require.NotPanics(t, func() {
		sdk.NewDecCoinFromDec(syscoinTestDenom1, sdk.NewDec(5))
	})
	require.NotPanics(t, func() {
		sdk.NewDecCoinFromDec(syscoinTestDenom1, sdk.ZeroDec())
	})
	require.Panics(t, func() {
		sdk.NewDecCoinFromDec(strings.ToUpper(syscoinTestDenom1), sdk.NewDec(5))
	})
	require.Panics(t, func() {
		sdk.NewDecCoinFromDec(syscoinTestDenom1, sdk.NewDec(-5))
	})
}

func TestDecCoinIsPositive(t *testing.T) {
	dc := sdk.NewInt64DecCoin(syscoinTestDenom1, 5)
	require.True(t, dc.IsPositive())

	dc = sdk.NewInt64DecCoin(syscoinTestDenom1, 0)
	require.False(t, dc.IsPositive())
}

func TestIsEqualCoin(t *testing.T) {
	cases := []struct {
		inputOne sdk.SysCoin
		inputTwo sdk.SysCoin
		expected bool
		panics   bool
	}{
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), true, false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom2, 1), false, true},
		{sdk.NewInt64Coin("stake", 1), sdk.NewInt64Coin("stake", 10), false, false},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsEqual(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsEqual(tc.inputTwo)
			require.Equal(t, tc.expected, res, "coin equality relation is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestCoinIsValid(t *testing.T) {
	cases := []struct {
		coin       sdk.SysCoin
		expectPass bool
	}{
		{sdk.SysCoin{Denom: syscoinTestDenom1, Amount: sdk.NewDec(-1)}, false},
		{sdk.SysCoin{Denom: syscoinTestDenom1, Amount: sdk.NewDec(0)}, true},
		{sdk.SysCoin{Denom: syscoinTestDenom1, Amount: sdk.NewDec(1)}, true},
		{sdk.SysCoin{Denom: "Atom", Amount: sdk.NewDec(1)}, false},
		{sdk.SysCoin{Denom: "a", Amount: sdk.NewDec(1)}, true},
		{sdk.SysCoin{Denom: "a very long coin denom", Amount: sdk.NewDec(1)}, false},
		{sdk.SysCoin{Denom: "atOm", Amount: sdk.NewDec(1)}, false},
		{sdk.SysCoin{Denom: "     ", Amount: sdk.NewDec(1)}, false},
	}

	for i, tc := range cases {
		require.Equal(t, tc.expectPass, tc.coin.IsValid(), "unexpected result for IsValid, tc #%d", i)
	}
}

func TestAddCoin(t *testing.T) {
	cases := []struct {
		inputOne    sdk.SysCoin
		inputTwo    sdk.SysCoin
		expected    sdk.SysCoin
		shouldPanic bool
	}{
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 2), false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 0), sdk.NewInt64Coin(syscoinTestDenom1, 1), false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom2, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), true},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Add(tc.inputTwo) })
		} else {
			res := tc.inputOne.Add(tc.inputTwo)
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestSubCoin(t *testing.T) {
	cases := []struct {
		inputOne    sdk.SysCoin
		inputTwo    sdk.SysCoin
		expected    sdk.SysCoin
		shouldPanic bool
	}{
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom2, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), true},
		{sdk.NewInt64Coin(syscoinTestDenom1, 10), sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 9), false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 5), sdk.NewInt64Coin(syscoinTestDenom1, 3), sdk.NewInt64Coin(syscoinTestDenom1, 2), false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 5), sdk.NewInt64Coin(syscoinTestDenom1, 0), sdk.NewInt64Coin(syscoinTestDenom1, 5), false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 5), sdk.SysCoin{}, true},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Sub(tc.inputTwo) })
		} else {
			res := tc.inputOne.Sub(tc.inputTwo)
			require.Equal(t, tc.expected, res, "difference of coins is incorrect, tc #%d", tcIndex)
		}
	}

	tc := struct {
		inputOne sdk.SysCoin
		inputTwo sdk.SysCoin
		expected int64
	}{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), 0}
	res := tc.inputOne.Sub(tc.inputTwo)
	require.Equal(t, tc.expected, res.Amount.Int64())
}

func TestIsGTECoin(t *testing.T) {
	cases := []struct {
		inputOne sdk.SysCoin
		inputTwo sdk.SysCoin
		expected bool
		panics   bool
	}{
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), true, false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 2), sdk.NewInt64Coin(syscoinTestDenom1, 1), true, false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom2, 1), false, true},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsGTE(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsGTE(tc.inputTwo)
			require.Equal(t, tc.expected, res, "coin GTE relation is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestIsLTCoin(t *testing.T) {
	cases := []struct {
		inputOne sdk.SysCoin
		inputTwo sdk.SysCoin
		expected bool
		panics   bool
	}{
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), false, false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 2), sdk.NewInt64Coin(syscoinTestDenom1, 1), false, false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 0), sdk.NewInt64Coin(syscoinTestDenom2, 1), false, true},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom2, 1), false, true},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 1), false, false},
		{sdk.NewInt64Coin(syscoinTestDenom1, 1), sdk.NewInt64Coin(syscoinTestDenom1, 2), true, false},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsLT(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsLT(tc.inputTwo)
			require.Equal(t, tc.expected, res, "coin LT relation is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestCoinIsZero(t *testing.T) {
	coin := sdk.NewInt64Coin(syscoinTestDenom1, 0)
	res := coin.IsZero()
	require.True(t, res)

	coin = sdk.NewInt64Coin(syscoinTestDenom1, 1)
	res = coin.IsZero()
	require.False(t, res)
}
