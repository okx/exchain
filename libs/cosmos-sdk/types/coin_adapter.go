package types

func CoinAdapterToCoin(adapter CoinAdapter) Coin {
	return Coin{
		Denom:  adapter.Denom,
		Amount: NewDecFromBigIntWithPrec(adapter.Amount.BigInt(), Precision),
	}
}

func CoinToCoinAdapter(adapter Coin) CoinAdapter {
	return CoinAdapter{
		Denom:  adapter.Denom,
		Amount: NewIntFromBigInt(adapter.Amount.BigInt()),
	}
}

func CoinAdaptersToCoins(adapters CoinAdapters) Coins {
	var coins Coins = make([]Coin, 0, len(adapters))
	for i, _ := range adapters {
		coins = append(coins, CoinAdapterToCoin(adapters[i]))
	}
	return coins
}

func CoinsToCoinAdapters(coins Coins) CoinAdapters {
	//Note:
	// `var adapters CoinAdapters = make([]CoinAdapter, 0)`
	// The code above if invalid.
	// []CoinAdapter{} and nil are different in json format which can make different signBytes.
	var adapters CoinAdapters
	for i, _ := range coins {
		adapters = append(adapters, CoinToCoinAdapter(coins[i]))
	}
	return adapters
}

func removeZeroCoinAdapters(coins CoinAdapters) CoinAdapters {
	for i := 0; i < len(coins); i++ {
		if coins[i].IsZero() {
			break
		} else if i == len(coins)-1 {
			return coins
		}
	}
	var result []CoinAdapter
	if len(coins) > 0 {
		result = make([]CoinAdapter, 0, len(coins)-1)
	}

	for _, coin := range coins {
		if !coin.IsZero() {
			result = append(result, coin)
		}
	}
	return result
}

func (coins CoinAdapters) IsZero() bool {
	for _, coin := range coins {
		if !coin.IsZero() {
			return false
		}
	}
	return true
}
