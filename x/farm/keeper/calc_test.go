package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
)

func TestCalculateAmountYieldedBetween(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"
	type testCase struct {
		curRewards     types.PoolCurrentRewards
		endBlockHeight int64
		yieldedInfos   types.YieldedTokenInfos
		expectedFunc   func(testCase, func() (types.FarmPool, sdk.SysCoins))
	}

	expectSuccess := func(test testCase, testFunc func() (types.FarmPool, sdk.SysCoins)) {
		pool, yieldedTokens := testFunc()
		oldRemaining := test.yieldedInfos[0].RemainingAmount
		newRemaining := pool.YieldedTokenInfos[0].RemainingAmount
		require.Equal(t, yieldedTokens, sdk.SysCoins{oldRemaining.Sub(newRemaining)})
	}

	expectNotYield := func(test testCase, testFunc func() (types.FarmPool, sdk.SysCoins)) {
		pool, yieldedTokens := testFunc()
		require.True(t, yieldedTokens.IsZero())
		require.Equal(t, test.yieldedInfos, pool.YieldedTokenInfos)
	}

	tests := []testCase{
		{
			curRewards:     types.NewPoolCurrentRewards(100, 1, sdk.SysCoins{}),
			endBlockHeight: 120,
			yieldedInfos: types.YieldedTokenInfos{
				types.NewYieldedTokenInfo(
					sdk.NewDecCoin("xxb", sdk.ZeroInt()),
					0,
					sdk.NewDec(0),
				),
			},
			expectedFunc: expectNotYield,
		},
		{
			curRewards:     types.NewPoolCurrentRewards(60, 1, sdk.SysCoins{}),
			endBlockHeight: 100,
			yieldedInfos: types.YieldedTokenInfos{
				types.NewYieldedTokenInfo(
					sdk.NewDecCoin("xxb", sdk.NewInt(1000)),
					100,
					sdk.NewDec(10),
				),
			},
			expectedFunc: expectNotYield,
		},
		{
			curRewards:     types.NewPoolCurrentRewards(120, 1, sdk.SysCoins{}),
			endBlockHeight: 100,
			yieldedInfos: types.YieldedTokenInfos{
				types.NewYieldedTokenInfo(
					sdk.NewDecCoin("xxb", sdk.NewInt(1000)),
					70,
					sdk.NewDec(10),
				),
			},
			expectedFunc: expectNotYield,
		},
		{
			curRewards:     types.NewPoolCurrentRewards(60, 1, sdk.SysCoins{}),
			endBlockHeight: 100,
			yieldedInfos: types.YieldedTokenInfos{
				types.NewYieldedTokenInfo(
					sdk.NewDecCoin("xxb", sdk.NewInt(1000)),
					70,
					sdk.NewDec(10),
				),
			},
			expectedFunc: expectSuccess,
		},
		{
			curRewards:     types.NewPoolCurrentRewards(70, 1, sdk.SysCoins{}),
			endBlockHeight: 100,
			yieldedInfos: types.YieldedTokenInfos{
				types.NewYieldedTokenInfo(
					sdk.NewDecCoin("xxb", sdk.NewInt(100)),
					60,
					sdk.NewDec(10),
				),
			},
			expectedFunc: expectSuccess,
		},
	}

	for _, test := range tests {
		ctx = ctx.WithBlockHeight(test.endBlockHeight)
		keeper.SetPoolCurrentRewards(ctx, poolName, test.curRewards)
		pool := types.FarmPool{
			Name:              poolName,
			YieldedTokenInfos: types.YieldedTokenInfos{test.yieldedInfos[0]},
		}
		wrappedTestFunc := func() (types.FarmPool, sdk.SysCoins) {
			return keeper.CalculateAmountYieldedBetween(ctx, pool)
		}
		test.expectedFunc(test, wrappedTestFunc)
	}
}

func TestIncrementReferenceCount(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"

	expectSuccess := func(testFunc func()) {
		testFunc()
	}

	expectPanic := func(testFunc func()) {
		require.Panics(t, testFunc)
	}

	tests := []struct {
		period         uint64
		referenceCount uint16
		expectedFunc   func(func())
	}{
		{0, 0, expectSuccess},
		{1, 1, expectSuccess},
		{2, 2, expectSuccess},
		{0, 3, expectPanic},
	}

	for _, test := range tests {
		his := types.NewPoolHistoricalRewards(sdk.SysCoins{}, test.referenceCount)
		keeper.SetPoolHistoricalRewards(ctx, poolName, test.period, his)
		wrappedTestFunc := func() {
			keeper.incrementReferenceCount(ctx, poolName, test.period)
			newHis := keeper.GetPoolHistoricalRewards(ctx, poolName, test.period)
			require.Equal(t, test.referenceCount+1, newHis.ReferenceCount)
		}
		test.expectedFunc(wrappedTestFunc)
	}
}

func TestDecrementReferenceCount(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"

	type testCase struct {
		period         uint64
		referenceCount uint16
		expectedFunc   func(testCase, func())
	}

	expectDeleted := func(test testCase, testFunc func()) {
		testFunc()
		require.Panics(t, func() {
			keeper.GetPoolHistoricalRewards(ctx, poolName, test.period)
		})
	}

	expectSuccess := func(test testCase, testFunc func()) {
		testFunc()
		newHis := keeper.GetPoolHistoricalRewards(ctx, poolName, test.period)
		require.Equal(t, test.referenceCount-1, newHis.ReferenceCount)
	}

	expectPanic := func(test testCase, testFunc func()) {
		require.Panics(t, testFunc)
	}

	tests := []testCase{
		{0, 0, expectPanic},
		{1, 1, expectDeleted},
		{2, 2, expectSuccess},
		{0, 3, expectSuccess},
	}

	for _, test := range tests {
		his := types.NewPoolHistoricalRewards(sdk.SysCoins{}, test.referenceCount)
		keeper.SetPoolHistoricalRewards(ctx, poolName, test.period, his)
		wrappedTestFunc := func() {
			keeper.decrementReferenceCount(ctx, poolName, test.period)
		}
		test.expectedFunc(test, wrappedTestFunc)
	}
}

func TestCalculateLockRewardsBetween(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"

	type testCase struct {
		startPeriod  uint64
		startRatio   sdk.SysCoins
		endPeriod    uint64
		endRatio     sdk.SysCoins
		amount       sdk.SysCoin
		expectedFunc func(testCase, func() sdk.SysCoins)
	}

	expectSuccess := func(test testCase, testFunc func() sdk.SysCoins) {
		rewards := testFunc()
		require.Equal(t, rewards, test.endRatio.Sub(test.startRatio).MulDecTruncate(test.amount.Amount))
	}

	expectPanic := func(test testCase, testFunc func() sdk.SysCoins) {
		require.Panics(t, func() {
			testFunc()
		})
	}

	tests := []testCase{
		{
			startPeriod: 0,
			startRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(10)),
			},
			endPeriod: 1,
			endRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(100)),
			},
			amount:       sdk.NewDecCoin("xxb", sdk.NewInt(10)),
			expectedFunc: expectSuccess,
		},
		{
			startPeriod: 0,
			startRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(10)),
			},
			endPeriod: 1,
			endRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(7)),
			},
			amount:       sdk.NewDecCoin("xxb", sdk.NewInt(10)),
			expectedFunc: expectPanic,
		},
		{
			startPeriod: 0,
			startRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(10)),
			},
			endPeriod: 1,
			endRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(100)),
			},
			amount:       sdk.SysCoin{"xxb", sdk.NewDec(-1)},
			expectedFunc: expectPanic,
		},
		{
			startPeriod: 1,
			startRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(10)),
			},
			endPeriod: 0,
			endRatio: sdk.SysCoins{
				sdk.NewDecCoin("wwb", sdk.NewInt(10)),
				sdk.NewDecCoin("okt", sdk.NewInt(100)),
			},
			amount:       sdk.NewDecCoin("xxb", sdk.NewInt(10)),
			expectedFunc: expectPanic,
		},
	}

	for _, test := range tests {
		startHis := types.NewPoolHistoricalRewards(test.startRatio, 1)
		keeper.SetPoolHistoricalRewards(ctx, poolName, test.startPeriod, startHis)
		endHis := types.NewPoolHistoricalRewards(test.endRatio, 1)
		keeper.SetPoolHistoricalRewards(ctx, poolName, test.endPeriod, endHis)

		wrappedTestFunc := func() sdk.SysCoins {
			return keeper.calculateLockRewardsBetween(ctx, poolName, test.startPeriod, test.endPeriod, test.amount)
		}
		test.expectedFunc(test, wrappedTestFunc)
	}
}

func TestIncrementPoolPeriod(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"

	type testCase struct {
		curRewards    types.PoolCurrentRewards
		preHisRewards types.PoolHistoricalRewards
		valueLocked   sdk.SysCoin
		yieldedTokens sdk.SysCoins
		expectedFunc  func(testCase, func() uint64)
	}

	expectSuccess := func(test testCase, testFunc func() uint64) {
		period := testFunc()
		// previous historical rewards is deleted
		require.Panics(t, func() {
			keeper.GetPoolHistoricalRewards(ctx, poolName, test.curRewards.Period-1)
		})
		// return the current period ended just now
		require.Equal(t, test.curRewards.Period, period)

		var currentRatio sdk.SysCoins
		if test.valueLocked.IsZero() {
			currentRatio = sdk.SysCoins{}
		} else {
			currentRatio = test.curRewards.Rewards.Add2(test.yieldedTokens).QuoDecTruncate(test.valueLocked.Amount)
		}

		// create new current period rewards
		newHis := keeper.GetPoolHistoricalRewards(ctx, poolName, test.curRewards.Period)
		require.Equal(t, test.preHisRewards.CumulativeRewardRatio.Add2(currentRatio),
			newHis.CumulativeRewardRatio,
		)

		newCurr := keeper.GetPoolCurrentRewards(ctx, poolName)
		require.Equal(
			t,
			test.curRewards.Period+1,
			newCurr.Period,
		)
	}

	tests := []testCase{
		{
			curRewards: types.NewPoolCurrentRewards(100, 1, sdk.SysCoins{}),
			preHisRewards: types.NewPoolHistoricalRewards(
				sdk.SysCoins{sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10))}, 1,
			),
			valueLocked:   sdk.NewDecCoinFromDec("wwb", sdk.NewDec(100)),
			yieldedTokens: sdk.SysCoins{sdk.NewDecCoin("yyb", sdk.NewInt(100))},
			expectedFunc:  expectSuccess,
		},
		{
			curRewards: types.NewPoolCurrentRewards(100, 1, sdk.SysCoins{}),
			preHisRewards: types.NewPoolHistoricalRewards(
				sdk.SysCoins{sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10))}, 1,
			),
			valueLocked:   sdk.NewDecCoinFromDec("wwb", sdk.NewDec(0)),
			yieldedTokens: sdk.SysCoins{sdk.NewDecCoin("yyb", sdk.NewInt(100))},
			expectedFunc:  expectSuccess,
		},
	}

	for _, test := range tests {
		keeper.SetPoolHistoricalRewards(ctx, poolName, test.curRewards.Period-1, test.preHisRewards)
		keeper.SetPoolCurrentRewards(ctx, poolName, test.curRewards)

		wrappedTestFunc := func() uint64 {
			return keeper.IncrementPoolPeriod(ctx, poolName, test.valueLocked, test.yieldedTokens)
		}
		test.expectedFunc(test, wrappedTestFunc)
	}
}

func TestUpdateLockInfo(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"

	type testCase struct {
		lockInfo      types.LockInfo
		changeAmount  sdk.Dec
		isSetLockInfo bool
		expectedFunc  func(testCase, func())
	}

	expectPanic := func(test testCase, testFunc func()) {
		require.Panics(t, func() {
			testFunc()
		})
	}

	expectDelete := func(test testCase, testFunc func()) {
		testFunc()
		found := keeper.HasLockInfo(ctx, test.lockInfo.Owner, poolName)
		require.False(t, found)
		found = keeper.HasAddressInFarmPool(ctx, poolName, test.lockInfo.Owner)
		require.False(t, found)
	}

	expectUpdate := func(test testCase, testFunc func()) {
		testFunc()
		lockInfo, found := keeper.GetLockInfo(ctx, test.lockInfo.Owner, poolName)
		require.True(t, found)
		require.Equal(t, lockInfo.Amount.Amount, test.lockInfo.Amount.Amount.Add(test.changeAmount))
	}

	tests := []testCase{
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				100, 1,
			),
			changeAmount:  sdk.NewDec(0),
			isSetLockInfo: false,
			expectedFunc:  expectPanic,
		},
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				100, 1,
			),
			changeAmount:  sdk.NewDec(-10),
			isSetLockInfo: true,
			expectedFunc:  expectDelete,
		},
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				100, 1,
			),
			changeAmount:  sdk.NewDec(5),
			isSetLockInfo: true,
			expectedFunc:  expectUpdate,
		},
	}

	for _, test := range tests {
		keeper.SetPoolHistoricalRewards(
			ctx, poolName, 1,
			types.NewPoolHistoricalRewards(sdk.SysCoins{}, 1),
		)
		keeper.SetPoolCurrentRewards(
			ctx, poolName, types.NewPoolCurrentRewards(ctx.BlockHeight(), 2, sdk.SysCoins{}),
		)
		if test.isSetLockInfo {
			keeper.SetLockInfo(ctx, test.lockInfo)
		}
		wrappedTestFunc := func() {
			keeper.UpdateLockInfo(ctx, test.lockInfo.Owner, poolName, test.changeAmount)
		}
		test.expectedFunc(test, wrappedTestFunc)
	}
}

func TestWithdrawRewards(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "poolName"

	type testCase struct {
		lockInfo       types.LockInfo
		totalLocked    sdk.SysCoin
		yieldedAmount  sdk.SysCoins
		isSetLockInfo  bool
		isSetModuleAcc bool
		expectedFunc   func(testCase, func() (sdk.SysCoins, sdk.Error))
	}

	expectError := func(test testCase, testFunc func() (sdk.SysCoins, sdk.Error)) {
		rewards, err := testFunc()
		require.Equal(t, sdk.SysCoins(nil), rewards)
		require.NotNil(t, err)
	}

	expectSuccess := func(test testCase, testFunc func() (sdk.SysCoins, sdk.Error)) {
		rewards, err := testFunc()
		require.Equal(
			t,
			test.yieldedAmount.QuoDecTruncate(test.totalLocked.Amount).MulDecTruncate(test.lockInfo.Amount.Amount),
			rewards,
		)
		require.Nil(t, err)
		require.Panics(t, func() {
			keeper.GetPoolHistoricalRewards(ctx, poolName, test.lockInfo.ReferencePeriod)
		})
	}

	expectNoRewards := func(test testCase, testFunc func() (sdk.SysCoins, sdk.Error)) {
		rewards, err := testFunc()
		require.Equal(t, sdk.SysCoins(nil), rewards)
		require.Nil(t, err)
		require.Panics(t, func() {
			keeper.GetPoolHistoricalRewards(ctx, poolName, test.lockInfo.ReferencePeriod)
		})
	}

	tests := []testCase{
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				100, 1,
			),
			totalLocked:    sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			yieldedAmount:  sdk.NewDecCoinsFromDec("wwb", sdk.NewDec(100)),
			isSetLockInfo:  false,
			isSetModuleAcc: false,
			expectedFunc:   expectError,
		},
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				100, 1,
			),
			totalLocked:    sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			yieldedAmount:  sdk.NewDecCoinsFromDec("wwb", sdk.NewDec(100)),
			isSetLockInfo:  true,
			isSetModuleAcc: false,
			expectedFunc:   expectError,
		},
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				100, 1,
			),
			totalLocked:    sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			yieldedAmount:  sdk.NewDecCoinsFromDec("wwb", sdk.NewDec(100)),
			isSetLockInfo:  true,
			isSetModuleAcc: true,
			expectedFunc:   expectSuccess,
		},
		{
			lockInfo: types.NewLockInfo(
				Addrs[0], poolName, sdk.NewDecCoinFromDec("xxb", sdk.NewDec(10)),
				120, 1,
			),
			totalLocked:    sdk.NewDecCoinFromDec("xxb", sdk.NewDec(100)),
			yieldedAmount:  sdk.NewDecCoinsFromDec("wwb", sdk.NewDec(100)),
			isSetLockInfo:  true,
			isSetModuleAcc: true,
			expectedFunc:   expectNoRewards,
		},
	}

	for _, test := range tests {
		ctx = ctx.WithBlockHeight(120)
		keeper.SetPoolHistoricalRewards(
			ctx, poolName, 1,
			types.NewPoolHistoricalRewards(sdk.SysCoins{}, 2),
		)
		keeper.SetPoolCurrentRewards(
			ctx, poolName, types.NewPoolCurrentRewards(ctx.BlockHeight(), 2, sdk.SysCoins{}),
		)
		if test.isSetLockInfo {
			keeper.SetLockInfo(ctx, test.lockInfo)
		}
		if test.isSetModuleAcc {
			yieldModuleAcc := keeper.supplyKeeper.GetModuleAccount(ctx, types.YieldFarmingAccount)
			err := yieldModuleAcc.SetCoins(
				test.yieldedAmount.MulDecTruncate(test.totalLocked.Amount).MulDecTruncate(test.lockInfo.Amount.Amount),
			)
			require.Nil(t, err)
			keeper.supplyKeeper.SetModuleAccount(ctx, yieldModuleAcc)
		}
		wrappedTestFunc := func() (sdk.SysCoins, sdk.Error) {
			return keeper.WithdrawRewards(ctx, poolName, test.totalLocked, test.yieldedAmount, test.lockInfo.Owner)
		}
		test.expectedFunc(test, wrappedTestFunc)
	}
}
