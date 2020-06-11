package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
)

const TestBasePooledToken = "xxb"
const TestBasePooledToken2 = "yyb"
const TestBasePooledToken3 = "zzb"
const TestQuotePooledToken = common.NativeToken
const TestSwapTokenPairName = TestBasePooledToken + "_" + TestQuotePooledToken

func GetTestSwapTokenPair() SwapTokenPair {
	return SwapTokenPair{
		QuotePooledCoin: sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(0)),
		BasePooledCoin:  sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(0)),
		PoolTokenName:   PoolTokenPrefix + TestBasePooledToken,
	}
}