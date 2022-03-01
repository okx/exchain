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