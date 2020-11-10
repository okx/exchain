package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// swap business type
	SwapBusinessTypeCreate = "create"
	SwapBusinessTypeAdd    = "add"
	SwapBusinessTypeSwap   = "swap"
)

// nolint
type QuerySwapTokensParams struct {
	BusinessType  string `json:"business_type"`
	Address       string `json:"address"`
	BaseTokenName string `json:"base_token_name"`
}

// NewQuerySwapTokensParams creates a new instance of QueryDexFeesParams
func NewQuerySwapTokensParams(businessType string, address string, baseTokenName string) QuerySwapTokensParams {
	return QuerySwapTokensParams{
		BusinessType:  businessType,
		Address:       address,
		BaseTokenName: baseTokenName,
	}
}

// nolint
type QuerySwapBuyInfoParams struct {
	SellTokenAmount string `json:"sell_token_amount"`
	BuyToken        string `json:"buy_token"`
}

// NewQuerySwapBuyInfoParams creates a new instance of QuerySwapBuyInfoParams
func NewQuerySwapBuyInfoParams(sellTokenAmount string, buyToken string) QuerySwapBuyInfoParams {
	return QuerySwapBuyInfoParams{
		SellTokenAmount: sellTokenAmount,
		BuyToken:        buyToken,
	}
}

// nolint
type QuerySwapLiquidityInfoParams struct {
	Address       string `json:"address"`
	TokenPairName string `json:"token_pair_name"`
}

// NewQuerySwapBuyInfoParams creates a new instance of QuerySwapLiquidityInfoParams
func NewQuerySwapLiquidityInfoParams(address string, tokenPairName string) QuerySwapLiquidityInfoParams {
	return QuerySwapLiquidityInfoParams{
		Address:       address,
		TokenPairName: tokenPairName,
	}
}

// nolint
type QuerySwapAddInfoParams struct {
	QuoteTokenAmount string `json:"quote_token_amount"`
	BaseToken        string `json:"base_token"`
}

// NewQuerySwapAddInfoParams creates a new instance of QuerySwapAddInfoParams
func NewQuerySwapAddInfoParams(quoteTokenAmount string, baseToken string) QuerySwapAddInfoParams {
	return QuerySwapAddInfoParams{
		QuoteTokenAmount: quoteTokenAmount,
		BaseToken:        baseToken,
	}
}

type SwapToken struct {
	Symbol    string  `json:"symbol"`
	Available sdk.Dec `json:"available"`
}

func NewSwapToken(symbol string, available sdk.Dec) SwapToken {
	return SwapToken{
		Symbol:    symbol,
		Available: available,
	}
}

type SwapTokens []SwapToken

type SwapTokensResponse struct {
	NativeToken string     `json:"native_token"`
	Tokens      SwapTokens `json:"tokens"`
}

func (swapTokens SwapTokens) Len() int { return len(swapTokens) }

func (swapTokens SwapTokens) Less(i, j int) bool {
	return swapTokens[i].Available.GT(swapTokens[j].Available)
}

func (swapTokens SwapTokens) Swap(i, j int) {
	swapTokens[i], swapTokens[j] = swapTokens[j], swapTokens[i]
}

type SwapBuyInfo struct {
	BuyAmount   sdk.Dec `json:"buy_amount"`
	Price       sdk.Dec `json:"price"`
	PriceImpact sdk.Dec `json:"price_impact"`
	Fee         string  `json:"fee"`
	Route       string  `json:"route"`
}

type SwapLiquidityInfo struct {
	BasePooledCoin  sdk.SysCoin `json:"base_pooled_coin"`
	QuotePooledCoin sdk.SysCoin `json:"quote_pooled_coin"`
	PoolTokenCoin   sdk.SysCoin `json:"pool_token_coin"`
	PoolTokenRatio  sdk.Dec     `json:"pool_token_ratio"`
}

type SwapAddInfo struct {
	BaseTokenAmount sdk.Dec `json:"base_token_amount"`
	PoolShare       sdk.Dec `json:"pool_share"`
}

type QueryBuyAmountParams struct {
	SoldToken  sdk.SysCoin
	TokenToBuy string
}
