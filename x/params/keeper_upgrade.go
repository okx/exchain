package params

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/params/types"
)

// ClaimReadyForUpgrade tells Keeper that someone has get ready for the upgrade.
// cb could be nil if there's no code to be execute when the upgrade is take effective.
// NOTE: This method could only be called at initialize phase, and
// CAN NOT be called when hanlding a tx.
func (keeper *Keeper) ClaimReadyForUpgrade(name string, cb func(types.UpgradeInfo)) {
	keeper.upgradeCache.ClaimReadyForUpgrade(name, cb)
}

func (keeper *Keeper) IsUpgradeEffective(ctx sdk.Context, name string) bool {
	b, err := keeper.upgradeCache.IsUpgradeEffective(ctx, name)
	if err != nil {
		return false
	}
	return b
}

func (keeper *Keeper) IsUpgradeEffective2(store store.KVStore, name string) bool {
	b, err := keeper.upgradeCache.IsUpgradeEffective2(store, name)
	if err != nil {
		return false
	}
	return b
}

func (keeper *Keeper) queryReadyForUpgrade(name string) ([]func(types.UpgradeInfo), bool) {
	return keeper.upgradeCache.QueryReadyForUpgrade(name)
}

func (keeper *Keeper) readUpgradeInfoFromStore(ctx sdk.Context, name string) (types.UpgradeInfo, error) {
	return keeper.upgradeCache.ReadUpgradeInfoFromStore(ctx, name)
}

func (keeper Keeper) iterateAllUpgradeInfo(ctx sdk.Context, cb func(info types.UpgradeInfo) (stop bool)) sdk.Error {
	return keeper.upgradeCache.IterateAllUpgradeInfo(ctx, cb)
}

func (keeper *Keeper) writeUpgradeInfo(ctx sdk.Context, info types.UpgradeInfo, forceCover bool) sdk.Error {
	return keeper.upgradeCache.WriteUpgradeInfo(ctx, info, forceCover)
}

func (keeper *Keeper) isUpgradeExist(ctx sdk.Context, name string) bool {
	return keeper.upgradeCache.IsUpgradeExist(ctx, name)
}
