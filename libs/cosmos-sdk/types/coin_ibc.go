package types

type CoinAdapters []CoinAdapter

// NewCoin returns a new coin with a denomination and amount. It will panic if
// the amount is negative or if the denomination is invalid.
func NewCoinAdapter(denom string, amount Int) CoinAdapter {
	coin := CoinAdapter{
		Denom:  denom,
		Amount: amount,
	}

	if err := coin.Validate(); err != nil {
		panic(err)
	}

	return coin
}

func (cas CoinAdapters) IsAnyNil() bool {
	for _, coin := range cas {
		if coin.Amount.IsNil() {
			return true
		}
	}

	return false
}

func (cas CoinAdapters) IsAnyNegative() bool {
	for _, coin := range cas {
		if coin.Amount.IsNegative() {
			return true
		}
	}

	return false
}

// ParseCoinsNormalized will parse out a list of coins separated by commas, and normalize them by converting to smallest
// unit. If the parsing is successuful, the provided coins will be sanitized by removing zero coins and sorting the coin
// set. Lastly a validation of the coin set is executed. If the check passes, ParseCoinsNormalized will return the
// sanitized coins.
// Otherwise it will return an error.
// If an empty string is provided to ParseCoinsNormalized, it returns nil Coins.
// ParseCoinsNormalized supports decimal coins as inputs, and truncate them to int after converted to smallest unit.
// Expected format: "{amount0}{denomination},...,{amountN}{denominationN}"
func ParseCoinsNormalized(coinStr string) (Coins, error) {
	coins, err := ParseDecCoins(coinStr)
	if err != nil {
		return Coins{}, err
	}
	return NormalizeCoins(coins), nil
}

// ParseCoinNormalized parses and normalize a cli input for one coin type, returning errors if invalid or on an empty string
// as well.
// Expected format: "{amount}{denomination}"
func ParseCoinNormalized(coinStr string) (coin Coin, err error) {
	decCoin, err := ParseDecCoin(coinStr)
	if err != nil {
		return Coin{}, err
	}

	coin, _ = NormalizeDecCoin(decCoin).TruncateDecimal()
	return coin, nil
}
