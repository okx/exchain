package types

import (
	"time"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/supply"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/token"
	tokentypes "github.com/okex/exchain/x/token/types"
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

func SetTestTokens(ctx sdk.Context, tokenKeeper token.Keeper, supplyKeeper supply.Keeper, addr sdk.AccAddress, coins sdk.DecCoins) error {
	for _, coin := range coins {
		name := coin.Denom
		tokenKeeper.NewToken(ctx, tokentypes.Token{"", name, name,name, coin.Amount, 1,addr,true})
	}
	err := supplyKeeper.MintCoins(ctx, tokentypes.ModuleName, coins)
	if err != nil {
		return err
	}
	err = supplyKeeper.SendCoinsFromModuleToAccount(ctx, tokentypes.ModuleName, addr, coins)
	if err != nil {
		return err
	}
	return nil
}

func CreateTestMsgs(addr sdk.AccAddress) []sdk.Msg {
	return []sdk.Msg{
		NewMsgCreateExchange(TestBasePooledToken, TestQuotePooledToken, addr),
		NewMsgCreateExchange(TestBasePooledToken2, TestQuotePooledToken, addr),
		NewMsgAddLiquidity(sdk.ZeroDec(),
			sdk.NewDecCoin(TestBasePooledToken, sdk.OneInt()), sdk.NewDecCoin(TestQuotePooledToken, sdk.OneInt()),
			time.Now().Add(time.Hour).Unix(), addr),
		NewMsgAddLiquidity(sdk.ZeroDec(),
			sdk.NewDecCoin(TestBasePooledToken2, sdk.OneInt()), sdk.NewDecCoin(TestQuotePooledToken, sdk.OneInt()),
			time.Now().Add(time.Hour).Unix(), addr),
		NewMsgRemoveLiquidity(sdk.OneDec(),
			sdk.NewDecCoin(TestBasePooledToken2, sdk.OneInt()), sdk.NewDecCoin(TestQuotePooledToken, sdk.OneInt()),
			time.Now().Add(time.Hour).Unix(), addr),
	}
}