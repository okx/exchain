package keeper

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okexchain/x/farm/types"
)

const custom = "custom"

func getQueriedParams(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier,
) (params types.Params) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryParameters}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{types.QueryParameters}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &params))

	return
}

func getQueriedPool(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, poolName string,
) (pool types.FarmPool) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryPool}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryPoolParams(poolName)),
	}

	bz, err := querier(ctx, []string{types.QueryPool}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &pool))
	return
}

func getQueriedPools(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier,
) (pools types.FarmPools) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryPools}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryPoolsParams(1, 10)),
	}

	bz, err := querier(ctx, []string{types.QueryPools}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &pools))
	return
}

func getQueriedEarnings(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, poolName string, addr sdk.AccAddress,
) (earnings types.Earnings) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryEarnings}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryPoolAccountParams(poolName, addr)),
	}

	bz, err := querier(ctx, []string{types.QueryEarnings}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &earnings))

	return
}

func getQueriedLockInfo(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier,
	poolName string, addr sdk.AccAddress,
) (lockInfo types.LockInfo) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryLockInfo}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryPoolAccountParams(poolName, addr)),
	}

	bz, err := querier(ctx, []string{types.QueryLockInfo}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &lockInfo))

	return
}

func getQueriedWhitelist(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier,
) (whiteList types.PoolNameList) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryWhitelist}, "/"),
	}

	bz, err := querier(ctx, []string{types.QueryWhitelist}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &whiteList))

	return
}

func getQueriedAccount(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, addr sdk.AccAddress,
) (poolNameList types.PoolNameList) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryAccount}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryAccountParams(addr)),
	}

	bz, err := querier(ctx, []string{types.QueryAccount}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(bz, &poolNameList))

	return
}

func getQueriedAccountsLockedTo(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, poolName string,
) (addrList types.AccAddrList) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryAccountsLockedTo}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryPoolParams(poolName)),
	}

	cp, err := querier(ctx, []string{types.QueryAccountsLockedTo}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(cp, &addrList))

	return
}

func getQueriedPoolNum(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier,
) (poolNum types.PoolNum) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryPoolNum}, "/"),
		Data: []byte{},
	}

	cp, err := querier(ctx, []string{types.QueryPoolNum}, query)
	require.Nil(t, err)
	require.Nil(t, cdc.UnmarshalJSON(cp, &poolNum))

	return
}

func TestQueries(t *testing.T) {
	cdc := codec.New()
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	ctx, mockKeeper := GetKeeper(t)
	querier := NewQuerier(mockKeeper.Keeper)

	pools, lockInfos := initPoolsAndLockInfos(t, ctx, mockKeeper)

	retParams := getQueriedParams(t, ctx, cdc, querier)
	require.Equal(t, types.DefaultParams(), retParams)

	retPool := getQueriedPool(t, ctx, cdc, querier, pools[0].Name)
	require.Equal(t, pools[0], retPool)

	retPools := getQueriedPools(t, ctx, cdc, querier)
	require.Equal(t, pools, retPools)

	retLockInfo := getQueriedLockInfo(t, ctx, cdc, querier, pools[0].Name, Addrs[0])
	require.Equal(t, lockInfos[0], retLockInfo)

	whiteList := getQueriedWhitelist(t, ctx, cdc, querier)
	require.Equal(t, 1, len(whiteList))
	require.Equal(t, pools[0].Name, whiteList[0])

	retLockedPools := getQueriedAccount(t, ctx, cdc, querier, Addrs[0])
	require.Equal(t, types.PoolNameList{pools[0].Name, pools[1].Name}, retLockedPools)

	retLockedAddrs := getQueriedAccountsLockedTo(t, ctx, cdc, querier, pools[0].Name)
	require.Equal(t, types.AccAddrList{Addrs[0], Addrs[1]}, retLockedAddrs)

	retPoolNum := getQueriedPoolNum(t, ctx, cdc, querier)
	require.Equal(t, len(pools), int(retPoolNum.Number))

	// test query earnings
	ctx = ctx.WithBlockHeight(120)
	retEarnings := getQueriedEarnings(t, ctx, cdc, querier, pools[0].Name, Addrs[0])
	yieldAmount := pools[0].YieldedTokenInfos[0].AmountYieldedPerBlock.
		MulInt64(ctx.BlockHeight() - pools[0].YieldedTokenInfos[0].StartBlockHeightToYield)
	cur := mockKeeper.Keeper.GetPoolCurrentRewards(ctx, pools[0].Name)
	cur.Rewards = cur.Rewards.Add(
		sdk.DecCoins{sdk.NewDecCoinFromDec(pools[0].YieldedTokenInfos[0].RemainingAmount.Denom, yieldAmount)},
	)
	referHis := mockKeeper.Keeper.GetPoolHistoricalRewards(ctx, pools[0].Name, lockInfos[0].ReferencePeriod)
	newRatio := referHis.CumulativeRewardRatio.Add(cur.Rewards.QuoDecTruncate(pools[0].TotalValueLocked.Amount))
	expectedAmount := newRatio.Sub(referHis.CumulativeRewardRatio).MulDecTruncate(lockInfos[0].Amount.Amount)
	require.Equal(t, expectedAmount, retEarnings.AmountYielded)

	// test not existed query path
	bz, err := querier(ctx, []string{"xxxx"}, abci.RequestQuery{})
	require.NotNil(t, err)
	require.Nil(t, bz)
}
