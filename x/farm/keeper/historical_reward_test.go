package keeper

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPoolCurrentReward(t *testing.T) {
	ctx, k := GetKeeper(t)
	cdc := codec.New()

	poolNames := []string{"pool1", "pool2"}
	for _, poolName := range poolNames {
		poolCur := types.NewPoolCurrentRewards(100, 3, sdk.DecCoins{})
		k.Keeper.SetPoolCurrentRewards(ctx, poolName, poolCur)
		poolHis1 := types.NewPoolHistoricalRewards(sdk.DecCoins{}, 1)
		k.Keeper.SetPoolHistoricalRewards(ctx, poolName, 1, poolHis1)
		poolHis2 := types.NewPoolHistoricalRewards(sdk.DecCoins{}, 1)
		k.Keeper.SetPoolHistoricalRewards(ctx, poolName, 2, poolHis2)
	}

	require.True(t, k.Keeper.HasPoolCurrentRewards(ctx, poolNames[0]))
	require.True(t, k.Keeper.HasPoolCurrentRewards(ctx, poolNames[1]))

	var curs []types.PoolCurrentRewards
	k.Keeper.IterateAllPoolCurrentRewards(ctx, func(poolName string, rewards types.PoolCurrentRewards) bool {
		curs = append(curs, rewards)
		return false
	})
	require.Equal(t, len(poolNames), len(curs))

	var historicals []types.PoolHistoricalRewards
	k.Keeper.IteratePoolHistoricalRewards(
		ctx, poolNames[0],
		func(store sdk.KVStore, key []byte, value []byte) (stop bool) {
			var rewards types.PoolHistoricalRewards
			cdc.MustUnmarshalBinaryLengthPrefixed(value, &rewards)
			historicals = append(historicals, rewards)
			return false
		},
	)
	require.Equal(t, 2, len(historicals))

	historicals = make([]types.PoolHistoricalRewards, 0)
	k.Keeper.IterateAllPoolHistoricalRewards(
		ctx,
		func(poolName string, period uint64, rewards types.PoolHistoricalRewards) (stop bool) {
			historicals = append(historicals, rewards)
			return false
		},
	)
	require.Equal(t, 4, len(historicals))

	k.Keeper.DeletePoolCurrentRewards(ctx, poolNames[0])
	k.Keeper.DeletePoolCurrentRewards(ctx, poolNames[1])
	require.False(t, k.Keeper.HasPoolCurrentRewards(ctx, poolNames[0]))
	require.False(t, k.Keeper.HasPoolCurrentRewards(ctx, poolNames[1]))
}

func TestGetPoolHistoricalRewardsPoolNamePeriod(t *testing.T) {
	period := uint64(10)
	poolName := common.GetFixedLengthRandomString(120)
	key := types.GetPoolHistoricalRewardsKey(poolName, period)
	require.Panics(t, func() { GetPoolHistoricalRewardsPoolNamePeriod(key) })

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(period))
	key = append(types.PoolHistoricalRewardsPrefix, append([]byte("pool"), b...)...)
	require.Panics(t, func() { GetPoolHistoricalRewardsPoolNamePeriod(key) })

	poolName = "pool"
	key = types.GetPoolHistoricalRewardsKey(poolName, period)
	retPoolName, retPeriod := GetPoolHistoricalRewardsPoolNamePeriod(key)
	require.Equal(t, poolName, retPoolName)
	require.Equal(t, period, retPeriod)
}
