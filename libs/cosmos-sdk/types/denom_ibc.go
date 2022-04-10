package types

import (
	"fmt"
)

// baseDenom is the denom of smallest unit registered
var baseDenom string

// GetBaseDenom returns the denom of smallest unit registered
func GetBaseDenom() (string, error) {
	if baseDenom == "" {
		return "", fmt.Errorf("no denom is registered")
	}
	return baseDenom, nil
}

// ConvertDecCoin attempts to convert a decimal coin to a given denomination. If the given
// denomination is invalid or if neither denomination is registered, an error
// is returned.
func ConvertDecCoin(coin DecCoin, denom string) (DecCoin, error) {
	if err := ValidateDenom(denom); err != nil {
		return DecCoin{}, err
	}

	srcUnit, ok := GetDenomUnit(coin.Denom)
	if !ok {
		return DecCoin{}, fmt.Errorf("source denom not registered: %s", coin.Denom)
	}

	dstUnit, ok := GetDenomUnit(denom)
	if !ok {
		return DecCoin{}, fmt.Errorf("destination denom not registered: %s", denom)
	}

	if srcUnit.Equal(dstUnit) {
		return NewDecCoinFromDec(denom, coin.Amount), nil
	}

	return NewDecCoinFromDec(denom, coin.Amount.Mul(srcUnit).Quo(dstUnit)), nil
}

// NormalizeCoin try to convert a coin to the smallest unit registered,
// returns original one if failed.
func NormalizeCoin(coin Coin) Coin {
	base, err := GetBaseDenom()
	if err != nil {
		return coin
	}
	newCoin, err := ConvertCoin(coin, base)
	if err != nil {
		return coin
	}
	return newCoin
}

// NormalizeDecCoin try to convert a decimal coin to the smallest unit registered,
// returns original one if failed.
func NormalizeDecCoin(coin DecCoin) DecCoin {
	base, err := GetBaseDenom()
	if err != nil {
		return coin
	}
	newCoin, err := ConvertDecCoin(coin, base)
	if err != nil {
		return coin
	}
	return newCoin
}

// NormalizeCoins normalize and truncate a list of decimal coins
func NormalizeCoins(coins []DecCoin) Coins {
	if coins == nil {
		return nil
	}
	result := make([]Coin, 0, len(coins))

	for _, coin := range coins {
		newCoin, _ := NormalizeDecCoin(coin).TruncateDecimal()
		result = append(result, newCoin)
	}

	return result
}
