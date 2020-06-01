package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Account struct {
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
