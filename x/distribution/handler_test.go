package distribution

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) TestHandlerWithdrawDelegatorReward() {
	testCases := []struct {
		title    string
		dochange func(ctx sdk.Context, dk Keeper)
		errors   [3]sdk.Error
	}{
		{
			"change distribution type",
			func(ctx sdk.Context, dk Keeper) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)
			},
			[3]sdk.Error{types.ErrUnknownDistributionMsgType(), types.ErrCodeEmptyDelegationDistInfo(), nil},
		},
		{
			"set withdraw reward disable",
			func(ctx sdk.Context, dk Keeper) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)

				proposalWithdrawReward := types.NewWithdrawRewardEnabledProposal("title", "description", false)
				keeper.HandleWithdrawRewardEnabledProposal(ctx, dk, proposalWithdrawReward)
				require.Equal(suite.T(), false, dk.GetWithdrawRewardEnabled(ctx))
			},
			[3]sdk.Error{types.ErrUnknownDistributionMsgType(), types.ErrCodeEmptyDelegationDistInfo(), types.ErrCodeDisabledWithdrawRewards()},
		},
		{
			"no change distribution type",
			func(ctx sdk.Context, dk Keeper) {

			},
			[3]sdk.Error{types.ErrUnknownDistributionMsgType(), types.ErrUnknownDistributionMsgType(), types.ErrUnknownDistributionMsgType()},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			ctx, _, dk, sk, _ := keeper.CreateTestInputDefault(suite.T(), false, 10)
			handler := NewHandler(dk)
			delAddr1 := keeper.TestDelAddrs[0]
			valAddr1 := keeper.TestValAddrs[0]

			valOpAddrs := []sdk.ValAddress{valAddr1}

			msg := NewMsgWithdrawDelegatorReward(delAddr1, valAddr1)
			_, err := handler(ctx, msg)
			require.Equal(suite.T(), tc.errors[0], err)

			tc.dochange(ctx, dk)

			// no deposit and add shares
			_, err = handler(ctx, msg)
			require.Equal(suite.T(), tc.errors[1], err)

			// deposit and add shares
			keeper.DoDeposit(suite.T(), ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
			keeper.DoAddShares(suite.T(), ctx, sk, delAddr1, valOpAddrs)

			_, err = handler(ctx, msg)
			require.Equal(suite.T(), tc.errors[2], err)
		})
	}

}

type allocationParam struct {
	totalPower int64
	isVote     []bool
	isJailed   []bool
	fee        sdk.SysCoins
}

func createVotes(ctx sdk.Context, sk staking.Keeper, test allocationParam) []abci.VoteInfo {
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

func (suite *HandlerSuite) TestHandlerWithdrawValidatorCommission() {
	testCases := []struct {
		title            string
		doAllocateTokens func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper)
		dochange         func(ctx sdk.Context, dk Keeper)
		errors           [2]sdk.Error
	}{
		{
			"normal, no change distribution type",
			func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper) {
				feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
				require.NotNil(suite.T(), feeCollector)
				err := feeCollector.SetCoins(sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(100))}})
				require.NoError(suite.T(), err)
				ak.SetAccount(ctx, feeCollector)
				allocationParam := allocationParam{
					10,
					[]bool{true, true, true, true}, []bool{false, false, false, false},
					nil,
				}
				votes := createVotes(ctx, sk, allocationParam)
				dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
				require.Nil(suite.T(), err)
			},
			func(ctx sdk.Context, dk Keeper) {},
			[2]sdk.Error{types.ErrNoValidatorCommission(), nil},
		},
		{
			"no allocate tokens, no change distribution type",
			func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper) {

			},
			func(ctx sdk.Context, dk Keeper) {},
			[2]sdk.Error{types.ErrNoValidatorCommission(), types.ErrNoValidatorCommission()},
		},
		{
			"normal, change distribution type",
			func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper) {
				feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
				require.NotNil(suite.T(), feeCollector)
				err := feeCollector.SetCoins(sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(100))}})
				require.NoError(suite.T(), err)
				ak.SetAccount(ctx, feeCollector)
				allocationParam := allocationParam{
					10,
					[]bool{true, true, true, true}, []bool{false, false, false, false},
					nil,
				}
				votes := createVotes(ctx, sk, allocationParam)
				dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
				require.Nil(suite.T(), err)
			},
			func(ctx sdk.Context, dk Keeper) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)
			},
			[2]sdk.Error{types.ErrNoValidatorCommission(), nil},
		},
		{
			"no allocate tokens, change distribution type",
			func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)
			},
			func(ctx sdk.Context, dk Keeper) {},
			[2]sdk.Error{types.ErrNoValidatorCommission(), types.ErrNoValidatorCommission()},
		},
		{
			"normal, no impact when set withdraw reward disable",
			func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper) {
				feeCollector := supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
				require.NotNil(suite.T(), feeCollector)
				err := feeCollector.SetCoins(sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(100))}})
				require.NoError(suite.T(), err)
				ak.SetAccount(ctx, feeCollector)
				allocationParam := allocationParam{
					10,
					[]bool{true, true, true, true}, []bool{false, false, false, false},
					nil,
				}
				votes := createVotes(ctx, sk, allocationParam)
				dk.AllocateTokens(ctx, 1, keeper.TestConsAddrs[0], votes)
				require.Nil(suite.T(), err)
			},
			func(ctx sdk.Context, dk Keeper) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)

				proposalWithdrawReward := types.NewWithdrawRewardEnabledProposal("title", "description", false)
				keeper.HandleWithdrawRewardEnabledProposal(ctx, dk, proposalWithdrawReward)
				require.Equal(suite.T(), false, dk.GetWithdrawRewardEnabled(ctx))
			},
			[2]sdk.Error{types.ErrNoValidatorCommission(), nil},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			ctx, ak, dk, sk, supplyKeeper := keeper.CreateTestInputDefault(suite.T(), false, 10)
			handler := NewHandler(dk)
			valAddr1 := keeper.TestValAddrs[0]

			msg := NewMsgWithdrawValidatorCommission(valAddr1)

			_, err := handler(ctx, msg)
			require.Equal(suite.T(), tc.errors[0], err)

			staking.EndBlocker(ctx, sk)
			tc.dochange(ctx, dk)
			tc.doAllocateTokens(ctx, ak, dk, sk, supplyKeeper)
			_, err = handler(ctx, msg)
			require.Equal(suite.T(), tc.errors[1], err)
		})
	}
}
