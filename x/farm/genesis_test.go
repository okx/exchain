package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
	"testing"
)


func TestInitAndExportGenesis(t *testing.T) {
	defaultGenesisState := types.DefaultGenesisState()
	// init
	tCtx := initEnvironment(t)
	// create pool
	poolMsg := createPool(t, tCtx)
	farmPools := tCtx.k.GetFarmPools(tCtx.ctx)
	defaultGenesisState.Pools = farmPools
	defaultGenesisState.PoolHistoricalRewards = []types.PoolHistoricalRewardsRecord{
		{
			PoolName: poolMsg.PoolName,
			Period: 1,
			Rewards: types.PoolHistoricalRewards{},
		},
	}
	defaultGenesisState.LockInfos = []types.LockInfo{
		{
			Owner: poolMsg.Owner,
			PoolName: poolMsg.PoolName,
			Amount: sdk.NewDecCoinFromDec(poolMsg.MinLockedAmount.Denom, sdk.NewDec(1)),
			StartBlockHeight: 10,
			ReferencePeriod: 1,
		},
	}
	defaultGenesisState.PoolCurrentRewards = []types.PoolCurrentRewardsRecord{
		{
			PoolName: poolMsg.PoolName,
			Rewards: types.PoolCurrentRewards{},
		},
	}
	defaultGenesisState.WhiteList = types.PoolNameList{
		poolMsg.PoolName,
	}

	ctx, mk := keeper.GetKeeper(t)
	k := mk.Keeper
	InitGenesis(ctx, k, defaultGenesisState)
	exportedGenesis := ExportGenesis(ctx, k)
	require.Equal(t, defaultGenesisState, exportedGenesis)

}
