package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var DefaultTokenPairDeposit = sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(0))

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
	Deposits         sdk.DecCoin    `json:"deposits"`
	BlockHeight      int64          `json:"block_height"`
}

func (tp *TokenPair) Name() string {
	return fmt.Sprintf("%s_%s", tp.BaseAssetSymbol, tp.QuoteAssetSymbol)
}

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

type TokenPairs []*TokenPair

func (tp TokenPairs) Len() int { return len(tp) }

func (tp TokenPairs) Less(i, j int) bool { return tp[i].IsGT(tp[j]) }

func (tp TokenPairs) Swap(i, j int) { tp[i], tp[j] = tp[j], tp[i] }
