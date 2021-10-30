package types

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

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

type SwapBuyInfo struct {
	BuyAmount   sdk.Dec `json:"buy_amount"`
	Price       sdk.Dec `json:"price"`
	PriceImpact sdk.Dec `json:"price_impact"`
	Fee         string  `json:"fee"`
	Route       string  `json:"route"`
}

type SwapAddInfo struct {
	BaseTokenAmount sdk.Dec `json:"base_token_amount"`
	PoolShare       sdk.Dec `json:"pool_share"`
	Liquidity       sdk.Dec `json:"liquidity"`
}

type QueryBuyAmountParams struct {
	SoldToken  sdk.SysCoin
	TokenToBuy string
}
