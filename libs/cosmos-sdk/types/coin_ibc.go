package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var (
	// Denominations can be 3 ~ 128 characters long and support letters, followed by either
	// a letter, a number or a separator ('/').
	ibcReDnmString = `[a-zA-Z][a-zA-Z0-9/-]{2,127}`
	ibcReDecAmt    = `[[:digit:]]+(?:\.[[:digit:]]+)?|\.[[:digit:]]+`
	ibcReSpc       = `[[:space:]]*`
	ibcReDnm       *regexp.Regexp
	ibcReDecCoin   *regexp.Regexp
)
var ibcCoinDenomRegex = DefaultCoinDenomRegex

func init() {
	SetIBCCoinDenomRegex(DefaultIBCCoinDenomRegex)
}

func IBCParseDecCoin(coinStr string) (coin DecCoin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := ibcReDecCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return DecCoin{}, fmt.Errorf("invalid decimal coin expression: %s", coinStr)
	}

	amountStr, denomStr := matches[1], matches[2]

	amount, err := NewDecFromStr(amountStr)
	if err != nil {
		return DecCoin{}, errors.Wrap(err, fmt.Sprintf("failed to parse decimal coin amount: %s", amountStr))
	}

	if err := ValidateDenom(denomStr); err != nil {
		return DecCoin{}, fmt.Errorf("invalid denom cannot contain upper case characters or spaces: %s", err)
	}

	return NewDecCoinFromDec(denomStr, amount), nil
}

// DefaultCoinDenomRegex returns the default regex string
func DefaultIBCCoinDenomRegex() string {
	return ibcReDnmString
}

func SetIBCCoinDenomRegex(reFn func() string) {
	ibcCoinDenomRegex = reFn

	ibcReDnm = regexp.MustCompile(fmt.Sprintf(`^%s$`, ibcCoinDenomRegex()))
	ibcReDecCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, ibcReDecAmt, ibcReSpc, ibcCoinDenomRegex()))
}

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
func (cs CoinAdapter) ToCoin() Coin {
	if cs.Denom == DefaultIbcWei {
		cs.Denom = DefaultBondDenom
	}
	return CoinAdapterToCoin(cs)
}
func (cs CoinAdapters) ToCoins() Coins {
	ret := make([]Coin, 0)
	for _, v := range cs {
		ret = append(ret, v.ToCoin())
	}
	return ret
}

func (cas CoinAdapters) IsAnyNegative() bool {
	for _, coin := range cas {
		if coin.Amount.IsNegative() {
			return true
		}
	}

	return false
}

func (cas CoinAdapters) IsAnyNil() bool {
	for _, coin := range cas {
		if coin.Amount.IsNil() {
			return true
		}
	}

	return false
}

func (coins CoinAdapters) Add(coinsB ...CoinAdapter) CoinAdapters {
	return coins.safeAdd(coinsB)
}
func (coins CoinAdapters) safeAdd(coinsB CoinAdapters) CoinAdapters {
	// probably the best way will be to make Coins and interface and hide the structure
	// definition (type alias)
	if !coins.isSorted() {
		panic("Coins (self) must be sorted")
	}
	if !coinsB.isSorted() {
		panic("Wrong argument: coins must be sorted")
	}

	sum := ([]CoinAdapter)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)

	for {
		if indexA == lenA {
			if indexB == lenB {
				// return nil coins if both sets are empty
				return sum
			}

			// return set B (excluding zero coins) if set A is empty
			return append(sum, removeZeroCoinAdapters(coinsB[indexB:])...)
		} else if indexB == lenB {
			// return set A (excluding zero coins) if set B is empty
			return append(sum, removeZeroCoinAdapters(coins[indexA:])...)
		}

		coinA, coinB := coins[indexA], coinsB[indexB]

		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1: // coin A denom < coin B denom
			if !coinA.IsZero() {
				sum = append(sum, coinA)
			}

			indexA++

		case 0: // coin A denom == coin B denom
			res := coinA.Add(coinB)
			if !res.IsZero() {
				sum = append(sum, res)
			}

			indexA++
			indexB++

		case 1: // coin A denom > coin B denom
			if !coinB.IsZero() {
				sum = append(sum, coinB)
			}

			indexB++
		}
	}
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

// IsValid calls Validate and returns true when the Coins are sorted, have positive amount, with a
// valid and unique denomination (i.e no duplicates).
func (coins CoinAdapters) IsValid() bool {
	return coins.Validate() == nil
}

func (coins CoinAdapters) Validate() error {
	switch len(coins) {
	case 0:
		return nil

	case 1:
		if err := ValidateDenom(coins[0].Denom); err != nil {
			return err
		}
		if coins[0].Amount.IsNil() {
			return fmt.Errorf("coin %s amount is nil", coins[0])
		}
		if !coins[0].IsPositive() {
			return fmt.Errorf("coin %s amount is not positive", coins[0])
		}
		return nil

	default:
		// check single coin case
		if err := (CoinAdapters{coins[0]}).Validate(); err != nil {
			return err
		}

		lowDenom := coins[0].Denom
		seenDenoms := make(map[string]bool)
		seenDenoms[lowDenom] = true

		for _, coin := range coins[1:] {
			if seenDenoms[coin.Denom] {
				return fmt.Errorf("duplicate denomination %s", coin.Denom)
			}
			if err := ValidateDenom(coin.Denom); err != nil {
				return err
			}
			if coin.Denom <= lowDenom {
				return fmt.Errorf("denomination %s is not sorted", coin.Denom)
			}
			if !coin.IsPositive() {
				return fmt.Errorf("coin %s amount is not positive", coin.Denom)
			}

			// we compare each coin against the last denom
			lowDenom = coin.Denom
			seenDenoms[coin.Denom] = true
		}

		return nil
	}
}

func (coins CoinAdapters) isSorted() bool {
	for i := 1; i < len(coins); i++ {
		if coins[i-1].Denom > coins[i].Denom {
			return false
		}
	}
	return true
}

func (coins CoinAdapters) String() string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin.String())
	}
	return out[:len(out)-1]
}

// IsAllPositive returns true if there is at least one coin and all currencies
// have a positive value.
func (coins CoinAdapters) IsAllPositive() bool {
	if len(coins) == 0 {
		return false
	}

	for _, coin := range coins {
		if !coin.IsPositive() {
			return false
		}
	}

	return true
}

type coinAdaptersJSON CoinAdapters

// MarshalJSON implements a custom JSON marshaller for the Coins type to allow
// nil Coins to be encoded as an empty array.
func (coins CoinAdapters) MarshalJSON() ([]byte, error) {
	if coins == nil {
		return json.Marshal(coinAdaptersJSON(CoinAdapters{}))
	}

	return json.Marshal(coinAdaptersJSON(coins))
}

func (coins CoinAdapters) Copy() CoinAdapters {
	copyCoins := make(CoinAdapters, len(coins))

	for i, coin := range coins {
		copyCoins[i] = coin
	}

	return copyCoins
}

func ConvWei2TOkt(adapters CoinAdapters) (CoinAdapters, error) {
	copyAdapters := adapters.Copy()
	for index, _ := range copyAdapters {
		if copyAdapters[index].Denom == DefaultIbcWei {
			copyAdapters[index].Denom = DefaultBondDenom
		} else if strings.ToLower(copyAdapters[index].Denom) == DefaultBondDenom {
			return nil, errors.Wrap(errors.ErrInvalidCoins, "not support okt denom")
		}
	}
	return copyAdapters, nil
}

// AmountOf returns the amount of a denom from coins
func (coins CoinAdapters) AmountOf(denom string) Int {
	mustValidateDenom(denom)
	return coins.AmountOfNoDenomValidation(denom)
}

// AmountOfNoDenomValidation returns the amount of a denom from coins
// without validating the denomination.
func (coins CoinAdapters) AmountOfNoDenomValidation(denom string) Int {
	switch len(coins) {
	case 0:
		return ZeroInt()

	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return coin.Amount
		}
		return ZeroInt()

	default:
		// Binary search the amount of coins remaining
		midIdx := len(coins) / 2 // 2:1, 3:1, 4:2
		coin := coins[midIdx]
		switch {
		case denom < coin.Denom:
			return coins[:midIdx].AmountOfNoDenomValidation(denom)
		case denom == coin.Denom:
			return coin.Amount
		default:
			return coins[midIdx+1:].AmountOfNoDenomValidation(denom)
		}
	}
}

func (coins CoinAdapters) Sub(coinsB CoinAdapters) CoinAdapters {
	diff, hasNeg := coins.SafeSub(coinsB)
	if hasNeg {
		panic("negative coin amount")
	}

	return diff
}

func (coins CoinAdapters) SafeSub(coinsB CoinAdapters) (CoinAdapters, bool) {
	diff := coins.safeAdd(coinsB.negative())
	return diff, diff.IsAnyNegative()
}
func (coins CoinAdapters) negative() CoinAdapters {
	res := make([]CoinAdapter, 0, len(coins))

	for _, coin := range coins {
		res = append(res, CoinAdapter{
			Denom:  coin.Denom,
			Amount: coin.Amount.Neg(),
		})
	}

	return res
}

// AddAmount adds an amount to the Coin.
func (coin CoinAdapter) AddAmount(amount Int) CoinAdapter {
	return CoinAdapter{coin.Denom, coin.Amount.Add(amount)}
}
