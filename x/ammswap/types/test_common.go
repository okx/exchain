package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/token"
	tokentypes "github.com/okex/okexchain/x/token/types"
)

// nolint
const TestBasePooledToken = "aab"
const TestBasePooledToken2 = "ccb"
const TestBasePooledToken3 = "ddb"
const TestBasePooledToken4 = "kkb"
const TestBasePooledToken5 = "ffb"
const TestQuotePooledToken = common.NativeToken
const TestQuotePooledToken2 = TestBasePooledToken5
const TestSwapTokenPairName = TestBasePooledToken + "_" + TestQuotePooledToken

// GetTestSwapTokenPair just for test
func GetTestSwapTokenPair() SwapTokenPair {
	return SwapTokenPair{
		QuotePooledCoin: sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(1)),
		BasePooledCoin:  sdk.NewDecCoinFromDec(TestBasePooledToken, sdk.NewDec(1)),
		PoolTokenName:   GetPoolTokenName(TestBasePooledToken, TestQuotePooledToken),
	}
}

func GetTestSwapTokenPairWithLargeLiquidity() SwapTokenPair {
	return SwapTokenPair{
		QuotePooledCoin: sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(10000)),
		BasePooledCoin:  sdk.NewDecCoinFromDec(TestBasePooledToken2, sdk.NewDec(100)),
		PoolTokenName:   GetPoolTokenName(TestBasePooledToken2, TestQuotePooledToken),
	}
}

func GetTestSwapTokenPairWithZeroLiquidity() SwapTokenPair {
	return SwapTokenPair{
		QuotePooledCoin: sdk.NewDecCoinFromDec(TestQuotePooledToken, sdk.NewDec(0)),
		BasePooledCoin:  sdk.NewDecCoinFromDec(TestBasePooledToken3, sdk.NewDec(0)),
		PoolTokenName:   GetPoolTokenName(TestBasePooledToken3, TestQuotePooledToken),
	}
}

func CreateTestMsgs(addr sdk.AccAddress) []sdk.Msg {
	return []sdk.Msg{
		NewMsgCreateExchange(TestBasePooledToken4, TestQuotePooledToken, addr),
		NewMsgCreateExchange(TestBasePooledToken, TestQuotePooledToken2, addr),
		NewMsgAddLiquidity(sdk.ZeroDec(),
			sdk.NewDecCoin(TestBasePooledToken4, sdk.OneInt()), sdk.NewDecCoin(TestQuotePooledToken, sdk.OneInt()),
			time.Now().Add(time.Hour).Unix(), addr),
		NewMsgAddLiquidity(sdk.ZeroDec(),
			sdk.NewDecCoin(TestBasePooledToken, sdk.OneInt()), sdk.NewDecCoin(TestQuotePooledToken2, sdk.OneInt()),
			time.Now().Add(time.Hour).Unix(), addr),
		NewMsgRemoveLiquidity(sdk.OneDec(),
			sdk.NewDecCoin(TestBasePooledToken, sdk.OneInt()), sdk.NewDecCoin(TestQuotePooledToken2, sdk.OneInt()),
			time.Now().Add(time.Hour).Unix(), addr),
	}
}

func SetTestTokens(ctx sdk.Context, tokenKeeper token.Keeper, supplyKeeper supply.Keeper, addr sdk.AccAddress) error {
	balance := 100
	coins, err := sdk.ParseDecCoins(fmt.Sprintf("%d%s,%d%s,%d%s,%d%s,%d%s,%d%s",
		balance, TestQuotePooledToken, balance, TestBasePooledToken, balance, TestBasePooledToken2, balance, TestBasePooledToken3,
		balance, TestBasePooledToken4, balance, TestBasePooledToken5))
	if err != nil {
		return err
	}

	for _, coin := range coins {
		name := coin.Denom
		tokenKeeper.NewToken(ctx, tokentypes.Token{"", name, name,name, coin.Amount, 1,addr,true})
	}
	err = supplyKeeper.MintCoins(ctx, tokentypes.ModuleName, coins)
	if err != nil {
		return err
	}
	err = supplyKeeper.SendCoinsFromModuleToAccount(ctx, tokentypes.ModuleName, addr, coins)
	if err != nil {
		return err
	}
	return nil
}