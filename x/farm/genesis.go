package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

// InitGenesis initialize default parameters and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	var yieldModuleAccHoldings sdk.DecCoins
	var moduleAccHoldings sdk.DecCoins

	for _, pool := range data.Pools {
		moduleAccHoldings = moduleAccHoldings.Add(sdk.DecCoins{pool.TotalValueLocked})
		moduleAccHoldings = moduleAccHoldings.Add(sdk.DecCoins{pool.DepositAmount})
		yieldModuleAccHoldings = yieldModuleAccHoldings.Add(pool.TotalAccumulatedRewards)
		k.SetFarmPool(ctx, pool)
	}

	for _, lockInfo := range data.LockInfos {
		k.SetLockInfo(ctx, lockInfo)
	}

	for _, historical := range data.PoolHistoricalRewards {
		k.SetPoolHistoricalRewards(ctx, historical.PoolName, historical.Period, historical.Rewards)
	}

	for _, current := range data.PoolCurrentRewards {
		k.SetPoolCurrentRewards(ctx, current.PoolName, current.Rewards)
	}

	k.SetParams(ctx, data.Params)

	moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(moduleAccHoldings); err != nil {
			panic(err)
		}
		k.SupplyKeeper().SetModuleAccount(ctx, moduleAcc)
	}

	yieldMoudleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.YieldFarmingAccount)
	if yieldMoudleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.YieldFarmingAccount))
	}
	if yieldMoudleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(yieldModuleAccHoldings); err != nil {
			panic(err)
		}
		k.SupplyKeeper().SetModuleAccount(ctx, yieldMoudleAcc)
	}

	mintModuleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.MintFarmingAccount)
	if mintModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.MintFarmingAccount))
	}


}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data types.GenesisState) {
	pools := k.GetFarmPools(ctx)

	lockInfos := make([]types.LockInfo, 0)
	k.IterateAllLockInfos(ctx,
		func(lockInfo types.LockInfo) (stop bool) {
			lockInfos = append(lockInfos, lockInfo)
			return false
		},
	)

	allHistoricalRewards := make([]types.PoolHistoricalRewardsRecord, 0)
	k.IterateAllPoolHistoricalRewards(ctx,
		func(poolName string, period uint64, rewards types.PoolHistoricalRewards) (stop bool) {
			allHistoricalRewards = append(allHistoricalRewards, types.PoolHistoricalRewardsRecord{
				PoolName: poolName,
				Period:   period,
				Rewards:  rewards,
			})
			return false
		},
	)

	allCurRewards := make([]types.PoolCurrentRewardsRecord, 0)
	k.IterateAllPoolCurrentRewards(ctx,
		func(poolName string, rewards types.PoolCurrentRewards) (stop bool) {
			allCurRewards = append(allCurRewards, types.PoolCurrentRewardsRecord{
				PoolName: poolName,
				Rewards:  rewards,
			})
			return false
		},
	)

	params := k.GetParams(ctx)

	return types.NewGenesisState(pools, lockInfos, allHistoricalRewards, allCurRewards, params)
}
