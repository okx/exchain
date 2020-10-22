package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	swap "github.com/okex/okexchain/x/ammswap"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/okex/okexchain/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestHandleMsgCreatePool(t *testing.T) {
	// init
	ctx, mk := keeper.GetKeeper(t)
	k := mk.Keeper

	var blockHeight int64 = 10
	ctx = ctx.WithBlockHeight(blockHeight)
	BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: blockHeight}}, k)

	testBaseTokenName := swaptypes.TestBasePooledToken
	testQuoteTokenName := swaptypes.TestBasePooledToken2
	testQuoteTokenName2 := swaptypes.TestBasePooledToken3
	testYieldTokenName := swaptypes.TestBasePooledToken

	token.NewTestToken(t, ctx, mk.TokenKeeper, mk.BankKeeper, testBaseTokenName, keeper.Addrs)
	token.NewTestToken(t, ctx, mk.TokenKeeper, mk.BankKeeper, testQuoteTokenName, keeper.Addrs)
	token.NewTestToken(t, ctx, mk.TokenKeeper, mk.BankKeeper, testQuoteTokenName2, keeper.Addrs)


	var initPoolTokenAmount int64 = 100
	testBaseToken := sdk.NewDecCoinFromDec(testBaseTokenName, sdk.NewDec(initPoolTokenAmount))
	testQuoteToken := sdk.NewDecCoinFromDec(testQuoteTokenName, sdk.NewDec(initPoolTokenAmount))
	testAddr := keeper.Addrs[0]
	testSwapTokenPair := swap.NewTestSwapTokenPairWithInitLiquidity(t, ctx, mk.SwapKeeper, testBaseToken, testQuoteToken, testAddr)

	acc := mk.AccKeeper.GetAccount(ctx, keeper.Addrs[0])
	fmt.Println(acc)

	handler := NewHandler(k)

	// create pool
	createPoolMsg := types.NewMsgCreatePool(testAddr, "abc", testSwapTokenPair.PoolTokenName, testYieldTokenName)
	result := handler(ctx, createPoolMsg)
	require.True(t, result.IsOK())


}
