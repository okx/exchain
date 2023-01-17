package params

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/params/types"
)

var (
	upgradeInfoPreifx = []byte("upgrade")
)

// ClaimReadyForUpgrade tells Keeper that someone has get ready for the upgrade.
// cb could be nil if there's no code to be execute when the upgrade is take effective.
// NOTE: for the same name, if this method is called more than once time, the callback of
// latest call will cover all of previous ones.
func (keeper *Keeper) ClaimReadyForUpgrade(name string, cb func(types.UpgradeInfo)) {
	if _, ok := keeper.upgradeReadyMap[name]; ok {
		keeper.logger.Error("more than one guys ready for the same upgrade, the front one will be cover", "upgrade name", name)
	}
	keeper.upgradeReadyMap[name] = cb
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

func (keeper *Keeper) queryReadyForUpgrade(name string) (func(types.UpgradeInfo), bool) {
	cb, ok := keeper.upgradeReadyMap[name]
	return cb, ok
}

func (keeper *Keeper) readUpgradeInfo(ctx sdk.Context, name string) (types.UpgradeInfo, error) {
	info, exist := keeper.readUpgradeInfoFromCache(name)
	if exist {
		return info, nil
	}

	info, err := keeper.readUpgradeInfoFromStore(ctx, name)
	if err != nil {
		return info, err
	}

	keeper.writeUpgradeInfoToCache(info)
	return info, nil
}

func (keeper Keeper) iterateAllUpgradeInfo(ctx sdk.Context, cb func(info types.UpgradeInfo) (stop bool)) sdk.Error {
	store := keeper.getUpgradeStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		data := iterator.Value()

		var info types.UpgradeInfo
		keeper.cdc.MustUnmarshalJSON(data, &info)

		if stop := cb(info); stop {
			break
		}
	}

	return nil
}

func (keeper *Keeper) writeUpgradeInfo(ctx sdk.Context, info types.UpgradeInfo, forceCover bool) sdk.Error {
	if err := keeper.writeUpgradeInfoToStore(ctx, info, forceCover); err != nil {
		return err
	}

	keeper.writeUpgradeInfoToCache(info)
	return nil
}

func (keeper *Keeper) readUpgradeInfoFromStore(ctx sdk.Context, name string) (types.UpgradeInfo, sdk.Error) {
	store := keeper.getUpgradeStore(ctx)

	data := store.Get([]byte(name))
	if len(data) == 0 {
		err := fmt.Errorf("upgrade '%s' is not exist", name)
		keeper.Logger(ctx).Info(err.Error())
		return types.UpgradeInfo{}, err
	}

	var info types.UpgradeInfo
	keeper.cdc.MustUnmarshalJSON(data, &info)
	return info, nil
}

func (keeper *Keeper) isUpgradeExist(ctx sdk.Context, name string) bool {
	store := keeper.getUpgradeStore(ctx)
	return store.Has([]byte(name))
}

func (keeper *Keeper) writeUpgradeInfoToStore(ctx sdk.Context, info types.UpgradeInfo, forceCover bool) sdk.Error {
	key := []byte(info.Name)

	store := keeper.getUpgradeStore(ctx)
	if !forceCover && store.Has(key) {
		keeper.Logger(ctx).Error("upgrade proposal name has been exist", "proposal name", info.Name)
		return sdk.ErrInternal(fmt.Sprintf("upgrade proposal name '%s' has been exist", info.Name))
	}

	data := keeper.cdc.MustMarshalJSON(info)
	store.Set(key, data)

	return nil
}

func (keeper *Keeper) getUpgradeStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(keeper.storeKey), upgradeInfoPreifx)
}

func (keeper *Keeper) readUpgradeInfoFromCache(name string) (types.UpgradeInfo, bool) {
	info, ok := keeper.upgradeInfoCache[name]
	return info, ok
}

func (keeper *Keeper) writeUpgradeInfoToCache(info types.UpgradeInfo) {
	keeper.upgradeInfoCache[info.Name] = info
}

func isUpgradeEffective(ctx sdk.Context, info types.UpgradeInfo) bool {
	return info.Status == types.UpgradeStatusEffective && uint64(ctx.BlockHeight()) >= info.EffectiveHeight
}
