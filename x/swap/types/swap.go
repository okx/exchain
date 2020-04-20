package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	token "github.com/okex/okchain/x/token/types"
	"github.com/okex/okchain/x/common"
	"strings"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const TestBasePooledToken = "xxb"
const TestQuotePooledToken = common.NativeToken
const TestSwapTokenPairName = TestBasePooledToken + "_" + TestQuotePooledToken
//
//type SwapTokenPairPools struct {
//	SwapTokenPairs []SwapTokenPair `json:"swap_token_pairs"`
//}

type SwapTokenPair struct {
	QuotePooledCoin sdk.DecCoin `json:"quote_pooled_coin"`
	BasePooledCoin  sdk.DecCoin `json:"base_pooled_coin"`
	PoolTokenName   string      `json:"pool_token_name"` //the name of poolToken
}

func NewSwapTokenPair(quotePooledCoin sdk.DecCoin, basePooledCoin sdk.DecCoin, poolTokenName string) *SwapTokenPair {
	swapTokenPair := &SwapTokenPair{
		QuotePooledCoin: quotePooledCoin,
		BasePooledCoin:  basePooledCoin,
		PoolTokenName:   poolTokenName,
	}
	return swapTokenPair
}

// implement fmt.Stringer
func (s SwapTokenPair) String() string {
	return strings.TrimSpace(fmt.Sprintf(`QuotePooledCoin: %s
BasePooledCoin: %s
PoolTokenName: %s`, s.QuotePooledCoin.String(), s.BasePooledCoin.String(), s.PoolTokenName))
}

func (s SwapTokenPair) TokenPairName() string {
	return s.BasePooledCoin.Denom + "_" + s.QuotePooledCoin.Denom
}

func GetTestSwapTokenPair() SwapTokenPair {
	return SwapTokenPair{
		QuotePooledCoin: sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(0)),
		BasePooledCoin:  sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(0)),
		PoolTokenName:   "spt",
	}
}

func InitPoolToken(poolTokenName string) token.Token {
	return token.Token{
		Description:         poolTokenName,
		Symbol:              poolTokenName,
		OriginalSymbol:      poolTokenName,
		WholeName:           poolTokenName,
		OriginalTotalSupply: sdk.NewDec(0),
		TotalSupply:         sdk.NewDec(0),
		Owner:               supply.NewModuleAddress(ModuleName),
		Mintable:            true,
	}
}
