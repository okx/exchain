package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMaxCanBorrow(t *testing.T) {
	denom := sdk.DefaultBondDenom
	maxLeverage := sdk.NewDec(5)
	account := Account{
		Available: sdk.NewDecCoinsFromDec(denom, sdk.NewDec(10)),
		Borrowed:  sdk.NewDecCoinsFromDec(denom, sdk.NewDec(6)),
		Interest:  sdk.NewDecCoinsFromDec(denom, sdk.NewDec(1)),
	}
	maxCanBorrow := account.MaxCanBorrow(denom, maxLeverage)
	expected := sdk.NewDecCoinFromDec(denom, sdk.NewDec(6))
	require.Equal(t, expected, maxCanBorrow)

	availableAmount, err := sdk.NewDecFromStr("10.20")
	require.Nil(t, err)
	account.Available = sdk.NewDecCoinsFromDec(denom, availableAmount)
	expectedAmount, err := sdk.NewDecFromStr("6.8")
	require.Nil(t, err)
	expected = sdk.NewDecCoinFromDec(denom, expectedAmount)
	maxCanBorrow = account.MaxCanBorrow(denom, maxLeverage)
	require.Equal(t, expected, maxCanBorrow)

}

func TestMarginRatio(t *testing.T) {
	baseSmbol := "xxb"
	quoteSmbol := sdk.DefaultBondDenom

	latestPrice := sdk.NewDec(5000)
	baseBorrow := sdk.NewDecCoinFromDec(baseSmbol, sdk.NewDec(4))
	baseInterestAmount, err := sdk.NewDecFromStr("0.5")
	require.Nil(t, err)
	baseInterest := sdk.NewDecCoinFromDec(baseSmbol, baseInterestAmount)

	quoteTotal := sdk.NewDecCoinFromDec(quoteSmbol, sdk.NewDec(40000))
	expected, err := sdk.NewDecFromStr("0.875")
	require.Nil(t, err)

	account := Account{
		Available: sdk.DecCoins{quoteTotal},
		Borrowed:  sdk.DecCoins{baseBorrow},
		Interest:  sdk.DecCoins{baseInterest},
	}

	marginRatio := account.MarginRatio(baseSmbol, quoteSmbol, latestPrice)
	require.Equal(t, expected, marginRatio)
}
