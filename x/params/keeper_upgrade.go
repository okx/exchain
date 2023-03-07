package params

import (
	"fmt"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/params/types"
)

// ClaimReadyForUpgrade tells Keeper that someone has get ready for the upgrade.
// cb could be nil if there's no code to be execute when the upgrade is take effective.
// NOTE: This method could only be called at initialize phase, and
// CAN NOT be called when hanlding a tx.
func (keeper *Keeper) ClaimReadyForUpgrade(name string, cb func(types.UpgradeInfo)) {
	keeper.upgradeCache.ClaimReadyForUpgrade(name, cb)
}

func (keeper *Keeper) IsUpgradeEffective(ctx sdk.Context, name string) bool {
	_, err := keeper.GetEffectiveUpgradeInfo(ctx, name)
	return err == nil
}

func (keeper *Keeper) GetEffectiveUpgradeInfo(ctx sdk.Context, name string) (types.UpgradeInfo, error) {
	info, err := keeper.readUpgradeInfo(ctx, name)
	if err != nil {
		return types.UpgradeInfo{}, err
	}

	if !isUpgradeEffective(ctx, info) {
		keeper.Logger(ctx).Debug("upgrade is not effective", "name", name)
		return types.UpgradeInfo{}, fmt.Errorf("upgrade '%s' is not effective", name)
	}

	keeper.Logger(ctx).Debug("upgrade is effective", "name", name)
	return info, nil
}

func (keeper *Keeper) queryReadyForUpgrade(name string) ([]func(types.UpgradeInfo), bool) {
	return keeper.upgradeCache.QueryReadyForUpgrade(name)
}

func (keeper *Keeper) readUpgradeInfo(ctx sdk.Context, name string) (types.UpgradeInfo, error) {
	return keeper.upgradeCache.ReadUpgradeInfo(ctx, name)
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

func isUpgradeEffective(ctx sdk.Context, info types.UpgradeInfo) bool {
	return info.Status == types.UpgradeStatusEffective && uint64(ctx.BlockHeight()) >= info.EffectiveHeight
}

func (keeper *Keeper) ApplyEffectiveUpgrade(ctx sdk.Context) error {
	return keeper.iterateAllUpgradeInfo(ctx, func(info types.UpgradeInfo) (stop bool) {
		if info.Status == types.UpgradeStatusEffective {
			if cbs, ready := keeper.queryReadyForUpgrade(info.Name); ready {
				for _, cb := range cbs {
					if cb != nil {
						cb(info)
					}
				}
			}
		}
		return false
	})
}
