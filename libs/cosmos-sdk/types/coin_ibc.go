package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"regexp"
	"strings"
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
