package types

import (
	"fmt"
	"strings"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

// DefaultTokenPairDeposit defines default deposit of token pair
var DefaultTokenPairDeposit = sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(0))

// TokenPair represents token pair object
type TokenPair struct {
	BaseAssetSymbol  string         `json:"base_asset_symbol"`
	QuoteAssetSymbol string         `json:"quote_asset_symbol"`
	InitPrice        sdk.Dec        `json:"price"`
	MaxPriceDigit    int64          `json:"max_price_digit"`
	MaxQuantityDigit int64          `json:"max_size_digit"`
	MinQuantity      sdk.Dec        `json:"min_trade_size"`
	ID               uint64         `json:"token_pair_id"`
	Delisting        bool           `json:"delisting"`
	Owner            sdk.AccAddress `json:"owner"`
	Deposits         sdk.SysCoin    `json:"deposits"`
	BlockHeight      int64          `json:"block_height"`
}

// Name returns name of token pair
func (tp *TokenPair) Name() string {
	return fmt.Sprintf("%s_%s", tp.BaseAssetSymbol, tp.QuoteAssetSymbol)
}

// IsGT returns true if the token pair is greater than the other one
// 1. compare deposits
// 2. compare block height
// 3. compare name
func (tp *TokenPair) IsGT(other *TokenPair) bool {
	if tp.Deposits.IsLT(other.Deposits) {
		return false
	}

	if !tp.Deposits.IsEqual(other.Deposits) {
		return true
	}

	if tp.BlockHeight < other.BlockHeight {
		return true
	}

	if tp.BlockHeight > other.BlockHeight {
		return false
	}

	return strings.Compare(tp.BaseAssetSymbol, other.BaseAssetSymbol) < 0 || strings.Compare(tp.QuoteAssetSymbol, other.QuoteAssetSymbol) < 0
}

// TokenPairs represents token pair slice, support sorting
type TokenPairs []*TokenPair

// Len Implements Sort
func (tp TokenPairs) Len() int { return len(tp) }

// Less Implements Sort
func (tp TokenPairs) Less(i, j int) bool { return tp[i].IsGT(tp[j]) }

// Swap Implements Sort
func (tp TokenPairs) Swap(i, j int) { tp[i], tp[j] = tp[j], tp[i] }
