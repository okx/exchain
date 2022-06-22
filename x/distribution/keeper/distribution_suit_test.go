package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	custom                 = "custom"
	valOpAddrs             []sdk.ValAddress
	valConsPks             []crypto.PubKey
	valConsAddrs           []sdk.ConsAddress
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
)

type delegatorVoteParam struct {
	delegators                []sdk.AccAddress
	beforeDelegatorsRewards   []string
	afterDelegatorsRewards    []string
	beforeValidatorCommission []string
	afterValidatorCommission  []string
}

//func getTestCases() []delegatorVoteParam {
//	return []delegatorVoteParam{
//		{ //test the case when fee is zero
//			10,
//			[]bool{true, true, true, true}, []bool{false, false, false, false},
//			nil,
//		},
//	}
//}

func initEnv(t *testing.T) {
	valOpAddrs, valConsPks, valConsAddrs = GetTestAddrs()
	ctx, ak, dk, sk, _ = CreateTestInputDefault(t, false, 1000)
	oneValOpAddrs = []sdk.ValAddress{valOpAddr1}
	twoValOpAddrs = []sdk.ValAddress{valOpAddr1, valOpAddr2}
	threeValOpAddrs = []sdk.ValAddress{valOpAddr1, valOpAddr2, valOpAddr3}
	fourValOpAddrs = []sdk.ValAddress{valOpAddr1, valOpAddr2, valOpAddr3, valOpAddr4}
	expectReferenceCount += len(valOpAddrs)
	require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	testAllocationParam := testAllocationParam{
		10,
		[]bool{true, true, true, true}, []bool{false, false, false, false},
		nil,
	}
	votes = createTestVotes(ctx, sk, testAllocationParam)
}

func TestDistributionTypeSuit(t *testing.T) {
	//initEnv(t)
	//
	//setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	//
	//dk.AllocateTokens(ctx, 1, valConsAddr1, votes)
	//decCommission, _ := sdk.NewDecFromStr("24.5")
	//decCommissionRemain, _ := sdk.NewDecFromStr("0.5")
	//decOutstanding, _ := sdk.NewDecFromStr("24.5")
	//decOutstandingRemain, _ := sdk.NewDecFromStr("0.5")
	//decCommunity, _ := sdk.NewDecFromStr("2")
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))
	//
	//dk.WithdrawValidatorCommission(ctx, valOpAddr1)
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommissionRemain}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstandingRemain}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))
	//
	////next block, first delegator vote
	//
	//doDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//doAddShares(t, ctx, sk, delAddr1, oneValOpAddrs)
	//expectReferenceCount += len(oneValOpAddrs)
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//
	////next block, second delegator vote
	//doDeposit(t, ctx, sk, delAddr2, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//doAddShares(t, ctx, sk, delAddr2, oneValOpAddrs)
	//expectReferenceCount += len(oneValOpAddrs)
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//
	//setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	//dk.AllocateTokens(ctx, 1, valConsAddr1, votes)
	//decCommission, _ = sdk.NewDecFromStr("25")
	//decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	//decOutstanding, _ = sdk.NewDecFromStr("25")
	//decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	//decCommunity, _ = sdk.NewDecFromStr("4")
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))
	//
	//dk.WithdrawValidatorCommission(ctx, valOpAddr1)
	////require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommissionRemain}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	////require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstandingRemain}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))
	//
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//
	////next block, third delegator vote
	//doDeposit(t, ctx, sk, delAddr3, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//doAddShares(t, ctx, sk, delAddr3, oneValOpAddrs)
	//expectReferenceCount += len(oneValOpAddrs)
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//
	//setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	//dk.AllocateTokens(ctx, 1, valConsAddr1, votes)
	//decCommission, _ = sdk.NewDecFromStr("24.5")
	//decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	//decOutstanding, _ = sdk.NewDecFromStr("24.5")
	//decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	//decCommunity, _ = sdk.NewDecFromStr("6")
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))
	//
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//
	////next block, fourth delegator vote
	//doDeposit(t, ctx, sk, delAddr4, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//doAddShares(t, ctx, sk, delAddr4, oneValOpAddrs)
	//expectReferenceCount += 1
	//require.Equal(t, uint64(expectReferenceCount), dk.GetValidatorHistoricalReferenceCount(ctx))
	//
	//setTestFees(t, ctx, dk, ak, blockRewardValueTokens)
	//dk.AllocateTokens(ctx, 1, valConsAddr1, votes)
	//decCommission, _ = sdk.NewDecFromStr("49")
	//decCommissionRemain, _ = sdk.NewDecFromStr("0.5")
	//decOutstanding, _ = sdk.NewDecFromStr("49")
	//decOutstandingRemain, _ = sdk.NewDecFromStr("0.5")
	//decCommunity, _ = sdk.NewDecFromStr("8")
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommission}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decOutstanding}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
	//require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: decCommunity}}, dk.GetFeePoolCommunityCoins(ctx))
	//
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//
	////next block, first delegator withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr1)
	//
	//querier := NewQuerier(dk)
	//delRewards := getQueriedDelegatorTotalRewards(t, ctx, dk.cdc, querier, delAddr1)
	////require.Equal(t, types.QueryDelegatorTotalRewardsResponse{}, delRewards)
	//fmt.Println(delRewards)
	//
	////next block, second delegator withdraw
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr1)
	//
	////next block, third delegator withdraw
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr1)
	//
	////next block, fourth delegator withdraw
	//staking.EndBlocker(ctx, sk)
	//ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr1)
	//
	////check offchain type
	//require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)
	//
	////change val commission
	//newRate, _ := sdk.NewDecFromStr("0.5")
	//ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	//doEditValidator(t, ctx, sk, valOpAddr1, &newRate)
	//
	////check val commission and distribution type
	//val := sk.Validator(ctx, valOpAddr1)
	//require.Equal(t, newRate, val.GetCommission())
	//require.Equal(t, dk.GetDistributionType(ctx), types.DistributionTypeOffChain)
	//
	////change del1 to proxy1
	//doRegProxy(t, ctx, sk, delAddr1, true)
	//
	////bind del2 and del3 to proxy1
	//doBindProxy(t, ctx, sk, delAddr2, delAddr1)
	//doBindProxy(t, ctx, sk, delAddr3, delAddr1)
	//
	////del1 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr4)
	//
	////del2 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr4)
	//
	////del3 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr4)
	//
	////del4 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr4)
	//
	////del1 deposit
	//doDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//
	////del2 deposit
	//doDeposit(t, ctx, sk, delAddr2, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//
	////del3 deposit
	//doDeposit(t, ctx, sk, delAddr3, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//
	////del4 deposit
	//doDeposit(t, ctx, sk, delAddr4, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	//
	////del1 add shares 1 vals
	//doAddShares(t, ctx, sk, delAddr1, oneValOpAddrs)
	//
	////del2 add shares 2 vals
	//doAddShares(t, ctx, sk, delAddr2, twoValOpAddrs)
	//
	////del3 add shares 3 vals
	//doAddShares(t, ctx, sk, delAddr3, threeValOpAddrs)
	//
	////del4 add shares 4 vals
	//doAddShares(t, ctx, sk, delAddr4, fourValOpAddrs)
	//
	////del1 unbind
	//doUnBindProxy(t, ctx, sk, delAddr1)
	//
	////proxy unreg
	//doRegProxy(t, ctx, sk, delAddr1, false)
	//
	////val1 destroy
	//doDestroyValidator(t, ctx, sk, valAccAddr1)
	//
	////de1 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr4)
	//
	////de2 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr4)
	//
	////de3 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr3, valOpAddr4)
	//
	////de4 withdraw
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr1)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr2)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr3)
	//dk.WithdrawDelegationRewards(ctx, delAddr4, valOpAddr4)
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
