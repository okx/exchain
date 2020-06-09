package poolswap

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okchain/x/poolswap/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
	"time"
)

func TestHandleMsgCreateExchange(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	testToken := types.InitPoolToken(types.TestBasePooledToken)

	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	handler := NewHandler(keeper)
	msg := types.NewMsgCreateExchange(testToken.Symbol, addrKeysSlice[0].Address)


	// test case1: token is not exist
	result := handler(ctx, msg)
	require.NotNil(t, result.Log)

	mapp.tokenKeeper.NewToken(ctx, testToken)

	// test case2: success
	result = handler(ctx, msg)
	require.Equal(t, "", result.Log)

	// check account balance
	acc := mapp.AccountKeeper.GetAccount(ctx, addrKeysSlice[0].Address)
	expectCoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.MustNewDecFromStr("100")),
		sdk.NewDecCoinFromDec(types.TestBasePooledToken, sdk.MustNewDecFromStr("100")),
	}
	require.EqualValues(t, expectCoins.String(), acc.GetCoins().String())

	expectSwapTokenPair := types.GetTestSwapTokenPair()
	swapTokenPair, err := keeper.GetSwapTokenPair(ctx, types.TestSwapTokenPairName)
	require.Nil(t, err)
	require.EqualValues(t, expectSwapTokenPair, swapTokenPair)

	// test case3: swapTokenPair already exists
	result = handler(ctx, msg)
	require.NotNil(t, result.Log)
}

func TestHandleMsgAddLiquidity(t *testing.T) {
	mapp, addrKeysSlice := getMockAppWithBalance(t, 1, 100000000)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	testToken := types.InitPoolToken(types.TestBasePooledToken)

	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	handler := NewHandler(keeper)
	msg := types.NewMsgCreateExchange(testToken.Symbol, addrKeysSlice[0].Address)
	mapp.tokenKeeper.NewToken(ctx, testToken)

	result := handler(ctx, msg)
	require.Equal(t, "", result.Log)

	minLiquidity := sdk.NewDec(1)
	maxBaseAmount := sdk.NewDecCoinFromDec(types.TestBasePooledToken, sdk.NewDec(1000000))
	quoteAmount := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.NewDec(100000))
	deadLine := time.Now().Unix()
	addr := addrKeysSlice[0].Address

	tests := []struct {
		testCase string
		minLiquidity sdk.Dec
		maxBaseAmount sdk.DecCoin
		quoteAmount sdk.DecCoin
		deadLine int64
		addr sdk.AccAddress
		exceptResultCode sdk.CodeType
	}{
		{"success", minLiquidity, maxBaseAmount, quoteAmount, deadLine, addr, 0},
	}

	for _, testCase := range tests {
		addLiquidityMsg := types.NewMsgAddLiquidity(testCase.minLiquidity, testCase.maxBaseAmount, testCase.quoteAmount, testCase.deadLine, testCase.addr)
		result = handler(ctx, addLiquidityMsg)
		require.Equal(t, testCase.exceptResultCode, result.Code)
	}
}


func TestBenchmark(t *testing.T) {
	mapp, addrKeysSlice := getMockAppWithBalance(t, 1, 100000000)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10).WithBlockTime(time.Now())
	testToken := types.InitPoolToken(types.TestBasePooledToken)

	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	handler := NewHandler(keeper)
	msg := types.NewMsgCreateExchange(testToken.Symbol, addrKeysSlice[0].Address)

	mapp.tokenKeeper.NewToken(ctx, testToken)

	result := handler(ctx, msg)
	require.Equal(t, "", result.Log)

	minLiquidity := sdk.NewDec(1)
	maxBaseAmount := sdk.NewDecCoinFromDec(types.TestBasePooledToken, sdk.NewDec(1000000))
	quoteAmount := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.NewDec(100000))
	deadLine := time.Now().Unix()
	addr := addrKeysSlice[0].Address

	addLiquidityMsg := types.NewMsgAddLiquidity(minLiquidity, maxBaseAmount, quoteAmount, deadLine, addr)
	result = handler(ctx, addLiquidityMsg)
	require.Equal(t, "", result.Log)


	minBoughtTokenAmount := sdk.NewDecCoinFromDec(types.TestBasePooledToken, sdk.NewDec(1))
	deadLine = time.Now().Unix()
	soldTokenAmountList := [3]int{200, 500, 1000}

	for _, soldTokenAmountInt := range soldTokenAmountList {
		startTime := time.Now()
		soldTokenAmount := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.NewDec(int64(soldTokenAmountInt)))
		tokenToNativeTokenMsg := types.NewMsgTokenToNativeToken(soldTokenAmount, minBoughtTokenAmount, deadLine, addr, addr)
		result = handler(ctx, tokenToNativeTokenMsg)
		require.Equal(t, "", result.Log)
		fmt.Println(time.Since(startTime))
	}

}

