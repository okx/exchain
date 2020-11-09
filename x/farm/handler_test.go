package farm

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	swap "github.com/okex/okexchain/x/ammswap"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/okex/okexchain/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
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

type verificationFunc func(t *testing.T, tCtx *testContext, result sdk.Result, testCase testCaseItem, preCoins, afterCoins sdk.SysCoins, preData interface{})

var verification verificationFunc = func(t *testing.T, context *testContext, result sdk.Result, testCase testCaseItem, preCoins, afterCoins sdk.SysCoins, preData interface{}) {
	if result.Code != testCase.expectedCode {
		fmt.Println(result.Log)
	}
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
		addrList := msg.GetSigners()
		addr := addrList[0]
		preCoins := tCtx.k.TokenKeeper().GetCoins(tCtx.ctx, addr)
		result := tCtx.handler(tCtx.ctx, msg)
		afterCoins := tCtx.k.TokenKeeper().GetCoins(tCtx.ctx, addr)
		testCase.verification(t, tCtx, result, testCase, preCoins, afterCoins, preData)
	}
}

func testCaseCombinationTest(t *testing.T, testCaseList []testCaseItem) {
	tCtx := initEnvironment(t)
	for _, testCase := range testCaseList {
		fmt.Println(testCase.caseName)
		preData := testCase.preExec(t, tCtx)
		msg := testCase.getMsg(tCtx, preData)
		addrList := msg.GetSigners()
		addr := addrList[0]
		preCoins := tCtx.k.TokenKeeper().GetCoins(tCtx.ctx, addr)
		result := tCtx.handler(tCtx.ctx, msg)
		afterCoins := tCtx.k.TokenKeeper().GetCoins(tCtx.ctx, addr)
		testCase.verification(t, tCtx, result, testCase, preCoins, afterCoins, preData)
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
	testSwapTokenPair := swap.NewTestSwapTokenPairWithInitLiquidity(t, ctx, mk.SwapKeeper, testBaseToken, testQuoteToken, keeper.Addrs)

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

var normalGetCreatePoolMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
	testSwapTokenPair := tCtx.swapTokenPairs[0]
	testYieldTokenName := testSwapTokenPair.BasePooledCoin.Denom
	owner := tCtx.tokenOwner
	poolName := "abc"
	minLockAmount := sdk.NewDecCoinFromDec(testSwapTokenPair.PoolTokenName, sdk.ZeroDec())
	createPoolMsg := types.NewMsgCreatePool(owner, poolName, minLockAmount, testYieldTokenName)
	return createPoolMsg
}

var normalGetDestroyPoolMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
	createPoolMsg := preData.(types.MsgCreatePool)
	addr := createPoolMsg.Owner
	poolName := createPoolMsg.PoolName
	destroyPoolMsg := types.NewMsgDestroyPool(addr, poolName)
	return destroyPoolMsg
}

var normalGetProvideMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
	createPoolMsg := preData.(types.MsgCreatePool)
	poolName := createPoolMsg.PoolName
	address := createPoolMsg.Owner
	amount := sdk.NewDecCoinFromDec(createPoolMsg.YieldedSymbol, sdk.NewDec(10))
	amountYieldedPerBlock := sdk.NewDec(1)
	startBlockHeight := tCtx.ctx.BlockHeight() + 1
	provideMsg := types.NewMsgProvide(poolName, address, amount, amountYieldedPerBlock, startBlockHeight)
	return provideMsg
}

var normalGetLockMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
	createPoolMsg := preData.(types.MsgCreatePool)
	poolName := createPoolMsg.PoolName
	address := createPoolMsg.Owner
	amount := sdk.NewDecCoinFromDec(createPoolMsg.MinLockAmount.Denom, sdk.NewDec(1))
	lockMsg := types.NewMsgLock(poolName, address, amount)
	return lockMsg
}

var normalGetUnlockMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
	createPoolMsg := preData.(types.MsgCreatePool)
	poolName := createPoolMsg.PoolName
	address := createPoolMsg.Owner
	amount := sdk.NewDecCoinFromDec(createPoolMsg.MinLockAmount.Denom, sdk.NewDec(1))
	unlockMsg := types.NewMsgUnlock(poolName, address, amount)
	return unlockMsg
}

var normalGetClaimMsg getMsgFunc = func(tCtx *testContext, preData interface{}) sdk.Msg {
	createPoolMsg := preData.(types.MsgCreatePool)
	claimMsg := types.NewMsgClaim(createPoolMsg.PoolName, createPoolMsg.Owner)
	return claimMsg
}

func createPool(t *testing.T, tCtx *testContext) types.MsgCreatePool {
	createPoolMsg := normalGetCreatePoolMsg(tCtx, nil).(types.MsgCreatePool)
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
	destroyPoolMsg := normalGetDestroyPoolMsg(tCtx, createPoolMsg)
	result := tCtx.handler(tCtx.ctx, destroyPoolMsg)
	require.True(t, result.IsOK())
	_, found = k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
	require.False(t, found)
}

func provide(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) types.MsgProvide {
	provideMsg := normalGetProvideMsg(tCtx, createPoolMsg)
	result := tCtx.handler(tCtx.ctx, provideMsg)
	require.True(t, result.IsOK())
	return provideMsg.(types.MsgProvide)
}

func lock(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) types.MsgLock {
	lockMsg := normalGetLockMsg(tCtx, createPoolMsg)
	result := tCtx.handler(tCtx.ctx, lockMsg)
	if !result.IsOK() {
		fmt.Println(result.Log)
	}
	require.True(t, result.IsOK())
	return lockMsg.(types.MsgLock)
}

func unlock(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	unlockMsg := normalGetUnlockMsg(tCtx, createPoolMsg)
	result := tCtx.handler(tCtx.ctx, unlockMsg)
	require.True(t, result.IsOK())
}

func claim(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	claimMsg := normalGetClaimMsg(tCtx, createPoolMsg)
	result := tCtx.handler(tCtx.ctx, claimMsg)
	require.True(t, result.IsOK())
}

func setWhite(t *testing.T, tCtx *testContext, createPoolMsg types.MsgCreatePool) {
	setWhiteMsg := types.NewMsgSetWhite(createPoolMsg.PoolName, createPoolMsg.Owner)
	result := tCtx.handler(tCtx.ctx, setWhiteMsg)
	require.True(t, result.IsOK())
}

func TestHandlerMsgCreatePool(t *testing.T) {
	preExec := func(t *testing.T, tCtx *testContext) interface{} {
		return nil
	}

	tests := []testCaseItem{
		{
			caseName:     "success",
			preExec:      preExec,
			getMsg:       normalGetCreatePoolMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. create again after destroying",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				createPoolMsg := createPool(t, tCtx)

				provide(t, tCtx, createPoolMsg)

				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				claim(t, tCtx, createPoolMsg)

				unlock(t, tCtx, createPoolMsg)

				destroyPool(t, tCtx, createPoolMsg)

				return nil
			},
			getMsg:       normalGetCreatePoolMsg,
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
				createPoolMsg := normalGetCreatePoolMsg(tCtx, preData).(types.MsgCreatePool)
				createPoolMsg.MinLockAmount = sdk.NewDecCoinFromDec(tCtx.nonExistTokenName[0], sdk.ZeroDec())
				return createPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeTokenNotExist,
		},
		{
			caseName: "failed. yield token does not exists",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil).(types.MsgCreatePool)
				createPoolMsg.YieldedSymbol = tCtx.nonExistTokenName[0]
				return createPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeTokenNotExist,
		},
		{
			caseName: "failed. insufficient fee coins",
			preExec: func(t *testing.T, context *testContext) interface{} {
				params := context.k.GetParams(context.ctx)
				params.CreatePoolFee = sdk.NewDecCoinFromDec(context.nonExistTokenName[0], sdk.NewDec(1))
				context.k.SetParams(context.ctx, params)
				return nil
			},
			getMsg:       normalGetCreatePoolMsg,
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
			getMsg:       normalGetCreatePoolMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
	}
	testCaseTest(t, tests)
}

func TestHandlerMsgDestroyPool(t *testing.T) {
	preExec := func(t *testing.T, tCtx *testContext) interface{} {
		// create pool
		createPoolMsg := createPool(t, tCtx)
		return createPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName:     "success",
			preExec:      preExec,
			getMsg:       normalGetDestroyPoolMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. Farm pool does not exist",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil)
				return createPoolMsg
			},
			getMsg:       normalGetDestroyPoolMsg,
			verification: verification,
			expectedCode: types.CodeNoFarmPoolFound,
		},
		{
			caseName: "failed. the address isn't the owner of pool",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				destroyPoolMsg := normalGetDestroyPoolMsg(tCtx, preData).(types.MsgDestroyPool)
				destroyPoolMsg.Owner = tCtx.addrList[0]
				return destroyPoolMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. insufficient fee coins",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// modify params
				pools, found := tCtx.k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
				require.True(t, found)
				pools.DepositAmount = sdk.NewDecCoinFromDec(tCtx.nonExistTokenName[0], sdk.NewDec(1))
				tCtx.k.SetFarmPool(tCtx.ctx, pools)
				return createPoolMsg
			},
			getMsg:       normalGetDestroyPoolMsg,
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
			getMsg:       normalGetDestroyPoolMsg,
			verification: verification,
			expectedCode: types.CodePoolNotFinished,
		},
		{
			caseName: "success. destroy after providing",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				return createPoolMsg
			},
			getMsg:       normalGetDestroyPoolMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. destroy after claiming",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				claim(t, tCtx, createPoolMsg)

				unlock(t, tCtx, createPoolMsg)

				return createPoolMsg
			},
			getMsg:       normalGetDestroyPoolMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. insufficient rewards coins",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				pool, found := tCtx.k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
				require.True(t, found)
				updatedPool, _ := tCtx.k.CalculateAmountYieldedBetween(tCtx.ctx, pool)

				err := tCtx.k.SupplyKeeper().SendCoinsFromModuleToAccount(tCtx.ctx, YieldFarmingAccount, createPoolMsg.Owner, updatedPool.TotalAccumulatedRewards)
				require.Nil(t, err)

				return createPoolMsg
			},
			getMsg:       normalGetDestroyPoolMsg,
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

				// lock
				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				return createPoolMsg
			},
			getMsg:       normalGetDestroyPoolMsg,
			verification: verification,
			expectedCode: types.CodePoolNotFinished,
		},
	}
	testCaseTest(t, tests)
}

func TestHandlerMsgProvide(t *testing.T) {
	var preExec preExecFunc = func(t *testing.T, tCtx *testContext) interface{} {
		// create pool
		createPoolMsg := createPool(t, tCtx)
		return createPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName:     "success",
			preExec:      preExec,
			getMsg:       normalGetProvideMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. The start height to yield is less than current height",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				provideMsg := normalGetProvideMsg(tCtx, preData).(types.MsgProvide)
				provideMsg.StartHeightToYield = 0
				return provideMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. Farm pool does not exist",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil)
				return createPoolMsg
			},
			getMsg:       normalGetProvideMsg,
			verification: verification,
			expectedCode: types.CodeNoFarmPoolFound,
		},
		{
			caseName: "failed. The coin name should be %s, not %s",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				provideMsg := normalGetProvideMsg(tCtx, preData).(types.MsgProvide)
				provideMsg.Amount = sdk.NewDecCoinFromDec(tCtx.nonExistTokenName[0], provideMsg.Amount.Amount)
				return provideMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. The remaining amount is %s, so it's not enable to provide token repeatedly util amount become zero",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)
				return createPoolMsg
			},
			getMsg:       normalGetProvideMsg,
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "insufficient amount",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				provideMsg := normalGetProvideMsg(tCtx, preData).(types.MsgProvide)
				provideMsg.Amount = sdk.NewDecCoinFromDec(provideMsg.Amount.Denom, sdk.NewDec(1000000000))
				return provideMsg
			},
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
	}

	testCaseTest(t, tests)
}

func TestHandlerMsgLock(t *testing.T) {
	var preExec preExecFunc = func(t *testing.T, tCtx *testContext) interface{} {
		// create pool
		createPoolMsg := createPool(t, tCtx)

		// provide
		provide(t, tCtx, createPoolMsg)

		return createPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName:     "success",
			preExec:      preExec,
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. Farm pool does not exist",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil)
				return createPoolMsg
			},
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: types.CodeNoFarmPoolFound,
		},
		{
			caseName: "failed. The coin name should be %s, not %s",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				lockMsg := normalGetLockMsg(tCtx, preData).(types.MsgLock)
				lockMsg.Amount.Denom = tCtx.nonExistTokenName[0]
				return lockMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. lock amount %s must be greater than the pool`s min lock amount %s",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil).(types.MsgCreatePool)
				createPoolMsg.MinLockAmount.Amount = sdk.NewDec(math.MaxInt64)
				result := tCtx.handler(tCtx.ctx, createPoolMsg)
				require.True(t, result.IsOK())

				// provide
				provide(t, tCtx, createPoolMsg)

				return createPoolMsg
			},
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "success. has lockInfo",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				// lock
				lock(t, tCtx, createPoolMsg)

				return createPoolMsg
			},
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. withdraw failed",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provideMsg := provide(t, tCtx, createPoolMsg)

				// lock
				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				err := tCtx.k.SupplyKeeper().SendCoinsFromModuleToAccount(tCtx.ctx, types.YieldFarmingAccount, provideMsg.Address, sdk.NewCoins(provideMsg.Amount))
				require.Nil(t, err)
				return createPoolMsg
			},
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
		{
			caseName: "failed. insufficient coins",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				lockMsg := normalGetLockMsg(tCtx, preData).(types.MsgLock)
				lockMsg.Amount.Amount = sdk.NewDec(1000000)
				return lockMsg
			},
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
	}

	testCaseTest(t, tests)
}

func TestHandlerMsgUnlock(t *testing.T) {
	var preExec preExecFunc = func(t *testing.T, tCtx *testContext) interface{} {
		// create pool
		createPoolMsg := createPool(t, tCtx)

		// provide
		provide(t, tCtx, createPoolMsg)

		// lock
		lock(t, tCtx, createPoolMsg)

		return createPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName:     "success",
			preExec:      preExec,
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. the addr doesn't have any lock infos",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)
				return createPoolMsg
			},
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: types.CodeInvalidLockInfo,
		},
		{
			caseName: "failed. The coin name should be %s, not %s",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				unlockMsg := normalGetUnlockMsg(tCtx, preData).(types.MsgUnlock)
				unlockMsg.Amount.Denom = tCtx.nonExistTokenName[0]
				return unlockMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. The actual amount %s is less than %s",
			preExec:  preExec,
			getMsg: func(tCtx *testContext, preData interface{}) sdk.Msg {
				unlockMsg := normalGetUnlockMsg(tCtx, preData).(types.MsgUnlock)
				unlockMsg.Amount.Amount = unlockMsg.Amount.Amount.Add(sdk.NewDec(1))
				return unlockMsg
			},
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. remain lock amount %s is less than pool`s min lock amount %s",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil).(types.MsgCreatePool)
				createPoolMsg.MinLockAmount.Amount = sdk.NewDec(2)
				result := tCtx.handler(tCtx.ctx, createPoolMsg)
				require.True(t, result.IsOK())

				// provide
				provide(t, tCtx, createPoolMsg)

				// lock
				lockMsg := normalGetLockMsg(tCtx, createPoolMsg).(types.MsgLock)
				lockMsg.Amount.Amount = sdk.NewDec(2)
				result = tCtx.handler(tCtx.ctx, lockMsg)
				require.True(t, result.IsOK())

				return createPoolMsg
			},
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: types.CodeInvalidInput,
		},
		{
			caseName: "failed. Farm pool %s does not exist",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				preData := preExec(t, tCtx).(types.MsgCreatePool)
				tCtx.k.DeleteFarmPool(tCtx.ctx, preData.PoolName)
				return preData
			},
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: types.CodeNoFarmPoolFound,
		},
		{
			caseName: "failed. withdraw failed",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provideMsg := provide(t, tCtx, createPoolMsg)

				// lock
				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				err := tCtx.k.SupplyKeeper().SendCoinsFromModuleToAccount(tCtx.ctx, types.YieldFarmingAccount, provideMsg.Address, sdk.NewCoins(provideMsg.Amount))
				require.Nil(t, err)
				return createPoolMsg
			},
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
		{
			caseName: "failed. insufficient coins from module account",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				// lock
				lockMsg := lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				err := tCtx.k.SupplyKeeper().SendCoinsFromModuleToAccount(tCtx.ctx, ModuleName, lockMsg.Address, sdk.NewCoins(lockMsg.Amount))
				require.Nil(t, err)
				return createPoolMsg
			},
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
		{
			caseName: "success. lock and unlock without provide before",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// lock
				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)
				return createPoolMsg
			},
			getMsg:       normalGetUnlockMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
	}

	testCaseTest(t, tests)
}

func TestHandlerMsgClaim(t *testing.T) {
	var preExec preExecFunc = func(t *testing.T, tCtx *testContext) interface{} {
		// create pool
		createPoolMsg := createPool(t, tCtx)

		// provide
		provide(t, tCtx, createPoolMsg)

		// lock
		lock(t, tCtx, createPoolMsg)

		return createPoolMsg
	}
	tests := []testCaseItem{
		{
			caseName:     "success. claim after providing at the same block height",
			preExec:      preExec,
			getMsg:       normalGetClaimMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. claim after providing at the lower block height",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provide(t, tCtx, createPoolMsg)

				// lock
				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 2)

				return createPoolMsg
			},
			getMsg: normalGetClaimMsg,
			verification: func(t *testing.T, tCtx *testContext, result sdk.Result, testCase testCaseItem, preCoins, afterCoins sdk.SysCoins, preData interface{}) {
				if result.Code != testCase.expectedCode {
					fmt.Println(result.Log)
				}
				require.Equal(t, testCase.expectedCode, result.Code)
				createPoolMsg := preData.(types.MsgCreatePool)
				diffCoins := afterCoins.Sub(preCoins)
				actualDec := diffCoins.AmountOf(createPoolMsg.YieldedSymbol)
				require.Equal(t, sdk.NewDec(1), actualDec)
			},
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "failed. Farm pool %s does not exist",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil)
				return createPoolMsg
			},
			getMsg:       normalGetClaimMsg,
			verification: verification,
			expectedCode: types.CodeNoFarmPoolFound,
		},
		{
			caseName: "failed. withdraw failed",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				// create pool
				createPoolMsg := createPool(t, tCtx)

				// provide
				provideMsg := provide(t, tCtx, createPoolMsg)

				// lock
				lock(t, tCtx, createPoolMsg)

				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1000)

				err := tCtx.k.SupplyKeeper().SendCoinsFromModuleToAccount(tCtx.ctx, types.YieldFarmingAccount, provideMsg.Address, sdk.NewCoins(provideMsg.Amount))
				require.Nil(t, err)
				return createPoolMsg
			},
			getMsg:       normalGetClaimMsg,
			verification: verification,
			expectedCode: sdk.CodeInsufficientCoins,
		},
	}

	testCaseTest(t, tests)
}

func TestHandlerMsgSetWhite(t *testing.T) {
	// init
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	// setWhite
	setWhite(t, tCtx, createPoolMsg)
}

func TestNewHandler(t *testing.T) {
	// init
	tCtx := initEnvironment(t)
	msg := swaptypes.NewMsgCreateExchange(tCtx.swapTokenPairs[0].BasePooledCoin.Denom, tCtx.swapTokenPairs[0].QuotePooledCoin.Denom, tCtx.tokenOwner)
	result := tCtx.handler(tCtx.ctx, msg)
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func TestHandlerMultiLockAtOneBlockHeight(t *testing.T) {
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 10)
	// provide
	provide(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)
	// lock
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[2]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[3]
	lock(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)

	createPoolMsg.Owner = tCtx.addrList[4]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[5]
	lock(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)

	createPoolMsg.Owner = tCtx.addrList[6]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[7]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	lock(t, tCtx, createPoolMsg)

	curPeriodRewards := tCtx.k.GetPoolCurrentRewards(tCtx.ctx, createPoolMsg.PoolName)
	fmt.Println(string(types.ModuleCdc.MustMarshalJSON(curPeriodRewards)))
	//var period uint64
	//for period = 0;period < curPeriodRewards.Period;period++ {
	//	historyPeriodRewards := tCtx.k.GetPoolHistoricalRewards(tCtx.ctx, createPoolMsg.PoolName, period)
	//	fmt.Println("period:", period)
	//	fmt.Println(string(types.ModuleCdc.MustMarshalJSON(historyPeriodRewards)))
	//}
	tCtx.k.IterateAllLockInfos(tCtx.ctx, func(lockInfo types.LockInfo) (stop bool) {

		fmt.Println(lockInfo.String())
		return false
	})

}

func TestHandlerMultiLockAtOneBlockHeight2(t *testing.T) {
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 10)
	// provide
	provide(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)
	// lock
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[0]
	lock(t, tCtx, createPoolMsg)

	curPeriodRewards := tCtx.k.GetPoolCurrentRewards(tCtx.ctx, createPoolMsg.PoolName)
	fmt.Println(string(types.ModuleCdc.MustMarshalJSON(curPeriodRewards)))
	//var period uint64
	//for period = 0;period < curPeriodRewards.Period;period++ {
	//	historyPeriodRewards := tCtx.k.GetPoolHistoricalRewards(tCtx.ctx, createPoolMsg.PoolName, period)
	//	fmt.Println("period:", period)
	//	fmt.Println(string(types.ModuleCdc.MustMarshalJSON(historyPeriodRewards)))
	//}
	tCtx.k.IterateAllLockInfos(tCtx.ctx, func(lockInfo types.LockInfo) (stop bool) {

		fmt.Println(lockInfo.String())
		return false
	})

}

func TestHandlerMultiLockAndUnlock(t *testing.T) {
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 10)
	// provide
	provide(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)
	// lock
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[2]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[3]
	lock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[4]
	lock(t, tCtx, createPoolMsg)

	tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 4)

	createPoolMsg.Owner = tCtx.tokenOwner
	unlock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[1]
	unlock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[2]
	unlock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[3]
	unlock(t, tCtx, createPoolMsg)

	createPoolMsg.Owner = tCtx.addrList[4]
	unlock(t, tCtx, createPoolMsg)

	curPeriodRewards := tCtx.k.GetPoolCurrentRewards(tCtx.ctx, createPoolMsg.PoolName)
	fmt.Println(string(types.ModuleCdc.MustMarshalJSON(curPeriodRewards)))
	numHistoricalRewards := 0
	tCtx.k.IteratePoolHistoricalRewards(tCtx.ctx, createPoolMsg.PoolName,
		func(store sdk.KVStore, key []byte, value []byte) (stop bool) {
			var rewards types.PoolHistoricalRewards
			types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(value, &rewards)
			fmt.Println(string(key), rewards)
			numHistoricalRewards++
			return false
		})
	require.Equal(t, 1, numHistoricalRewards)
	numLockInfo := 0
	tCtx.k.IterateAllLockInfos(tCtx.ctx, func(lockInfo types.LockInfo) (stop bool) {
		numLockInfo++
		fmt.Println(lockInfo.String())
		return false
	})
	require.Equal(t, 0, numLockInfo)
}

func TestHandlerRandom(t *testing.T) {
	tCtx := initEnvironment(t)

	// create pool
	createPoolMsg := createPool(t, tCtx)
	for i := 0; i < 10000; i++ {
		var msg sdk.Msg
		judge := rand.Intn(5)
		switch judge {
		case 0:
			msg = normalGetProvideMsg(tCtx, createPoolMsg)
		case 1:
			msg = normalGetCreatePoolMsg(tCtx, createPoolMsg)
		case 2:
			msg = normalGetLockMsg(tCtx, createPoolMsg)
		case 3:
			msg = normalGetUnlockMsg(tCtx, createPoolMsg)
		case 4:
			msg = normalGetClaimMsg(tCtx, createPoolMsg)
		case 5:
			msg = normalGetDestroyPoolMsg(tCtx, createPoolMsg)
		}
		ctx, writeCache := tCtx.ctx.CacheContext()
		res := tCtx.handler(ctx, msg)
		if !res.Code.IsOK() {
			fmt.Println(res.Log)
		} else {
			writeCache()
		}
		tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + int64(rand.Intn(2)))
	}
}

func TestHandlerCheckCombination(t *testing.T) {
	var preExec preExecFunc = func(t *testing.T, tCtx *testContext) interface{} {
		return normalGetCreatePoolMsg(tCtx, nil)
	}
	tests := []testCaseItem{
		{
			caseName:     "success. create pool",
			preExec:      preExec,
			getMsg:       normalGetCreatePoolMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. provide",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1)
				return normalGetCreatePoolMsg(tCtx, nil)
			},
			getMsg:       normalGetProvideMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. lock address 1",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1)
				return normalGetCreatePoolMsg(tCtx, nil)
			},
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. lock address 2",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1)
				createPoolMsg := normalGetCreatePoolMsg(tCtx, nil).(types.MsgCreatePool)
				createPoolMsg.Owner = tCtx.addrList[0]
				return createPoolMsg
			},
			getMsg:       normalGetLockMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName:     "success. claim address 1",
			preExec:      preExec,
			getMsg:       normalGetClaimMsg,
			verification: verification,
			expectedCode: sdk.CodeOK,
		},
		{
			caseName: "success. unlock address 1",
			preExec: func(t *testing.T, tCtx *testContext) interface{} {
				tCtx.ctx = tCtx.ctx.WithBlockHeight(tCtx.ctx.BlockHeight() + 1)
				return normalGetCreatePoolMsg(tCtx, nil)
			},
			getMsg: normalGetUnlockMsg,
			verification: func(t *testing.T, tCtx *testContext, result sdk.Result, testCase testCaseItem, preCoins, afterCoins sdk.SysCoins, preData interface{}) {
				verification(t, tCtx, result, testCase, preCoins, afterCoins, preData)
				createPoolMsg := preData.(types.MsgCreatePool)

				// check current rewards
				curPeriodRewards := tCtx.k.GetPoolCurrentRewards(tCtx.ctx, createPoolMsg.PoolName)
				var expectedCurrentPeriod uint64 = 6
				require.Equal(t, expectedCurrentPeriod, curPeriodRewards.Period)
				require.Equal(t, tCtx.ctx.BlockHeight(), curPeriodRewards.StartBlockHeight)

				// check the number of historicalRewards
				numHistoricalRewards := 0
				tCtx.k.IteratePoolHistoricalRewards(tCtx.ctx, createPoolMsg.PoolName,
					func(store sdk.KVStore, key []byte, value []byte) (stop bool) {
						var rewards types.PoolHistoricalRewards
						types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(value, &rewards)
						numHistoricalRewards++
						return false
					})
				require.Equal(t, 2, numHistoricalRewards)

				// check the number of lockInfo
				numLockInfo := 0
				tCtx.k.IterateAllLockInfos(tCtx.ctx, func(lockInfo types.LockInfo) (stop bool) {
					numLockInfo++
					return false
				})
				require.Equal(t, int(1), numLockInfo)

				// check historical rewards of period 3
				expectedRatio := sdk.NewDecCoinsFromDec(createPoolMsg.YieldedSymbol, sdk.NewDec(1))
				rewards := tCtx.k.GetPoolHistoricalRewards(tCtx.ctx, createPoolMsg.PoolName, 3)
				require.Equal(t, uint16(1), rewards.ReferenceCount)
				require.Equal(t, expectedRatio.String(), rewards.CumulativeRewardRatio.String())

				// check historical rewards of period 5
				expectedRatio = sdk.NewDecCoinsFromDec(createPoolMsg.YieldedSymbol, sdk.NewDecWithPrec(15, 1))
				rewards = tCtx.k.GetPoolHistoricalRewards(tCtx.ctx, createPoolMsg.PoolName, 5)
				require.Equal(t, uint16(1), rewards.ReferenceCount)
				require.Equal(t, expectedRatio.String(), rewards.CumulativeRewardRatio.String())

				// check farm pool
				pool, found := tCtx.k.GetFarmPool(tCtx.ctx, createPoolMsg.PoolName)
				params := tCtx.k.GetParams(tCtx.ctx)
				require.True(t, found)
				require.Equal(t, params.CreatePoolDeposit.String(), pool.DepositAmount.String())
				expectedLockedValue := sdk.NewDecCoinFromDec(createPoolMsg.MinLockAmount.Denom, sdk.NewDec(1))
				require.Equal(t, expectedLockedValue.String(), pool.TotalValueLocked.String())
				require.Equal(t, 1, len(pool.YieldedTokenInfos))
				expectedRemainingAmount := sdk.NewDecCoinFromDec(createPoolMsg.YieldedSymbol, sdk.NewDec(8))
				require.Equal(t, expectedRemainingAmount.String(), pool.YieldedTokenInfos[0].RemainingAmount.String())
				expectedTotalAccumulatedRewards := sdk.NewDecCoinsFromDec(createPoolMsg.YieldedSymbol, sdk.NewDecWithPrec(5, 1))
				require.Equal(t, expectedTotalAccumulatedRewards.String(), pool.TotalAccumulatedRewards.String())

			},
			expectedCode: sdk.CodeOK,
		},
	}

	testCaseCombinationTest(t, tests)

}
