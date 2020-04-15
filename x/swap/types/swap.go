package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"strings"
)

type SwapTokenPair struct {
	Quote       string  `json:"value"`        //swap token
	QuoteAmount sdk.Dec `json:"quote_amount"` //token pool amount
	BaseAmount  sdk.Dec `json:"base_amount"`  //native token pool amount
	PoolToken   string  `json:"pool_token"`   //pool token address
}

func NewSwapTokenPair(quote string, quoteAmount sdk.Dec, baseAmount sdk.Dec, poolToken string) *SwapTokenPair {
	swapTokenPair := &SwapTokenPair{
		Quote:       quote,
		QuoteAmount: quoteAmount,
		BaseAmount:  baseAmount,
		PoolToken:   poolToken,
	}
	return swapTokenPair
}

// implement fmt.Stringer
func (s SwapTokenPair) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Quote: %s
QuoteAmount: %s
Base: %s
BaseAmount: %s
PoolToken: %s`, s.Quote, s.QuoteAmount, common.NativeToken, s.BaseAmount, s.PoolToken))
}
