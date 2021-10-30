package keeper

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/supply/exported"
	"github.com/okex/exchain/x/staking/types"
)

// GetBondedPool returns the bonded tokens pool's module account
func (k Keeper) GetBondedPool(ctx sdk.Context) (bondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.BondedPoolName)
}

// GetNotBondedPool returns the not bonded tokens pool's module account
func (k Keeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.NotBondedPoolName)
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdk.SysCoin) {

	coins := tokens.ToCoins()
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// TotalBondedTokens total staking tokens supply which is bonded
// TODO:No usages found in project files,remove it later
func (k Keeper) TotalBondedTokens(ctx sdk.Context) sdk.Dec {
	bondedPool := k.GetBondedPool(ctx)
	return bondedPool.GetCoins().AmountOf(k.BondDenom(ctx))
}

// StakingTokenSupply staking tokens from the total supply
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdk.Dec {
	return k.supplyKeeper.GetSupplyByDenom(ctx, k.BondDenom(ctx))
}

// BondedRatio the fraction of the staking tokens which are currently bonded
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	bondedPool := k.GetBondedPool(ctx)

	stakeSupply := k.StakingTokenSupply(ctx)
	if stakeSupply.IsPositive() {
		return bondedPool.GetCoins().AmountOf(k.BondDenom(ctx)).Quo(stakeSupply)
	}
	return sdk.ZeroDec()
}
