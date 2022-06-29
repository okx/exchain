package distribution

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var (
	custom                 = "custom"
	oneValOpAddrs          []sdk.ValAddress
	twoValOpAddrs          []sdk.ValAddress
	threeValOpAddrs        []sdk.ValAddress
	fourValOpAddrs         []sdk.ValAddress
	expectReferenceCount   = 0
	ctx                    sdk.Context
	ak                     auth.AccountKeeper
	dk                     Keeper
	sk                     staking.Keeper
	supplyKeeper           types.SupplyKeeper
	blockRewardValueTokens = sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(100))}}
	votes                  []abci.VoteInfo
	stdCommissionRate      sdk.Dec
	newCommissionRate      sdk.Dec
)

type delegatorVoteParam struct {
	delegators                []sdk.AccAddress
	beforeDelegatorsRewards   []string
	afterDelegatorsRewards    []string
	beforeValidatorCommission []string
	afterValidatorCommission  []string
}

type testAllocationParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        sdk.SysCoins
}

func initEnv(t *testing.T) {
	//keeper.TestValAddrs, valConsPks, valConsAddrs = keeper.GetTestAddrs()
	ctx, ak, dk, sk, supplyKeeper = keeper.CreateTestInputDefault(t, false, 1000)
	oneValOpAddrs = []sdk.ValAddress{keeper.TestValAddrs[0]}
	twoValOpAddrs = []sdk.ValAddress{keeper.TestValAddrs[1], keeper.TestValAddrs[1]}
	threeValOpAddrs = []sdk.ValAddress{keeper.TestValAddrs[0], keeper.TestValAddrs[1], keeper.TestValAddrs[2]}
	fourValOpAddrs = []sdk.ValAddress{keeper.TestValAddrs[0], keeper.TestValAddrs[1], keeper.TestValAddrs[2], keeper.TestValAddrs[3]}
	stdCommissionRate, _ = sdk.NewDecFromStr("1.0")
	newCommissionRate, _ = sdk.NewDecFromStr("0.5")
	//expectReferenceCount += len(keeper.TestValAddrs)
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
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

func TestDistributionSuit(t *testing.T) {
	initEnv(t)
	firstOffChain(t)
	//firstOnChain(t)
	//secondOffline(t)
	//secondOnline(t)
}

func firstOffChain(t *testing.T) {
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, second delegator vote
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, third delegator vote
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, fourth delegator vote
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, first delegator withdraw
	newBlockAndAllocateReward(t)
	_, err := dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	require.NotNil(t, err)

	//next block, second delegator withdraw
	newBlockAndAllocateReward(t)
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	require.NotNil(t, err)

	//next block, third delegator withdraw
	newBlockAndAllocateReward(t)
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	require.NotNil(t, err)

	//next block, fourth delegator withdraw
	newBlockAndAllocateReward(t)
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	require.NotNil(t, err)
	//check offchain type
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)
	//change val commission
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	//keeper.DoEditValidator(t, ctx, sk, keeper.TestValAddrs[0], newCommissionRate)
	//check val commission and distribution type
	val := sk.Validator(ctx, keeper.TestValAddrs[0])
	require.Equal(t, stdCommissionRate, val.GetCommission())
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//next block
	newBlockAndAllocateReward(t)
	//change del1 to proxy1
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], true)

	//bind del2 and del3 to proxy1
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[1], keeper.TestDelAddrs[0])
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[2], keeper.TestDelAddrs[0])

	//next block
	newBlockAndAllocateReward(t)
	//del1 and 4 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del1 add shares 1 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)

	//del4 add shares 4 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], fourValOpAddrs)

	//del2 unbind
	keeper.DoUnBindProxy(t, ctx, sk, keeper.TestDelAddrs[1])

	//withdraw error
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	require.NotNil(t, err)
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	require.NotNil(t, err)
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	require.NotNil(t, err)
	_, err = dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	require.NotNil(t, err)

	commission1 := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0])
	commission2 := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[1])
	commission3 := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[2])
	commission4 := dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[3])
	community := dk.GetFeePoolCommunityCoins(ctx)

	mulCount, _ := sdk.NewDecFromStr(fmt.Sprintf("%d", ctx.BlockHeight()-1))
	require.Equal(t, blockRewardValueTokens.MulDec(mulCount), commission1.Add(commission2[0]).Add(commission3[0]).Add(commission4[0]).Add(community[0]))
}

func firstOnChain(t *testing.T) {
	//set onchain
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)

	//set diff type, first
	proposal := makeChangeDistributionTypeProposal(types.DistributionTypeOnChain)
	hdlr := NewDistributionProposalHandler(dk)
	require.NoError(t, hdlr(ctx, &proposal))
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOnChain)

	newBlockAndAllocateReward(t)

	//next block, first delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, second delegator vote
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, third delegator vote
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, fourth delegator vote
	newBlockAndAllocateReward(t)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	//next block, first delegator withdraw
	newBlockAndAllocateReward(t)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])

	//querier := NewQuerier(dk)
	//delRewards := getQueriedDelegatorTotalRewards(t, ctx, dk.cdc, querier, keeper.TestDelAddrs[0])
	//require.Equal(t, types.QueryDelegatorTotalRewardsResponse{}, delRewards)
	//fmt.Println(delRewards)

	//next block, second delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])

	//next block, third delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])

	//next block, fourth delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])

	//check offchain type
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//change val commission
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	keeper.DoEditValidator(t, ctx, sk, keeper.TestValAddrs[0], newRate)

	//check val commission and distribution type
	val := sk.Validator(ctx, keeper.TestValAddrs[0])
	require.Equal(t, newRate, val.GetCommission())
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//change del1 to proxy1
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], true)

	//bind del2 and del3 to proxy1
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[1], keeper.TestDelAddrs[0])
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[2], keeper.TestDelAddrs[0])

	//del1 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[3])

	//del2 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[3])

	//del3 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[3])

	//del4 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[3])

	//del1 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del2 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del3 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del4 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del1 add shares 1 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)

	//del2 add shares 2 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], twoValOpAddrs)

	//del3 add shares 3 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], threeValOpAddrs)

	//del4 add shares 4 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], fourValOpAddrs)

	//del1 unbind
	keeper.DoUnBindProxy(t, ctx, sk, keeper.TestDelAddrs[0])

	//proxy unreg
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], false)

	//val1 destroy
	keeper.DoDestroyValidator(t, ctx, sk, keeper.TestValAccAddrs[0])

	//de1 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[3])

	//de2 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[3])

	//de3 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[3])

	//de4 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[3])
}

func secondOnline(t *testing.T) {
	initEnv(t)

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)

	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ := sdk.NewDecFromStr("24.5")
	decCommissionRemain, _ := sdk.NewDecFromStr("0.5")
	decOutstanding, _ := sdk.NewDecFromStr("24.5")
	decOutstandingRemain, _ := sdk.NewDecFromStr("0.5")
	decCommunity, _ := sdk.NewDecFromStr("2")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	dk.WithdrawValidatorCommission(ctx, keeper.TestValAddrs[0])
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommissionRemain}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstandingRemain}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	//next block, first delegator vote

	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, second delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ = sdk.NewDecFromStr("25")
	decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	decOutstanding, _ = sdk.NewDecFromStr("25")
	decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	decCommunity, _ = sdk.NewDecFromStr("4")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	dk.WithdrawValidatorCommission(ctx, keeper.TestValAddrs[0])
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, third delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ = sdk.NewDecFromStr("24.5")
	decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	decOutstanding, _ = sdk.NewDecFromStr("24.5")
	decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	decCommunity, _ = sdk.NewDecFromStr("6")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, fourth delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ = sdk.NewDecFromStr("49")
	decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	decOutstanding, _ = sdk.NewDecFromStr("49")
	decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	decCommunity, _ = sdk.NewDecFromStr("8")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, first delegator withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])

	//querier := NewQuerier(dk)
	//delRewards := getQueriedDelegatorTotalRewards(t, ctx, dk.cdc, querier, keeper.TestDelAddrs[0])
	//require.Equal(t, types.QueryDelegatorTotalRewardsResponse{}, delRewards)
	//fmt.Println(delRewards)

	//next block, second delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])

	//next block, third delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])

	//next block, fourth delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])

	//check offchain type
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//change val commission
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	keeper.DoEditValidator(t, ctx, sk, keeper.TestValAddrs[0], newRate)

	//check val commission and distribution type
	val := sk.Validator(ctx, keeper.TestValAddrs[0])
	require.Equal(t, newRate, val.GetCommission())
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//change del1 to proxy1
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], true)

	//bind del2 and del3 to proxy1
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[1], keeper.TestDelAddrs[0])
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[2], keeper.TestDelAddrs[0])

	//del1 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[3])

	//del2 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[3])

	//del3 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[3])

	//del4 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[3])

	//del1 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del2 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del3 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del4 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del1 add shares 1 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)

	//del2 add shares 2 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], twoValOpAddrs)

	//del3 add shares 3 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], threeValOpAddrs)

	//del4 add shares 4 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], fourValOpAddrs)

	//del1 unbind
	keeper.DoUnBindProxy(t, ctx, sk, keeper.TestDelAddrs[0])

	//proxy unreg
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], false)

	//val1 destroy
	keeper.DoDestroyValidator(t, ctx, sk, keeper.TestValAccAddrs[0])

	//de1 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[3])

	//de2 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[3])

	//de3 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[3])

	//de4 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[3])
}

func secondOffline(t *testing.T) {
	initEnv(t)
	newBlockAndAllocateReward(t)

	decCommission, _ := sdk.NewDecFromStr("24.5")
	decCommissionRemain, _ := sdk.NewDecFromStr("0.5")
	decOutstanding, _ := sdk.NewDecFromStr("24.5")
	decOutstandingRemain, _ := sdk.NewDecFromStr("0.5")
	decCommunity, _ := sdk.NewDecFromStr("2")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	dk.WithdrawValidatorCommission(ctx, keeper.TestValAddrs[0])
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommissionRemain}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstandingRemain}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	//next block, first delegator vote

	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, second delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ = sdk.NewDecFromStr("25")
	decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	decOutstanding, _ = sdk.NewDecFromStr("25")
	decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	decCommunity, _ = sdk.NewDecFromStr("4")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	dk.WithdrawValidatorCommission(ctx, keeper.TestValAddrs[0])
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, third delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], oneValOpAddrs)
	expectReferenceCount += len(oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ = sdk.NewDecFromStr("24.5")
	decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	decOutstanding, _ = sdk.NewDecFromStr("24.5")
	decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	decCommunity, _ = sdk.NewDecFromStr("6")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, fourth delegator vote
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], oneValOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))

	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	decCommission, _ = sdk.NewDecFromStr("49")
	decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	decOutstanding, _ = sdk.NewDecFromStr("49")
	decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	decCommunity, _ = sdk.NewDecFromStr("8")
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, keeper.TestValAddrs[0]))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))

	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	//next block, first delegator withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])

	//querier := NewQuerier(dk)
	//delRewards := getQueriedDelegatorTotalRewards(t, ctx, dk.cdc, querier, keeper.TestDelAddrs[0])
	//require.Equal(t, types.QueryDelegatorTotalRewardsResponse{}, delRewards)
	//fmt.Println(delRewards)

	//next block, second delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])

	//next block, third delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])

	//next block, fourth delegator withdraw
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])

	//check offchain type
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//change val commission
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	keeper.DoEditValidator(t, ctx, sk, keeper.TestValAddrs[0], newRate)

	//check val commission and distribution type
	val := sk.Validator(ctx, keeper.TestValAddrs[0])
	require.Equal(t, newRate, val.GetCommission())
	require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)

	//change del1 to proxy1
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], true)

	//bind del2 and del3 to proxy1
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[1], keeper.TestDelAddrs[0])
	keeper.DoBindProxy(t, ctx, sk, keeper.TestDelAddrs[2], keeper.TestDelAddrs[0])

	//del1 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[3])

	//del2 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[3])

	//del3 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[3])

	//del4 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[3])

	//del1 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del2 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[1], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del3 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[2], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del4 deposit
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

	//del1 add shares 1 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[0], oneValOpAddrs)

	//del2 add shares 2 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[1], twoValOpAddrs)

	//del3 add shares 3 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[2], threeValOpAddrs)

	//del4 add shares 4 vals
	keeper.DoAddShares(t, ctx, sk, keeper.TestDelAddrs[3], fourValOpAddrs)

	//del1 unbind
	keeper.DoUnBindProxy(t, ctx, sk, keeper.TestDelAddrs[0])

	//proxy unreg
	keeper.DoRegProxy(t, ctx, sk, keeper.TestDelAddrs[0], false)

	//val1 destroy
	keeper.DoDestroyValidator(t, ctx, sk, keeper.TestValAccAddrs[0])

	//de1 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[0], keeper.TestValAddrs[3])

	//de2 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[1], keeper.TestValAddrs[3])

	//de3 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[2], keeper.TestValAddrs[3])

	//de4 withdraw
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[0])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[1])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[2])
	dk.WithdrawDelegationRewards(ctx, keeper.TestDelAddrs[3], keeper.TestValAddrs[3])
}

func getQueriedDelegatorTotalRewards(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, delegatorAddr sdk.AccAddress) (response types.QueryDelegatorTotalRewardsResponse) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryDelegatorTotalRewards}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryDelegatorParams(delegatorAddr)),
	}

	bz, err := querier(ctx, []string{types.QueryDelegatorTotalRewards}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &response))

	return
}

func newBlockAndAllocateReward(t *testing.T) {
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
	keeper.DoDeposit(t, ctx, sk, keeper.TestDelAddrs[3], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
}

func setTestFees(t *testing.T, ctx sdk.Context, k Keeper, ak auth.AccountKeeper, fees sdk.SysCoins) {
	feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	require.NotNil(t, feeCollector)
	err := feeCollector.SetCoins(fees)
	require.NoError(t, err)
	ak.SetAccount(ctx, feeCollector)
}
