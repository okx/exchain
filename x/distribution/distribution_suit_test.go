package distribution

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	ctx                    sdk.Context
	ak                     auth.AccountKeeper
	dk                     Keeper
	sk                     staking.Keeper
	supplyKeeper           types.SupplyKeeper
	blockRewardValueTokens = sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(100))}}
	votes                  []abci.VoteInfo
	depositCoin            = sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100))
)

type testAllocationParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        sdk.SysCoins
}

func allocateTokens(t *testing.T) {
	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
}

func initEnv(t *testing.T, validatorCount int64) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, ak, _, dk, sk, _, supplyKeeper = keeper.CreateTestInputAdvanced(t, false, 1000, communityTax)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	dk.SetInitAllocateValidator(ctx, true)

	h := staking.NewHandler(sk)
	valOpAddrs, valConsPks, _ := keeper.GetTestAddrs()

	// create four validators
	for i := int64(0); i < validatorCount; i++ {
		msg := staking.NewMsgCreateValidator(valOpAddrs[i], valConsPks[i],
			staking.Description{}, keeper.NewTestSysCoin(i+1, 0))
		_, e := h(ctx, msg)
		require.Nil(t, e)
		require.True(t, dk.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]).IsZero())
	}
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	testAllocationParam := testAllocationParam{
		10,
		[]bool{true, true, true, true}, []bool{false, false, false, false},
		nil,
	}
	votes = createTestVotes(ctx, sk, testAllocationParam)
}

func createTestVotes(ctx sdk.Context, sk staking.Keeper, test testAllocationParam) []abci.VoteInfo {
	var votes []abci.VoteInfo
	for i := int64(0); i < int64(len(test.isVote)); i++ {
		if test.isJailed[i] {
			sk.Jail(ctx, keeper.TestConsAddrs[i])
		}
		abciVal := abci.Validator{Address: keeper.TestConsAddrs[i], Power: i + 1}
		if test.isVote[i] {
			votes = append(votes, abci.VoteInfo{Validator: abciVal, SignedLastBlock: true})
		}
	}
	return votes
}

type DistributionSuite struct {
	suite.Suite
}

func TestDistributionOnline(t *testing.T) {
	suite.Run(t, new(DistributionSuite))
}

func getDecCoins(value string) sdk.SysCoins {
	if value == "0" {
		var decCoins sdk.SysCoins
		return decCoins
	}

	dec, _ := sdk.NewDecFromStr(value)
	return sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: dec}}
}

func (suite *DistributionSuite) TestNormal() {
	testCases := []struct {
		title                string
		valCount             int64
		beforeCommissionDec  [4]string
		beforeOutstandingDec [4]string
		afterCommissionDec   [4]string
		afterOutstandingDec  [4]string
		decCommunity         string
		chainType            uint32
	}{
		{
			"1 validator onchain",
			int64(1),
			[4]string{"98"},
			[4]string{"98"},
			[4]string{"0"},
			[4]string{"0"},
			"2",
			1,
		},
		{
			"2 validator onchain",
			int64(2),
			[4]string{"49", "49"},
			[4]string{"49", "49"},
			[4]string{"0", "0"},
			[4]string{"0", "0"},
			"2",
			1,
		},
		{
			"3 validator onchain",
			int64(3),
			[4]string{"32.666666666666666633", "32.666666666666666633", "32.666666666666666633"},
			[4]string{"32.666666666666666633", "32.666666666666666633", "32.666666666666666633"},
			[4]string{"0.666666666666666633", "0.666666666666666633", "0.666666666666666633"},
			[4]string{"0.666666666666666633", "0.666666666666666633", "0.666666666666666633"},
			"2.000000000000000101",
			1,
		},
		{
			"4 validator onchain",
			int64(4),
			[4]string{"24.5", "24.5", "24.5", "24.5"},
			[4]string{"24.5", "24.5", "24.5", "24.5"},
			[4]string{"0.5", "0.5", "0.5", "0.5"},
			[4]string{"0.5", "0.5", "0.5", "0.5"},
			"2",
			1,
		},
		{
			"1 validator offchain",
			int64(1),
			[4]string{"98"},
			[4]string{"98"},
			[4]string{"0"},
			[4]string{"0"},
			"2",
			0,
		},
		{
			"2 validator offchain",
			int64(2),
			[4]string{"49", "49"},
			[4]string{"49", "49"},
			[4]string{"0", "0"},
			[4]string{"0", "0"},
			"2",
			0,
		},
		{
			"3 validator offchain",
			int64(3),
			[4]string{"32.666666666666666633", "32.666666666666666633", "32.666666666666666633"},
			[4]string{"32.666666666666666633", "32.666666666666666633", "32.666666666666666633"},
			[4]string{"0.666666666666666633", "0.666666666666666633", "0.666666666666666633"},
			[4]string{"0.666666666666666633", "0.666666666666666633", "0.666666666666666633"},
			"2.000000000000000101",
			0,
		},
		{
			"4 validator offchain",
			int64(4),
			[4]string{"24.5", "24.5", "24.5", "24.5"},
			[4]string{"24.5", "24.5", "24.5", "24.5"},
			[4]string{"0.5", "0.5", "0.5", "0.5"},
			[4]string{"0.5", "0.5", "0.5", "0.5"},
			"2",
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount)
			dk.SetDistributionType(ctx, tc.chainType)
			allocateTokens(suite.T())

			for i := int64(0); i < tc.valCount; i++ {
				require.Equal(suite.T(), getDecCoins(tc.beforeCommissionDec[i]), dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.beforeOutstandingDec[i]), dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.decCommunity), dk.GetFeePoolCommunityCoins(ctx))

				dk.WithdrawValidatorCommission(ctx, keeper.TestValAddrs[i])

				require.Equal(suite.T(), getDecCoins(tc.afterCommissionDec[i]), dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.afterOutstandingDec[i]), dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.decCommunity), dk.GetFeePoolCommunityCoins(ctx))
			}
		})
	}
}

func (suite *DistributionSuite) TestDelegator() {
	testCases := []struct {
		title      string
		valCount   int64
		delCount   int64
		decRewards [4][4]string
		rate       string
		chainType  uint32
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"48"}},
			"0.5",
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
		},
		{
			"1 delegator，2 validator, onchain",
			2,
			1,
			[4][4]string{{"24", "24"}},
			"0.5",
			1,
		},
		{
			"1 delegator，2 validator, offchain",
			2,
			1,
			[4][4]string{{"0", "0"}},
			"0.5",
			0,
		},
		{
			"1 delegator，4 validator, onchain",
			4,
			1,
			[4][4]string{{"12", "12", "12", "12"}},
			"0.5",
			1,
		},
		{
			"1 delegator，4 validator, offchain",
			4,
			1,
			[4][4]string{{"0", "0", "0", "0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，1 validator, onchain",
			1,
			2,
			[4][4]string{{"24"}, {"24"}},
			"0.5",
			1,
		},
		{
			"2 delegator，1 validator, offchain",
			1,
			2,
			[4][4]string{{"0"}, {"0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，2 validator, onchain",
			2,
			2,
			[4][4]string{{"12", "12"}, {"12", "12"}},
			"0.5",
			1,
		},
		{
			"2 delegator，2 validator, offchain",
			2,
			2,
			[4][4]string{{"0", "0"}, {"0", "0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，4 validator, onchain",
			4,
			2,
			[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
			"0.5",
			1,
		},
		{
			"2 delegator，4 validator, offchain",
			4,
			2,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			"0.5",
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount)
			dk.SetDistributionType(ctx, tc.chainType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
			}

			allocateTokens(suite.T())

			beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decRewards[i][j]))
				}
			}

			afterValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			//withdraw again
			staking.EndBlocker(ctx, sk)
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				}
			}
			require.Equal(suite.T(), beforeValCommission, afterValCommission)

			//allocate and withdraw agagin
			allocateTokens(suite.T())
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decRewards[i][j]))
				}
			}
		})
	}
}

func (suite *DistributionSuite) TestProxy() {
	testCases := []struct {
		title        string
		valCount     int64
		proxyCount   int64
		proxyRewards [4][4]string
		rate         string
		chainType    uint32
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"48"}},
			"0.5",
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
		},
		{
			"1 delegator，2 validator, onchain",
			2,
			1,
			[4][4]string{{"24", "24"}},
			"0.5",
			1,
		},
		{
			"1 delegator，2 validator, offchain",
			2,
			1,
			[4][4]string{{"0", "0"}},
			"0.5",
			0,
		},
		{
			"1 delegator，4 validator, onchain",
			4,
			1,
			[4][4]string{{"12", "12", "12", "12"}},
			"0.5",
			1,
		},
		{
			"1 delegator，4 validator, offchain",
			4,
			1,
			[4][4]string{{"0", "0", "0", "0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，1 validator, onchain",
			1,
			2,
			[4][4]string{{"24"}, {"24"}},
			"0.5",
			1,
		},
		{
			"2 delegator，1 validator, offchain",
			1,
			2,
			[4][4]string{{"0"}, {"0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，2 validator, onchain",
			2,
			2,
			[4][4]string{{"12", "12"}, {"12", "12"}},
			"0.5",
			1,
		},
		{
			"2 delegator，2 validator, offchain",
			2,
			2,
			[4][4]string{{"0", "0"}, {"0", "0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，4 validator, onchain",
			4,
			2,
			[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
			"0.5",
			1,
		},
		{
			"2 delegator，4 validator, offchain",
			4,
			2,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			"0.5",
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount)
			dk.SetDistributionType(ctx, tc.chainType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			staking.EndBlocker(ctx, sk)
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			for i := int64(0); i < tc.proxyCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], depositCoin)
				keeper.DoRegProxy(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], true)

				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoBindProxy(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestProxyAddrs[i])
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], keeper.TestValAddrs[0:tc.valCount])
			}

			//test withdraw
			testProxyWithdraw(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

			//proxy withdraw again, delegator withdraw again
			testProxyWithdrawAgain(suite, tc.valCount, tc.proxyCount)

			// UnBindProxy
			for i := int64(0); i < tc.proxyCount; i++ {
				keeper.DoUnBindProxy(suite.T(), ctx, sk, keeper.TestDelAddrs[i])
			}

			allocateTokens(suite.T())

			for i := int64(0); i < tc.proxyCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.proxyRewards[i][j]))
				}
			}

			//proxy withdraw again, delegator withdraw again
			testProxyWithdrawAgain(suite, tc.valCount, tc.proxyCount)

			//bind again, will withdraw to proxy
			testProxyBindAgain(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

			//del deposit will do withdraw
			testProxyDelDepositAgain(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

			//proxy deposit will do withdraw
			testProxyProxyDepositAgain(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

			allocateTokens(suite.T())
		})
	}
}

func testProxyWithdraw(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	//test withdraw
	allocateTokens(suite.T())
	beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < valCount; i++ {
		beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
	}

	for i := int64(0); i < proxyCount; i++ {
		for j := int64(0); j < valCount; j++ {
			beforeAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
			dk.WithdrawDelegationRewards(ctx, keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
			afterAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
			require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(proxyRewards[i][j]))
		}
	}

	afterValCommission := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < valCount; i++ {
		afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
	}
}

func testProxyWithdrawAgain(suite *DistributionSuite, valCount int64, proxyCount int64) {
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	for i := int64(0); i < proxyCount; i++ {
		for j := int64(0); j < valCount; j++ {
			beforeAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
			dk.WithdrawDelegationRewards(ctx, keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
			afterAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
			require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
		}
	}

	//delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	for i := int64(0); i < proxyCount; i++ {
		for j := int64(0); j < valCount; j++ {
			beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
			dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
			afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
			require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
		}
	}
}

func testProxyBindAgain(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	allocateTokens(suite.T())
	beforeProxyAccount := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccount[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
		for j := int64(0); j < valCount; j++ {
			beforeProxyAccount[i] = beforeProxyAccount[i].Add2(getDecCoins(proxyRewards[i][j]))
		}
	}
	for i := int64(0); i < proxyCount; i++ {
		keeper.DoBindProxy(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestProxyAddrs[i])
	}
	afterProxyAccount := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < proxyCount; i++ {
		afterProxyAccount[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
	}
	require.Equal(suite.T(), beforeProxyAccount, afterProxyAccount)
}

func testProxyDelDepositAgain(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	allocateTokens(suite.T())
	beforeProxyAccount := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccount[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
		for j := int64(0); j < valCount; j++ {
			beforeProxyAccount[i] = beforeProxyAccount[i].Add2(getDecCoins(proxyRewards[i][j]))
		}
	}
	for i := int64(0); i < proxyCount; i++ {
		keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
	}
	afterProxyAccount := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < proxyCount; i++ {
		afterProxyAccount[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
	}
	require.Equal(suite.T(), beforeProxyAccount, afterProxyAccount)
}

func testProxyProxyDepositAgain(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	allocateTokens(suite.T())
	beforeProxyAccount := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccount[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
		for j := int64(0); j < valCount; j++ {
			beforeProxyAccount[i] = beforeProxyAccount[i].Add2(getDecCoins(proxyRewards[i][j]))
		}
	}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccount[i] = beforeProxyAccount[i].Sub(sdk.NewCoins(depositCoin))
		keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], depositCoin)
	}
	afterProxyAccount := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < proxyCount; i++ {
		afterProxyAccount[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
	}
	require.Equal(suite.T(), beforeProxyAccount, afterProxyAccount)
}

func (suite *DistributionSuite) TestDestroyValidator() {
	testCases := []struct {
		title                string
		valCount             int64
		beforeCommissionDec  [4]string
		beforeOutstandingDec [4]string
		afterCommissionDec   [4]string
		afterOutstandingDec  [4]string
		decCommunity         string
		chainType            uint32
	}{
		{},
	}
	_ = testCases
}

func (suite *DistributionSuite) TestWithdraw() {
	testCases := []struct {
		title      string
		valCount   int64
		delCount   int64
		decRewards [4][4]string
		rate       string
		chainType  uint32
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"48"}},
			"0.5",
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
		},
		{
			"1 delegator，1 validator, onchain",
			2,
			1,
			[4][4]string{{"24", "24"}},
			"0.5",
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			2,
			1,
			[4][4]string{{"0", "0"}},
			"0.5",
			0,
		},
		{
			"1 delegator，4 validator, onchain",
			4,
			1,
			[4][4]string{{"12", "12", "12", "12"}},
			"0.5",
			1,
		},
		{
			"1 delegator，4 validator, offchain",
			4,
			1,
			[4][4]string{{"0", "0", "0", "0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，1 validator, onchain",
			1,
			2,
			[4][4]string{{"24"}, {"24"}},
			"0.5",
			1,
		},
		{
			"2 delegator，1 validator, offchain",
			1,
			2,
			[4][4]string{{"0"}, {"0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，2 validator, onchain",
			2,
			2,
			[4][4]string{{"12", "12"}, {"12", "12"}},
			"0.5",
			1,
		},
		{
			"2 delegator，2 validator, offchain",
			2,
			2,
			[4][4]string{{"0", "0"}, {"0", "0"}},
			"0.5",
			0,
		},
		{
			"2 delegator，4 validator, onchain",
			4,
			2,
			[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
			"0.5",
			1,
		},
		{
			"2 delegator，4 validator, offchain",
			4,
			2,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			"0.5",
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount)
			dk.SetDistributionType(ctx, tc.chainType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
			}

			allocateTokens(suite.T())

			beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				rewards := sdk.SysCoins{}
				for j := int64(0); j < tc.valCount; j++ {
					rewards = rewards.Add2(getDecCoins(tc.decRewards[i][j]))

				}
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), rewards)
			}

			afterValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}
			require.Equal(suite.T(), beforeValCommission, afterValCommission)

			//withdraw again
			staking.EndBlocker(ctx, sk)
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				}
			}

			//allocate and withdraw again, do nothing
			allocateTokens(suite.T())
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				}
			}
		})
	}
}

func (suite *DistributionSuite) TestWithdrawValidator() {
	testCases := []struct {
		title      string
		valCount   int64
		delCount   int64
		decRewards [4][4]string
		rate       string
		chainType  uint32
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"24"}},
			"0.5",
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount)
			dk.SetDistributionType(ctx, tc.chainType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
			}

			for j := int64(0); j < tc.valCount; j++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestValAccAddrs[j], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestValAccAddrs[j], keeper.TestValAddrs[0:tc.valCount])
			}

			allocateTokens(suite.T())

			beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			//withdraw validator
			for j := int64(0); j < tc.valCount; j++ {
				keeper.DoDestroyValidator(suite.T(), ctx, sk, keeper.TestValAccAddrs[j])
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestValAccAddrs[j], depositCoin)
			}

			staking.EndBlocker(ctx, sk)

			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				rewards := sdk.SysCoins{}
				for j := int64(0); j < tc.valCount; j++ {
					rewards = rewards.Add2(getDecCoins(tc.decRewards[i][j]))

				}
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), rewards)
			}

			afterValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}
			require.Equal(suite.T(), beforeValCommission, afterValCommission)

			//withdraw again
			staking.EndBlocker(ctx, sk)
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				}
			}

			//allocate and withdraw again, do nothing
			staking.EndBlocker(ctx, sk)
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				}
			}
		})
	}
}

func setTestFees(t *testing.T, ctx sdk.Context, k Keeper, ak auth.AccountKeeper, fees sdk.SysCoins) {
	feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	require.NotNil(t, feeCollector)
	err := feeCollector.SetCoins(fees)
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
}
