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

type getMsgFunc func(tCtx *testContext, preData interface{}) sdk.Msg

type preExecFunc func(t *testing.T, tCtx *testContext) interface{}

type verificationFunc func(t *testing.T, tCtx *testContext, result sdk.Result, testCase testCaseItem)

var verification verificationFunc = func(t *testing.T, context *testContext, result sdk.Result, testCase testCaseItem) {
	require.Equal(t, testCase.expectedCode, result.Code)
}

type testCaseItem struct {
	caseName     string           // the name of the case
	preExec      preExecFunc      // function "preExec" executes the code before executing the specific handler to be tested
	getMsg       getMsgFunc       // function "getMsg" returns a sdk.Msg for testing, this msg will be tested by executing the function "handler"
	verification verificationFunc // function "verification" Verifies that the test results are the same as expected
	expectedCode sdk.CodeType     // expectedCode represents the expected code in the test result
}

func testCaseTest(t *testing.T, testCaseList []testCaseItem) {
	for _, testCase := range testCaseList {
		fmt.Println(testCase.caseName)
		tCtx := initEnvironment(t)
		preData := testCase.preExec(t, tCtx)
		msg := testCase.getMsg(tCtx, preData)
		result := tCtx.handler(tCtx.ctx, msg)
		testCase.verification(t, tCtx, result, testCase)
	}
}

func initEnvironment(t *testing.T) *testContext {
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

	return &testContext{
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

func createPool(t *testing.T, tCtx *testContext) types.MsgCreatePool {
	// create pool
	testSwapTokenPair := tCtx.swapTokenPairs[0]
	testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
	owner := tCtx.tokenOwner
	poolName := "abc"
	createPoolMsg := types.NewMsgCreatePool(owner, poolName, testSwapTokenPair.PoolTokenName, testYieldTokenName)
	result := tCtx.handler(tCtx.ctx, createPoolMsg)
	require.True(t, result.IsOK())

	k := tCtx.k
	_, found := k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
	require.True(t, found)
	return createPoolMsg
}

func destroyPool(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	k := tCtx.k
	_, found := k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
	require.True(t, found)
	destroyPoolMsg := types.NewMsgDestroyPool(createPoolMsg.Owner, createPoolMsg.PoolName)
	result := tCtx.handler(tCtx.ctx, destroyPoolMsg)
	require.True(t, result.IsOK())
	_, found = k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
	require.False(t, found)
}

func provide(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	poolName := createPoolMsg.PoolName
	address := createPoolMsg.Owner
	amount := sdk.NewDecCoinFromDec(createPoolMsg.YieldedSymbol, sdk.NewDec(10))
	amountYieldedPerBlock := sdk.NewDec(1)
	startBlockHeight := tCtx.ctx.BlockHeight() + 1
	provideMsg := types.NewMsgProvide(poolName, address, amount, amountYieldedPerBlock, startBlockHeight)
	result := tCtx.handler(tCtx.ctx, provideMsg)
	require.True(t, result.IsOK())
}

func lock(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	poolName := createPoolMsg.PoolName
	address := createPoolMsg.Owner
	amount := sdk.NewDecCoinFromDec(createPoolMsg.LockedSymbol, sdk.NewDec(1))
	lockMsg := types.NewMsgLock(poolName, address, amount)
	result := tCtx.handler(tCtx.ctx, lockMsg)
	require.True(t, result.IsOK())
}

func unlock(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	poolName := createPoolMsg.PoolName
	address := createPoolMsg.Owner
	amount := sdk.NewDecCoinFromDec(createPoolMsg.LockedSymbol, sdk.NewDec(1))
	unlockMsg := types.NewMsgUnlock(poolName, address, amount)
	result := tCtx.handler(tCtx.ctx, unlockMsg)
	require.True(t, result.IsOK())
}

func claim(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	claimMsg := types.NewMsgClaim(createPoolMsg.PoolName, createPoolMsg.Owner)
	result := tCtx.handler(tCtx.ctx, claimMsg)
	require.True(t, result.IsOK())
}

func setWhite(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	setWhiteMsg := types.NewMsgSetWhite(createPoolMsg.PoolName, createPoolMsg.Owner)
	result := tCtx.handler(tCtx.ctx, setWhiteMsg)
	require.True(t, result.IsOK())
}

func TestHandleMsgCreatePool(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPool(t, tCtx)
}

func TestHandlerMsgCreatePoolInvalid(t *testing.T) {
	preExec := func(t *testing.T, tCtx *testContext) interface{} {
		return nil
	}
	var normalGetMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
		testSwapTokenPair := tCtx.swapTokenPairs[0]
		testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
		owner := tCtx.tokenOwner
		poolName := "abc"
		createPoolMsg := types.NewMsgCreatePool(owner, poolName, testSwapTokenPair.PoolTokenName, testYieldTokenName)
		return createPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName:     "success",
			preExec:      preExec,
			getMsg:       normalGetMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. farm pool already exists",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				return createPool(t, tCtx)
			},
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				createPoolMsg := preData.(types.MsgCreatePool)
				return createPoolMsg
			},
			verification: verification,
			expectedCode: types.CodePoolAlreadyExist,
		},
		{
			caseName: "failed. lock token does not exists",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				testSwapTokenPair := tCtx.swapTokenPairs[0]
				testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
				lockSymbol := tCtx.nonExistTokenName[0]
				owner := tCtx.tokenOwner
				poolName := "abc"
				createPoolMsg := types.NewMsgCreatePool(owner, poolName, lockSymbol, testYieldTokenName)
				return createPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeTokenNotExist,
		},
		{
			caseName: "failed. the addr isn't the owner of the token",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				testSwapTokenPair := tCtx.swapTokenPairs[0]
				testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
				owner := tCtx.addrList[0]
				poolName := "abc"
				createPoolMsg := types.NewMsgCreatePool(owner, poolName, testSwapTokenPair.PoolTokenName, testYieldTokenName)
				return createPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. yield token does not exists",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				testSwapTokenPair := tCtx.swapTokenPairs[0]
				testYieldTokenName := tCtx.nonExistTokenName[0]
				owner := tCtx.addrList[0]
				poolName := "abc"
				createPoolMsg := types.NewMsgCreatePool(owner, poolName, testSwapTokenPair.PoolTokenName, testYieldTokenName)
				return createPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. insufficient fee coins",
			preExec: func(t *testing.T, context *testContext) interface{} {
				params := context.k.GetParams(context.ctx)
				params.CreatePoolFee = sdk.NewDecCoinFromDec(context.nonExistTokenName[0], sdk.NewDec(1))
				context.k.SetParams(context.ctx, params)
				return nil
			},
			getMsg:       normalGetMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientFee,
		},
		{
			caseName: "failed. insufficient coins",
			preExec: func(t *testing.T, context *testContext) interface{} {
				params := context.k.GetParams(context.ctx)
				params.CreatePoolDeposit = sdk.NewDecCoinFromDec(context.nonExistTokenName[0], sdk.NewDec(1))
				context.k.SetParams(context.ctx, params)
				return nil
			},
			getMsg:       normalGetMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
	}
	testCaseTest(t, tests)
}

func TestHandlerMsgDestroyPool(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// destroy pool
	destroyPool(t, tCtx, createPoolMsg)
}

func TestHandlerMsgDestroyPoolInvalid(t *testing.T) {
	preExec := func(t *testing.T, tCtx *testContext) interface{} {
		// create pool
		createPoolMsg := createPool(t, tCtx)
		return createPoolMsg
	}
	var normalGetMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
		createPoolMsg := preData.(types.MsgCreatePool)
		addr := createPoolMsg.Owner
		poolName := createPoolMsg.PoolName
		destroyPoolMsg := types.NewMsgDestroyPool(addr, poolName)
		return destroyPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName: "success",
			preExec:  preExec,
			getMsg: normalGetMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. Farm pool does not exist",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				testSwapTokenPair := tCtx.swapTokenPairs[0]
				testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
				owner := tCtx.tokenOwner
				poolName := "abc"
				createPoolMsg := types.NewMsgCreatePool(owner, poolName, testSwapTokenPair.PoolTokenName, testYieldTokenName)
				return createPoolMsg
			},
			getMsg: normalGetMsg,
			verification: verification,
			expectedCode: types.CodeInvalidFarmPool,
		},
		{
			caseName: "failed. the address isn't the owner of pool",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				createPoolMsg := preData.(types.MsgCreatePool)
				addr := tCtx.addrList[0]
				poolName := createPoolMsg.PoolName
				destroyPoolMsg := types.NewMsgDestroyPool(addr, poolName)
				return destroyPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "insufficient fee coins",
			preExec:  func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// modify params
				pools, found := tCtx.k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
				require.True(t, found)
				pools.DepositAmount = sdk.NewDecCoinFromDec(tCtx.nonExistTokenName[0], sdk.NewDec(1))
				tCtx.k.SetFarmPool(tCtx.ctx, pools)
				return createPoolMsg
			},
			getMsg: normalGetMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
		{
			caseName: "failed. the pool is not finished and can not be destroyed",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				return createPoolMsg
			},
			getMsg: normalGetMsg,
			verification: verification,
			expectedCode: types.CodePoolNotFinished,
		},
		{
			caseName: "insufficient rewards coins",
			preExec:  preExec,
			getMsg: normalGetMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
	}
	testCaseTest(t, tests)
}

func TestHandlerMsgProvide(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	ceatePoolMsg := createPool(t, tCtx)

	// provide
	provide(t, tCtx, ceatePoolMsg)
}

func TestHandlerMsgLock(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// provide
	provide(t, tCtx, createPoolMsg)

	// lock
	lock(t, tCtx, createPoolMsg)
}

func TestHandlerMsgUnlock(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// provide
	provide(t, tCtx, createPoolMsg)

	// lock
	lock(t, tCtx, createPoolMsg)

	// unlock
	unlock(t, tCtx, createPoolMsg)
}

func TestHandlerMsgClaim(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// provide
	provide(t, tCtx, createPoolMsg)

	// lock
	lock(t, tCtx, createPoolMsg)

	// claim
	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 2)
	claim(t, tCtx, createPoolMsg)
}

func TestHandlerMsgSetWhite(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// setWhite
	setWhite(t, tCtx, createPoolMsg)
}
