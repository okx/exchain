package types

import (
	"fmt"
	"math/big"
	"strings"
)

// String provides a human-readable representation of a coin
func (coin CoinAdapter) String() string {
	return fmt.Sprintf("%v%v", coin.Amount, coin.Denom)
}

// Validate returns an error if the Coin has a negative amount or if
// the denom is invalid.
func (coin CoinAdapter) Validate() error {
	if err := ValidateDenom(coin.Denom); err != nil {
		return err
	}

	if coin.Amount.IsNegative() {
		return fmt.Errorf("negative coin amount: %v", coin.Amount)
	}

	return nil
}

// IsValid returns true if the Coin has a non-negative amount and the denom is valid.
func (coin CoinAdapter) IsValid() bool {
	return coin.Validate() == nil
}

// IsZero returns if this represents no money
func (coin CoinAdapter) IsZero() bool {
	return coin.Amount.IsZero()
}

// IsGTE returns true if they are the same type and the receiver is
// an equal or greater value
func (coin CoinAdapter) IsGTE(other CoinAdapter) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return !coin.Amount.LT(other.Amount)
}

// IsLT returns true if they are the same type and the receiver is
// a smaller value
func (coin CoinAdapter) IsLT(other CoinAdapter) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return coin.Amount.LT(other.Amount)
}

// IsEqual returns true if the two sets of Coins have the same value
func (coin CoinAdapter) IsEqual(other CoinAdapter) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return coin.Amount.Equal(other.Amount)
}

// Add adds amounts of two coins with same denom. If the coins differ in denom then
// it panics.
func (coin CoinAdapter) Add(coinB CoinAdapter) CoinAdapter {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, coinB.Denom))
	}

	return CoinAdapter{coin.Denom, coin.Amount.Add(coinB.Amount)}
}

// Sub subtracts amounts of two coins with same denom. If the coins differ in denom
// then it panics.
func (coin CoinAdapter) Sub(coinB CoinAdapter) CoinAdapter {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, coinB.Denom))
	}

	res := CoinAdapter{coin.Denom, coin.Amount.Sub(coinB.Amount)}
	if res.IsNegative() {
		panic("negative coin amount")
	}

	return res
}

// IsPositive returns true if coin amount is positive.
//
// TODO: Remove once unsigned integers are used.
func (coin CoinAdapter) IsPositive() bool {
	return coin.Amount.Sign() == 1
}

// IsNegative returns true if the coin amount is negative and false otherwise.
//
// TODO: Remove once unsigned integers are used.
func (coin CoinAdapter) IsNegative() bool {
	return coin.Amount.Sign() == -1
}

// Unmarshal implements the gogo proto custom type interface.
func (i *Int) Unmarshal(data []byte) error {
	if len(data) == 0 {
		i = nil
		return nil
	}

	if i.i == nil {
		i.i = new(big.Int)
	}

	if err := i.i.UnmarshalText(data); err != nil {
		return err
	}

	if i.i.BitLen() > maxBitLen {
		return fmt.Errorf("integer out of range; got: %d, max: %d", i.i.BitLen(), maxBitLen)
	}

	return nil
}

// Size implements the gogo proto custom type interface.
func (i *Int) Size() int {
	bz, _ := i.Marshal()
	return len(bz)
}

// Marshal implements the gogo proto custom type interface.
func (i Int) Marshal() ([]byte, error) {
	if i.i == nil {
		i.i = new(big.Int)
	}
	return i.i.MarshalText()
}

// MarshalTo implements the gogo proto custom type interface.
func (i *Int) MarshalTo(data []byte) (n int, err error) {
	if i.i == nil {
		i.i = new(big.Int)
	}
	if len(i.i.Bytes()) == 0 {
		copy(data, []byte{0x30})
		return 1, nil
	}

	bz, err := i.Marshal()
	if err != nil {
		return 0, err
	}

	copy(data, bz)
	return len(bz), nil
}

func NewIBCCoin(denom string, amount interface{}) DecCoin {
	switch amount := amount.(type) {
	case Int:
		return NewIBCDecCoin(denom, amount)
	case Dec:
		return NewDecCoinFromDec(denom, amount)
	default:
		panic("Invalid amount")
	}
}

func NewIBCDecCoin(denom string, amount Int) DecCoin {
	if err := validateIBCCotin(denom, amount); err != nil {
		panic(err)
	}

	return DecCoin{
		Denom:  denom,
		Amount: amount.ToDec(),
	}
}

func validateIBCCotin(denom string, amount Int) error {
	if !validIBCCoinDenom(denom) {
		return fmt.Errorf("invalid denom: %s", denom)
	}

	if amount.IsNegative() {
		return fmt.Errorf("negative coin amount: %v", amount)
	}

	return nil
}

func validIBCCoinDenom(denom string) bool {
	return ibcReDnm.MatchString(denom)
}

// NewCoins constructs a new coin set.
func NewIBCCoins(coins ...Coin) Coins {
	// remove zeroes
	newCoins := removeZeroCoins(Coins(coins))
	if len(newCoins) == 0 {
		return Coins{}
	}

	newCoins.Sort()

	// detect duplicate Denoms
	if dupIndex := findDup(newCoins); dupIndex != -1 {
		panic(fmt.Errorf("find duplicate denom: %s", newCoins[dupIndex]))
	}
	newCoins.IsValid()
	if !ValidCoins(newCoins) {
		panic(fmt.Errorf("invalid coin set: %s", newCoins))
	}

	return newCoins
}
func ValidCoins(coins Coins) bool {
	switch len(coins) {
	case 0:
		return true

	case 1:
		if !validIBCCoinDenom(coins[0].Denom) {
			return false
		}
		return coins[0].IsPositive()

	default:
		// check single coin case
		if !ValidCoins(DecCoins{coins[0]}) {
			return false
		}

		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if strings.ToLower(coin.Denom) != coin.Denom {
				return false
			}
			if coin.Denom <= lowDenom {
				return false
			}
			if !coin.IsPositive() {
				return false
			}

			// we compare each coin against the last denom
			lowDenom = coin.Denom
		}

		return true
	}
}
