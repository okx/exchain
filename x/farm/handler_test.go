package farm

import (
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

type testContext struct {
	ctx               sdk.Context
	k                 Keeper
	swapTokenPairs    []swaptypes.SwapTokenPair
	tokenOwner        sdk.AccAddress
	nonPairTokenName  []string
	nonExistTokenName []string
	addrList          []sdk.AccAddress // 1000 okt per address
	handler           sdk.Handler
}

func initEnvironment(t *testing.T) testContext {
	// init
	ctx, mk := keeper.GetKeeper(t)
	k := mk.Keeper

	var blockHeight int64 = 10
	ctx = ctx.WithBlockHeight(blockHeight)
	BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: blockHeight}}, k)

	testBaseTokenName := swaptypes.TestBasePooledToken
	testQuoteTokenName := swaptypes.TestBasePooledToken2
	testQuoteTokenName2 := swaptypes.TestBasePooledToken3
	nonExistTokenName := "fff"

	token.NewTestToken(t, ctx, mk.TokenKeeper, mk.BankKeeper, testBaseTokenName, keeper.Addrs)
	token.NewTestToken(t, ctx, mk.TokenKeeper, mk.BankKeeper, testQuoteTokenName, keeper.Addrs)
	token.NewTestToken(t, ctx, mk.TokenKeeper, mk.BankKeeper, testQuoteTokenName2, keeper.Addrs)

	var initPoolTokenAmount int64 = 100
	testBaseToken := sdk.NewDecCoinFromDec(testBaseTokenName, sdk.NewDec(initPoolTokenAmount))
	testQuoteToken := sdk.NewDecCoinFromDec(testQuoteTokenName, sdk.NewDec(initPoolTokenAmount))
	testAddr := keeper.Addrs[0]
	testSwapTokenPair := swap.NewTestSwapTokenPairWithInitLiquidity(t, ctx, mk.SwapKeeper, testBaseToken, testQuoteToken, testAddr)

	//acc := mk.AccKeeper.GetAccount(ctx, Addrs[0])
	//fmt.Println(acc)

	handler := NewHandler(k)

	return testContext{
		ctx:               ctx,
		k:                 k,
		swapTokenPairs:    []swap.SwapTokenPair{testSwapTokenPair},
		tokenOwner:        testAddr,
		nonPairTokenName:  []string{testQuoteTokenName2},
		nonExistTokenName: []string{nonExistTokenName},
		addrList:          keeper.Addrs[1:],
		handler:           handler,
	}
}

func createPool(t *testing.T, tCtx testContext) types.MsgCreatePool {
	// create pool
	testSwapTokenPair := tCtx.swapTokenPairs[0]
	testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
	owner := tCtx.tokenOwner
	poolName := "abc"
	createPoolMsg := types.NewMsgCreatePool(owner, poolName, testSwapTokenPair.PoolTokenName, testYieldTokenName)
	result := tCtx.handler(tCtx.ctx, createPoolMsg)
	require.True(t, result.IsOK())
	return createPoolMsg
}

func destroyPool(t *testing.T, tCtx testContext, createPoolMsg types.MsgCreatePool) {
	destroyPoolMsg := types.NewMsgDestroyPool(createPoolMsg.Owner, createPoolMsg.PoolName)
	result := tCtx.handler(tCtx.ctx, destroyPoolMsg)
	require.True(t, result.IsOK())
}

func provide(t *testing.T, tCtx testContext, msgCreatePool types.MsgCreatePool) {
	poolName := msgCreatePool.PoolName
	address := msgCreatePool.Owner
	amount := sdk.NewDecCoinFromDec(msgCreatePool.YieldedSymbol, sdk.NewDec(10))
	amountYieldedPerBlock := sdk.NewDec(1)
	startBlockHeight := tCtx.ctx.BlockHeight() + 1
	provideMsg := types.NewMsgProvide(poolName, address, amount, amountYieldedPerBlock, startBlockHeight)
	result := tCtx.handler(tCtx.ctx, provideMsg)
	require.True(t, result.IsOK())
}

func lock(t *testing.T, tCtx testContext, msgCreatePool types.MsgCreatePool) {
	poolName := msgCreatePool.PoolName
	address := msgCreatePool.Owner
	amount := sdk.NewDecCoinFromDec(msgCreatePool.LockedSymbol, sdk.NewDec(1))
	lockMsg := types.NewMsgLock(poolName, address, amount)
	result := tCtx.handler(tCtx.ctx, lockMsg)
	require.True(t, result.IsOK())
}

func unlock(t *testing.T, tCtx testContext, msgCreatePool types.MsgCreatePool) {
	poolName := msgCreatePool.PoolName
	address := msgCreatePool.Owner
	amount := sdk.NewDecCoinFromDec(msgCreatePool.LockedSymbol, sdk.NewDec(1))
	unlockMsg := types.NewMsgUnlock(poolName, address, amount)
	result := tCtx.handler(tCtx.ctx, unlockMsg)
	require.True(t, result.IsOK())
}

func claim(t *testing.T, tCtx testContext, msgCreatePool types.MsgCreatePool) {
	claimMsg := types.NewMsgClaim(msgCreatePool.PoolName, msgCreatePool.Owner)
	result := tCtx.handler(tCtx.ctx, claimMsg)
	require.True(t, result.IsOK())
}

func TestHandleMsgCreatePool(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPool(t, tCtx)
}

func TestHandlerMsgDestroyPool(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// destroy pool
	destroyPool(t, tCtx, createPoolMsg)
}

func TestHandlerMsgProvide(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	msgCreatePool := createPool(t, tCtx)

	// provide
	provide(t, tCtx, msgCreatePool)
}

func TestHandlerMsgLock(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	msgCreatePool := createPool(t, tCtx)

	// provide
	provide(t, tCtx, msgCreatePool)

	// lock
	lock(t, tCtx, msgCreatePool)
}

func TestHandlerMsgUnlock(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	msgCreatePool := createPool(t, tCtx)

	// provide
	provide(t, tCtx, msgCreatePool)

	// lock
	lock(t, tCtx, msgCreatePool)

	// unlock
	unlock(t, tCtx, msgCreatePool)
}

func TestHandlerMsgClaim(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	msgCreatePool := createPool(t, tCtx)

	// provide
	provide(t, tCtx, msgCreatePool)

	// lock
	lock(t, tCtx, msgCreatePool)

	// claim
	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 2)
	claim(t, tCtx, msgCreatePool)
}
