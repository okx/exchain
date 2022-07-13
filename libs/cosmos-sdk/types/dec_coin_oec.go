package types

// SafeSub performs the same arithmetic as Sub but returns a boolean if any
// negative coin amount was returned.
func (coins DecCoins) fastSafeSubOneCoin(coinsB DecCoins) (DecCoins, bool) {

	diff := coins.safeAdd(coinsB.negative())
	return diff, diff.IsAnyNegative()
}
