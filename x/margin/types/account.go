package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Account struct {
	Product   string       `json:"product"`
	Available sdk.DecCoins `json:"available"`
	Locked    sdk.DecCoins `json:"locked"`
	Borrowed  sdk.DecCoins `json:"borrowed"`
	Interest  sdk.DecCoins `json:"interest"`
}

func (account Account) MaxCanBorrow(denom string, maxLeverage sdk.Dec) sdk.DecCoin {
	available := account.Available.AmountOf(denom)
	borrowed := account.Borrowed.AmountOf(denom)
	interest := account.Interest.AmountOf(denom)

	canBorrow := available.Sub(borrowed).Sub(interest).Mul(maxLeverage.Sub(sdk.NewDec(1))).Sub(borrowed)
	return sdk.NewDecCoinFromDec(denom, canBorrow)
}

func (account Account) MarginRatio(baseSymbol string, quoteSymbol string, latestPrice sdk.Dec) sdk.Dec {
	quoteTotal := account.Available.AmountOf(quoteSymbol).Add(account.Locked.AmountOf(quoteSymbol))
	quoteBorrow := account.Borrowed.AmountOf(quoteSymbol)
	quoteInterest := account.Interest.AmountOf(quoteSymbol)
	baseTotal := account.Available.AmountOf(baseSymbol).Add(account.Locked.AmountOf(baseSymbol))
	baseBorrow := account.Borrowed.AmountOf(baseSymbol)
	baseInterest := account.Interest.AmountOf(baseSymbol)

	quoteRemain := quoteTotal.Sub(quoteBorrow).Sub(quoteInterest)
	baseRemain := baseTotal.Sub(baseBorrow).Sub(baseInterest)

	numerator := quoteRemain.Quo(latestPrice).Add(baseRemain)
	denominator := quoteBorrow.Quo(latestPrice).Add(baseBorrow)
	return numerator.Quo(denominator)
}
