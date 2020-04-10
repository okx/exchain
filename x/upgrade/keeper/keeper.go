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

// create a new upgrade keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, pk ProtocolKeeper, sk StakingKeeper, ck BankKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		key,
		cdc,
		pk,
		sk,
		ck,
		paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// get app upgrade config
func (k Keeper) GetAppUpgradeConfig(ctx sdk.Context) (proto.AppUpgradeConfig, bool) {
	return k.protocolKeeper.GetUpgradeConfig(ctx)
}

// clear upgrade config
func (k Keeper) ClearUpgradeConfig(ctx sdk.Context) {
	k.protocolKeeper.ClearUpgradeConfig(ctx)
}

// get validator by consAddr
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool) {
	return k.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
}

// set current version
func (k Keeper) SetCurrentVersion(ctx sdk.Context, currentVersion uint64) {
	k.protocolKeeper.SetCurrentVersion(ctx, currentVersion)
}

//set last failed version
func (k Keeper) SetLastFailedVersion(ctx sdk.Context, lastFailedVersion uint64) {
	k.protocolKeeper.SetLastFailedVersion(ctx, lastFailedVersion)
}

// set signal
func (k Keeper) SetSignal(ctx sdk.Context, protocol uint64, address string) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes, err := k.cdc.MarshalBinaryLengthPrefixed(true)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetSignalKey(protocol, address), cmsgBytes)
}

// get signal
func (k Keeper) GetSignal(ctx sdk.Context, protocol uint64, address string) bool {
	kvStore := ctx.KVStore(k.storeKey)
	flagBytes := kvStore.Get(GetSignalKey(protocol, address))
	if flagBytes != nil {
		var flag bool
		err := k.cdc.UnmarshalBinaryLengthPrefixed(flagBytes, &flag)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}

// remove signal
func (k Keeper) DeleteSignal(ctx sdk.Context, protocol uint64, address string) bool {
	if ok := k.GetSignal(ctx, protocol, address); ok {
		kvStore := ctx.KVStore(k.storeKey)
		kvStore.Delete(GetSignalKey(protocol, address))
		return true
	}
	return false
}

// cleanup signals
func (k Keeper) ClearSignals(ctx sdk.Context, protocol uint64) {
	kvStore := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(kvStore, GetSignalPrefixKey(protocol))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		kvStore.Delete(iterator.Key())
	}
}

// add new version info
func (k Keeper) AddNewVersionInfo(ctx sdk.Context, versionInfo types.VersionInfo) {
	kvStore := ctx.KVStore(k.storeKey)

	versionInfoBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetProposalIDKey(versionInfo.UpgradeInfo.ProposalID), versionInfoBytes)

	proposalIDBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo.UpgradeInfo.ProposalID)
	if err != nil {
		panic(err)
	}

	if versionInfo.Success {
		kvStore.Set(GetSuccessVersionKey(versionInfo.UpgradeInfo.ProtocolDef.Version), proposalIDBytes)
	} else {
		kvStore.Set(GetFailedVersionKey(versionInfo.UpgradeInfo.ProtocolDef.Version, versionInfo.UpgradeInfo.ProposalID), proposalIDBytes)
	}
}

// get iterate bonded validators by power
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, fn)
}

// get current version
func (k Keeper) GetCurrentVersion(ctx sdk.Context) uint64 {
	return k.protocolKeeper.GetCurrentVersion(ctx)
}

func (k Keeper) SetParams(ctx sdk.Context, params types.UpgradeParams) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// get inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.UpgradeParams) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// just 4 test
// getter
func (k Keeper) GetProtocolKeeper() ProtocolKeeper {
	return k.protocolKeeper
}
