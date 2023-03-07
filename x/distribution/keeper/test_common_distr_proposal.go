package keeper

import (
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/crypto"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func DoCreateValidator(t *testing.T, ctx sdk.Context, sk staking.Keeper, valAddr sdk.ValAddress, valConsPk crypto.PubKey) {
	sh := staking.NewHandler(sk)
	msg := staking.NewMsgCreateValidator(valAddr, valConsPk, staking.Description{}, NewTestSysCoin(1, 0))
	res, err := sh(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func DoEditValidator(t *testing.T, ctx sdk.Context, sk staking.Keeper, valAddr sdk.ValAddress, newRate sdk.Dec) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgEditValidatorCommissionRate(valAddr, newRate)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoWithdraw(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, amount sdk.SysCoin) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgWithdraw(delAddr, amount)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoDestroyValidator(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgDestroyValidator(delAddr)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoDeposit(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, amount sdk.SysCoin) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgDeposit(delAddr, amount)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoDepositWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, amount sdk.SysCoin, err error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgDeposit(delAddr, amount)
	_, e := h(ctx, msg)
	require.Equal(t, err, e)
}

func DoAddShares(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgAddShares(delAddr, valAddrs)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoAddSharesWithError(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress, err error) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgAddShares(delAddr, valAddrs)
	_, e := h(ctx, msg)
	require.Equal(t, err, e)
}

func DoRegProxy(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, reg bool) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgRegProxy(delAddr, reg)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoBindProxy(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress, proxyAddr sdk.AccAddress) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgBindProxy(delAddr, proxyAddr)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func DoUnBindProxy(t *testing.T, ctx sdk.Context, sk staking.Keeper, delAddr sdk.AccAddress) {
	h := staking.NewHandler(sk)
	msg := staking.NewMsgUnbindProxy(delAddr)
	_, e := h(ctx, msg)
	require.Nil(t, e)
}

func GetQueriedDelegationRewards(t *testing.T, ctx sdk.Context, querier sdk.Querier,
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) (rewards sdk.DecCoins) {
	bz, err := amino.MarshalJSON(types.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr))
	require.NoError(t, err)

	ctx, _ = ctx.CacheContext()
	result, err := querier(ctx, []string{types.QueryDelegationRewards}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	err = amino.UnmarshalJSON(result, &rewards)
	require.NoError(t, err)

	return rewards
}

func GetQueriedDelegationTotalRewards(t *testing.T, ctx sdk.Context, querier sdk.Querier,
	delegatorAddr sdk.AccAddress) types.QueryDelegatorTotalRewardsResponse {

	params := types.NewQueryDelegatorParams(delegatorAddr)
	bz, err := amino.MarshalJSON(params)
	require.NoError(t, err)

	ctx, _ = ctx.CacheContext()
	result, err := querier(ctx, []string{types.QueryDelegatorTotalRewards}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	var response types.QueryDelegatorTotalRewardsResponse
	err = amino.UnmarshalJSON(result, &response)
	require.NoError(t, err)

	return response
}
