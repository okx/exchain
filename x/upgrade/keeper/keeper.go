package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/params"

	//"github.com/okex/okchain/x/staking/exported"
	"github.com/okex/okchain/x/staking/exported"

	//stakingtypes "github.com/okex/okchain/x/staking/types"
	stakingtypes "github.com/okex/okchain/x/staking/types"
	"github.com/okex/okchain/x/upgrade/types"
)

// Keeper is the keeper struct of the upgrade store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// The ValidatorSet to get information about validators
	//protocolKeeper proto.ProtocolKeeper
	protocolKeeper ProtocolKeeper
	stakingKeeper  StakingKeeper
	bankKeeper     BankKeeper
	paramSpace     params.Subspace
}

// NewKeeper creates a new upgrade keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, pk ProtocolKeeper, sk StakingKeeper, ck BankKeeper,
	paramSpace params.Subspace) Keeper {
	return Keeper{
		key,
		cdc,
		pk,
		sk,
		ck,
		paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// GetAppUpgradeConfig gets app upgrade config
func (k Keeper) GetAppUpgradeConfig(ctx sdk.Context) (proto.AppUpgradeConfig, bool) {
	return k.protocolKeeper.GetUpgradeConfig(ctx)
}

// ClearUpgradeConfig clears upgrade config
func (k Keeper) ClearUpgradeConfig(ctx sdk.Context) {
	k.protocolKeeper.ClearUpgradeConfig(ctx)
}

// GetValidatorByConsAddr gets validator by its consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator,
	found bool) {
	return k.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
}

// SetCurrentVersion sets current version to store
func (k Keeper) SetCurrentVersion(ctx sdk.Context, currentVersion uint64) {
	k.protocolKeeper.SetCurrentVersion(ctx, currentVersion)
}

// SetLastFailedVersion sets last failed version to store
func (k Keeper) SetLastFailedVersion(ctx sdk.Context, lastFailedVersion uint64) {
	k.protocolKeeper.SetLastFailedVersion(ctx, lastFailedVersion)
}

// SetSignal sets signal for upgrade
func (k Keeper) SetSignal(ctx sdk.Context, protocol uint64, address string) {
	kvStore := ctx.KVStore(k.storeKey)
	kvStore.Set(types.GetSignalKey(protocol, address), k.cdc.MustMarshalBinaryLengthPrefixed(true))
}

// GetSignal gets signal
func (k Keeper) GetSignal(ctx sdk.Context, protocol uint64, address string) bool {
	kvStore := ctx.KVStore(k.storeKey)
	flagBytes := kvStore.Get(types.GetSignalKey(protocol, address))
	if flagBytes != nil {
		var flag bool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(flagBytes, &flag)
		return true
	}
	return false
}

// DeleteSignal removes signal
func (k Keeper) DeleteSignal(ctx sdk.Context, protocol uint64, address string) bool {
	if ok := k.GetSignal(ctx, protocol, address); ok {
		kvStore := ctx.KVStore(k.storeKey)
		kvStore.Delete(types.GetSignalKey(protocol, address))
		return true
	}
	return false
}

// ClearSignals cleans up signals
func (k Keeper) ClearSignals(ctx sdk.Context, protocol uint64) {
	kvStore := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(kvStore, types.GetSignalPrefixKey(protocol))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		kvStore.Delete(iterator.Key())
	}
}

// AddNewVersionInfo adds new version info
func (k Keeper) AddNewVersionInfo(ctx sdk.Context, versionInfo types.VersionInfo) {
	kvStore := ctx.KVStore(k.storeKey)

	versionInfoBytes := k.cdc.MustMarshalBinaryLengthPrefixed(versionInfo)
	kvStore.Set(types.GetProposalIDKey(versionInfo.UpgradeInfo.ProposalID), versionInfoBytes)
	proposalIDBytes := k.cdc.MustMarshalBinaryLengthPrefixed(versionInfo.UpgradeInfo.ProposalID)

	if versionInfo.Success {
		kvStore.Set(types.GetSuccessVersionKey(versionInfo.UpgradeInfo.ProtocolDef.Version), proposalIDBytes)
	} else {
		kvStore.Set(types.GetFailedVersionKey(versionInfo.UpgradeInfo.ProtocolDef.Version,
			versionInfo.UpgradeInfo.ProposalID), proposalIDBytes)
	}
}

// IterateBondedValidatorsByPower iterates bonded validators by power
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context,
	fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, fn)
}

// GetCurrentVersion gets current version
func (k Keeper) GetCurrentVersion(ctx sdk.Context) uint64 {
	return k.protocolKeeper.GetCurrentVersion(ctx)
}

// SetParams sets upgrade params to store
func (k Keeper) SetParams(ctx sdk.Context, params types.UpgradeParams) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.UpgradeParams) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// GetProtocolKeeper gets proto keeper
func (k Keeper) GetProtocolKeeper() ProtocolKeeper {
	return k.protocolKeeper
}
