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
	var adapters CoinAdapters
	for i, _ := range coins {
		adapters = append(adapters, CoinToCoinAdapter(coins[i]))
	}
	return adapters
}
