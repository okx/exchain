package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
)

// nolint
const TestBasePooledToken = "aab"
const TestBasePooledToken2 = "ccb"
const TestBasePooledToken3 = "ddb"
const TestQuotePooledToken = common.NativeToken
const TestSwapTokenPairName = TestBasePooledToken + "_" + TestQuotePooledToken

// GetTestSwapTokenPair just for test
func GetTestSwapTokenPair() SwapTokenPair {
	return SwapTokenPair{
		QuotePooledCoin: sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(0)),
		BasePooledCoin:  sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(0)),
		PoolTokenName:   GetPoolTokenName(TestBasePooledToken, TestQuotePooledToken),
	}
}
