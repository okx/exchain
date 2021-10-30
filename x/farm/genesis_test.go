package farm

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
	"github.com/stretchr/testify/require"
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
			Period:   1,
			Rewards:  types.PoolHistoricalRewards{},
		},
	}
	defaultGenesisState.LockInfos = []types.LockInfo{
		{
			Owner:            poolMsg.Owner,
			PoolName:         poolMsg.PoolName,
			Amount:           sdk.NewDecCoinFromDec(poolMsg.MinLockAmount.Denom, sdk.NewDec(1)),
			StartBlockHeight: 10,
			ReferencePeriod:  1,
		},
	}
	defaultGenesisState.PoolCurrentRewards = []types.PoolCurrentRewardsRecord{
		{
			PoolName: poolMsg.PoolName,
			Rewards:  types.PoolCurrentRewards{},
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
