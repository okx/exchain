package distribution

import (
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/crypto"
	"github.com/okx/okbchain/x/distribution/keeper"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/staking"
	stakingtypes "github.com/okx/okbchain/x/staking/types"
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
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)
			},
			[4]sdk.Error{types.ErrCodeEmptyDelegationDistInfo(), types.ErrCodeEmptyDelegationDistInfo(), nil, nil},
		},
		{
			"set withdraw reward disable",
			func(ctx sdk.Context, dk Keeper) {
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)

				proposalWithdrawReward := types.NewWithdrawRewardEnabledProposal("title", "description", false)
				keeper.HandleWithdrawRewardEnabledProposal(ctx, dk, proposalWithdrawReward)
				require.Equal(suite.T(), false, dk.GetWithdrawRewardEnabled(ctx))
			},
			[4]sdk.Error{types.ErrCodeEmptyDelegationDistInfo(), types.ErrCodeDisabledWithdrawRewards(),
				stakingtypes.ErrCodeDisabledOperate(), types.ErrCodeDisabledWithdrawRewards()},
		},
		{
			"empty delegation",
			func(ctx sdk.Context, dk Keeper) {

			},
			[4]sdk.Error{types.ErrCodeEmptyDelegationDistInfo(), types.ErrCodeEmptyDelegationDistInfo(), nil, nil},
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

			msg2 := NewMsgWithdrawDelegatorAllRewards(delAddr1)
			_, err = handler(ctx, msg2)
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
				proposal := types.NewChangeDistributionTypeProposal("change distri type", "", types.DistributionTypeOnChain)
				keeper.HandleChangeDistributionTypeProposal(ctx, dk, proposal)
				require.Equal(suite.T(), dk.GetDistributionType(ctx), types.DistributionTypeOnChain)
			},
			[2]sdk.Error{types.ErrNoValidatorCommission(), nil},
		},
		{
			"no allocate tokens, change distribution type",
			func(ctx sdk.Context, ak auth.AccountKeeper, dk Keeper, sk staking.Keeper, supplyKeeper types.SupplyKeeper) {
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
		enable      bool
		expectError error
	}

	testCases := []struct {
		title   string
		execute func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int)
		params  []param
	}{
		{
			"create val",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoCreateValidatorWithError(suite.T(), *ctx, sk, keeper.TestValAddrs[0], nil, p.expectError)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, stakingtypes.ErrValidatorOwnerExists()},
			},
		},
		{
			"disable edit val",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoEditValidatorWithError(suite.T(), *ctx, sk, keeper.TestValAddrs[0], sdk.NewDec(0), p.expectError)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, stakingtypes.ErrCommissionUpdateTime()},
			},
		},
		{
			"disable destroy val",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoDestroyValidatorWithError(suite.T(), *ctx, sk, keeper.TestValAccAddrs[i], p.expectError)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, nil},
			},
		},
		{
			"disable withdraw DoAddShares",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoAddSharesWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[0], valOpAddrs, p.expectError)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, nil},
			},
		},
		{
			"disable withdraw DoRegProxy",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoRegProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i], true, p.expectError)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, nil},
			},
		},
		{
			"disable withdraw DoWithdraw",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i], sdk.NewCoin(sk.BondDenom(*ctx),
					sdk.NewInt(100)), p.expectError)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, nil},
			},
		},
		{
			"disable withdraw DoBindProxy",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoBindProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i+1], keeper.TestDelAddrs[0], p.expectError)
				dk.SetWithdrawRewardEnabled(*ctx, true)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, nil},
			},
		},
		{
			"disable withdraw DoUnBindProxy",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				DoBindProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i+1], keeper.TestDelAddrs[0], nil)
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoUnBindProxyWithError(suite.T(), *ctx, sk, keeper.TestDelAddrs[i+1], p.expectError)
				dk.SetWithdrawRewardEnabled(*ctx, true)
			},
			[]param{
				{false, stakingtypes.ErrCodeDisabledOperate()},
				{true, nil},
			},
		},
		{
			"disable withdraw address",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoSetWithdrawAddressWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[i], p.expectError)
			},
			[]param{
				{false, types.ErrCodeDisabledWithdrawRewards()},
				{true, nil},
			},
		},
		{
			"disable withdraw validator commission",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawValidatorCommissionWithError(suite.T(), *ctx, dk, keeper.TestValAddrs[0], p.expectError)
			},
			[]param{
				{false, types.ErrCodeDisabledWithdrawRewards()},
				{true, types.ErrNoValidatorCommission()},
			},
		},
		{
			"disable set withdraw address",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoSetWithdrawAddressWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[i], p.expectError)
			},
			[]param{
				{false, types.ErrCodeDisabledWithdrawRewards()},
				{true, nil},
			},
		},
		{
			"disable set withdraw delegator reward",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawDelegatorRewardWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[0], keeper.TestValAddrs[0], p.expectError)
			},
			[]param{
				{false, types.ErrCodeDisabledWithdrawRewards()},
				{true, nil},
			},
		},
		{
			"disable set withdraw delegator all reward",
			func(ctx *sdk.Context, dk Keeper, sk staking.Keeper, valOpAddrs []sdk.ValAddress, p param, i int) {
				dk.SetWithdrawRewardEnabled(*ctx, p.enable)
				DoWithdrawDelegatorAllRewardWithError(suite.T(), *ctx, dk, keeper.TestDelAddrs[0], p.expectError)
			},
			[]param{
				{false, types.ErrCodeDisabledWithdrawRewards()},
				{true, nil},
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
	// No reg proxy
	//h := staking.NewHandler(sk)
	//msg := staking.NewMsgRegProxy(delAddr, reg)
	//_, e := h(ctx, msg)
	//require.Equal(t, expectError, e)
}

func DoBindProxyWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, proxyAddr sdk.AccAddress, expectError error) {
	// No reg proxy
	//h := staking.NewHandler(sk)
	//msg := staking.NewMsgBindProxy(delAddr, proxyAddr)
	//_, e := h(ctx, msg)
	//require.Equal(t, expectError, e)
}

func DoUnBindProxyWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, expectError error) {
	// No reg proxy
	//h := staking.NewHandler(sk)
	//msg := staking.NewMsgUnbindProxy(delAddr)
	//_, e := h(ctx, msg)
	//require.Equal(t, expectError, e)
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

func DoWithdrawDelegatorAllRewardWithError(t *testing.T, ctx sdk.Context, dk Keeper, delAddr sdk.AccAddress, expectError error) {
	h := NewHandler(dk)
	msg := NewMsgWithdrawDelegatorAllRewards(delAddr)
	_, e := h(ctx, msg)
	require.Equal(t, expectError, e)
}
