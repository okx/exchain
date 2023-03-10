package distribution

import (
	"testing"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/distribution/keeper"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	feePoolBefore, _ := dk.GetFeePool(ctx).CommunityPool.TruncateDecimal()
	setTestFees(t, ctx, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	feePoolAfter, _ := dk.GetFeePool(ctx).CommunityPool.TruncateDecimal()
	require.Equal(t, feePoolBefore.Add2(getDecCoins("2")), feePoolAfter)
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
}

func initEnv(t *testing.T, validatorCount int64, newVersion bool) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, ak, _, dk, sk, _, supplyKeeper = keeper.CreateTestInputAdvanced(t, false, 1000, communityTax)
	if newVersion {
		dk.SetInitExistedValidatorFlag(ctx, true)
	}

	h := staking.NewHandler(sk)
	valOpAddrs, valConsPks, _ := keeper.GetTestAddrs()

	// create four validators
	for i := int64(0); i < validatorCount; i++ {
		msg := staking.NewMsgCreateValidator(valOpAddrs[i], valConsPks[i],
			staking.Description{}, keeper.NewTestSysCoin(i+1, 0))
		_, e := h(ctx, msg)
		require.Nil(t, e)
		if newVersion {
			require.True(t, dk.GetValidatorAccumulatedCommission(ctx, valOpAddrs[i]).IsZero())
		} else {
			require.Panics(t, func() {
				dk.GetValidatorOutstandingRewards(ctx, valOpAddrs[i])
			})
		}
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

func TestDistributionSuit(t *testing.T) {
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
		distrType            uint32
		remainReferenceCount uint64
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
			2,
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
			3,
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
			4,
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
			1,
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
			2,
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
			3,
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
			4,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount, true)
			dk.SetDistributionType(ctx, tc.distrType)
			allocateTokens(suite.T())

			for i := int64(0); i < tc.valCount; i++ {
				require.Equal(suite.T(), getDecCoins(tc.beforeCommissionDec[i]), dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.beforeOutstandingDec[i]), dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.decCommunity), dk.GetFeePoolCommunityCoins(ctx))

				dk.WithdrawValidatorCommission(ctx, keeper.TestValAddrs[i])

				require.Equal(suite.T(), getDecCoins(tc.afterCommissionDec[i]), dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.afterOutstandingDec[i]), dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]))
				require.Equal(suite.T(), getDecCoins(tc.decCommunity), dk.GetFeePoolCommunityCoins(ctx))

				truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
				truncatedCommission, _ := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
				require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
			}

			require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
		})
	}
}

func (suite *DistributionSuite) TestDelegator() {
	testCases := []struct {
		title                string
		valCount             int64
		delCount             int64
		decRewards           [4][4]string
		rate                 string
		distrType            uint32
		remainReferenceCount uint64
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"48"}},
			"0.5",
			1,
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
			1,
		},
		{
			"1 delegator，2 validator, onchain",
			2,
			1,
			[4][4]string{{"24", "24"}},
			"0.5",
			1,
			2,
		},
		{
			"1 delegator，2 validator, offchain",
			2,
			1,
			[4][4]string{{"0", "0"}},
			"0.5",
			0,
			2,
		},
		{
			"1 delegator，4 validator, onchain",
			4,
			1,
			[4][4]string{{"12", "12", "12", "12"}},
			"0.5",
			1,
			4,
		},
		{
			"1 delegator，4 validator, offchain",
			4,
			1,
			[4][4]string{{"0", "0", "0", "0"}},
			"0.5",
			0,
			4,
		},
		{
			"2 delegator，1 validator, onchain",
			1,
			2,
			[4][4]string{{"24"}, {"24"}},
			"0.5",
			1,
			1,
		},
		{
			"2 delegator，1 validator, offchain",
			1,
			2,
			[4][4]string{{"0"}, {"0"}},
			"0.5",
			0,
			1,
		},
		{
			"2 delegator，2 validator, onchain",
			2,
			2,
			[4][4]string{{"12", "12"}, {"12", "12"}},
			"0.5",
			1,
			2,
		},
		{
			"2 delegator，2 validator, offchain",
			2,
			2,
			[4][4]string{{"0", "0"}, {"0", "0"}},
			"0.5",
			0,
			2,
		},
		{
			"2 delegator，4 validator, onchain",
			4,
			2,
			[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
			"0.5",
			1,
			4,
		},
		{
			"2 delegator，4 validator, offchain",
			4,
			2,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			"0.5",
			0,
			4,
		},
		{
			"1 delegator，4 validator, onchain, rate 0",
			4,
			1,
			[4][4]string{{"24", "24", "24", "24"}},
			"0",
			1,
			4,
		},
		{
			"1 delegator，4 validator, onchain, rate 1",
			4,
			1,
			[4][4]string{{"0", "0", "0", "0"}},
			"1",
			1,
			4,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount, true)
			dk.SetDistributionType(ctx, tc.distrType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
				delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
			}

			allocateTokens(suite.T())

			beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					queryRewards := keeper.GetQueriedDelegationRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					queryRewards, _ = queryRewards.TruncateWithPrec(int64(0))
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					require.Nil(suite.T(), err)
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decRewards[i][j]))
					require.Equal(suite.T(), queryRewards, getDecCoins(tc.decRewards[i][j]))
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

					truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[j]).TruncateDecimal()
					truncatedCommission, _ := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[j]).TruncateDecimal()
					require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
				}
			}
			require.Equal(suite.T(), beforeValCommission, afterValCommission)

			//allocate and withdraw agagin
			allocateTokens(suite.T())
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					queryRewards := keeper.GetQueriedDelegationRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					queryRewards, _ = queryRewards.TruncateWithPrec(int64(0))
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decRewards[i][j]))
					require.Equal(suite.T(), queryRewards, getDecCoins(tc.decRewards[i][j]))
				}
			}

			//withdraw token
			allocateTokens(suite.T())
			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				rewards := sdk.SysCoins{}
				for j := int64(0); j < tc.valCount; j++ {
					rewards = rewards.Add2(getDecCoins(tc.decRewards[i][j]))
				}
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), rewards)
				for j := int64(0); j < tc.valCount; j++ {
					require.False(suite.T(), dk.HasDelegatorStartingInfo(ctx, keeper.TestValAddrs[j], keeper.TestDelAddrs[i]))
				}
			}

			//withdraw again
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)
				}
			}

			require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
		})
	}
}

func (suite *DistributionSuite) TestProxy() {
	/*
		testCases := []struct {
			title                string
			valCount             int64
			proxyCount           int64
			proxyRewards         [4][4]string
			rate                 string
			distrType            uint32
			remainReferenceCount uint64
		}{
			{
				"1 proxy，1 validator, onchain",
				1,
				1,
				[4][4]string{{"48"}},
				"0.5",
				1,
				1,
			},
			{
				"1 proxy，1 validator, offchain",
				1,
				1,
				[4][4]string{{"0"}},
				"0.5",
				0,
				1,
			},
			{
				"1 proxy，2 validator, onchain",
				2,
				1,
				[4][4]string{{"24", "24"}},
				"0.5",
				1,
				2,
			},
			{
				"1 proxy，2 validator, offchain",
				2,
				1,
				[4][4]string{{"0", "0"}},
				"0.5",
				0,
				2,
			},
			{
				"1 proxy，4 validator, onchain",
				4,
				1,
				[4][4]string{{"12", "12", "12", "12"}},
				"0.5",
				1,
				4,
			},
			{
				"1 proxy，4 validator, offchain",
				4,
				1,
				[4][4]string{{"0", "0", "0", "0"}},
				"0.5",
				0,
				4,
			},
			{
				"2 proxy，1 validator, onchain",
				1,
				2,
				[4][4]string{{"24"}, {"24"}},
				"0.5",
				1,
				1,
			},
			{
				"2 proxy，1 validator, offchain",
				1,
				2,
				[4][4]string{{"0"}, {"0"}},
				"0.5",
				0,
				1,
			},
			{
				"2 proxy，2 validator, onchain",
				2,
				2,
				[4][4]string{{"12", "12"}, {"12", "12"}},
				"0.5",
				1,
				2,
			},
			{
				"2 proxy，2 validator, offchain",
				2,
				2,
				[4][4]string{{"0", "0"}, {"0", "0"}},
				"0.5",
				0,
				2,
			},
			{
				"2 proxy，4 validator, onchain",
				4,
				2,
				[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
				"0.5",
				1,
				4,
			},
			{
				"2 proxy，4 validator, offchain",
				4,
				2,
				[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
				"0.5",
				0,
				4,
			},
		}

		for _, tc := range testCases {
			suite.Run(tc.title, func() {
				initEnv(suite.T(), tc.valCount, true)
				dk.SetDistributionType(ctx, tc.distrType)

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
					delegator := sk.Delegator(ctx, keeper.TestProxyAddrs[i])
					require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
				}

				//test withdraw rewards
				testProxyWithdrawRewards(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

				//proxy withdraw rewards again, delegator withdraw reards again
				testProxyWithdrawRewardsAgain(suite, tc.valCount, tc.proxyCount)

				// UnBindProxy
				for i := int64(0); i < tc.proxyCount; i++ {
					keeper.DoUnBindProxy(suite.T(), ctx, sk, keeper.TestDelAddrs[i])
				}

				allocateTokens(suite.T())

				for i := int64(0); i < tc.proxyCount; i++ {
					for j := int64(0); j < tc.valCount; j++ {
						queryRewards := keeper.GetQueriedDelegationRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
						queryRewards, _ = queryRewards.TruncateWithPrec(int64(0))
						beforeAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
						_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
						require.Nil(suite.T(), err)
						afterAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
						require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.proxyRewards[i][j]))
						require.Equal(suite.T(), queryRewards, getDecCoins(tc.proxyRewards[i][j]))
					}
				}

				//proxy withdraw rewards again, delegator withdraw rewards again
				testProxyWithdrawRewardsAgain(suite, tc.valCount, tc.proxyCount)

				//bind proxy again
				testProxyBindAgain(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

				//delegator deposit to proxy
				testProxyDelDepositAgain(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

				//proxy deposit again
				testProxyProxyDepositAgain(suite, tc.valCount, tc.proxyCount, tc.proxyRewards)

				//withdraw token
				allocateTokens(suite.T())
				for i := int64(0); i < tc.proxyCount; i++ {
					beforeAccountCoins := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
					keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin.Add(depositCoin))
					keeper.DoRegProxy(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], false)
					keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], depositCoin.Add(depositCoin))
					afterAccountCoins := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
					rewards := sdk.SysCoins{}
					for j := int64(0); j < tc.valCount; j++ {
						rewards = rewards.Add2(getDecCoins(tc.proxyRewards[i][j]))
					}
					require.Equal(suite.T(), afterAccountCoins.Sub(beforeAccountCoins), rewards)
					for j := int64(0); j < tc.valCount; j++ {
						require.False(suite.T(), dk.HasDelegatorStartingInfo(ctx, keeper.TestValAddrs[j], keeper.TestProxyAddrs[i]))
					}
				}

				//proxy withdraw rewards again
				for i := int64(0); i < tc.proxyCount; i++ {
					for j := int64(0); j < tc.valCount; j++ {
						_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
						require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)
					}
				}

				//delegator withdraw rewards again
				for i := int64(0); i < tc.proxyCount; i++ {
					for j := int64(0); j < tc.valCount; j++ {
						_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
						require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)
					}
				}

				require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
			})
		}
	*/
}

func testProxyWithdrawRewards(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	allocateTokens(suite.T())
	beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < valCount; i++ {
		beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
	}

	for i := int64(0); i < proxyCount; i++ {
		for j := int64(0); j < valCount; j++ {
			queryRewards := keeper.GetQueriedDelegationRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
			queryRewards, _ = queryRewards.TruncateWithPrec(int64(0))
			beforeAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
			dk.WithdrawDelegationRewards(ctx, keeper.TestProxyAddrs[i], keeper.TestValAddrs[j])
			afterAccount := ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
			require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(proxyRewards[i][j]))
			require.Equal(suite.T(), queryRewards, getDecCoins(proxyRewards[i][j]))
		}
	}

	afterValCommission := [4]types.ValidatorAccumulatedCommission{}
	for i := int64(0); i < valCount; i++ {
		afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
		truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
		truncatedCommission, _ := afterValCommission[i].TruncateDecimal()
		require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
	}
	require.Equal(suite.T(), beforeValCommission, afterValCommission)
}

func testProxyWithdrawRewardsAgain(suite *DistributionSuite, valCount int64, proxyCount int64) {
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
	beforeProxyAccountCoins := [4]sdk.SysCoins{}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccountCoins[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
		for j := int64(0); j < valCount; j++ {
			beforeProxyAccountCoins[i] = beforeProxyAccountCoins[i].Add2(getDecCoins(proxyRewards[i][j]))
		}
	}
	for i := int64(0); i < proxyCount; i++ {
		keeper.DoBindProxy(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestProxyAddrs[i])
	}
	afterProxyAccountCoins := [4]sdk.SysCoins{}
	for i := int64(0); i < proxyCount; i++ {
		afterProxyAccountCoins[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
	}
	require.Equal(suite.T(), beforeProxyAccountCoins, afterProxyAccountCoins)
}

func testProxyDelDepositAgain(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	allocateTokens(suite.T())
	beforeProxyAccountCoins := [4]sdk.SysCoins{}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccountCoins[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
		for j := int64(0); j < valCount; j++ {
			beforeProxyAccountCoins[i] = beforeProxyAccountCoins[i].Add2(getDecCoins(proxyRewards[i][j]))
		}
	}
	for i := int64(0); i < proxyCount; i++ {
		keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
	}
	afterProxyAccountCoins := [4]sdk.SysCoins{}
	for i := int64(0); i < proxyCount; i++ {
		afterProxyAccountCoins[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
	}
	require.Equal(suite.T(), beforeProxyAccountCoins, afterProxyAccountCoins)
}

func testProxyProxyDepositAgain(suite *DistributionSuite, valCount int64, proxyCount int64, proxyRewards [4][4]string) {
	allocateTokens(suite.T())
	beforeProxyAccountCoins := [4]sdk.SysCoins{}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccountCoins[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
		for j := int64(0); j < valCount; j++ {
			beforeProxyAccountCoins[i] = beforeProxyAccountCoins[i].Add2(getDecCoins(proxyRewards[i][j]))
		}
	}
	for i := int64(0); i < proxyCount; i++ {
		beforeProxyAccountCoins[i] = beforeProxyAccountCoins[i].Sub(sdk.NewCoins(depositCoin))
		keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestProxyAddrs[i], depositCoin)
	}
	afterProxyAccountCoins := [4]sdk.SysCoins{}
	for i := int64(0); i < proxyCount; i++ {
		afterProxyAccountCoins[i] = ak.GetAccount(ctx, keeper.TestProxyAddrs[i]).GetCoins()
	}
	require.Equal(suite.T(), beforeProxyAccountCoins, afterProxyAccountCoins)
}

func (suite *DistributionSuite) TestWithdraw() {
	testCases := []struct {
		title                string
		valCount             int64
		delCount             int64
		decRewards           [4][4]string
		rate                 string
		distrType            uint32
		remainReferenceCount uint64
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"48"}},
			"0.5",
			1,
			1,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
			1,
		},
		{
			"1 delegator，2 validator, onchain",
			2,
			1,
			[4][4]string{{"24", "24"}},
			"0.5",
			1,
			2,
		},
		{
			"1 delegator，2 validator, offchain",
			2,
			1,
			[4][4]string{{"0", "0"}},
			"0.5",
			0,
			2,
		},
		{
			"1 delegator，4 validator, onchain",
			4,
			1,
			[4][4]string{{"12", "12", "12", "12"}},
			"0.5",
			1,
			4,
		},
		{
			"1 delegator，4 validator, offchain",
			4,
			1,
			[4][4]string{{"0", "0", "0", "0"}},
			"0.5",
			0,
			4,
		},
		{
			"2 delegator，1 validator, onchain",
			1,
			2,
			[4][4]string{{"24"}, {"24"}},
			"0.5",
			1,
			1,
		},
		{
			"2 delegator，1 validator, offchain",
			1,
			2,
			[4][4]string{{"0"}, {"0"}},
			"0.5",
			0,
			1,
		},
		{
			"2 delegator，2 validator, onchain",
			2,
			2,
			[4][4]string{{"12", "12"}, {"12", "12"}},
			"0.5",
			1,
			2,
		},
		{
			"2 delegator，2 validator, offchain",
			2,
			2,
			[4][4]string{{"0", "0"}, {"0", "0"}},
			"0.5",
			0,
			2,
		},
		{
			"2 delegator，4 validator, onchain",
			4,
			2,
			[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
			"0.5",
			1,
			4,
		},
		{
			"2 delegator，4 validator, offchain",
			4,
			2,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			"0.5",
			0,
			4,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount, true)
			dk.SetDistributionType(ctx, tc.distrType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
				delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
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
				truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
				truncatedCommission, _ := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
				require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
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
					_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				}
			}

			require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
		})
	}
}

func (suite *DistributionSuite) TestWithdrawAllRewards() {
	testCases := []struct {
		title                string
		valCount             int64
		delCount             int64
		decRewards           [4]string
		rate                 string
		distrType            uint32
		remainReferenceCount uint64
		addShares            bool
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4]string{"0"},
			"0.5",
			1,
			1,
			false,
		},
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4]string{"48"},
			"0.5",
			1,
			1,
			true,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4]string{"0"},
			"0.5",
			0,
			1,
			true,
		},
		{
			"1 delegator，2 validator, onchain",
			2,
			1,
			[4]string{"48"},
			"0.5",
			1,
			2,
			true,
		},
		{
			"1 delegator，2 validator, offchain",
			2,
			1,
			[4]string{"0"},
			"0.5",
			0,
			2,
			true,
		},
		{
			"1 delegator，4 validator, onchain",
			4,
			1,
			[4]string{"48"},
			"0.5",
			1,
			4,
			true,
		},
		{
			"1 delegator，4 validator, offchain",
			4,
			1,
			[4]string{"0"},
			"0.5",
			0,
			4,
			true,
		},
		{
			"2 delegator，1 validator, onchain",
			1,
			2,
			[4]string{"24", "24"},
			"0.5",
			1,
			1,
			true,
		},
		{
			"2 delegator，1 validator, offchain",
			1,
			2,
			[4]string{"0", "0"},
			"0.5",
			0,
			1,
			true,
		},
		{
			"2 delegator，2 validator, onchain",
			2,
			2,
			[4]string{"24", "24"},
			"0.5",
			1,
			2,
			true,
		},
		{
			"2 delegator，2 validator, offchain",
			2,
			2,
			[4]string{"0", "0"},
			"0.5",
			0,
			2,
			true,
		},
		{
			"2 delegator，4 validator, onchain",
			4,
			2,
			[4]string{"24", "24"},
			"0.5",
			1,
			4,
			true,
		},
		{
			"2 delegator，4 validator, offchain",
			4,
			2,
			[4]string{"0", "0"},
			"0.5",
			0,
			4,
			true,
		},
		{
			"1 delegator，4 validator, onchain, rate 0",
			4,
			1,
			[4]string{"96"},
			"0",
			1,
			4,
			true,
		},
		{
			"1 delegator，4 validator, onchain, rate 1",
			4,
			1,
			[4]string{"0"},
			"1",
			1,
			4,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount, true)
			dk.SetDistributionType(ctx, tc.distrType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				if tc.addShares {
					keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
					delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
					require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
				}
			}

			allocateTokens(suite.T())

			beforeValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				beforeValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				response := keeper.GetQueriedDelegationTotalRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestDelAddrs[i])
				var queryRewards sdk.SysCoins
				for _, v := range response.Rewards {
					coins, _ := v.Reward.TruncateWithPrec(int64(0))
					queryRewards = queryRewards.Add2(coins)
				}
				err := dk.WithdrawDelegationAllRewards(ctx, keeper.TestDelAddrs[i])
				if tc.addShares {
					require.Nil(suite.T(), err)
				} else {
					require.Equal(suite.T(), types.ErrCodeEmptyDelegationVoteValidator(), err)
				}

				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decRewards[i]))
				require.Equal(suite.T(), queryRewards, getDecCoins(tc.decRewards[i]))
			}

			afterValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
			}

			//withdraw again
			staking.EndBlocker(ctx, sk)
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				err := dk.WithdrawDelegationAllRewards(ctx, keeper.TestDelAddrs[i])
				if tc.addShares {
					require.Nil(suite.T(), err)
				} else {
					require.Equal(suite.T(), types.ErrCodeEmptyDelegationVoteValidator(), err)
				}
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins("0"))
				truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]).TruncateDecimal()
				truncatedCommission, _ := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]).TruncateDecimal()
				if tc.addShares {
					require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
				} else {
					require.Equal(suite.T(), truncatedOutstanding, truncatedCommission.QuoDec(sdk.MustNewDecFromStr(tc.rate)))
				}
			}
			require.Equal(suite.T(), beforeValCommission, afterValCommission)

			//allocate and withdraw agagin
			allocateTokens(suite.T())
			for i := int64(0); i < tc.delCount; i++ {
				response := keeper.GetQueriedDelegationTotalRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestDelAddrs[i])
				var queryRewards sdk.SysCoins
				for _, v := range response.Rewards {
					coins, _ := v.Reward.TruncateWithPrec(int64(0))
					queryRewards = queryRewards.Add2(coins)
				}
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				err := dk.WithdrawDelegationAllRewards(ctx, keeper.TestDelAddrs[i])
				if tc.addShares {
					require.Nil(suite.T(), err)
				} else {
					require.Equal(suite.T(), types.ErrCodeEmptyDelegationVoteValidator(), err)
				}
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decRewards[i]))
				require.Equal(suite.T(), queryRewards, getDecCoins(tc.decRewards[i]))
			}

			//withdraw token
			allocateTokens(suite.T())
			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				rewards := getDecCoins(tc.decRewards[i])
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), rewards)
				for j := int64(0); j < tc.valCount; j++ {
					require.False(suite.T(), dk.HasDelegatorStartingInfo(ctx, keeper.TestValAddrs[j], keeper.TestDelAddrs[i]))
				}
			}

			//withdraw again
			for i := int64(0); i < tc.delCount; i++ {
				err := dk.WithdrawDelegationAllRewards(ctx, keeper.TestDelAddrs[i])
				require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)
			}

			require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
		})
	}
}

func (suite *DistributionSuite) TestDestroyValidator() {
	testCases := []struct {
		title                string
		valCount             int64
		delCount             int64
		decRewards           [4][4]string
		rate                 string
		distrType            uint32
		remainReferenceCount uint64
	}{
		{
			"1 delegator，1 validator, onchain",
			1,
			1,
			[4][4]string{{"24"}},
			"0.5",
			1,
			0,
		},
		{
			"1 delegator，1 validator, offchain",
			1,
			1,
			[4][4]string{{"0"}},
			"0.5",
			0,
			0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount, true)
			dk.SetDistributionType(ctx, tc.distrType)

			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
				delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
			}

			for j := int64(0); j < tc.valCount; j++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestValAccAddrs[j], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestValAccAddrs[j], keeper.TestValAddrs[0:tc.valCount])
				delegator := sk.Delegator(ctx, keeper.TestValAccAddrs[j])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
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
				truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
				truncatedCommission, _ := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i]).TruncateDecimal()
				require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
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

			hook := dk.Hooks()
			for j := int64(0); j < tc.valCount; j++ {
				hook.AfterValidatorRemoved(ctx, nil, keeper.TestValAddrs[j])
				require.True(suite.T(), dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[j]).IsZero())
				require.Panics(suite.T(), func() {
					require.True(suite.T(), dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[j]).IsZero())
				})
			}

			require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
		})
	}
}

func setTestFees(t *testing.T, ctx sdk.Context, ak auth.AccountKeeper, fees sdk.SysCoins) {
	feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	require.NotNil(t, feeCollector)
	err := feeCollector.SetCoins(fees)
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
}

func (suite *DistributionSuite) TestUpgrade() {
	testCases := []struct {
		title                   string
		valCount                int64
		delCount                int64
		decBeforeUpgradeRewards [4][4]string
		decAfterUpgradeRewards  [4][4]string
		rate                    string
		remainReferenceCount    uint64
	}{
		{
			"1 delegator，1 validator",
			1,
			1,
			[4][4]string{{"0"}},
			[4][4]string{{"48"}},
			"0.5",
			2,
		},
		{
			"1 delegator，2 validator",
			2,
			1,
			[4][4]string{{"0", "0"}},
			[4][4]string{{"24", "24"}},
			"0.5",
			4,
		},
		{
			"2 delegator，1 validator",
			1,
			2,
			[4][4]string{{"0"}, {"0"}},
			[4][4]string{{"24"}, {"24"}},
			"0.5",
			2,
		},
		{
			"2 delegator，2 validator",
			2,
			2,
			[4][4]string{{"0", "0"}, {"0", "0"}},
			[4][4]string{{"12", "12"}, {"12", "12"}},
			"0.5",
			4,
		},
		{
			"2 delegator，4 validator",
			4,
			2,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			[4][4]string{{"6", "6", "6", "6"}, {"6", "6", "6", "6"}},
			"0.5",
			8,
		},
		{
			"4 delegator，4 validator",
			4,
			4,
			[4][4]string{{"0", "0", "0", "0"}, {"0", "0", "0", "0"}, {"0", "0", "0", "0"}, {"0", "0", "0", "0"}},
			[4][4]string{{"3", "3", "3", "3"}, {"3", "3", "3", "3"}, {"3", "3", "3", "3"}, {"3", "3", "3", "3"}},
			"0.5",
			8,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), tc.valCount, false)
			ctx.SetBlockTime(time.Now())
			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:tc.valCount])
				delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
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
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decBeforeUpgradeRewards[i][j]))
				}
			}

			afterValCommission := [4]types.ValidatorAccumulatedCommission{}
			for i := int64(0); i < tc.valCount; i++ {
				afterValCommission[i] = dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[i])
				require.Panics(suite.T(), func() {
					require.True(suite.T(), dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[i]).IsZero())
				})
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
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decBeforeUpgradeRewards[i][j]))
				}
			}

			// upgrade
			proposal := makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
			hdlr := NewDistributionProposalHandler(dk)
			require.NoError(suite.T(), hdlr(ctx, &proposal))
			queryDistrType := dk.GetDistributionType(ctx)
			require.Equal(suite.T(), queryDistrType, types.DistributionTypeOnChain)
			staking.EndBlocker(ctx, sk)
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)

			//set rate
			newRate, _ := sdk.NewDecFromStr(tc.rate)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			for i := int64(0); i < tc.valCount; i++ {
				keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[i], newRate)
			}

			allocateTokens(suite.T())

			//withdraw reward
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					queryRewards := keeper.GetQueriedDelegationRewards(suite.T(), ctx, NewQuerier(dk), keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					queryRewards, _ = queryRewards.TruncateWithPrec(int64(0))
					beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
					require.Equal(suite.T(), afterAccount.Sub(beforeAccount), getDecCoins(tc.decAfterUpgradeRewards[i][j]))
					require.Equal(suite.T(), queryRewards, getDecCoins(tc.decAfterUpgradeRewards[i][j]))
				}
			}

			//withdraw reward again
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

			allocateTokens(suite.T())

			//withdraw token
			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				rewards := sdk.SysCoins{}
				for j := int64(0); j < tc.valCount; j++ {
					rewards = rewards.Add2(getDecCoins(tc.decAfterUpgradeRewards[i][j]))
				}
				require.Equal(suite.T(), afterAccount.Sub(beforeAccount), rewards)
				for j := int64(0); j < tc.valCount; j++ {
					require.False(suite.T(), dk.HasDelegatorStartingInfo(ctx, keeper.TestValAddrs[j], keeper.TestDelAddrs[i]))
				}
			}

			//withdraw again
			for i := int64(0); i < tc.delCount; i++ {
				for j := int64(0); j < tc.valCount; j++ {
					_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[j])
					require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)

					truncatedOutstanding, _ := dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[j]).TruncateDecimal()
					truncatedCommission, _ := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[j]).TruncateDecimal()
					require.Equal(suite.T(), truncatedOutstanding, truncatedCommission)
				}
			}

			require.Equal(suite.T(), tc.remainReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
		})
	}
}

func allocateVariateTokens(t *testing.T, blockRewards string) {
	feePoolBefore, _ := dk.GetFeePool(ctx).CommunityPool.TruncateDecimal()
	VariateBlockRewards := sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.MustNewDecFromStr(blockRewards)}}
	setTestFees(t, ctx, ak, VariateBlockRewards)
	dk.SetCommunityTax(ctx, sdk.MustNewDecFromStr("0"))
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes[0:1])
	feePoolAfter, _ := dk.GetFeePool(ctx).CommunityPool.TruncateDecimal()
	require.Equal(t, feePoolBefore, feePoolAfter)
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
}

func (suite *DistributionSuite) TestTruncateWithPrecWithdraw() {
	testCases := []struct {
		title        string
		precision    int64
		delCount     int64
		depositCoins [4]string
		blockRewards string
		decRewards   [4]string
	}{
		{
			"1 delegator, precision 18, reward 1",
			18,
			1,
			[4]string{"100"},
			"1",
			[4]string{"0.990099009900990000"},
		},
		{
			"2 delegator, precision 18, reward 1",
			18,
			2,
			[4]string{"100", "200"},
			"1",
			[4]string{"0.332225913621262400", "0.664451827242524800"},
		},
		{
			"3 delegator, precision 18, reward 1",
			18,
			3,
			[4]string{"100", "200", "300"},
			"1",
			[4]string{"0.166389351081530700", "0.332778702163061400", "0.499168053244592100"},
		},
		{
			"4 delegator, precision 18, reward 1",
			18,
			4,
			[4]string{"100", "200", "300", "400"},
			"1",
			[4]string{"0.099900099900099900", "0.199800199800199800", "0.299700299700299700", "0.399600399600399600"},
		},
		{
			"1 delegator, precision 5, reward 1",
			5,
			1,
			[4]string{"100"},
			"1",
			[4]string{"0.99009"},
		},
		{
			"2 delegator, precision 18, reward 1",
			5,
			2,
			[4]string{"100", "200"},
			"1",
			[4]string{"0.33222", "0.66445"},
		},
		{
			"3 delegator, precision 18, reward 1",
			5,
			3,
			[4]string{"100", "200", "300"},
			"1",
			[4]string{"0.16638", "0.33277", "0.49916"},
		},
		{
			"4 delegator, precision 18, reward 1",
			5,
			4,
			[4]string{"100", "200", "300", "400"},
			"1",
			[4]string{"0.09990", "0.19980", "0.29970", "0.39960"},
		},
		{
			"4 delegator, precision 0, reward 1",
			0,
			4,
			[4]string{"100", "200", "300", "400"},
			"1",
			[4]string{"0", "0", "0", "0"},
		},
		{
			"4 delegator, precision 18, reward 100",
			18,
			4,
			[4]string{"100", "200", "300", "400"},
			"100",
			[4]string{"9.990009990009990000", "19.980019980019980000", "29.970029970029970000", "39.960039960039960000"},
		},
		{
			"4 delegator, precision 10, reward 100",
			10,
			4,
			[4]string{"100", "200", "300", "400"},
			"100",
			[4]string{"9.9900099900", "19.9800199800", "29.9700299700", "39.9600399600"},
		},
		{
			"4 delegator, precision 1, reward 100",
			1,
			4,
			[4]string{"100", "200", "300", "400"},
			"100",
			[4]string{"9.9", "19.9", "29.9", "39.9"},
		},
		{
			"4 delegator, precision 1, reward 100",
			0,
			4,
			[4]string{"100", "200", "300", "400"},
			"100",
			[4]string{"9", "19", "29", "39"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			initEnv(suite.T(), 1, true)
			dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
			dk.SetRewardTruncatePrecision(ctx, tc.precision)
			ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
			keeper.DoEditValidator(suite.T(), ctx, sk, keeper.TestValAddrs[0], sdk.MustNewDecFromStr("0"))
			// UTC Time: 2000/1/1 00:00:00
			blockTimestampEpoch := int64(946684800)
			ctx.SetBlockTime(time.Unix(blockTimestampEpoch, 0))

			//deposit, add shares, withdraw msg in one block
			allocateVariateTokens(suite.T(), tc.blockRewards)
			staking.EndBlocker(ctx, sk)
			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], getDecCoins(tc.depositCoins[i])[0])
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:1])
				delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())

				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[0])
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), getDecCoins("0"), afterAccount.Sub(beforeAccount))

				beforeAccount = ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], getDecCoins(tc.depositCoins[i])[0])
				afterAccount = ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), getDecCoins("0"), afterAccount.Sub(beforeAccount))
			}

			allocateVariateTokens(suite.T(), tc.blockRewards)
			staking.EndBlocker(ctx, sk)
			//nomal
			for i := int64(0); i < tc.delCount; i++ {
				keeper.DoDeposit(suite.T(), ctx, sk, keeper.TestDelAddrs[i], getDecCoins(tc.depositCoins[i])[0])
				keeper.DoAddShares(suite.T(), ctx, sk, keeper.TestDelAddrs[i], keeper.TestValAddrs[0:1])
				delegator := sk.Delegator(ctx, keeper.TestDelAddrs[i])
				require.False(suite.T(), delegator.GetLastAddedShares().IsZero())
			}

			allocateVariateTokens(suite.T(), tc.blockRewards)
			staking.EndBlocker(ctx, sk)
			//withdraw reward
			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			feePoolBefore := dk.GetFeePool(ctx).CommunityPool
			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[i], keeper.TestValAddrs[0])
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), getDecCoins(tc.decRewards[i]), afterAccount.Sub(beforeAccount))
			}
			feePoolEnd := dk.GetFeePool(ctx).CommunityPool
			diff := feePoolEnd.Sub(feePoolBefore)
			if tc.precision == sdk.Precision {
				require.True(suite.T(), diff.IsZero())
			} else {
				require.False(suite.T(), diff.IsZero())
			}

			// withdraw
			allocateVariateTokens(suite.T(), tc.blockRewards)
			staking.EndBlocker(ctx, sk)
			for i := int64(0); i < tc.delCount; i++ {
				beforeAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				keeper.DoWithdraw(suite.T(), ctx, sk, keeper.TestDelAddrs[i], depositCoin)
				afterAccount := ak.GetAccount(ctx, keeper.TestDelAddrs[i]).GetCoins()
				require.Equal(suite.T(), getDecCoins(tc.decRewards[i]), afterAccount.Sub(beforeAccount))
			}
		})
	}
}
