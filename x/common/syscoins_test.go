package common

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Coins tests

var (
	syscoinsTestDenom1 = "eos"
	syscoinsTestDenom2 = "okt"
)

func TestIsZeroCoins(t *testing.T) {
	cases := []struct {
		inputOne sdk.SysCoins
		expected bool
	}{
		{sdk.SysCoins{}, true},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0)}, true},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 0)}, true},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 1)}, false},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 1)}, false},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsZero()
		require.Equal(t, tc.expected, res)
	}
}

func TestEqualCoins(t *testing.T) {
	cases := []struct {
		inputOne sdk.SysCoins
		inputTwo sdk.SysCoins
		expected bool
		panics   bool
	}{
		{sdk.SysCoins{}, sdk.SysCoins{}, true, false},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0)}, sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0)}, true, false},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 1)}, sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 1)}, true, false},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0)}, sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom2, 0)}, false, true},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0)}, sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 1)}, false, false},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0)}, sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 1)}, false, false},
		{sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 1)}, sdk.SysCoins{sdk.NewInt64Coin(syscoinsTestDenom1, 0), sdk.NewInt64Coin(syscoinsTestDenom2, 1)}, true, false},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsEqual(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsEqual(tc.inputTwo)
			require.Equal(t, tc.expected, res, "Equality is differed from exported. tc #%d, expected %b, actual %b.", tcIndex, tc.expected, res)
		}
	}
}

func TestAddCoins(t *testing.T) {
	cases := []struct {
		inputOne sdk.SysCoins
		inputTwo sdk.SysCoins
		expected sdk.SysCoins
		shouldPanic bool
	}{
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}, {syscoinsTestDenom2, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}, {syscoinsTestDenom2, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}, {syscoinsTestDenom2, sdk.NewDec(2)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(0)}, {syscoinsTestDenom2, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(0)}, {syscoinsTestDenom2, sdk.NewDec(0)}}, sdk.SysCoins{{syscoinsTestDenom2, sdk.NewDec(1)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}}, sdk.SysCoins{{syscoinsTestDenom2, sdk.NewDec(0)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}, {syscoinsTestDenom2, sdk.NewDec(2)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}, {syscoinsTestDenom2, sdk.NewDec(2)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(0)}, {syscoinsTestDenom2, sdk.NewDec(0)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(0)}, {syscoinsTestDenom2, sdk.NewDec(0)}}, sdk.SysCoins(nil), false},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Add2(tc.inputTwo)})
		} else {
			res := tc.inputOne.Add2(tc.inputTwo)
			assert.True(t, res.IsValid())
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestSubCoins(t *testing.T) {
	testCases := []struct {
		inputOne    sdk.SysCoins
		inputTwo    sdk.SysCoins
		expected    sdk.SysCoins
		shouldPanic bool
	}{
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}}, sdk.SysCoins{{syscoinsTestDenom2, sdk.NewDec(0)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom2, sdk.NewDec(0)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}, {syscoinsTestDenom2, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom2, sdk.NewDec(1)}}, false},
		{sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(1)}, {syscoinsTestDenom2, sdk.NewDec(1)}}, sdk.SysCoins{{syscoinsTestDenom1, sdk.NewDec(2)}}, sdk.SysCoins{}, true},
	}

	for tcIndex, tc := range testCases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Sub(tc.inputTwo) })
		} else {
			res := tc.inputOne.Sub(tc.inputTwo)
			assert.True(t, res.IsValid())
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestSortDecCoins(t *testing.T) {
	good := sdk.SysCoins{
		{"gas", sdk.NewDec(1)},
		{"mineral", sdk.NewDec(1)},
		{"tree", sdk.NewDec(1)},
	}
	mixedCase1 := sdk.SysCoins{
		{"gAs", sdk.NewDec(1)},
		{"MineraL", sdk.NewDec(1)},
		{"TREE", sdk.NewDec(1)},
	}
	mixedCase2 := sdk.SysCoins{
		{"gAs", sdk.NewDec(1)},
		{"mineral", sdk.NewDec(1)},
	}
	mixedCase3 := sdk.SysCoins{
		{"gAs", sdk.NewDec(1)},
	}
	empty := sdk.NewCoins()
	badSort1 := sdk.SysCoins{
		{"tree", sdk.NewDec(1)},
		{"gas", sdk.NewDec(1)},
		{"mineral", sdk.NewDec(1)},
	}

	// both are after the first one, but the second and third are in the wrong order
	badSort2 := sdk.SysCoins{
		{"gas", sdk.NewDec(1)},
		{"tree", sdk.NewDec(1)},
		{"mineral", sdk.NewDec(1)},
	}
	badAmt := sdk.SysCoins{
		{"gas", sdk.NewDec(1)},
		{"tree", sdk.NewDec(0)},
		{"mineral", sdk.NewDec(1)},
	}
	dup := sdk.SysCoins{
		{"gas", sdk.NewDec(1)},
		{"gas", sdk.NewDec(1)},
		{"mineral", sdk.NewDec(1)},
	}
	neg := sdk.SysCoins{
		{"gas", sdk.NewDec(-1)},
		{"mineral", sdk.NewDec(1)},
	}

	assert.True(t, good.IsValid(), "Coins are valid")
	assert.False(t, mixedCase1.IsValid(), "Coins denoms contain upper case characters")
	assert.False(t, mixedCase2.IsValid(), "First Coins denoms contain upper case characters")
	assert.False(t, mixedCase3.IsValid(), "Single denom in Coins contains upper case characters")
	assert.True(t, good.IsAllPositive(), "Expected coins to be positive: %v", good)
	assert.False(t, empty.IsAllPositive(), "Expected coins to not be positive: %v", empty)
	assert.True(t, good.IsAllGTE(empty), "Expected %v to be >= %v", good, empty)
	assert.False(t, good.IsAllLT(empty), "Expected %v to be < %v", good, empty)
	assert.True(t, empty.IsAllLT(good), "Expected %v to be < %v", empty, good)
	assert.False(t, badSort1.IsValid(), "Coins are not sorted")
	assert.False(t, badSort2.IsValid(), "Coins are not sorted")
	assert.False(t, badAmt.IsValid(), "Coins cannot include 0 amounts")
	assert.False(t, dup.IsValid(), "Duplicate coin")
	assert.False(t, neg.IsValid(), "Negative first-denom coin")
}

func TestCoinsGT(t *testing.T) {
	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	assert.False(t, sdk.SysCoins{}.IsAllGT(sdk.SysCoins{}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom2, two}}))
}

func TestCoinsLT(t *testing.T) {
	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	assert.False(t, sdk.SysCoins{}.IsAllLT(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllLT(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom1, two}, {syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{}.IsAllLT(sdk.SysCoins{{syscoinsTestDenom1, one}}))
}

func TestCoinsLTE(t *testing.T) {
	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	assert.True(t, sdk.SysCoins{}.IsAllLTE(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllLTE(sdk.SysCoins{}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{}.IsAllLTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
}

func TestParse(t *testing.T) {
	cases := []struct {
		input    string
		valid    bool  // if false, we expect an error on parse
		expected sdk.SysCoins // if valid is true, make sure this is returned
	}{
		{"", true, nil},
		{"1foo", true, sdk.SysCoins{{"foo", sdk.NewDec(1)}}},
		{"10bar", true, sdk.SysCoins{{"bar", sdk.NewDec(10)}}},
		{"99bar,1foo", true, sdk.SysCoins{{"bar", sdk.NewDec(99)}, {"foo", sdk.NewDec(1)}}},
		{"98 bar , 1 foo  ", true, sdk.SysCoins{{"bar", sdk.NewDec(98)}, {"foo", sdk.NewDec(1)}}},
		{"  55\t \t bling\n", true, sdk.SysCoins{{"bling", sdk.NewDec(55)}}},
		{"2foo, 97 bar", true, sdk.SysCoins{{"bar", sdk.NewDec(97)}, {"foo", sdk.NewDec(2)}}},
		{"5 mycoin,", false, nil},             // no empty coins in a list
		{"2 3foo, 97 bar", false, nil},        // 3foo is invalid coin name
		{"11me coin, 12you coin", false, nil}, // no spaces in coin names
		//{"1.2btc", true, sdk.Coins{{"btc", sdk.NewDec(1.2)}}},                // amount must be integer
		{"5foo-bar", false, nil},              // once more, only letters in coin name
	}

	for tcIndex, tc := range cases {
		res, err := sdk.ParseCoins(tc.input)
		if !tc.valid {
			require.NotNil(t, err, "%s: %#v. tc #%d", tc.input, res, tcIndex)
		} else if assert.Nil(t, err, "%s: %+v", tc.input, err) {
			require.Equal(t, tc.expected, res, "coin parsing was incorrect, tc #%d", tcIndex)
		}
	}
}

func TestSortCoins(t *testing.T) {
	good := sdk.SysCoins{
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("mineral", 1),
		sdk.NewInt64Coin("tree", 1),
	}
	empty := sdk.SysCoins{
		sdk.NewInt64Coin("gold", 0),
	}
	badSort1 := sdk.SysCoins{
		sdk.NewInt64Coin("tree", 1),
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("mineral", 1),
	}
	badSort2 := sdk.SysCoins{ // both are after the first one, but the second and third are in the wrong order
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("tree", 1),
		sdk.NewInt64Coin("mineral", 1),
	}
	badAmt := sdk.SysCoins{
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("tree", 0),
		sdk.NewInt64Coin("mineral", 1),
	}
	dup := sdk.SysCoins{
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("mineral", 1),
	}

	cases := []struct {
		coins         sdk.SysCoins
		before, after bool // valid before/after sort
	}{
		{good, true, true},
		{empty, false, false},
		{badSort1, false, true},
		{badSort2, false, true},
		{badAmt, false, false},
		{dup, false, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.before, tc.coins.IsValid(), "coin validity is incorrect before sorting, tc #%d", tcIndex)
		tc.coins.Sort()
		require.Equal(t, tc.after, tc.coins.IsValid(), "coin validity is incorrect after sorting, tc #%d", tcIndex)
	}
}

func TestAmountOf(t *testing.T) {
	case0 := sdk.SysCoins{}
	case1 := sdk.SysCoins{
		sdk.NewInt64Coin("gold", 0),
	}
	case2 := sdk.SysCoins{
		sdk.NewInt64Coin("gas", 1),
		sdk.NewInt64Coin("mineral", 1),
		sdk.NewInt64Coin("tree", 1),
	}
	case3 := sdk.SysCoins{
		sdk.NewInt64Coin("mineral", 1),
		sdk.NewInt64Coin("tree", 1),
	}
	case4 := sdk.SysCoins{
		sdk.NewInt64Coin("gas", 8),
	}

	cases := []struct {
		coins           sdk.SysCoins
		amountOf        int64
		amountOfSpace   int64
		amountOfGAS     int64
		amountOfMINERAL int64
		amountOfTREE    int64
	}{
		{case0, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000},
		{case1, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000},
		{case2, 0.0000000000000000000, 0.000000000000000000, 1.000000000000000000, 1.000000000000000000, 1.000000000000000000},
		{case3, 0.000000000000000000, 0.000000000000000000, 0.000000000000000000, 1.000000000000000000, 1.000000000000000000},
		{case4, 0.000000000000000000, 0.000000000000000000, 8.000000000000000000, 0.000000000000000000, 0.000000000000000000},
	}

	for _, tc := range cases {
		assert.Equal(t, sdk.NewDec(tc.amountOfGAS), tc.coins.AmountOf("gas"))
		assert.Equal(t, sdk.NewDec(tc.amountOfMINERAL), tc.coins.AmountOf("mineral"))
		assert.Equal(t, sdk.NewDec(tc.amountOfTREE), tc.coins.AmountOf("tree"))
	}

	assert.Panics(t, func() { cases[0].coins.AmountOf("Invalid") })
}

func TestCoinsIsAnyGTE(t *testing.T) {
	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	assert.False(t, sdk.SysCoins{}.IsAnyGTE(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAnyGTE(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, two}, {syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, two}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom2, two}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom2, one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{{"xxx", one}, {"yyy", one}}.IsAnyGTE(sdk.SysCoins{{syscoinsTestDenom2, one}, {"ccc", one}, {"yyy", one}, {"zzz", one}}))
}

func TestCoinsIsAllGT(t *testing.T) {
	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	assert.False(t, sdk.SysCoins{}.IsAllGT(sdk.SysCoins{}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{}))
	assert.False(t, sdk.SysCoins{}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, two}, {syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, two}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom2, two}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom2, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{"xxx", one}, {"yyy", one}}.IsAllGT(sdk.SysCoins{{syscoinsTestDenom2, one}, {"ccc", one}, {"yyy", one}, {"zzz", one}}))
}

func TestCoinsIsAllGTE(t *testing.T) {
	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	assert.True(t, sdk.SysCoins{}.IsAllGTE(sdk.SysCoins{}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGTE(sdk.SysCoins{}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, two}, {syscoinsTestDenom2, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, two}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom2, two}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom2, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.True(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}))
	assert.False(t, sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom1, one}, {syscoinsTestDenom2, two}}))
	assert.False(t, sdk.SysCoins{{"xxx", one}, {"yyy", one}}.IsAllGTE(sdk.SysCoins{{syscoinsTestDenom2, one}, {"ccc", one}, {"yyy", one}, {"zzz", one}}))
}

func TestNewCoins(t *testing.T) {
	tenatom := sdk.NewInt64Coin("atom", 10)
	tenbtc := sdk.NewInt64Coin("btc", 10)
	zeroeth := sdk.NewInt64Coin("eth", 0)
	tests := []struct {
		name      string
		coins     sdk.SysCoins
		want      sdk.SysCoins
		wantPanic bool
	}{
		{"empty args", []sdk.SysCoin{}, sdk.SysCoins{}, false},
		{"one coin", []sdk.SysCoin{tenatom}, sdk.SysCoins{tenatom}, false},
		{"sort after create", []sdk.SysCoin{tenbtc, tenatom}, sdk.SysCoins{tenatom, tenbtc}, false},
		{"sort and remove zeroes", []sdk.SysCoin{zeroeth, tenbtc, tenatom}, sdk.SysCoins{tenatom, tenbtc}, false},
		{"panic on dups", []sdk.SysCoin{tenatom, tenatom}, sdk.SysCoins{}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() { sdk.NewCoins(tt.coins...) })
				return
			}
			got := sdk.NewCoins(tt.coins...)
			require.True(t, got.IsEqual(tt.want))
		})
	}
}

func TestCoinsIsAnyGT(t *testing.T) {
	twoAtom := sdk.NewInt64Coin("atom", 2)
	fiveAtom := sdk.NewInt64Coin("atom", 5)
	threeEth := sdk.NewInt64Coin("eth", 3)
	sixEth := sdk.NewInt64Coin("eth", 6)
	twoBtc := sdk.NewInt64Coin("btc", 2)

	require.False(t, sdk.SysCoins{}.IsAnyGT(sdk.SysCoins{}))

	require.False(t, sdk.SysCoins{fiveAtom}.IsAnyGT(sdk.SysCoins{}))
	require.False(t, sdk.SysCoins{}.IsAnyGT(sdk.SysCoins{fiveAtom}))
	require.True(t, sdk.SysCoins{fiveAtom}.IsAnyGT(sdk.SysCoins{twoAtom}))
	require.False(t, sdk.SysCoins{twoAtom}.IsAnyGT(sdk.SysCoins{fiveAtom}))

	require.True(t, sdk.SysCoins{twoAtom, sixEth}.IsAnyGT(sdk.SysCoins{twoBtc, fiveAtom, threeEth}))
	require.False(t, sdk.SysCoins{twoBtc, twoAtom, threeEth}.IsAnyGT(sdk.SysCoins{fiveAtom, sixEth}))
	require.False(t, sdk.SysCoins{twoAtom, sixEth}.IsAnyGT(sdk.SysCoins{twoBtc, fiveAtom}))
}

func TestMarshalJSONCoins(t *testing.T) {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)

	testCases := []struct {
		name      string
		input     sdk.SysCoins
		strOutput string
	}{
		{"nil coins", nil, `[]`},
		{"empty coins", sdk.SysCoins{}, `[]`},
		{"non-empty coins", sdk.NewCoins(sdk.NewInt64Coin("foo", 50)), `[{"denom":"foo","amount":"50.000000000000000000"}]`},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			bz, err := cdc.MarshalJSON(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.strOutput, string(bz))

			var newCoins sdk.SysCoins
			require.NoError(t, cdc.UnmarshalJSON(bz, &newCoins))

			if tc.input.Empty() {
				require.Nil(t, newCoins)
			} else {
				require.Equal(t, tc.input, newCoins)
			}
		})
	}
}

func TestCoinsIsValid(t *testing.T) {
	testCases := []struct {
		input    sdk.Coins
		expected bool
	}{
		{sdk.SysCoins{}, true},
		{sdk.SysCoins{sdk.DecCoin{Denom: syscoinsTestDenom1, Amount: sdk.NewDec(5)}}, true},
		{sdk.SysCoins{sdk.DecCoin{Denom: syscoinsTestDenom1, Amount: sdk.NewDec(5)}, sdk.DecCoin{Denom: syscoinsTestDenom2, Amount: sdk.NewDec(100000)}}, true},
		{sdk.SysCoins{sdk.DecCoin{Denom: syscoinsTestDenom1, Amount: sdk.NewDec(-5)}}, false},
		{sdk.SysCoins{sdk.DecCoin{Denom: "AAA", Amount: sdk.NewDec(5)}}, false},
		{sdk.SysCoins{sdk.DecCoin{Denom: syscoinsTestDenom1, Amount: sdk.NewDec(5)}, sdk.DecCoin{Denom: "B", Amount: sdk.NewDec(100000)}}, false},
		{sdk.SysCoins{sdk.DecCoin{Denom: syscoinsTestDenom1, Amount: sdk.NewDec(5)}, sdk.DecCoin{Denom: syscoinsTestDenom2, Amount: sdk.NewDec(-100000)}}, false},
		{sdk.SysCoins{sdk.DecCoin{Denom: syscoinsTestDenom1, Amount: sdk.NewDec(-5)}, sdk.DecCoin{Denom: syscoinsTestDenom2, Amount: sdk.NewDec(100000)}}, false},
		{sdk.SysCoins{sdk.DecCoin{Denom: "AAA", Amount: sdk.NewDec(5)}, sdk.DecCoin{Denom: syscoinsTestDenom2, Amount: sdk.NewDec(100000)}}, false},
	}

	for i, tc := range testCases {
		res := tc.input.IsValid()
		require.Equal(t, tc.expected, res, "unexpected result for test case #%d, input: %v", i, tc.input)
	}
}

func TestDecCoinsString(t *testing.T) {
	testCases := []struct {
		input    sdk.SysCoins
		expected string
	}{
		{sdk.SysCoins{}, ""},
		{
			sdk.SysCoins{
				sdk.NewDecCoinFromDec("atom", sdk.NewDecWithPrec(5040000000000000000, sdk.Precision)),
				sdk.NewDecCoinFromDec("stake", sdk.NewDecWithPrec(4000000000000000, sdk.Precision)),
			},
			"5.040000000000000000atom,0.004000000000000000stake",
		},
	}

	for i, tc := range testCases {
		out := tc.input.String()
		require.Equal(t, tc.expected, out, "unexpected result for test case #%d, input: %v", i, tc.input)
	}
}

func TestDecCoinsIntersect(t *testing.T) {
	testCases := []struct {
		input1         string
		input2         string
		expectedResult string
	}{
		{"", "", ""},
		{"1.0stake", "", ""},
		{"1.0stake", "1.0stake", "1.0stake"},
		{"", "1.0stake", ""},
		{"1.0stake", "", ""},
		{"2.0stake,1.0trope", "1.9stake", "1.9stake"},
		{"2.0stake,1.0trope", "2.1stake", "2.0stake"},
		{"2.0stake,1.0trope", "0.9trope", "0.9trope"},
		{"2.0stake,1.0trope", "1.9stake,0.9trope", "1.9stake,0.9trope"},
		{"2.0stake,1.0trope", "1.9stake,0.9trope,20.0other", "1.9stake,0.9trope"},
		{"2.0stake,1.0trope", "1.0other", ""},
	}

	for i, tc := range testCases {
		in1, err := sdk.ParseDecCoins(tc.input1)
		require.NoError(t, err, "unexpected parse error in %v", i)
		in2, err := sdk.ParseDecCoins(tc.input2)
		require.NoError(t, err, "unexpected parse error in %v", i)
		exr, err := sdk.ParseDecCoins(tc.expectedResult)
		require.NoError(t, err, "unexpected parse error in %v", i)

		require.True(t, in1.Intersect(in2).IsEqual(exr), "in1.cap(in2) != exr in %v", i)
	}
}

func TestDecCoinsTruncateDecimal(t *testing.T) {
	decCoinA := sdk.NewDecCoinFromDec("bar", sdk.MustNewDecFromStr("5.41"))
	decCoinB := sdk.NewDecCoinFromDec("foo", sdk.MustNewDecFromStr("6.00"))

	testCases := []struct {
		input          sdk.DecCoins
		truncatedCoins sdk.SysCoins
		changeCoins    sdk.DecCoins
	}{
		{sdk.DecCoins{}, sdk.SysCoins(nil), sdk.DecCoins(nil)},
		{
			sdk.DecCoins{decCoinA, decCoinB},
			sdk.SysCoins{sdk.NewInt64Coin(decCoinA.Denom, 5), sdk.NewInt64Coin(decCoinB.Denom, 6)},
			sdk.DecCoins{sdk.NewDecCoinFromDec(decCoinA.Denom, sdk.MustNewDecFromStr("0.41"))},
		},
		{
			sdk.DecCoins{decCoinB},
			sdk.SysCoins{sdk.NewInt64Coin(decCoinB.Denom, 6)},
			sdk.DecCoins(nil),
		},
	}

	for i, tc := range testCases {
		truncatedCoins, changeCoins := tc.input.TruncateDecimal()
		require.Equal(
			t, tc.truncatedCoins, truncatedCoins,
			"unexpected truncated coins; tc #%d, input: %s", i, tc.input,
		)
		require.Equal(
			t, tc.changeCoins, changeCoins,
			"unexpected change coins; tc #%d, input: %s", i, tc.input,
		)
	}
}

func TestDecCoinsQuoDecTruncate(t *testing.T) {
	x := sdk.MustNewDecFromStr("1.00")
	y := sdk.MustNewDecFromStr("10000000000000000000.00")

	testCases := []struct {
		coins  sdk.DecCoins
		input  sdk.Dec
		result sdk.DecCoins
		panics bool
	}{
		{sdk.DecCoins{}, sdk.ZeroDec(), sdk.DecCoins(nil), true},
		{sdk.DecCoins{sdk.NewDecCoinFromDec("foo", x)}, y, sdk.DecCoins(nil), false},
		{sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)}, sdk.NewDec(2), sdk.DecCoins{sdk.NewDecCoinFromDec("foo", sdk.MustNewDecFromStr("2.5"))}, false},
	}

	for i, tc := range testCases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.coins.QuoDecTruncate(tc.input) })
		} else {
			res := tc.coins.QuoDecTruncate(tc.input)
			require.Equal(t, tc.result, res, "unexpected result; tc #%d, coins: %s, input: %s", i, tc.coins, tc.input)
		}
	}
}