package types

import "github.com/okex/exchain/libs/cosmos-sdk/types"

func CoinAdapterToCoin(adapter types.CoinAdapter) types.Coin {
	return types.Coin{
		Denom:  adapter.Denom,
		Amount: adapter.Amount.ToDec(),
	}
}

func CoinToCoinAdapter(adapter types.Coin) types.CoinAdapter {
	return types.CoinAdapter{
		Denom:  adapter.Denom,
		Amount: types.NewIntFromBigInt(adapter.Amount.BigInt()),
	}
}

func CoinAdaptersToCoins(adapters types.CoinAdapters) types.Coins {
	var coins types.Coins = make([]types.Coin, 0)
	for i, _ := range adapters {
		coins = append(coins, CoinAdapterToCoin(adapters[i]))
	}
	return coins
}

func CoinsToCoinAdapters(coins types.Coins) types.CoinAdapters {
	var adapters types.CoinAdapters
	for i, _ := range coins {
		adapters = append(adapters, CoinToCoinAdapter(coins[i]))
	}
	return adapters
}
