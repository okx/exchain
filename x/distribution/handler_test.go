package distribution

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	stakingtypes "github.com/okex/exchain/x/staking/types"
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
		errors   [4]sdk.Error
	}{
		{
			"change distribution type",
			func(ctx sdk.Context, dk Keeper) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)
			},
			[4]sdk.Error{types.ErrUnknownDistributionMsgType(), types.ErrCodeEmptyDelegationDistInfo(), nil, nil},
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
			[4]sdk.Error{types.ErrUnknownDistributionMsgType(), types.ErrCodeDisabledWithdrawRewards(),
				stakingtypes.ErrCodeDisabledOperate(), types.ErrCodeDisabledWithdrawRewards()},
		},
		{
			"no change distribution type",
			func(ctx sdk.Context, dk Keeper) {

			},
			[4]sdk.Error{types.ErrUnknownDistributionMsgType(), types.ErrUnknownDistributionMsgType(), nil, types.ErrUnknownDistributionMsgType()},
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
			keeper.DoDepositWithError(suite.T(), ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)), tc.errors[2])
			keeper.DoAddSharesWithError(suite.T(), ctx, sk, delAddr1, valOpAddrs, tc.errors[2])

			_, err = handler(ctx, msg)
			require.Equal(suite.T(), tc.errors[3], err)
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
			[2]sdk.Error{types.ErrNoValidatorCommission(), types.ErrCodeDisabledWithdrawRewards()},
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

func (suite *HandlerSuite) TestWithdrawDisabled() {
	type param struct {
		blockVersion int64
		enable       bool
		expectError  error
	}

	testCases := []struct {
		title   string
		execute func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int)
		params  []param
	}{
		{
			"create val",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoCreateValidatorWithError(suite.T(), *ctx, sk, keeper.TestValAddrs[0], nil, p.expectError)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, stakingtypes.ErrValidatorOwnerExists()},
				{0, false, stakingtypes.ErrValidatorOwnerExists()},
				{0, true, stakingtypes.ErrValidatorOwnerExists()},
			},
		},
		{
			"disable edit val",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoEditValidatorWithError(suite.T(), *ctx, sk, keeper.TestValAddrs[0], sdk.NewDec(0), p.expectError)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, stakingtypes.ErrCommissionUpdateTime()},
			},
		},
		{
			"disable destroy val",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoDestroyValidatorWithError(suite.T(), *ctx, sk, keeper.TestValAccAddrs[i], p.expectError)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, nil},
				{0, false, nil},
				{0, true, nil},
			},
		},
		{
			"disable withdraw DoAddShares",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoAddSharesWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[0], valOpAddrs, p.expectError)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, nil},
				{0, false, nil},
				{0, true, nil},
			},
		},
		{
			"disable withdraw DoRegProxy",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoRegProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i], true, p.expectError)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, nil},
				{0, false, nil},
				{0, true, nil},
			},
		},
		{
			"disable withdraw DoWithdraw",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i], sdk.NewCoin(sk.BondDenom(*ctx),
					sdk.NewInt(100)), p.expectError)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, nil},
				{0, false, nil},
				{0, true, nil},
			},
		},
		{
			"disable withdraw DoBindProxy",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoBindProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i+1], keeper.TestDelAddrs[0], p.expectError)
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				dk.SetWithdrawRewardEnabled(*ctx, true)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, nil},
				{0, false, nil},
			},
		},
		{
			"disable withdraw DoUnBindProxy",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				DoBindProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i+1], keeper.TestDelAddrs[0], nil)
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoUnBindProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i+1], p.expectError)
				tmtypes.UnittestOnlySetMilestoneVenus2Height(-1)
				dk.SetWithdrawRewardEnabled(*ctx, true)
			},
			[]param{
				{-1, false, stakingtypes.ErrCodeDisabledOperate()},
				{-1, true, nil},
				{0, false, nil},
			},
		},
		{
			"disable withdraw address",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoSetWithdrawAddressWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[i], p.expectError)
			},
			[]param{
				{-1, false, types.ErrCodeDisabledWithdrawRewards()},
				{-1, true, nil},
				{0, false, nil},
				{0, true, nil},
			},
		},
		{
			"disable withdraw validator",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawValidatorCommissionWithError(suite.T(), *ctx, dk, keeper.TestValAddrs[0], p.expectError)
			},
			[]param{
				{-1, false, types.ErrCodeDisabledWithdrawRewards()},
				{-1, true, types.ErrNoValidatorCommission()},
				{0, false, types.ErrNoValidatorCommission()},
				{0, true, types.ErrNoValidatorCommission()},
			},
		},
		{
			"disable set withdraw address",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoSetWithdrawAddressWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[i], p.expectError)
			},
			[]param{
				{-1, false, types.ErrCodeDisabledWithdrawRewards()},
				{-1, true, nil},
				{0, false, nil},
				{0, true, nil},
			},
		},
		{
			"disable set withdraw address",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawDelegatorRewardWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[0], keeper.TestValAddrs[0], p.expectError)
			},
			[]param{
				{-1, false, types.ErrCodeDisabledWithdrawRewards()},
				{-1, true, types.ErrUnknownDistributionMsgType()},
				{0, false, types.ErrUnknownDistributionMsgType()},
				{0, true, types.ErrUnknownDistributionMsgType()},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			communityTax := sdk.NewDecWithPrec(2, 2)
			ctx, _, _, dk, sk, _, _ := keeper.CreateTestInputAdvanced(suite.T(), false, 1000, communityTax)
			valOpAddrs, valConsPks, _ := keeper.GetTestAddrs()
			for i, _ := range valOpAddrs {
				keeper.DoCreateValidator(suite.T(), ctx, sk, valOpAddrs[i], valConsPks[i])
			}
			// end block to bond validator
			staking.EndBlocker(ctx, sk)
			//delegation
			for _, v := range keeper.TestDelAddrs {
				keeper.DoDeposit(suite.T(), ctx, sk, v, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
				keeper.DoAddShares(suite.T(), ctx, sk, v, valOpAddrs)
			}

			DoRegProxyWithError(suite.T(), ctx, sk, keeper.TestDelAddrs[0], true, nil)
			DoDepositWithError(suite.T(), ctx, sk, keeper.TestDelAddrs[0], sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)), nil)

			for i, p := range tc.params {
				tc.execute(&ctx, dk, sk, valOpAddrs, p, i)
			}

			proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
			keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
		})
	}
}

//func (suite *HandlerSuite) TestValidatorSubmitProposal() {
//
//	type param struct {
//		blockVersion int64
//		addr         sdk.AccAddress
//		expectError  error
//	}
//
//	testCases := []struct {
//		title   string
//		execute func(ctx *sdk.Context, dk Keeper, p param)
//		params  []param
//	}{
//		{
//			"ChangeDistributionTypeProposal",
//			func(ctx *sdk.Context, dk Keeper, p param) {
//				tmtypes.UnittestOnlySetMilestoneVenus2Height(p.blockVersion)
//				proposal := types.NewChangeDistributionTypeProposal("change type", "", types.DistributionTypeOnChain)
//				minDeposit := dk.GetMinDeposit(*ctx, proposal)
//				msg := govtypes.NewMsgSubmitProposal(proposal, minDeposit, p.addr)
//				dk.CheckMsgSubmitProposal(*ctx, msg)
//				keeper.HandleChangeDistributionTypeProposal(*ctx, dk, proposal)
//			},
//			[]param{
//				{-1, keeper.TestDelAddrs[0], types.ErrCodeDisabledWithdrawRewards()},
//				//{0, keeper.TestDelAddrs[0], types.ErrUnknownDistributionMsgType()},
//				//{-1, keeper.TestValAccAddrs[0], types.ErrCodeDisabledWithdrawRewards()},
//				//{0, keeper.TestValAccAddrs[0], types.ErrUnknownDistributionMsgType()},
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		suite.Run(tc.title, func() {
//			communityTax := sdk.NewDecWithPrec(2, 2)
//			ctx, _, _, dk, _, _, _ := CreateTestInputAdvancedA(suite.T(), false, 1000, communityTax)
//			for _, p := range tc.params {
//				tc.execute(&ctx, dk, p)
//			}
//
//			proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
//			keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
//		})
//	}
//}

func DoCreateValidatorWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, valAddr sdk.ValAddress, valConsPk crypto.PubKey, expectError error) {
	s := staking.NewHandler(sk)
	msg := staking.NewMsgCreateValidator(valAddr, valConsPk, staking.Description{}, keeper.NewTestSysCoin(1, 0))
	_, e := s(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoEditValidatorWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, valAddr sdk.ValAddress, newRate sdk.Dec, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgEditValidatorCommissionRate(valAddr, newRate)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoWithdrawWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, amount sdk.SysCoin, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgWithdraw(delAddr, amount)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoDestroyValidatorWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgDestroyValidator(delAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoDepositWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, amount sdk.SysCoin, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgDeposit(delAddr, amount)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoAddSharesWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgAddShares(delAddr, valAddrs)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoRegProxyWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, reg bool, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgRegProxy(delAddr, reg)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoBindProxyWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, proxyAddr sdk.AccAddress, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgBindProxy(delAddr, proxyAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoUnBindProxyWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, expectError error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgUnbindProxy(delAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoSetWithdrawAddressWithError(t *testing.T, ctx sdk.Context, dk Keeper, delAddr sdk.AccAddress, expectError error) {
	h := NewHandler(dk)
	msg := NewMsgSetWithdrawAddress(delAddr, delAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoWithdrawValidatorCommissionWithError(t *testing.T, ctx sdk.Context, dk Keeper, valAddr sdk.ValAddress, expectError error) {
	h := NewHandler(dk)
	msg := NewMsgWithdrawValidatorCommission(valAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

func DoWithdrawDelegatorRewardWithError(t *testing.T, ctx sdk.Context, dk Keeper, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress, expectError error) {
	h := NewHandler(dk)
	msg := NewMsgWithdrawDelegatorReward(delAddr, valAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}

//// CreateTestInputAdvanced hogpodge of all sorts of input required for testing
//func CreateTestInputAdvancedA(t *testing.T, isCheckTx bool, initPower int64, communityTax sdk.Dec) (
//	sdk.Context, auth.AccountKeeper, bank.Keeper, Keeper, staking.Keeper, params.Keeper, types.SupplyKeeper) {
//
//	//initTokens := sdk.TokensFromConsensusPower(initPower)
//
//	keyDistr := sdk.NewKVStoreKey(types.StoreKey)
//	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
//	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
//	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
//	keyMpt := sdk.NewKVStoreKey(mpt.StoreKey)
//	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
//	keyGov := sdk.NewKVStoreKey(gov.StoreKey)
//	keyParams := sdk.NewKVStoreKey(params.StoreKey)
//	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
//
//	db := dbm.NewMemDB()
//	ms := store.NewCommitMultiStore(db)
//
//	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
//	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(keyMpt, sdk.StoreTypeMPT, db)
//	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
//	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
//
//	err := ms.LoadLatestVersion()
//	require.Nil(t, err)
//
//	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
//	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
//	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
//
//	cdc := keeper.MakeTestCodec()
//	reg := types2.NewInterfaceRegistry()
//	cc := codec.NewProtoCodec(reg)
//	pro := codec.NewCodecProxy(cc, cdc)
//
//	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
//
//	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
//	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, keyMpt, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
//	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), nil)
//	maccPerms := map[string][]string{
//		auth.FeeCollectorName:     nil,
//		types.ModuleName:          nil,
//		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
//		staking.BondedPoolName:    {supply.Burner, supply.Staking},
//		gov.ModuleName:            {supply.Burner, supply.Staking},
//	}
//	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
//
//	sk := staking.NewKeeper(pro, keyStaking, supplyKeeper,
//		pk.Subspace(staking.DefaultParamspace))
//	sk.SetParams(ctx, staking.DefaultParams())
//
//	keeper := NewKeeper(cdc, keyDistr, pk.Subspace(types.DefaultParamspace), sk, supplyKeeper,
//		auth.FeeCollectorName, nil)
//
//	keeper.SetWithdrawAddrEnabled(ctx, true)
//	//initCoins := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens))
//	//totalSupply := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), initTokens.MulRaw(int64(len(keeper.TestAddrs)))))
//	//supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))
//
//	// fill all the addresses with some coins, set the loose pool tokens simultaneously
//	//for _, addr := range keeper.TestAddrs {
//	//	_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
//	//	require.Nil(t, err)
//	//}
//
//	//set module accounts
//	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
//	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//	supplyKeeper.SetModuleAccount(ctx, bondPool)
//	supplyKeeper.SetModuleAccount(ctx, supply.NewEmptyModuleAccount(types.ModuleName))
//	supplyKeeper.SetModuleAccount(ctx, supply.NewEmptyModuleAccount(gov.ModuleName))
//
//	// set the distribution hooks on staking
//	sk.SetHooks(keeper.Hooks())
//
//	// set genesis items required for distribution
//	keeper.SetFeePool(ctx, types.InitialFeePool())
//	keeper.SetCommunityTax(ctx, communityTax)
//
//	govRouter := gov.NewRouter().
//		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
//		AddRoute(distributiontypes.RouterKey, NewDistributionProposalHandler(keeper))
//
//	govProposalHandlerRouter := govkeeper.NewProposalHandlerRouter()
//	govProposalHandlerRouter.AddRoute(params.RouterKey, &keeper)
//
//	keys := sdk.NewKVStoreKeys(
//		auth.StoreKey, staking.StoreKey,
//		supply.StoreKey, distr.StoreKey,
//		gov.StoreKey, params.StoreKey,
//		order.OrderStoreKey, ammswap.StoreKey, farm.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
//		ibchost.StoreKey,
//		erc20.StoreKey,
//		mpt.StoreKey,
//		types.StoreKey,
//	)
//	govKeeper := gov.NewKeeper(
//		cdc, keys[govtypes.StoreKey], pk, pk.Subspace(govtypes.DefaultParamspace),
//		supplyKeeper, &sk, gov.DefaultParamspace, govRouter,
//		bankKeeper, govProposalHandlerRouter, auth.FeeCollectorName,
//	)
//
//	keeper.SetGovKeeper(govKeeper)
//
//	return ctx, accountKeeper, bankKeeper, keeper, sk, pk, supplyKeeper
//}
