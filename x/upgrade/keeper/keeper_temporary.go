package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"
)

// only for unittest

// deprecated
func (k Keeper) SetAppUpgradeConfig(ctx sdk.Context, proposalID, version, upgradeHeight uint64, software string) sdk.Error {
	if _, found := k.GetAppUpgradeConfig(ctx); found {
		return sdk.ErrInternal("there is an app upgrade config existing within the protocolKeeper. Only one entry is permitted")
	}

	appUpgradeConfig := proto.NewAppUpgradeConfig(proposalID, proto.NewProtocolDefinition(version, software, upgradeHeight, sdk.NewDecWithPrec(7, 1)))
	k.protocolKeeper.SetUpgradeConfig(ctx, appUpgradeConfig)
	return nil
}

// deprecated
func (k Keeper) getVersionInfoSuccessResult(ctx sdk.Context, version uint64) (proposalID uint64) {
	kvStore := ctx.KVStore(k.storeKey)
	bytes := kvStore.Get(GetSuccessVersionKey(version))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &proposalID)
	return
}

// deprecated
func (k Keeper) getVersionInfoFailResult(ctx sdk.Context, version uint64, proposalID uint64) (proposalIDRet uint64) {
	kvStore := ctx.KVStore(k.storeKey)
	bytes := kvStore.Get(GetFailedVersionKey(version, proposalID))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &proposalIDRet)
	return
}
