package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

// InitGenesis initialize default parameters and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	var yieldModuleAccHoldings sdk.SysCoins
	var moduleAccHoldings sdk.SysCoins

	for _, pool := range data.Pools {
		moduleAccHoldings = moduleAccHoldings.Add2(sdk.SysCoins{pool.TotalValueLocked})
		moduleAccHoldings = moduleAccHoldings.Add2(sdk.SysCoins{pool.DepositAmount})
		yieldModuleAccHoldings = yieldModuleAccHoldings.Add2(pool.TotalAccumulatedRewards)
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

	for _, poolName := range data.WhiteList {
		k.SetWhitelist(ctx, poolName)
	}

	k.SetParams(ctx, data.Params)

	// init module account
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

	yieldModuleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.YieldFarmingAccount)
	if yieldModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.YieldFarmingAccount))
	}
	if yieldModuleAcc.GetCoins().IsZero() {
		if err := yieldModuleAcc.SetCoins(yieldModuleAccHoldings); err != nil {
			panic(err)
		}
		k.SupplyKeeper().SetModuleAccount(ctx, yieldModuleAcc)
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

	whiteList := k.GetWhitelist(ctx)

	params := k.GetParams(ctx)

	return types.NewGenesisState(pools, lockInfos, allHistoricalRewards, allCurRewards, whiteList, params)
}
