package types

import (
	"fmt"
	"sync"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
)

var (
	upgradeInfoPreifx = []byte("upgrade")
	readyPrefix       = []byte("readyUpgrade")
)

type UpgradeCache struct {
	storeKey *sdk.KVStoreKey
	logger   log.Logger
	cdc      *codec.Codec

	readyLock        sync.Mutex
	upgradeReadyMap  map[string][]func(UpgradeInfo)
	infoLock         sync.Mutex
	upgradeInfoCache map[string]UpgradeInfo
}

func NewUpgreadeCache(storeKey *sdk.KVStoreKey, logger log.Logger, cdc *codec.Codec) *UpgradeCache {
	return &UpgradeCache{
		storeKey: storeKey,
		logger:   logger,
		cdc:      cdc,

		upgradeReadyMap:  make(map[string][]func(UpgradeInfo)),
		upgradeInfoCache: make(map[string]UpgradeInfo),
	}
}

func (uc *UpgradeCache) ReadUpgradeInfo(ctx sdk.Context, name string) (UpgradeInfo, error) {
	if ctx.UseParamCache() {
		info, exist := uc.readUpgradeInfo(name)
		if exist {
			return info, nil
		}
	}

	info, err := readUpgradeInfoFromStore(ctx, name, uc.storeKey, uc.cdc)
	if err != nil {
		return info, err
	}

	uc.writeUpgradeInfo(info)
	return info, nil
}

func (uc *UpgradeCache) ClaimReadyForUpgrade(name string, cb func(UpgradeInfo)) {
	uc.writeClaim(name, cb)
}

func (uc *UpgradeCache) QueryReadyForUpgrade(name string) ([]func(UpgradeInfo), bool) {
	return uc.readClaim(name)
}

func (uc *UpgradeCache) WriteUpgradeInfo(ctx sdk.Context, info UpgradeInfo, forceCover bool) sdk.Error {
	if err := writeUpgradeInfoToStore(ctx, info, forceCover, uc.storeKey, uc.cdc, uc.logger); err != nil {
		return err
	}

	// store is updated, remove the info from cache so
	// makeing ReadUpgradeInfo to re-read from store.
	uc.removeUpgradeInfo(info.Name)
	return nil
}

func (uc *UpgradeCache) IsUpgradeExist(ctx sdk.Context, name string) bool {
	store := getUpgradeStore(ctx, uc.storeKey)
	return store.Has([]byte(name))
}

func (uc *UpgradeCache) IterateAllUpgradeInfo(ctx sdk.Context, cb func(info UpgradeInfo) (stop bool)) sdk.Error {
	store := getUpgradeStore(ctx, uc.storeKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		data := iterator.Value()

		var info UpgradeInfo
		uc.cdc.MustUnmarshalJSON(data, &info)

		if stop := cb(info); stop {
			break
		}
	}

	return nil
}

func (uc *UpgradeCache) readUpgradeInfo(name string) (UpgradeInfo, bool) {
	uc.infoLock.Lock()
	defer uc.infoLock.Unlock()

	info, ok := uc.upgradeInfoCache[name]
	return info, ok
}

func (uc *UpgradeCache) removeUpgradeInfo(name string) {
	uc.infoLock.Lock()
	defer uc.infoLock.Unlock()

	delete(uc.upgradeInfoCache, name)
}

func (uc *UpgradeCache) writeUpgradeInfo(info UpgradeInfo) {
	uc.infoLock.Lock()
	defer uc.infoLock.Unlock()

	uc.upgradeInfoCache[info.Name] = info
}

func (uc *UpgradeCache) readClaim(name string) ([]func(UpgradeInfo), bool) {
	uc.readyLock.Lock()
	defer uc.readyLock.Unlock()

	cb, ok := uc.upgradeReadyMap[name]
	return cb, ok
}

func (uc *UpgradeCache) writeClaim(name string, cb func(UpgradeInfo)) {
	uc.readyLock.Lock()
	defer uc.readyLock.Unlock()

	readies, ok := uc.upgradeReadyMap[name]
	if !ok {
		uc.upgradeReadyMap[name] = []func(UpgradeInfo){cb}
	} else {
		uc.upgradeReadyMap[name] = append(readies, cb)
	}
}

func readUpgradeInfoFromStore(ctx sdk.Context, name string, skey *sdk.KVStoreKey, cdc *codec.Codec) (UpgradeInfo, sdk.Error) {
	store := getUpgradeStore(ctx, skey)

	data := store.Get([]byte(name))
	if len(data) == 0 {
		err := fmt.Errorf("upgrade '%s' is not exist", name)
		return UpgradeInfo{}, err
	}

	var info UpgradeInfo
	cdc.MustUnmarshalJSON(data, &info)
	return info, nil
}

func writeUpgradeInfoToStore(ctx sdk.Context, info UpgradeInfo, forceCover bool, skey *sdk.KVStoreKey, cdc *codec.Codec, logger log.Logger) sdk.Error {
	key := []byte(info.Name)

	store := getUpgradeStore(ctx, skey)
	if !forceCover && store.Has(key) {
		logger.Error("upgrade proposal name has been exist", "proposal name", info.Name)
		return sdk.ErrInternal(fmt.Sprintf("upgrade proposal name '%s' has been exist", info.Name))
	}

	data := cdc.MustMarshalJSON(info)
	store.Set(key, data)

	return nil
}

func getUpgradeStore(ctx sdk.Context, skey *sdk.KVStoreKey) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(skey), upgradeInfoPreifx)
}

func readReadyFromStore(ctx sdk.Context, name string, skey *sdk.KVStoreKey) ([]byte, bool) {
	store := getReadyStore(ctx, skey)
	data := store.Get([]byte(name))
	return data, len(data) != 0
}

func writeReadyToStore(ctx sdk.Context, name string, skey *sdk.KVStoreKey) {
	store := getReadyStore(ctx, skey)
	store.Set([]byte(name), []byte(name))
}

func getReadyStore(ctx sdk.Context, skey *sdk.KVStoreKey) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(skey), readyPrefix)
}
